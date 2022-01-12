package utils

func Recover(f func()) {
	if r := recover() ; r != nil {
		f()
	}
}

func CheckError(e error) {
	if e != nil {
		panic(e)
	}
}
