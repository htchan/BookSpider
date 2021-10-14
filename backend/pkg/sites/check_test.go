package sites

import (
	"testing"
	"os"
	"io"

	"golang.org/x/text/encoding/traditionalchinese"
	"github.com/htchan/BookSpider/internal/utils"
)

func init() {
	source, err := os.Open("../../test/site-test-data/ck101_template.db")
	utils.CheckError(err)
	destination, err := os.Create("../../test/site-test-data/ck101_check.db")
	utils.CheckError(err)
	io.Copy(destination, source)
	source.Close()
	destination.Close()
}

var testSite_check = Site{
	SiteName: "ck101",
	database: nil,
	meta: *meta,
	decoder: traditionalchinese.Big5.NewDecoder(),
	databaseLocation: "../../test/site-test-data/ck101_check.db",
	DownloadLocation: "./test_res/site-test-data/",
	MAX_THREAD_COUNT: 100,
}

func TestCheckEnd(t *testing.T) {
	testSite_check.CheckEnd()
	testSite_check.OpenDatabase()
	defer testSite_check.CloseDatabase()
	distinctEndCount, endCount := testSite_check.endCount()
	if distinctEndCount != 1 || endCount != 2 {
		t.Fatalf("Site.CheckEnd() make final results (%v, %v) but not (1, 2)",
			distinctEndCount, endCount)
	}
}
// As function used in Site.Check is just logging, here will not have any test for them