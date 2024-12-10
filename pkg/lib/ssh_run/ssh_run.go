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

func CreateNewSSHClient(meta SSHMeta) (*ssh.Client, error) {
	client, err := ssh.Dial("tcp", meta.Address, GetSSHConfig(meta.Username, meta.Password))
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %v", err)
	}
	return client, nil
}

func CreateClientFromClient(meta SSHMeta, client *ssh.Client) (*ssh.Client, error) {
	conn, err := client.Dial("tcp", meta.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %v", err)
	}

	clientConn, chans, reqs, err := ssh.NewClientConn(conn, meta.Address, GetSSHConfig(meta.Username, meta.Password))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to target: %v", err)
	}
	return ssh.NewClient(clientConn, chans, reqs), nil
}

func ConsumeSession(target SSHMeta, jumpers []SSHMeta, consumer func(*ssh.Session) (string, error)) (string, error) {
	client, err := (*ssh.Client)(nil), error(nil)
	for _, meta := range jumpers {
		if client == nil {
			client, err = CreateNewSSHClient(meta)
			if err != nil {
				return "", fmt.Errorf("failed to create new client: %v", err)
			}
			defer client.Close()
			continue
		}

		client, err = CreateClientFromClient(meta, client)
		if err != nil {
			return "", fmt.Errorf("failed to create new client: %v", err)
		}
		defer client.Close()
	}

	if client == nil {
		client, err = CreateNewSSHClient(target)
		if err != nil {
			return "", fmt.Errorf("failed to create new client: %v", err)
		}
		defer client.Close()
	} else {
		client, err = CreateClientFromClient(target, client)
		if err != nil {
			return "", fmt.Errorf("failed to create new client: %v", err)
		}
		defer client.Close()
	}

	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %v", err)
	}

	return consumer(session)
}
