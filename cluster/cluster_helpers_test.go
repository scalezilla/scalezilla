package cluster

import (
	"math/rand"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/scalezilla/scalezilla/logger"
)

func makeBasicCluster(randomPort bool) *Cluster {
	if randomPort {
		testHTTPPort = uint16(makeRandomInt(16000, 40000))
		testRaftGRPCPort += 1
		testGRPCPort += 1
	}

	config := ClusterInitialConfig{
		Logger:        logger.NewLogger(),
		BindAddress:   testBindAddress,
		HostIPAddress: testHostIPAddress,
		HTTPPort:      testHTTPPort,
		RaftGRPCPort:  testRaftGRPCPort,
		GRPCPort:      testGRPCPort,
	}
	z := NewCluster(config)
	z.buildAddressAndID()

	return z
}

func makeRandomInt(min int, max int) int {
	if min < 1 {
		min = 1
	}
	if max < 1 {
		max = 1
	}
	if min == 1 && max == 1 {
		return 1
	}
	if min > max {
		return 1
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return min + r.Intn(max-min)
}

func makeHTTPRequestRecorder(h http.Handler, method, uri string, header map[string]string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, uri, nil)
	if len(header) > 0 {
		for k, v := range header {
			req.Header.Add(k, v)
		}
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	return w
}
