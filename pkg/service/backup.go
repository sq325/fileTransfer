package service

type Backuper interface {
	Backup(ip, user, passwd, remoteFilePath, localPath, duration, timeLayout string) error
}

type backuper struct{}

func NewBackuper() Backuper {
	return &backuper{}
}

func (b backuper) Backup(ip, user, passwd, remoteFilePath, localPath, duration, timeLayout string) error {
	return nil
}
