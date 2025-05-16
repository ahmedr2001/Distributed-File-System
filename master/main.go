package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	pb "dfs/proto"

	"google.golang.org/grpc"
)

type FileRecord struct {
	Filename   string
	DataKeeper string
	Filepath   string
}

type DataKeeperInfo struct {
	ID           string
	IP           string
	Port         int32
	TransferPort string
	LastSeen     time.Time
	IsAlive      bool
}

type MasterServer struct {
	pb.UnimplementedMasterTrackerServer // like a default implementation for non-implemented methods

	mu          sync.RWMutex
	dataKeepers map[string]*DataKeeperInfo
	files       map[string][]*FileRecord

	replicationInterval time.Duration
	replicaCount        int
}

func NewMasterServer() *MasterServer {
	server := &MasterServer{
		dataKeepers:         make(map[string]*DataKeeperInfo),
		files:               make(map[string][]*FileRecord),
		replicationInterval: 10 * time.Second,
		replicaCount:        3,
	}

	go server.replicationChecker()
	go server.heartbeatChecker()

	return server
}

func (s *MasterServer) GetUploadNode(ctx context.Context, req *pb.GetUploadNodeRequest) (*pb.GetUploadNodeResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var aliveKeepers []*DataKeeperInfo
	for _, dk := range s.dataKeepers {
		if dk.IsAlive {
			aliveKeepers = append(aliveKeepers, dk)
		}
	}

	if len(aliveKeepers) == 0 {
		return nil, fmt.Errorf("no available data keepers")
	}

	selectedKeeper := aliveKeepers[rand.Intn(len(aliveKeepers))]

	return &pb.GetUploadNodeResponse{
		Ip:           selectedKeeper.IP,
		Port:         selectedKeeper.Port,
		TransferPort: selectedKeeper.TransferPort,
	}, nil
}

func (s *MasterServer) RegisterUpload(ctx context.Context, req *pb.RegisterUploadRequest) (*pb.RegisterUploadResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	fileRecord := &FileRecord{
		Filename:   req.Filename,
		DataKeeper: req.DataKeeperId,
		Filepath:   req.Filepath,
	}

	s.files[req.Filename] = append(s.files[req.Filename], fileRecord)

	log.Printf("Registered file upload: %s on data keeper %s at path %s", req.Filename, req.DataKeeperId, req.Filepath)

	return &pb.RegisterUploadResponse{
		Success: true,
		Message: "File registered successfully",
	}, nil
}

func (s *MasterServer) GetDownloadNodes(ctx context.Context, req *pb.GetDownloadNodesRequest) (*pb.GetDownloadNodesResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	fileRecords, exists := s.files[req.Filename]
	if !exists || len(fileRecords) == 0 {
		return nil, fmt.Errorf("file %s not found", req.Filename)
	}

	var nodes []*pb.DataKeeperInfo
	for _, fileRecord := range fileRecords {
		dataKeeper, exists := s.dataKeepers[fileRecord.DataKeeper]
		if exists && dataKeeper.IsAlive {
			nodes = append(nodes, &pb.DataKeeperInfo{
				Ip:           dataKeeper.IP,
				Port:         dataKeeper.Port,
				TransferPort: dataKeeper.TransferPort,
			})
		}
	}

	if len(nodes) == 0 {
		return nil, fmt.Errorf("no available nodes for file %s", req.Filename)
	}

	return &pb.GetDownloadNodesResponse{
		Nodes: nodes,
	}, nil
}

func (s *MasterServer) Heartbeat(ctx context.Context, req *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	dk, exists := s.dataKeepers[req.DataKeeperId]
	if !exists {
		dk = &DataKeeperInfo{
			ID:           req.DataKeeperId,
			IP:           req.Ip,
			Port:         req.Port,
			TransferPort: req.TransferPort,
		}
		s.dataKeepers[req.DataKeeperId] = dk
		log.Printf("New data keeper registered: %s at %s:%d", req.DataKeeperId, req.Ip, req.Port)
	}

	dk.LastSeen = time.Now()
	dk.IsAlive = true

	return &pb.HeartbeatResponse{
		Acknowledged: true,
	}, nil
}

func (s *MasterServer) NotifyReplication(ctx context.Context, req *pb.ReplicationRequest) (*pb.ReplicationResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sourceKeeper, exists := s.dataKeepers[req.SourceKeeperId]
	if !exists || !sourceKeeper.IsAlive {
		return nil, fmt.Errorf("source data keeper %s not available", req.SourceKeeperId)
	}

	destKeeper, exists := s.dataKeepers[req.DestKeeperId]
	if !exists || !destKeeper.IsAlive {
		return nil, fmt.Errorf("destination data keeper %s not available", req.DestKeeperId)
	}

	log.Printf("Notifying replication of %s from %s to %s", req.Filename, req.SourceKeeperId, req.DestKeeperId)

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", destKeeper.IP, destKeeper.Port), grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to destination data keeper: %v", err)
	}
	defer conn.Close()

	client := pb.NewDataKeeperClient(conn)
	replicateReq := &pb.ReplicateFileRequest{
		Filename:           req.Filename,
		Filepath:           req.Filepath,
		SourceIp:           sourceKeeper.IP,
		SourceTransferPort: sourceKeeper.TransferPort,
	}

	_, err = client.ReplicateFile(ctx, replicateReq)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate replication: %v", err)
	}

	return &pb.ReplicationResponse{
		Acknowledged: true,
	}, nil
}

func (s *MasterServer) replicationChecker() {
	ticker := time.NewTicker(s.replicationInterval)
	defer ticker.Stop()

	for range ticker.C {
		s.checkReplication()
	}
}

func (s *MasterServer) checkReplication() {
	s.mu.Lock()
	defer s.mu.Unlock()

	log.Println("Checking replication status...")

	for filename, fileRecords := range s.files {
		var aliveFileRecords []*FileRecord
		for _, fileRecord := range fileRecords {
			dk, exists := s.dataKeepers[fileRecord.DataKeeper]
			if exists && dk.IsAlive {
				aliveFileRecords = append(aliveFileRecords, fileRecord)
			}
		}

		if len(aliveFileRecords) < s.replicaCount {
			if len(aliveFileRecords) == 0 {
				log.Printf("Warning: File %s has no alive replicas", filename)
				continue
			}

			sourceFileRecord := aliveFileRecords[0]
			sourceKeeper := s.dataKeepers[sourceFileRecord.DataKeeper]

			var potentialDests []*DataKeeperInfo
			for id, dk := range s.dataKeepers {
				if !dk.IsAlive {
					continue
				}

				hasFile := false
				for _, fileRecord := range fileRecords {
					if fileRecord.DataKeeper == id {
						hasFile = true
						break
					}
				}

				if !hasFile {
					potentialDests = append(potentialDests, dk)
				}
			}

			for i := len(aliveFileRecords); i < s.replicaCount && i-len(aliveFileRecords) < len(potentialDests); i++ {
				destKeeper := potentialDests[i-len(aliveFileRecords)]

				go func(src, dest *DataKeeperInfo, fname, fpath string) {
					ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // can increase timeout to handle network ping error
					defer cancel()

					req := &pb.ReplicationRequest{
						SourceKeeperId: src.ID,
						DestKeeperId:   dest.ID,
						Filename:       fname,
						Filepath:       fpath,
					}

					_, err := s.NotifyReplication(ctx, req)
					if err != nil {
						log.Printf("Failed to initiate replication: %v", err)
					}
				}(sourceKeeper, destKeeper, filename, sourceFileRecord.Filepath)

				s.files[filename] = append(s.files[filename], &FileRecord{
					Filename:   filename,
					DataKeeper: destKeeper.ID,
					Filepath:   fmt.Sprintf("/storage/%s", filename),
				})

				log.Printf("Initiated replication of %s to %s", filename, destKeeper.ID)
			}
		}
	}
}

func (s *MasterServer) heartbeatChecker() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		s.checkHeartbeats()
	}
}

func (s *MasterServer) checkHeartbeats() {
	s.mu.Lock()
	defer s.mu.Unlock()

	threshold := time.Now().Add(-3 * time.Second)

	for id, dk := range s.dataKeepers {
		if dk.LastSeen.Before(threshold) {
			if dk.IsAlive {
				log.Printf("Data keeper %s is now offline", id)
				dk.IsAlive = false
			}
		}
	}
}

func main() {
	port := 50051
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterMasterTrackerServer(server, NewMasterServer())

	log.Printf("Master Tracker server started on port %d", port)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
