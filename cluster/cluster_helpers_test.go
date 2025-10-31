package cluster

import (
	"math/rand"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/jackc/fake"
	"github.com/scalezilla/scalezilla/logger"
)

type basicClusterConfig struct {
	randomPort, dev, isVoter bool
}

func makeBasicCluster(cfg basicClusterConfig) *Cluster {
	if cfg.randomPort {
		defaultHTTPPort = uint16(makeRandomInt(16000, 40000))
		defaultRaftGRPCPort += 1
		defaultGRPCPort += 1
	}

	config := ClusterInitialConfig{
		Logger:               logger.NewLogger(),
		BindAddress:          defaultBindAddress,
		HostIPAddress:        defaultHostIPAddress,
		HTTPPort:             defaultHTTPPort,
		RaftGRPCPort:         defaultRaftGRPCPort,
		GRPCPort:             defaultGRPCPort,
		TestRaftMetricPrefix: fake.CharactersN(50),
		Dev:                  cfg.dev,
	}
	z, _ := NewCluster(config)
	z.buildAddressAndID()

	if !cfg.isVoter {
		client := &Client{
			Raft: &RaftConfig{
				TimeMultiplier:    1,
				SnapshotInterval:  30 * time.Second,
				SnapshotThreshold: defaultSnapshotThreshold,
			},
		}
		z.config.Client = client
	}

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
