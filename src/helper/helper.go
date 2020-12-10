// This is a package provide helper functions.
package helper

import (
	"net/http"
	"io/ioutil"
	"regexp"
	"time"
	"os"
)

func CheckError(e error) {
	if (e != nil) {
		panic(e);
	}
}

func GetWeb(url string) (string) {
	client := http.Client{Timeout: 10*time.Second}
	resp, err := client.Get(url);
	if err != nil {
		return "";
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "";
	}
	resp.Body.Close();
	client.CloseIdleConnections()
	return string(body);
}

func Match(str, regex string) (bool) {
	return false;
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

func Contains(arr []int, target int) (bool) {
	for _, i := range arr {
		if (i == target) {
			return true
		}
	}
	return false
}

func Exists(path string) (bool) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}