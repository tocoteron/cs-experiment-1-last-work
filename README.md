# The last work of CS Experiment 1

## Prerequisites

- Go (v1.14)

## Getting Started

### Install golang

Ubuntu:
```sh
sudo apt install golang
```

### Setup GOPATH

Set your GOPATH environment.
If you don't set GOPATH, it is set to `$HOME/go`(Unix) or `%USERPROFILE%\go`(Windows) automatically.

### Setup program

```sh
cd $GOPATH/src # If it doesn't exist, you should make the directory
go get github.com/tokoroten-lab/cs-experiment-1-last-work
```

## Usage

Move to `$GOPATH/src/cs-experiment-1/part-3/last-work/`.

Run(use default port):

```sh
go run server.go
```

Run(use 80 port)

```sh
go run server.go -port=80
```