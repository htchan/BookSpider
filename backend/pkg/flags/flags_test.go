package flags

import (
	"flag"
	"testing"
)

func Test_Flags_Flags_Load(t *testing.T) {
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
			t.Errorf("flags.Load(%v) returns \nOperation=%v,\tSite=%v,\tId=%v\tMaxThreads=%v\n"+
				"but not\nOperation=%v,\tSite=%v,\tId=%v\tMaxThreads=%v\n",
				testcase.maxThreadsConfig,
				*actual.Operation, *actual.Site, *actual.Id, *actual.MaxThreads,
				testcase.expect.Operation, testcase.expect.Site,
				testcase.expect.Id, testcase.expect.MaxThreads)
		}
	}
}

func Test_Flags_Flags_IsEverything(t *testing.T) {
	flagOperation, flagSite, flagId := "operation", "test", 123

	t.Run("true if it is empty", func(t *testing.T) {
		f := Flags{}
		if !f.IsEverything() {
			t.Errorf("flags IsEverything return false for empty")
		}
	})

	t.Run("true if it provides operation", func(t *testing.T) {
		f := Flags{ Operation: &flagOperation }
		if !f.IsEverything() {
			t.Errorf("flags IsEverything return false for empty")
		}
	})

	t.Run("false if providing site, id", func(t *testing.T) {
		f := Flags{ Site: &flagSite, Id: &flagId }
		if f.IsEverything() {
			t.Errorf("flags IsEverything return true for providing site, id")
		}
	})

	t.Run("false if providing site", func(t *testing.T) {
		f := Flags{ Site: &flagSite }
		if f.IsEverything() {
			t.Errorf("flags IsEverything return true for providing site")
		}
	})
}

func Test_Flags_Flags_IsBook(t *testing.T) {
	flagSite, flagId, flagHashCode := "test", 123, "abc"

	t.Run("true if it provide site, id, hash", func(t *testing.T) {
		f := Flags{ Site: &flagSite, Id: &flagId, HashCode: &flagHashCode }
		if !f.IsBook() {
			t.Errorf("flags IsBook return false for site, id, hash")
		}
	})

	t.Run("true if it provide site, id", func(t *testing.T) {
		f := Flags{ Site: &flagSite, Id: &flagId }
		if !f.IsBook() {
			t.Errorf("flags IsBook return false for site, id")
		}
	})

	t.Run("false if missing site", func(t *testing.T) {
		f := Flags{ Id: &flagId, HashCode: &flagHashCode }
		if f.IsBook() {
			t.Errorf("flags IsBook return true for missing site")
		}
	})

	t.Run("false if missing id", func(t *testing.T) {
		f := Flags{ Site: &flagSite, HashCode: &flagHashCode }
		if f.IsBook() {
			t.Errorf("flags IsBook return true for missing id")
		}
	})
}

func Test_Flags_Flags_IsSite(t *testing.T) {
	flagSite, flagId, flagHashCode := "test", 123, "abc"
	t.Run("true if it only provide site", func(t *testing.T) {
		f := Flags{ Site: &flagSite }
		if !f.IsSite() {
			t.Errorf("flags IsSite return false for only site")
		}
	})

	t.Run("false if site not provide", func(t *testing.T) {
		f := Flags{}
		if f.IsSite() {
			t.Errorf("flags IsSite return true for empty flag")
		}
	})

	t.Run("false if id provided", func(t *testing.T) {
		f := Flags{ Site: &flagSite, Id: &flagId }
		if f.IsSite() {
			t.Errorf("flags IsBook return true for site with id")
		}
	})

	t.Run("false if hash code provided", func(t *testing.T) {
		f := Flags{ Site: &flagSite, HashCode: &flagHashCode }
		if f.IsSite() {
			t.Errorf("flags IsBook return true for site with hash code")
		}
	})
}

func Test_Flags_Flags_GetBookInfo(t *testing.T) {
	flagSite, flagId, flagHashCode := "test", 123, "abc"
	t.Run("success if it provide site, id, hash", func(t *testing.T) {
		f := Flags{ Site: &flagSite, Id: &flagId, HashCode: &flagHashCode }
		site, id, hash := f.GetBookInfo()
		if site != "test" || id != 123 || hash != 13368 {
			t.Errorf(
				"flags GetBookInfo return wrong result - site: %v, id: %v, hash: %v",
				site, id, hash)
		}
	})

	t.Run("success if it provide site, id", func(t *testing.T) {
		f := Flags{ Site: &flagSite, Id: &flagId }
		site, id, hash := f.GetBookInfo()
		if site != "test" || id != 123 || hash != -1 {
			t.Errorf(
				"flags GetBookInfo return wrong result - site: %v, id: %v, hash: %v",
				site, id, hash)
		}
	})

	t.Run("invalid data if missing site", func(t *testing.T) {
		f := Flags{ Id: &flagId, HashCode: &flagHashCode }
		site, id, hash := f.GetBookInfo()
		if site != "" || id != 123 || hash != 13368 {
			t.Errorf(
				"flags GetBookInfo return wrong result - site: %v, id: %v, hash: %v",
				site, id, hash)
		}
	})

	t.Run("invalid data if missing id", func(t *testing.T) {
		f := Flags{ Site: &flagSite, HashCode: &flagHashCode }
		site, id, hash := f.GetBookInfo()
		if site != "test" || id != -1 || hash != 13368 {
			t.Errorf(
				"flags GetBookInfo return wrong result - site: %v, id: %v, hash: %v",
				site, id, hash)
		}
	})
}

func Test_Flags_Flags_Valid(t *testing.T) {
	flagSite, flagId, flagHashCode := "test", 123, "abc"

	t.Run("true for valid book", func(t *testing.T) {
		f := Flags{ Site: &flagSite, Id: &flagId, HashCode: &flagHashCode }
		if !f.Valid() {
			t.Errorf("flags Valid return false for valid book")
		}
	})

	t.Run("true for valid site", func(t *testing.T) {
		f := Flags{ Site: &flagSite }
		if !f.Valid() {
			t.Errorf("flags Valid return false for valid site")
		}
	})

	t.Run("return false", func(t *testing.T) {
		f := Flags{ Id: &flagId, HashCode: &flagHashCode }
		if f.Valid() {
			t.Errorf("flags Valid return true for invalid arguments")
		}
	})
}