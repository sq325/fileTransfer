package service

import "path/filepath"

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
	pwd, _ := filepath.Abs("./")
	fn := filepath.Base(remoteFilePath)
	if err := d.Get(remoteIp, remoteUser, remotePasswd, remoteFilePath, pwd); err != nil {
		return err
	}
	if err := d.Put(clientIp, clientUser, clientPasswd, clientDir, filepath.Join(pwd, fn)); err != nil {
		return err
	}
	return nil
}
