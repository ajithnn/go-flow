# go-flow

Pipe WorkFlow manager through watch folders.Designed for use with processing workflows involving movement between different folders for different stages.

## Install

go get github.com/ajithnn/go-flow

## Intro

1. Folder scanner package - Scans a given folder to look for stable files supports windows and linux.
2. Stable is defined as files not being written.
3. Returns the filepaths which can be process through outChannel
4. Takes a white list of sub folders under the scan path and only scans those sub folders
5. Passes the stable files to a process pipeline. 
6. Process pipeline needs to be defined as mentioned below.

## Configuration

1. Configure pipes.json file, to include pipelines and their concurrency.
2. Pipelines are defined as types inside components package. Define a type with a process method, process method defines the entire flow for the pipe.
3. Pipeline capacity is the number of parallel pipes running through go routines. 
4. Pipeline type common - If one of the common pipes is at full capacity the others need to wait.
5. Pipeline type separate - Each pipe will get dedicated capacity.
6. Configure in components/asset.go file the TypeMap resgistry after defining a pipe.

## Run

go run scan_folder.go -logtostderr=true -v=2 <Inbox Path> <Comma separated whitelist of folders>

Example:
go run scan_folder.go -logtostderr=true -v=2 "./Inbox/" "Media,Meta"

## More to do

1. Tests to be written for scanner
2. Use automatically created channels and allow parallel pipes to be used.

