package site

func Backup(st *Site) error {
	return st.rp.Backup(st.StConf.BackupDirectory)
}
