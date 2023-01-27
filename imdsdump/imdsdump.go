// Dumps credentials from the IMDS
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// Used for deserializing the JSON
type Credentials struct {
	AccessKeyId     string
	SecretAccessKey string
	Token           string
	Expiration      time.Time
}

// Unlikely any of these will change to keep backward compatalbility
const (
	Host           = "http://169.254.169.254"
	CredsPath      = "/latest/meta-data/identity-credentials/ec2/security-credentials/ec2-instance"
	TokenPath      = "/latest/api/token"
	TokenTTLHeader = "X-aws-ec2-metadata-token-ttl-seconds"
	TokenTTL       = 21600
	TokenHeader    = "X-aws-ec2-metadata-token"
)

// Grab the IMDSv2 token
func getV2Token(ctx context.Context) (string, error) {
	buf := bytes.Buffer{}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, Host+TokenPath, nil)
	req.Header.Add(TokenTTLHeader, strconv.Itoa((TokenTTL)))
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// Dump the EC2 assigned role creds using the IMDS
//
// If a 401 status is returned it will assume IMDSv2 and will retrieve the token and retry the request
//
// Context is provided to control timeouts or other reasons to cancel
func DumpCreds(ctx context.Context) (*Credentials, error) {
	creds := &Credentials{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, Host+CredsPath, nil)
	if err != nil {
		return nil, err
	}
	// attempt direct call for IMDSv1
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	// using IMDSv2 which requires a token
	if resp.StatusCode == http.StatusUnauthorized {
		token, err := getV2Token(ctx)
		if err != nil {
			return nil, err
		}
		// add token to header and retry request
		req.Header.Add(TokenHeader, token)
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
	}
	// something else?!
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("[%v] %v", resp.Status, resp.Body)
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&creds)
	if err != nil {
		return nil, err
	}
	return creds, nil
}
