package helper
import (
	"fmt"
	"net/http"
	"sync"
	"context"
	"golang.org/x/sync/semaphore"
	"time"
)

func requestNet(client http.Client) (int) {
	// _, err := client.Get("http://localhost:10427")
	_, err := client.Get("https://reqres.in/api/users?page=2")
	if err != nil {
		log.Println(err)
		return 1
	} else {
		return 0
	}
}

func NetLinearRequest(n int) (errorCount int) {
	client := http.Client{Timeout: 10*time.Second}
	errorCount = 0
	for i := 0; i < n; i++ {
		errorCount += requestNet(client)
	}
	return
}

func NetGoRequest(n int) (errorCount int) {
	client := http.Client{Timeout: 10*time.Second}
	errorCount = 0
	var wait sync.WaitGroup
	s := semaphore.NewWeighted(int64(100))
	ctx := context.Background()
	for i := 0; i < n; i++ {
		s.Acquire(ctx, 1)
		wait.Add(1)
		go func() {
			errorCount += requestNet(client)
			wait.Done()
			s.Release(1)
		}()
	}
	wait.Wait()
	return
}