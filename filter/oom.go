package filter

import (
	"bufio"
	"bytes"
	"io"
)

type OOM struct{}

func (o OOM) Filter(reader io.Reader, results chan<- []byte) error {
	defer close(results)

	var result *bytes.Buffer
	recordResult := false
	buf := bufio.NewReader(reader)
	for {
		line, _ := buf.ReadBytes('\n')
		if bytes.Contains(line, []byte("invoked oom-killer")) {
			result = new(bytes.Buffer)
			recordResult = true
		}
		if recordResult {
			if _, err := result.Write(line); err != nil {
				return err
			}
		}
		if bytes.Contains(line, []byte("Killed process")) {
			results <- result.Bytes()
			recordResult = false
		}
	}
	return nil
}
