package service

import (
	"fmt"
	"os"
	"path/filepath"
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
	sshClient, err := newSshClient(ip, user, passwd)
	if err != nil {
		return err
	}
	defer sshClient.Close()

	sftpClient, err := newSftpClient(sshClient)
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

		var remoteFileName string
		{
			remoteFileInfo, err := remoteFile.Stat()
			if err != nil {
				return fmt.Errorf("failed to get remote file info: %w", err)
			}
			remoteFileName = remoteFileInfo.Name()
		}

		var localFilePath string
		// 如果localFilePath是文件夹
		{
			info, err := os.Stat(localPath)
			if err != nil {
				return fmt.Errorf("failed to get local file info: %w", err)
			}
			if info.IsDir() {
				localFilePath = filepath.Join(localPath, remoteFileName)
			} else {
				localFilePath = localPath
			}
		}
		if err := createLocalFile(remoteFile, localFilePath); err != nil {
			return fmt.Errorf("failed to create local file %s: %w", localFilePath, err)
		}
		remoteFile.Close()
	}

	return nil
}

func (t transfer) List(ip, user, passwd, remoteFilePath string) ([]string, error) {
	sshClient, err := newSshClient(ip, user, passwd)
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

func newSshClient(ip, user, passwd string) (*ssh.Client, error) {
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
	return sshClient, nil
}

func newSftpClient(sshClient *ssh.Client) (*sftp.Client, error) {
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create SFTP client: %w", err)
	}
	return sftpClient, nil
}
