#!/bin/bash

GOOS=windows GOARCH=amd64 go build -ldflags "-w -s" -trimpath -o bin/windows/client_x64.exe client/client.go
GOOS=windows GOARCH=386 go build -ldflags "-w -s" -trimpath -o bin/windows/tlsserver_x64.exe server/tlsserver.go
GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -trimpath -o bin/linux/client_x86.exe client/client.go
GOOS=linux GOARCH=386 go build -ldflags "-w -s" -trimpath -o bin/linux/tlsserver_x86.exe server/tlsserver.go