package client

import (
	"flag"
	"net/http"
	"os"
	"testing"

	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	leak := flag.Bool("leak", false, "check for memory leaks")
	flag.Parse()

	if *leak {
		goleak.VerifyTestMain(
			m,
			goleak.Cleanup(func(exitCode int) { http.DefaultClient.CloseIdleConnections() }),
		)
	} else {
		os.Exit(m.Run())
	}
}
