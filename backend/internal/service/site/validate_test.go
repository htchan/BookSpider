package site

import (
	"testing"

	"github.com/htchan/BookSpider/internal/config"
	"github.com/htchan/BookSpider/internal/mock"
	psql "github.com/htchan/BookSpider/internal/repo/psql"
)

func Test_Validate(t *testing.T) {
	t.Parallel()

	server := mock.MockSiteServer()

	psql.StubPsqlConn()
	db, err := psql.OpenDatabase("")
	if err != nil {
		t.Fatalf("error in open database: %v", err)
	}
	site := "valid/st"

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		db.Exec("delete from writers where id>0 and name like $1", site+"%")
		db.Exec("delete from errors where site=$1", site)
		server.Close()
	})

	st, err := NewSite(
		site,
		&config.BookConfig{},
		&config.SiteConfig{
			BookKey:         "test_book",
			MaxExploreError: 2,
		},
		&config.CircuitBreakerClientConfig{MaxThreads: 1}, nil, nil)
	if err != nil {
		t.Errorf("fail to init site: %v", err)
		return
	}

	stubData(st.rp, site)

	tests := []struct {
		name               string
		st                 *Site
		availabilityConfig config.AvailabilityConfig
		expectErr          bool
	}{
		{
			name: "works",
			st:   st,
			availabilityConfig: config.AvailabilityConfig{
				URL:   server.URL + "/error",
				Check: "error",
			},
			expectErr: false,
		},
		{
			name: "return err if check not found",
			st:   st,
			availabilityConfig: config.AvailabilityConfig{
				URL:   server.URL + "/error",
				Check: "not found",
			},
			expectErr: true,
		},
		{
			name: "return err if web fail to fetch",
			st:   st,
			availabilityConfig: config.AvailabilityConfig{
				URL:   "https://127.0.0.1:1/error",
				Check: "error",
			},
			expectErr: true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			st.StConf.AvailabilityConfig = test.availabilityConfig
			err := Validate(test.st)
			if (err != nil) != test.expectErr {
				t.Errorf("got error: %v; want error: %v", err, test.expectErr)
			}
		})
	}
}
