package site

func Check(st *Site) error {
	return st.rp.UpdateBooksStatus()
}
