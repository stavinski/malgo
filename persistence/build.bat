set GOARCH=amd64
go build -ldflags "-w -s" -trimpath -o bin/schtaskpoc_x64.exe main_windows.go

set GOARCH=386
go build -ldflags "-w -s" -trimpath -o bin/schtaskpoc_x86.exe main_windows.go