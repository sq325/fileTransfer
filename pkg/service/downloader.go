package service

import (
	"fmt"
	"os"
	"path/filepath"
)

type Downloader interface {
	Download(remoteIp, remoteUser, remotePasswd, remoteFilePath, clientIp, clientUser, clientPasswd, clientDir string) error
}

type downloader struct {
	Transfer
}

func NewDownloader() Downloader {
	return &downloader{NewTransfer()}
}

func (d downloader) Download(remoteIp, remoteUser, remotePasswd, remoteFilePath, clientIp, clientUser, clientPasswd, clientDir string) error {
	if isLocalIp(remoteIp) || isLocalIp(clientIp) {
		return fmt.Errorf("%s or %s is server ip", remoteIp, clientIp)
	}

	pwd, _ := filepath.Abs("./")
	fn := filepath.Base(remoteFilePath)
	serverFilePath := filepath.Join(pwd, fn)
	if err := d.Get(remoteIp, remoteUser, remotePasswd, remoteFilePath, pwd); err != nil {
		return err
	}
	if err := d.Put(clientIp, clientUser, clientPasswd, clientDir, serverFilePath); err != nil {
		return err
	}
	err := os.Remove(serverFilePath)
	if err != nil {
		return err
	}
	return nil
}
