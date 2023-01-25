package service

func (serv *ServiceImp) Backup() error {
	return serv.rpo.Backup(serv.conf.BackupDirectory)
}
