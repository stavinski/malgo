package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

type empty struct{}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %v <port> <remote_addr>\n", os.Args[0])
	os.Exit(1)
}

func main() {

	var (
		inPort  int
		outHost string
	)

	if len(os.Args) < 3 {
		usage()
	}

	inPort, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "[!] cannot convert %v to port: %q\n", os.Args[1], err)
		return
	}
	outHost = os.Args[2]

	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", inPort))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer listener.Close()
	fmt.Printf("[*] listening on '%v'\n", inPort)
	in, err := listener.Accept()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[!] error from client connection: %q\n", err)
		return
	}
	defer in.Close()
	fmt.Printf("[+] received connection '%v'\n", in.RemoteAddr())
	out, err := net.Dial("tcp", outHost)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[!] cannot connect to %v to port: %q\n", outHost, err)
		return
	}
	defer out.Close()
	fmt.Println("[+] setting up outbound connection")
	fmt.Printf("[+] now proxying: '%v' <=> '%v'\n", in.LocalAddr(), out.RemoteAddr())

	end := make(chan empty)
	copy := func(src, dst net.Conn) {
		_, err := io.Copy(dst, src)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		end <- empty{}
	}

	// just continously copy from one to other and vice versa
	go copy(out, in)
	go copy(in, out)

	<-end
}
