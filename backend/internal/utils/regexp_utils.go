package utils

import (
	"errors"
	"regexp"
)

func Match(str, regex string) bool {
	return false
}

func Search(str, regex string) (string, error) {
	re := regexp.MustCompile(regex)
	result := re.FindStringSubmatch(str)
	if len(result) > 1 {
		return result[1], nil
	}
	return "", errors.New("no result found")
}

func SearchAll(str, regex string) []string {
	re := regexp.MustCompile(regex)
	matches := re.FindAllStringSubmatch(str, -1)
	results := make([]string, len(matches))
	for i := range matches {
		results[i] = matches[i][1]
	}
	return results
}

func Contains(arr []int, target int) bool {
	for _, i := range arr {
		if i == target {
			return true
		}
	}
	return false
}
