package testutil

import (
	"bytes"
	"github.com/golang/glog"
	"io"
	"os"
	"path"
	"strings"
)

// Normalize a local path to point to an absolute path
func MakeTestPath(localPath string) string {
	return path.Join(os.Getenv("PWD"), "../../", localPath)
}

// Capture output generated by glog
func GetOutput(fn func()) (string, error) {
	old := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		return "", err
	}
	os.Stderr = w

	ch := make(chan string)
	// Need to drain the channel in a separate goroutine so that we don't deadlock
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		ch <- buf.String()
	}()

	fn()
	glog.Flush()
	// Restore the original stderr
	os.Stderr = old
	w.Close()

	return <-ch, nil
}

/*
 * Execute a function and return a list of strings representing the output generated by the function.
 * This function is specifically meant to capture the output generated by the glog() functions that
 * the functions use for output.  The first prefix number of characters are removed from each line. To
 * include the entire line, specify a prefix of zero.
 */
func GetOutputAsList(fn func(), prefix int) ([]string, error) {
	output, err := GetOutput(fn)
	if err != nil {
		return []string{}, err
	}
	strlist := strings.Split(output, "\n")
	for i := 0; i < len(strlist); i++ {
		if len(strlist[i]) > prefix {
			strlist[i] = strlist[i][prefix:]
		} else {
			strlist[i] = ""
		}
	}
	return strlist, err
}
