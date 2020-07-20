package podcast

import (
	"net/url"
)

func inSlice(a string, list []string) bool {
	for _, b := range list {
		//fmt.Printf("[%s] == [%s]\n", a, b)
		if b == a {
			return true
		}
	}
	return false
}

func inSliceInt(a int, list []int) bool {
	for _, b := range list {
		//fmt.Printf("[%s] == [%s]\n", a, b)
		if b == a {
			return true
		}
	}
	return false
}

func isValidURL(URL string) bool {
	u, err := url.Parse(URL)

	if err != nil {
		return false
	}

	if u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}
