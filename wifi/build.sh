#!/bin/bash

GOOS=windows GOARCH=amd64 go build -ldflags '-w -s' -trimpath -o wifi_pwds.exe 