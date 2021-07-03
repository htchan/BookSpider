// This is a package provide helper functions.
package helper

import (
	"net/http"
	"io/ioutil"
	"regexp"
	"strconv"
	"time"
	"os"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

func CheckError(e error) {
	if (e != nil) {
		panic(e);
	}
}

/* web related */
func getWeb(url string) (string) {
	client := http.Client{Timeout: 30*time.Second}
	resp, err := client.Get(url);
	if err != nil {
		return ""
	}
	if resp.StatusCode >= 300 {
		return strconv.Itoa(resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	resp.Body.Close();
	client.CloseIdleConnections()
	return string(body);
}

func GetWeb(url string, trial int, decoder *encoding.Decoder) (html string, i int) {
	for i = 0; i < 10; i++ {
		html = getWeb(url);
		if _, err := strconv.Atoi(html); err == nil || (len(html) == 0) {
			time.Sleep(time.Duration(i * i) * time.Second)
			continue
		}
		if (decoder != nil) {
			html, _, _ = transform.String(decoder, html)
		}
		break
	}
	return
}

/* regex relates */
func Match(str, regex string) (bool) {
	return false
}

func Search(str, regex string) (string) {
	re := regexp.MustCompile(regex);
	result := re.FindStringSubmatch(str);
	if(len(result) > 1) {
		return result[1]
	}
	return "error"
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