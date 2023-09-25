#!/bin/bash

GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -trimpath -o bin/linux/dump_x64 main.go imdsdump.go
GOOS=linux GOARCH=386 go build -ldflags "-w -s" -trimpath -o bin/linux/dump_x86 main.go imdsdump.go
GOOS=windows GOARCH=amd64 go build -ldflags "-w -s" -trimpath -o bin/win/dump_x64.exe main.go imdsdump.go
GOOS=windows GOARCH=386 go build -ldflags "-w -s" -trimpath -o bin/win/dump_x86.exe main.go imdsdump.go