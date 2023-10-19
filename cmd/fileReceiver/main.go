package main

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/sftp"
	"github.com/spf13/pflag"
	"golang.org/x/crypto/ssh"
)

var (
	ip         *string = pflag.String("ip", "", "remote host ip")
	user       *string = pflag.String("user", "", "remote host user")
	passwd     *string = pflag.String("passwd", "", "remote host user password")
	remotePath *string = pflag.String("remotePath", "", "remote file path")
	outFile    *string = pflag.String("localPath", "", "local file path")
)

func main() {
	pflag.Parse()

	// SFTP服务器地址和端口
	sftpHost := *ip
	sftpPort := 22

	// 登录凭据
	sftpUser := *user
	sftpPass := *passwd

	// 远程文件路径和本地文件路径
	remoteFilePath := *remotePath
	localFilePath := *outFile

	// 建立SSH连接
	sshConfig := &ssh.ClientConfig{
		User: sftpUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(sftpPass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	sshClient, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", sftpHost, sftpPort), sshConfig)
	if err != nil {
		fmt.Println("Failed to establish SSH connection:", err)
		return
	}
	defer sshClient.Close()

	// 建立SFTP会话
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		fmt.Println("Failed to create SFTP client:", err)
		return
	}
	defer sftpClient.Close()

	// 打开远程文件
	remoteFile, err := sftpClient.Open(remoteFilePath)
	if err != nil {
		fmt.Println("Failed to open remote file:", err)
		return
	}
	defer remoteFile.Close()

	// 创建本地文件
	localFile, err := os.Create(localFilePath)
	if err != nil {
		fmt.Println("Failed to create local file:", err)
		return
	}
	defer localFile.Close()

	// 从远程文件复制到本地文件
	_, err = io.Copy(localFile, remoteFile)
	if err != nil {
		fmt.Println("Failed to download file:", err)
		return
	}

	fmt.Println("File downloaded successfully!")
}
