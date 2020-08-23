package model

import (
	"testing"
	//"fmt"
)

func TestUpdate(t *testing.T) {
	sites := LoadSites("../test-resource/config/config.json")
	var testCases = []int{2, 2000, 200000}
	for _, site := range sites {
		success := false
		for _, testCase := range testCases {
			book := site.Book(testCase)
			check := book.Update()
			if check {
				t.Logf("success to get info of %v-book(%v)\n%v\n\n", site.SiteName, testCase, book.JsonString())
				success = true
			} else {
				t.Logf("error of get info of %v-book(%v)\n%v\n\n", site.SiteName, testCase, book.JsonString())
			}
		}
		if !success {
			t.Errorf("site <%v> fail in getting book %v", site.SiteName, testCases)
		}
	}
}
