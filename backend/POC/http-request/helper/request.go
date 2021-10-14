package utils

import (
	"github.com/parnurzeal/gorequest"
	"sync"
	"context"
	"golang.org/x/sync/semaphore"
)

func requestGo(req *gorequest.SuperAgent) (int) {
	// _, _, errs := req.Get("http://localhost:10427/").End()
	_, _, errs := req.Get("https://reqres.in/api/users?page=2").End()
	if errs != nil {
		return 1
	} else {
		return 0
	}
}

func GoreqLinearRequest(n int) (errorCount int) {
	errorCount = 0
	req := gorequest.New()
	for i := 0; i < n; i++ {
		errorCount += requestGo(req)
	}
	return
}

func GoreqGoRequest(n int) (errorCount int) {
	errorCount = 0
	req := gorequest.New()
	var wait sync.WaitGroup
	s := semaphore.NewWeighted(int64(100))
	ctx := context.Background()
	for i := 0; i < n; i++ {
		s.Acquire(ctx, 1)
		wait.Add(1)
		go func() {
			errorCount += requestGo(req)
			wait.Done()
			s.Release(1)
		}()
	}
	wait.Wait()
	return
}