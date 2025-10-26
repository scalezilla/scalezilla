package cluster

import "errors"

type failWriter struct {
	failOn int
	count  int
}

func (fw *failWriter) Write(p []byte) (int, error) {
	fw.count++
	if fw.count == fw.failOn {
		return 0, errors.New("forced write error")
	}
	return len(p), nil
}
