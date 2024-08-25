package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

func NewClientConfigForKey(keyFilePath, user string) (*ssh.ClientConfig, error) {
	key, err := os.ReadFile(keyFilePath)
	if err != nil {
		return nil, fmt.Errorf("unable to read private key: %w", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %w", err)
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			// use the PublicKeys method for remote authentication
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	return config, nil
}
