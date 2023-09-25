set GOARCH=amd64
go build -ldflags "-w -s" -trimpath -o bin\samdump_x64.exe
set GOARCH=386
go build -ldflags "-w -s" -trimpath -o bin\samdump_x86.exe
