// Simple TLS server to provide shell to a rshell client
package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %v [-port|-p] <cert> <key>", os.Args[0])
	os.Exit(1)
}

func main() {
	var (
		port     int
		certPath string
		keyPath  string
	)

	flag.IntVar(&port, "port", 4444, "Port to listen on")
	flag.IntVar(&port, "p", 4444, "Port to listen on")
	flag.Parse()

	if len(flag.Args()) < 2 {
		usage()
	}

	certPath, keyPath = flag.Arg(0), flag.Arg(1)
	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	keyPEM, err := os.ReadFile(keyPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	config := &tls.Config{Certificates: []tls.Certificate{cert}}
	conn, err := tls.Listen("tcp", fmt.Sprintf(":%v", port), config)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer conn.Close()
	fmt.Printf("[!] waiting for remote connection on %v\n", port)

	clientConn, err := conn.Accept()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer clientConn.Close()
	fmt.Printf("[+] connection from '%v'\n", clientConn.RemoteAddr())
	c := make(chan uint64)

	// Read from Reader and write to Writer until EOF
	copy := func(r io.ReadCloser, w io.WriteCloser) {
		defer func() {
			r.Close()
			w.Close()
		}()
		n, _ := io.Copy(w, r)
		c <- uint64(n)
	}

	go copy(clientConn, os.Stdout)
	go copy(os.Stdin, clientConn)

	p := <-c
	log.Printf("[%s]: Connection has been closed by remote peer, %d bytes has been received\n", clientConn.RemoteAddr(), p)
	p = <-c
	log.Printf("[%s]: Local peer has been stopped, %d bytes has been sent\n", clientConn.RemoteAddr(), p)
}
