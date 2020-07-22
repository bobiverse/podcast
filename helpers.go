package podcast

import (
	"bytes"
	"log"
	"net/url"
	"os/exec"
	"strings"
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

// Execute bash script
func runBash(name string, args ...string) ([]byte, []byte, error) {
	log.Printf("⍄  %s %s", name, strings.Join(args, " "))

	cmd := exec.Command(name, args...)
	bufErrOutput := &bytes.Buffer{}
	cmd.Stderr = bufErrOutput

	output, err := cmd.Output()
	return output, bufErrOutput.Bytes(), err
}

// local path to URL
func pathToURL(base, path string) string {
	s := strings.Trim(base, "./") + "/" + strings.Trim(path, "./")
	if isValidURL(s) {
		return s
	}
	return ""
}
