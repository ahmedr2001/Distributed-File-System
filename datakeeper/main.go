package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	pb "dfs/proto"

	"google.golang.org/grpc"
)

type DataKeeperServer struct {
	pb.UnimplementedDataKeeperServer

	id           string
	ip           string
	port         int
	transferPort int
	storagePath  string

	masterIP   string
	masterPort int
}

func NewDataKeeperServer(id, ip string, port, transferPort int, storagePath, masterIP string, masterPort int) *DataKeeperServer {
	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		err := os.MkdirAll(storagePath, 0755)
		if err != nil {
			log.Fatalf("Failed to create storage directory: %v", err)
		}
	}

	server := &DataKeeperServer{
		id:           id,
		ip:           ip,
		port:         port,
		transferPort: transferPort,
		storagePath:  storagePath,
		masterIP:     masterIP,
		masterPort:   masterPort,
	}

	go server.heartbeatSender()
	go server.startFileTransferServer()

	return server
}

func (s *DataKeeperServer) InitiateUpload(ctx context.Context, req *pb.InitiateUploadRequest) (*pb.InitiateUploadResponse, error) {
	log.Printf("Preparing to receive file: %s", req.Filename)

	return &pb.InitiateUploadResponse{
		Ready:        true,
		TransferPort: fmt.Sprintf("%d", s.transferPort),
	}, nil
}

func (s *DataKeeperServer) ReplicateFile(ctx context.Context, req *pb.ReplicateFileRequest) (*pb.ReplicateFileResponse, error) {
	log.Printf("Replicating file %s from %s:%s", req.Filename, req.SourceIp, req.SourceTransferPort)

	go func() {
		filePath := filepath.Join(s.storagePath, req.Filename)

		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", req.SourceIp, req.SourceTransferPort))
		if err != nil {
			log.Printf("Failed to connect to source data keeper: %v", err)
			return
		}
		defer conn.Close()

		_, err = conn.Write([]byte(fmt.Sprintf("REPLICATE::%s\n", req.Filename)))
		if err != nil {
			log.Printf("Failed to send replication request: %v", err)
			return
		}

		file, err := os.Create(filePath)
		if err != nil {
			log.Printf("Failed to create file: %v", err)
			return
		}
		defer file.Close()

		bytesWritten, err := io.Copy(file, conn)
		if err != nil {
			log.Printf("Failed to copy file data: %v", err)
			return
		}

		fileInfo, err := file.Stat()
		if err != nil {
			log.Printf("Failed to stat file after writing: %v", err)
			return
		}

		if fileInfo.Size() == 0 {
			log.Printf("Warning: Replicated file %s has zero size", req.Filename)
		} else {
			log.Printf("File %s replicated successfully (%d bytes)", req.Filename, bytesWritten)
		}

		grpcConn, err := grpc.Dial(fmt.Sprintf("%s:%d", s.masterIP, s.masterPort), grpc.WithInsecure())
		if err != nil {
			log.Printf("Failed to connect to master tracker: %v", err)
			return
		}
		defer grpcConn.Close()

		client := pb.NewMasterTrackerClient(grpcConn)
		registerReq := &pb.RegisterUploadRequest{
			Filename:     req.Filename,
			DataKeeperId: s.id,
			Filepath:     filePath,
		}

		_, err = client.RegisterUpload(context.Background(), registerReq)
		if err != nil {
			log.Printf("Failed to register replicated file: %v", err)
			return
		}

		log.Printf("Replicated file registered with master tracker")
	}()

	return &pb.ReplicateFileResponse{
		Success: true,
		Message: "Replication initiated",
	}, nil
}

func (s *DataKeeperServer) heartbeatSender() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", s.masterIP, s.masterPort), grpc.WithInsecure())
		if err != nil {
			log.Printf("Failed to connect to master tracker: %v", err)
			continue
		}

		client := pb.NewMasterTrackerClient(conn)
		req := &pb.HeartbeatRequest{
			DataKeeperId: s.id,
			Ip:           s.ip,
			Port:         int32(s.port),
			TransferPort: fmt.Sprintf("%d", s.transferPort),
		}

		_, err = client.Heartbeat(context.Background(), req)
		if err != nil {
			log.Printf("Failed to send heartbeat: %v", err)
		}

		conn.Close()
	}
}

func (s *DataKeeperServer) startFileTransferServer() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.transferPort))
	if err != nil {
		log.Fatalf("Failed to start file transfer server: %v", err)
	}

	log.Printf("File transfer server started on port %d", s.transferPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		go s.handleFileTransfer(conn)
	}
}

func (s *DataKeeperServer) handleFileTransfer(conn net.Conn) {
	defer conn.Close()
	log.Printf("New file transfer connection established")

	reader := bufio.NewReader(conn)
	commandLine, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Error reading command: %v", err)
		return
	}

	command := strings.TrimSpace(commandLine)
	parts := strings.Split(command, "::")

	if len(parts) < 2 {
		log.Printf("Invalid command format: %s", command)
		return
	}

	operation := parts[0]
	filename := parts[1]

	log.Printf("Received %s operation for file: %s", operation, filename)

	if operation == "UPLOAD" {
		filePath := filepath.Join(s.storagePath, filename)
		file, err := os.Create(filePath)
		if err != nil {
			log.Printf("Failed to create file: %v", err)
			return
		}
		defer file.Close()

		bytesWritten, err := io.Copy(file, reader)
		if err != nil {
			log.Printf("Failed to copy file data: %v", err)
			return
		}

		fileInfo, err := file.Stat()
		if err != nil {
			log.Printf("Failed to stat file after writing: %v", err)
			return
		}

		if fileInfo.Size() == 0 {
			log.Printf("Warning: File %s has zero size after transfer", filename)
		} else {
			log.Printf("File %s received successfully (%d bytes)", filename, bytesWritten)
		}

		masterConn, err := grpc.Dial(fmt.Sprintf("%s:%d", s.masterIP, s.masterPort), grpc.WithInsecure())
		if err != nil {
			log.Printf("Failed to connect to master tracker: %v", err)
			return
		}
		defer masterConn.Close()

		client := pb.NewMasterTrackerClient(masterConn)
		req := &pb.RegisterUploadRequest{
			Filename:     filename,
			DataKeeperId: s.id,
			Filepath:     filePath,
		}

		_, err = client.RegisterUpload(context.Background(), req)
		if err != nil {
			log.Printf("Failed to register file: %v", err)
			return
		}

		log.Printf("File %s registered with master tracker", filename)
	} else if operation == "DOWNLOAD" || operation == "REPLICATE" {
		filePath := filepath.Join(s.storagePath, filename)
		file, err := os.Open(filePath)
		if err != nil {
			log.Printf("Failed to open file for %s: %v", operation, err)
			return
		}
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			log.Printf("Failed to get file info: %v", err)
			return
		}
		fileSize := fileInfo.Size()

		if fileSize == 0 {
			log.Printf("Warning: Sending empty file %s for %s operation", filename, operation)
		}

		bytesSent, err := io.Copy(conn, file)
		if err != nil {
			log.Printf("Failed to send file: %v", err)
			return
		}

		log.Printf("File %s sent successfully for %s operation (%d bytes)", filename, operation, bytesSent)
	} else if operation == "SIZE" {
		filePath := filepath.Join(s.storagePath, filename)
		file, err := os.Open(filePath)
		if err != nil {
			log.Printf("Failed to open file for size check: %v", err)
			conn.Write([]byte("-1\n"))
			return
		}
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			log.Printf("Failed to get file info: %v", err)
			conn.Write([]byte("-1\n"))
			return
		}

		fileSize := fileInfo.Size()
		_, err = conn.Write([]byte(fmt.Sprintf("%d\n", fileSize)))
		if err != nil {
			log.Printf("Failed to send file size: %v", err)
			return
		}

		log.Printf("Sent size information for file %s: %d bytes", filename, fileSize)
	} else if operation == "DOWNLOAD_RANGE" {
		if len(parts) < 4 {
			log.Printf("Invalid range format: %s", command)
			return
		}

		filePath := filepath.Join(s.storagePath, filename)
		start, err := parsePositionArg(parts[2])
		if err != nil {
			log.Printf("Invalid start position: %v", err)
			return
		}

		end, err := parsePositionArg(parts[3])
		if err != nil {
			log.Printf("Invalid end position: %v", err)
			return
		}

		file, err := os.Open(filePath)
		if err != nil {
			log.Printf("Failed to open file for range download: %v", err)
			return
		}
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			log.Printf("Failed to get file info: %v", err)
			return
		}
		totalSize := fileInfo.Size()

		if start < 0 || end > totalSize || start >= end {
			log.Printf("Invalid range [%d:%d] for file size %d", start, end, totalSize)
			return
		}

		_, err = file.Seek(start, 0)
		if err != nil {
			log.Printf("Failed to seek to position %d: %v", start, err)
			return
		}

		limitReader := io.LimitReader(file, end-start)

		bytesSent, err := io.Copy(conn, limitReader)
		if err != nil {
			log.Printf("Failed to send file range: %v", err)
			return
		}

		log.Printf("File range %s[%d:%d] sent successfully (%d bytes of %d total)",
			filename, start, end, bytesSent, totalSize)
	} else {
		log.Printf("Unknown operation: %s", operation)
	}
}

func parsePositionArg(posStr string) (int64, error) {
	var pos int64
	_, err := fmt.Sscanf(posStr, "%d", &pos)
	if err != nil {
		return 0, fmt.Errorf("invalid position: %v", err)
	}
	return pos, nil
}

func main() {
	id := flag.String("id", "datakeeper1", "Data Keeper ID")
	ip := flag.String("ip", "localhost", "Data Keeper IP address")
	port := flag.Int("port", 50052, "Data Keeper gRPC port")
	transferPort := flag.Int("transfer-port", 50053, "Data Keeper file transfer port")
	storagePath := flag.String("storage", "./storage", "Path for file storage")
	masterIP := flag.String("master-ip", "localhost", "Master Tracker IP address")
	masterPort := flag.Int("master-port", 50051, "Master Tracker port")

	flag.Parse()

	server := grpc.NewServer()
	dataKeeperServer := NewDataKeeperServer(*id, *ip, *port, *transferPort, *storagePath, *masterIP, *masterPort)
	pb.RegisterDataKeeperServer(server, dataKeeperServer)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Data Keeper server %s started on %s:%d with transfer port %d", *id, *ip, *port, *transferPort)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
