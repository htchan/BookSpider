package flags

import (
	"flag"
	"testing"
)

func TestLoad(t *testing.T) {
	type Inputs struct {
		operation, site, id, maxThreads string
	}
	type Expects struct {
		Operation, Site string
		Id, MaxThreads  int
	}
	var testcases = []struct {
		input            Inputs
		maxThreadsConfig int
		expect           Expects
	}{
		{
			Inputs{"", "", "", ""},
			200,
			Expects{"", "", -1, 200},
		},
		{
			Inputs{"operation 1", "site 1", "123", "10"},
			200,
			Expects{"operation 1", "site 1", 123, 100},
		},
		{
			Inputs{"operation 1", "site 1", "abc", "def"},
			10,
			Expects{"operation 1", "site 1", 0, 10},
		},
	}

	actual := *NewFlags()

	for _, testcase := range testcases {
		if len(testcase.input.operation) != 0 {
			flag.Set("operation", testcase.input.operation)
		}
		if len(testcase.input.operation) != 0 {
			flag.Set("site", testcase.input.site)
		}
		if len(testcase.input.operation) != 0 {
			flag.Set("id", testcase.input.id)
		}
		if len(testcase.input.operation) != 0 {
			flag.Set("max-threads", testcase.input.maxThreads)
		}

		actual.Load(testcase.maxThreadsConfig)

		if *actual.Operation != testcase.expect.Operation || *actual.Site != testcase.expect.Site ||
			*actual.Id != testcase.expect.Id || *actual.MaxThreads != testcase.expect.MaxThreads {
			t.Fatalf("flags.Load(%v) returns \nOperation=%v,\tSite=%v,\tId=%v\tMaxThreads=%v\n"+
				"but not\nOperation=%v,\tSite=%v,\tId=%v\tMaxThreads=%v\n",
				testcase.maxThreadsConfig,
				*actual.Operation, *actual.Site, *actual.Id, *actual.MaxThreads,
				testcase.expect.Operation, testcase.expect.Site,
				testcase.expect.Id, testcase.expect.MaxThreads)
		}
	}
}
