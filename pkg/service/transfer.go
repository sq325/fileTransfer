package service

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type Transfer interface {
	Get(remoteIp, remoteUser, remotePasswd, remoteFilePath, srcDir string) error
	List(remoteIp, remoteUser, remotePasswd, remoteFilePath string) ([]string, error)
	Put(clientIp, ClientUser, ClientPasswd, clientDir, srcFilePath string) error
	HealthCheck() bool
}

type transfer struct{}

func NewTransfer() Transfer {
	return &transfer{}
}

func (t transfer) Get(remoteIp, remoteUser, remotePasswd, remoteFilePath, srcDir string) error {

	// srcDir must be dir
	if info, err := os.Stat(srcDir); err != nil {
		return fmt.Errorf("failed to get srcDir:%s info: %w", srcDir, err)
	} else if !info.IsDir() {
		return fmt.Errorf("%s must be dir", srcDir)
	}
	srcDir, _ = filepath.Abs(srcDir)

	// remoteFilePath muse be abs path
	if !filepath.IsAbs(remoteFilePath) {
		return fmt.Errorf("%s must be absolute path", remoteFilePath)
	}

	sshClient, err := newSshClient(remoteIp, remoteUser, remotePasswd)
	if err != nil {
		return err
	}
	defer sshClient.Close()

	sftpClient, err := newSftpClient(sshClient)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	ms, err := sftpClient.Glob(remoteFilePath)
	if err != nil {
		return fmt.Errorf("failed to Glob %s: %w", remoteFilePath, err)
	}

	for _, path := range ms {
		remoteFile, err := sftpClient.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open remote file: %w", err)
		}

		remoteFileName := filepath.Base(path)
		localFilePath := filepath.Join(srcDir, remoteFileName)
		if err := createLocalFile(remoteFile, localFilePath); err != nil {
			return fmt.Errorf("failed to create local file %s: %w", localFilePath, err)
		}
		remoteFile.Close()
	}

	return nil
}

func (t transfer) List(remoteIp, remoteUser, remotePasswd, remoteFilePath string) ([]string, error) {
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

func (t transfer) Put(clientIp, ClientUser, ClientPasswd, clientDir, srcFilePath string) error {
	srcFilePath, _ = filepath.Abs(srcFilePath)
	// check srcFilePath
	if info, err := os.Stat(srcFilePath); err != nil {
		return fmt.Errorf("failed to get local file info: %w", err)
	} else if info.IsDir() {
		return fmt.Errorf("%s must not be dir", srcFilePath)
	}

	sshClient, err := newSshClient(clientIp, ClientUser, ClientPasswd)
	if err != nil {
		return err
	}
	defer sshClient.Close()

	sftpClient, err := newSftpClient(sshClient)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	// check dstDir
	if info, err := sftpClient.Stat(clientDir); err != nil {
		return fmt.Errorf("failed to get client file info: %w", err)
	} else if !info.IsDir() {
		return fmt.Errorf("%s must be dir", clientDir)
	}

	// src glob
	ms, err := filepath.Glob(srcFilePath)
	if err != nil {
		return fmt.Errorf("failed to Glob %s: %w", srcFilePath, err)
	}

	for _, srcpath := range ms {
		fileName := filepath.Base(srcpath)
		clientFilePath := filepath.Join(clientDir, fileName)
		clientFile, err := sftpClient.Create(clientFilePath)
		if err != nil {
			return err
		}
		defer clientFile.Close()

		srcFile, err := os.Open(srcpath)
		if err != nil {
			return fmt.Errorf("failed %s to create local file: %w", srcpath, err)
		}
		defer srcFile.Close()

		_, err = io.Copy(clientFile, srcFile) // server -> client
		if err != nil {
			return err
		}
	}

	return nil
}

func (t transfer) HealthCheck() bool {
	return true
}

func createLocalFile(remoteFile *sftp.File, localFilePath string) error {
	// 如果文件存在则报错
	// if _, err := os.Stat(localFilePath); err == nil {
	// 	return fmt.Errorf("create file Error, %s already exists", localFilePath)
	// }
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
