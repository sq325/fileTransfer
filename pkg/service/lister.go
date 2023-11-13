package service

import "fmt"

type Lister interface {
	List(remoteIp, remoteUser, remotePasswd, remoteFilePath string) ([]string, error)
}

type lister struct{}

func NewLister() Lister {
	return &lister{}
}

func (l lister) List(remoteIp, remoteUser, remotePasswd, remoteFilePath string) ([]string, error) {
	sshClient, err := newSshClient(remoteIp, remoteUser, remotePasswd)
	if err != nil {
		return nil, err
	}
	defer sshClient.Close()

	sftpClient, err := newSftpClient(sshClient)
	if err != nil {
		return nil, err
	}

	paths, err := sftpClient.Glob(remoteFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to Glob %s: %w", remoteFilePath, err)
	}
	return paths, nil
}
