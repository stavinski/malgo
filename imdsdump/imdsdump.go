// Dumps credentials from the IMDS
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Used for deserializing the JSON
type Credentials struct {
	AccessKeyId     string
	SecretAccessKey string
	Token           string
	Expiration      time.Time
}

// Dump the EC2 assigned role creds using the IMDS
//
// Context is provided to control timeouts or other reasons to cancel
func DumpCreds(ctx context.Context) (*Credentials, error) {
	creds := &Credentials{}
	host := "http://169.254.169.254"
	creds_path := "/latest/meta-data/identity-credentials/ec2/security-credentials/ec2-instance"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, host+creds_path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[%v] %v", resp.Status, resp.Body)
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&creds)
	if err != nil {
		return nil, err
	}
	return creds, nil
}
