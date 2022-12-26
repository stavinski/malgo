GOARCH=amd64 go build -ldflags "-w -s" -trimpath -o bin\inject_x64.exe
GOARCH=386 go build -ldflags "-w -s" -trimpath -o bin\inject_x86.exe
