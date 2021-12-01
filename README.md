# EdgeStats Client

The [EdgeStats](https://www.edgestats.io) client. May be used with public EdgeStats server, or private [edgestats-server](https://github.com/edgestats/edgestats-server) and [edgestats-webui](https://github.com/edgestats/edgestats-webui) setup.
>
> Check your Theta Edge Node uptime stats without needing to view the edge node client GUI/CLI.

## Basic Setup (connect to public server)

### Download and run client
Client downloads are available [here](https://github.com/edgestats/edgestats-client/releases).
Once the client is downloaded, simply unzip the client and run (double click) the executable file, or run the following commands from command line:
```shell
cd <path/to/edgestats-client>
./edgestats-client-<OS>-<ARCH>
# example: ./edgestats-client-windows-amd64
```

Note, double clicking the executable file will open a new command line terminal

## Advanced Setup (connect to private server)

### Setup EdgeStats server
Instructions for setting up an EdgeStats server available [here](https://github.com/edgestats/edgestats-server).

### Clone repository
Execute the following commands to clone this repository:

```shell
git clone https://github.com/edgestats/edgestats-client
cd edgestats-client
```

### Install dependencies
Execute the following command to install the dependencies:

```shell
go mod tidy
```

### Build client from source
```shell
GOOS=<OS> GOARCH=<ARCH> go build -ldflags "-X 'github.com/edgestats/edgestats-client/data.apiAddr=<http://127.0.0.1:port>' -X 'github.com/edgestats/edgestats-client/data.apiKey=<your-api-key>'" -o ./build/edgestats-client-<OS>-<ARCH> ./cmd/main.go
# example: GOOS=windows GOARCH=amd64 go build -ldflags "-X 'github.com/edgestats/edgestats-client/data.apiAddr=http://127.0.0.1:8000' -X 'github.com/edgestats/edgestats-client/data.apiKey=thetaverse'" -o ./build/edgestats-client-windows-amd64.exe ./cmd/main.go
```

### Set environment variables
The following environment variable is optional (required for linux!) and sets the location of the Theta Edge Node log file that the EdgeStats client watches:

```shell
export LOG_FILEPATH=<path/to/edge-node-logs/log.log>
# example: export LOG_FILEPATH=~/Library/Logs/Theta\ Edge\ Node/log.log
```

### Setup EdgeStats webui
Instructions for setting up an EdgeStats webui available [here](https://github.com/edgestats/edgestats-webui).

## FAQs

### How does it work?
The EdgeStats client works as follows:
> The client watches Theta Edge Node log file
> 
> Scans edge node logs as they are written
> 
> Filters edge node logs relevant to uptime
> 
> Sends uptime logs to the EdgeStats server
> 
> Web UI displays the edge node uptime stats

### Important to note!
> An edge node's uptime stats will not be available unless the EdgeStats client was run previously
> 
> If an edge node is running but the EdgeStats client is not, the uptime stats are not being collected by the EdgeStats server
> 
> If either the EdgeStats client or edge node is not running, uptime stats are not being collected by the EdgeStats server

## LICENSE
Copyright (c) EdgeStats Authors