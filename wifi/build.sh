#!/bin/bash

GOOS=windows GOARCH=amd64 go build -ldflags "-w -s" -trimpath -o bin/wifi_pwds_x64.exe
GOOS=windows GOARCH=386 go build -ldflags "-w -s" -trimpath -o bin/wifi_pwds_x86.exe