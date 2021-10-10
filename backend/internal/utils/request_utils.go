package utils

import (
	"strconv"
	"time"

	"io/ioutil"
	"math/rand"
	"net/http"

	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

func getWeb(url string) string {
	client := http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
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
	resp.Body.Close()
	client.CloseIdleConnections()
	return string(body)
}

const MIN_SLEEP_MS_MULTIPLIER = 2000
const MAX_SLEEP_MS = 30000

func GetWeb(url string, trial int, decoder *encoding.Decoder) (html string, i int) {
	for i = 0; i < 10; i++ {
		html = getWeb(url)
		if _, err := strconv.Atoi(html); err == nil || (len(html) == 0) {
			minSleepMs := i * MIN_SLEEP_MS_MULTIPLIER
			time.Sleep(time.Duration(rand.Intn(MAX_SLEEP_MS-minSleepMs)+minSleepMs) * time.Millisecond)
			continue
		}
		if decoder != nil {
			html, _, _ = transform.String(decoder, html)
		}
		break
	}
	return
}
