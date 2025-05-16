package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	pb "dfs/proto"

	"google.golang.org/grpc"
)

type Client struct {
	masterIP   string
	masterPort int
}

func NewClient(masterIP string, masterPort int) *Client {
	return &Client{
		masterIP:   masterIP,
		masterPort: masterPort,
	}
}

func (c *Client) UploadFile(filePath string) error {
	filename := filepath.Base(filePath)

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %v", err)
	}
	fileSize := fileInfo.Size()

	log.Printf("Uploading file %s (%d bytes)", filename, fileSize)

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", c.masterIP, c.masterPort), grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("failed to connect to master tracker: %v", err)
	}
	defer conn.Close()

	client := pb.NewMasterTrackerClient(conn)
	req := &pb.GetUploadNodeRequest{
		Filename: filename,
	}
	resp, err := client.GetUploadNode(context.Background(), req)
	if err != nil {
		return fmt.Errorf("failed to get upload node: %v", err)
	}

	log.Printf("Got upload node: %s:%d with transfer port %s", resp.Ip, resp.Port, resp.TransferPort)

	dataKeeperConn, err := grpc.Dial(fmt.Sprintf("%s:%d", resp.Ip, resp.Port), grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("failed to connect to data keeper: %v", err)
	}
	defer dataKeeperConn.Close()

	dataKeeperClient := pb.NewDataKeeperClient(dataKeeperConn)
	initReq := &pb.InitiateUploadRequest{
		Filename: filename,
	}
	initResp, err := dataKeeperClient.InitiateUpload(context.Background(), initReq)
	if err != nil {
		return fmt.Errorf("failed to initiate upload: %v", err)
	}

	if !initResp.Ready {
		return fmt.Errorf("data keeper not ready for upload")
	}

	transferConn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", resp.Ip, resp.TransferPort))
	if err != nil {
		return fmt.Errorf("failed to connect to data keeper transfer port: %v", err)
	}
	defer transferConn.Close()

	uploadCmd := fmt.Sprintf("UPLOAD::%s\n", filename)
	_, err = transferConn.Write([]byte(uploadCmd))
	if err != nil {
		return fmt.Errorf("failed to send upload command: %v", err)
	}

	bytesWritten, err := io.Copy(transferConn, file)
	if err != nil {
		return fmt.Errorf("failed to upload file: %v", err)
	}

	log.Printf("File %s uploaded successfully (%d bytes)", filename, bytesWritten)

	time.Sleep(1 * time.Second)

	return nil
}

func (c *Client) DownloadFile(filename, outputPath string) error {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", c.masterIP, c.masterPort), grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("failed to connect to master tracker: %v", err)
	}
	defer conn.Close()

	client := pb.NewMasterTrackerClient(conn)

	req := &pb.GetDownloadNodesRequest{
		Filename: filename,
	}
	resp, err := client.GetDownloadNodes(context.Background(), req)
	if err != nil {
		return fmt.Errorf("failed to get download nodes: %v", err)
	}

	if len(resp.Nodes) == 0 {
		return fmt.Errorf("no nodes available for downloading file %s", filename)
	}

	log.Printf("Got %d download nodes for file %s", len(resp.Nodes), filename)

	node := resp.Nodes[rand.Intn(len(resp.Nodes))]

	transferConn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", node.Ip, node.TransferPort))
	if err != nil {
		return fmt.Errorf("failed to connect to data keeper transfer port: %v", err)
	}
	defer transferConn.Close()

	downloadCmd := fmt.Sprintf("DOWNLOAD::%s\n", filename)
	_, err = transferConn.Write([]byte(downloadCmd))
	if err != nil {
		return fmt.Errorf("failed to send download command: %v", err)
	}

	outFilePath := filepath.Join(outputPath, filename)
	outFile, err := os.Create(outFilePath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer outFile.Close()

	bytesRead, err := io.Copy(outFile, transferConn)
	if err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}

	fileInfo, err := outFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get downloaded file info: %v", err)
	}

	if fileInfo.Size() == 0 {
		log.Printf("Warning: Downloaded file %s has zero size", filename)
	}

	log.Printf("File %s downloaded successfully (%d bytes)", filename, bytesRead)

	return nil
}

func (c *Client) GetFileSize(filename string, node *pb.DataKeeperInfo) (int64, error) {
	transferConn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", node.Ip, node.TransferPort))
	if err != nil {
		return 0, fmt.Errorf("failed to connect to data keeper transfer port: %v", err)
	}
	defer transferConn.Close()

	sizeCmd := fmt.Sprintf("SIZE::%s\n", filename)
	_, err = transferConn.Write([]byte(sizeCmd))
	if err != nil {
		return 0, fmt.Errorf("failed to send size command: %v", err)
	}

	var fileSize int64
	reader := bufio.NewReader(transferConn)
	response, err := reader.ReadString('\n')
	if err != nil {
		return 0, fmt.Errorf("failed to read size response: %v", err)
	}

	_, err = fmt.Sscanf(strings.TrimSpace(response), "%d", &fileSize)
	if err != nil {
		return 0, fmt.Errorf("invalid size format: %v", err)
	}

	return fileSize, nil
}

func (c *Client) DownloadFileParallel(filename, outputPath string) error {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", c.masterIP, c.masterPort), grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("failed to connect to master tracker: %v", err)
	}
	defer conn.Close()

	client := pb.NewMasterTrackerClient(conn)

	req := &pb.GetDownloadNodesRequest{
		Filename: filename,
	}
	resp, err := client.GetDownloadNodes(context.Background(), req)
	if err != nil {
		return fmt.Errorf("failed to get download nodes: %v", err)
	}

	if len(resp.Nodes) == 0 {
		return fmt.Errorf("no nodes available for downloading file %s", filename)
	}

	log.Printf("Got %d download nodes for file %s", len(resp.Nodes), filename)

	var fileSize int64
	for _, node := range resp.Nodes {
		size, err := c.GetFileSize(filename, node)
		if err == nil {
			fileSize = size
			break
		}
		log.Printf("Failed to get file size from node %s: %v, trying next node", node.Ip, err)
	}

	if fileSize == 0 {
		log.Printf("Couldn't determine file size, falling back to regular download")
		return c.DownloadFile(filename, outputPath)
	}

	log.Printf("File size for %s: %d bytes", filename, fileSize)

	outFilePath := filepath.Join(outputPath, filename)
	outFile, err := os.Create(outFilePath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer outFile.Close()

	err = outFile.Truncate(fileSize) // allocate fileSize space
	if err != nil {
		log.Printf("Warning: Failed to pre-allocate file space: %v", err)
	}

	numNodes := len(resp.Nodes)
	numChunks := numNodes
	chunkSize := fileSize / int64(numChunks)
	lastChunkSize := chunkSize + (fileSize % int64(numChunks))

	var wg sync.WaitGroup
	errs := make(chan error, numChunks)

	for i := 0; i < numChunks; i++ {
		wg.Add(1)
		node := resp.Nodes[i]

		go func(chunkIndex int, node *pb.DataKeeperInfo) {
			defer wg.Done()

			start := int64(chunkIndex) * chunkSize
			end := start + chunkSize
			if chunkIndex == numChunks-1 {
				end = start + lastChunkSize
			}

			log.Printf("Downloading chunk %d from %s:%s (bytes %d-%d)",
				chunkIndex, node.Ip, node.TransferPort, start, end-1)

			transferConn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", node.Ip, node.TransferPort))
			if err != nil {
				errs <- fmt.Errorf("failed to connect to data keeper transfer port for chunk %d: %v",
					chunkIndex, err)
				return
			}
			defer transferConn.Close()

			rangeCmd := fmt.Sprintf("DOWNLOAD_RANGE::%s::%d::%d\n", filename, start, end)
			_, err = transferConn.Write([]byte(rangeCmd))
			if err != nil {
				errs <- fmt.Errorf("failed to send range download command for chunk %d: %v",
					chunkIndex, err)
				return
			}

			buffer := make([]byte, end-start)
			bytesRead, err := io.ReadAtLeast(transferConn, buffer, int(end-start))
			if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
				errs <- fmt.Errorf("failed to read chunk %d: %v", chunkIndex, err)
				return
			}

			_, err = outFile.WriteAt(buffer[:bytesRead], start)
			if err != nil {
				errs <- fmt.Errorf("failed to write chunk %d to file: %v", chunkIndex, err)
				return
			}

			log.Printf("Chunk %d downloaded successfully (%d bytes)", chunkIndex, bytesRead)
		}(i, node)
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			return err
		}
	}

	fileInfo, err := outFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get downloaded file info: %v", err)
	}

	if fileInfo.Size() == 0 {
		log.Printf("Warning: Parallel downloaded file %s has zero size", filename)
	} else if fileInfo.Size() != fileSize {
		log.Printf("Warning: Downloaded file size (%d) doesn't match expected size (%d)",
			fileInfo.Size(), fileSize)
	}

	log.Printf("File %s downloaded successfully in parallel (%d bytes)", filename, fileInfo.Size())
	return nil
}

func main() {
	masterIP := flag.String("master-ip", "localhost", "Master Tracker IP address")
	masterPort := flag.Int("master-port", 50051, "Master Tracker port")
	command := flag.String("command", "", "Command (upload or download)")
	filePath := flag.String("file", "", "File path for upload or filename for download")
	outputPath := flag.String("output", "./", "Output directory for downloads")
	parallel := flag.Bool("parallel", false, "Download in parallel (for download command)")

	flag.Parse()

	if *command == "" || *filePath == "" {
		flag.Usage()
		os.Exit(1)
	}

	client := NewClient(*masterIP, *masterPort)

	switch *command {
	case "upload":
		err := client.UploadFile(*filePath)
		if err != nil {
			log.Fatalf("Upload failed: %v", err)
		}
		log.Println("Upload completed successfully")

	case "download":
		var err error
		if *parallel {
			err = client.DownloadFileParallel(*filePath, *outputPath)
		} else {
			err = client.DownloadFile(*filePath, *outputPath)
		}
		if err != nil {
			log.Fatalf("Download failed: %v", err)
		}
		log.Println("Download completed successfully")

	default:
		log.Fatalf("Unknown command: %s", *command)
	}
}
