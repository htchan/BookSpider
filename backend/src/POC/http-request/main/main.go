package main
import (
	"fmt"
	"time"
	"runtime"

	"../helper"
)

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
func printMemStat() (uint64) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	// log.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	// log.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	return  m.Mallocs
}

func test(f func (int) (int), n int) () {
	alloc := printMemStat()
	startTime := time.Now()
	err := helper.NetLinearRequest(n)
	log.Println(time.Since(startTime))
	newAlloc := printMemStat()
	log.Println("Alloc = ", (newAlloc - alloc))
	log.Println("total error : ", err, "/", n, "\n")
}

func main() () {
	totalN := 999
	// go helper.StartServer()

	log.Println("net - Linear request")
	test(helper.NetLinearRequest, totalN)
	
	log.Println("goreq - Linear request")
	test(helper.GoreqLinearRequest, totalN)

	log.Println("net - go routine request")
	test(helper.NetGoRequest, totalN)

	log.Println("goreq - go routine request")
	test(helper.GoreqGoRequest, totalN)
}