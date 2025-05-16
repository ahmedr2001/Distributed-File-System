# Distributed File System

A simple distributed file system that supports reading and writing MP4 files while keeping files replicated for fault tolerance.

## Architecture

This distributed file system has two types of nodes:

1. **Master Tracker**: Manages the metadata and coordinates operations between clients and data keepers.
2. **Data Keeper**: Stores the actual file data and communicates with the master tracker via heartbeats.

## Run the code

1. Generate the gRPC code:
```
protoc --go_out=proto --go-grpc_out=proto proto/dfs.proto
```

2. Build the binaries:
```
go build -o bin/master.exe master/main.go
go build -o bin/datakeeper.exe datakeeper/main.go
go build -o bin/client.exe client/main.go
```

## Running the System

### 1. Start the Master Tracker

```
./bin/master
```

This will start the Master Tracker on port 50051.

### 2. Start multiple Data Keepers

```
./bin/datakeeper --id datakeeper1 --port 50052 --transfer-port 50053 --storage ./storage1
./bin/datakeeper --id datakeeper2 --port 50054 --transfer-port 50055 --storage ./storage2
./bin/datakeeper --id datakeeper3 --port 50056 --transfer-port 50057 --storage ./storage3
```

Each Data Keeper needs a unique ID, port, transfer port, and storage directory.

### 3. Upload a file

```
./bin/client --command upload --file /path/to/your/video.mp4
```

This will upload the file to one of the Data Keepers, which will then get replicated to other Data Keepers.

### 4. Download a file

```
./bin/client --command download --file video.mp4 --output ./downloads
```

This will download the file from one of the available Data Keepers.

For parallel download:

```
./bin/client --command download --file video.mp4 --output ./downloads --parallel
```
