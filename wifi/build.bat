set GOARCH=amd64
go build -ldflags "-w -s" -trimpath -o bin\wifi_pwds_x64.exe
GOARCH=386
go build -ldflags "-w -s" -trimpath -o bin\wifi_pwds_x86.exe
