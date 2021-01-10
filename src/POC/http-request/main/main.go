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
	// fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	// fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	return  m.Mallocs
}

func test(f func (int) (int), n int) () {
	alloc := printMemStat()
	startTime := time.Now()
	err := helper.NetLinearRequest(n)
	fmt.Println(time.Since(startTime))
	newAlloc := printMemStat()
	fmt.Println("Alloc = ", (newAlloc - alloc))
	fmt.Println("total error : ", err, "/", n, "\n")
}

func main() () {
	totalN := 999
	// go helper.StartServer()

	fmt.Println("net - Linear request")
	test(helper.NetLinearRequest, totalN)
	
	fmt.Println("goreq - Linear request")
	test(helper.GoreqLinearRequest, totalN)

	fmt.Println("net - go routine request")
	test(helper.NetGoRequest, totalN)

	fmt.Println("goreq - go routine request")
	test(helper.GoreqGoRequest, totalN)
}