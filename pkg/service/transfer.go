package service

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type Transfer interface {
	Get(ip, user, passwd, remoteFilePath, localPath string) error
	List(ip, user, passwd, remoteFilePath string) ([]string, error)
}

type transfer struct{}

func NewTransfer() Transfer {
	return &transfer{}
}

func (t transfer) Get(ip, user, passwd, remoteFilePath, localPath string) error {
	sftpClient, err := newSftpClient(ip, user, passwd)
	if err != nil {
		return err
	}

	paths, err := sftpClient.Glob(remoteFilePath)
	if err != nil {
		return fmt.Errorf("failed to Glob %s: %w", remoteFilePath, err)
	}

	localPath = strings.TrimSpace(localPath)
	for _, path := range paths {
		remoteFile, err := sftpClient.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open remote file: %w", err)
		}
		defer remoteFile.Close()

		var remoteFileName string
		{
			remoteFileInfo, err := remoteFile.Stat()
			if err != nil {
				return fmt.Errorf("failed to get remote file info: %w", err)
			}
			remoteFileName = remoteFileInfo.Name()
		}

		var localFilePath string
		{
			if strings.HasSuffix(localPath, "/") {
				localPath = localPath + remoteFileName
			} else {
				localPath = localPath + "/" + remoteFileName
			}
		}
		if err := createLocalFile(remoteFile, localFilePath); err != nil {
			return fmt.Errorf("failed to create local file: %w", err)
		}
	}

	return nil
}

func (t transfer) List(ip, user, passwd, remoteFilePath string) ([]string, error) {
	sftpClient, err := newSftpClient(ip, user, passwd)
	if err != nil {
		return nil, err
	}

	paths, err := sftpClient.Glob(remoteFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to Glob %s: %w", remoteFilePath, err)
	}
	return paths, nil
}

func createLocalFile(remoteFile *sftp.File, localFilePath string) error {
	// 如果文件存在则报错
	if _, err := os.Stat(localFilePath); err == nil {
		return fmt.Errorf("create file Error, %s already exists", localFilePath)
	}

	localFile, err := os.Create(localFilePath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %w", err)
	}
	defer localFile.Close()

	if _, err := remoteFile.WriteTo(localFile); err != nil {
		return fmt.Errorf("failed to write to local file: %w", err)
	}
	return nil
}

func newSftpClient(ip, user, passwd string) (*sftp.Client, error) {
	sftpPort := 22
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(passwd),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	sshClient, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", ip, sftpPort), sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to establish SSH connection: %w", err)
	}
	defer sshClient.Close()

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create SFTP client: %w", err)
	}
	return sftpClient, nil
}
