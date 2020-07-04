package helper

import (
	"net/http"
	"io/ioutil"
	"regexp"
)

func CheckError(e error) {
	if (e != nil) {
		panic(e);
	}
}

func GetWeb(url string) (string) {
	resp, err := http.Get(url);
	if err != nil {
		return "";
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "";
	}
	resp.Body.Close();
	return string(body);
}

func Search(str, regex string) (string) {
	re := regexp.MustCompile(regex);
	result := re.FindStringSubmatch(str);
	if(len(result) > 1) {
		return result[1]
	}
	return "error";
}

func SearchAll(str, regex string) ([]string) {
	re := regexp.MustCompile(regex);
	matches := re.FindAllStringSubmatch(str, -1);
	results := make([]string, len(matches));
	for i := range matches {
		results[i] = matches[i][1];
	}
	return results
}