package ssh_run

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

type SSHMeta struct {
	Address  string `json:"address"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type BodyParams struct {
	Command string    `json:"command"`
	Target  SSHMeta   `json:"target"`
	Jumpers []SSHMeta `json:"jumpers"`
}

func GetSSHConfig(username, password string) *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
}

func ConsumeSession(target SSHMeta, jumpers []SSHMeta, consumer func(*ssh.Session) (string, error)) (string, error) {
	client := (*ssh.Client)(nil)
	for _, meta := range jumpers {
		if client == nil {
			_client, err := ssh.Dial("tcp", meta.Address, GetSSHConfig(meta.Username, meta.Password))
			if err != nil {
				return "", fmt.Errorf("failed to dial: %v", err)
			}
			if err != nil {
				return "", err
			}
			defer _client.Close()
			client = _client
			continue
		}

		conn, err := client.Dial("tcp", meta.Address)
		if err != nil {
			return "", fmt.Errorf("failed to dial: %v", err)
		}

		clientConn, chans, reqs, err := ssh.NewClientConn(conn, meta.Address, GetSSHConfig(meta.Username, meta.Password))
		if err != nil {
			return "", fmt.Errorf("failed to connect to target: %v", err)
		}
		client = ssh.NewClient(clientConn, chans, reqs)
		defer client.Close()
	}

	if client == nil {
		_client, err := ssh.Dial("tcp", target.Address, GetSSHConfig(target.Username, target.Password))
		if err != nil {
			return "", fmt.Errorf("failed to dial: %v", err)
		}
		defer _client.Close()
		client = _client
	}

	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %v", err)
	}

	return consumer(session)
}
