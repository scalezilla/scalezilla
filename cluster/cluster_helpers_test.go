package cluster

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/jackc/fake"
	"github.com/scalezilla/scalezilla/logger"
	"github.com/stretchr/testify/assert"
)

type basicClusterConfig struct {
	randomPort, dev, isVoter bool
}

type sizedClusterConfig struct {
	port, voterSize, clientSize uint16
}

func makeBasicCluster(cfg basicClusterConfig) *Cluster {
	httpPort := defaultHTTPPort
	grpcPort := defaultGRPCPort
	raftGRPCPort := defaultRaftGRPCPort

	if cfg.randomPort {
		port := uint16(makeRandomInt(16000, 40000))
		httpPort = port
		grpcPort = httpPort + 1
		raftGRPCPort = grpcPort + 1
	}

	config := ClusterInitialConfig{
		Logger:               logger.NewLogger(),
		BindAddress:          defaultBindAddress,
		HostIPAddress:        defaultHostIPAddress,
		HTTPPort:             httpPort,
		GRPCPort:             grpcPort,
		RaftGRPCPort:         raftGRPCPort,
		TestRaftMetricPrefix: fake.CharactersN(50),
		Dev:                  cfg.dev,
	}
	config.Members = append(config.Members, fmt.Sprintf("%s:%d", defaultHostIPAddress, grpcPort+1))
	config.Members = append(config.Members, fmt.Sprintf("%s:%d", defaultHostIPAddress, grpcPort+2))
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

func makeSizedCluster(cfg sizedClusterConfig) (cluster []*Cluster) {
	/*
		20000
		20001
		20002

		20003
		20004
		20005

		20006
		20007
		20008

		20009
		20010
		20011

		20012
		20013
		20014

		20015
		20016
		20017
	*/

	if cfg.voterSize == 0 {
		cfg.voterSize = 3
	}

	if cfg.port == 0 {
		cfg.port = 20000
	}

	var voters []string
	baseHTTPPort, baseGRPCPort, baseRaftPort := cfg.port+uint16(0), cfg.port+1, cfg.port+2
	dbaseHTTPPort, dbaseGRPCPort, dbaseRaftPort := cfg.port+uint16(0), cfg.port+1, cfg.port+2
	for i := range cfg.voterSize {
		var (
			httpPort, grpcPort, raftGRPCPort uint16
			members                          []string
		)
		if i == 0 {
			httpPort = baseHTTPPort
			grpcPort = baseGRPCPort
			raftGRPCPort = baseRaftPort
		} else {
			httpPort = baseHTTPPort + 3*uint16(i)
			grpcPort = baseGRPCPort + 3*uint16(i)
			raftGRPCPort = baseRaftPort + 3*uint16(i)
		}

		for s := range cfg.voterSize {
			if s == i {
				continue
			}
			members = append(members, fmt.Sprintf("%s:%d", defaultHostIPAddress, dbaseGRPCPort+3*uint16(s)))
		}

		if i == 0 {
			voters = append(voters, fmt.Sprintf("%s:%d", defaultHostIPAddress, baseGRPCPort))
			voters = append(voters, members...)
		}

		config := ClusterInitialConfig{
			Logger:               logger.NewLogger(),
			BindAddress:          defaultBindAddress,
			HostIPAddress:        defaultHostIPAddress,
			HTTPPort:             httpPort,
			GRPCPort:             grpcPort,
			RaftGRPCPort:         raftGRPCPort,
			TestRaftMetricPrefix: fake.CharactersN(50),
			Dev:                  true,
			NodePool:             defaultNodePool,
		}

		z, err := NewCluster(config)
		if err != nil {
			log.Fatal(err)
		}
		z.buildAddressAndID()

		// unset dev requirements
		z.dev = false
		z.isVoter = true
		z.members = members

		server := &Server{
			Enabled: true,
			ClusterJoin: &ClusterJoin{
				InitialMembers: members,
			},
			Raft: &RaftConfig{
				TimeMultiplier:        1,
				SnapshotInterval:      30 * time.Second,
				SnapshotThreshold:     defaultSnapshotThreshold,
				BootstrapExpectedSize: uint64(cfg.voterSize),
			},
		}
		z.config.Server = server
		cluster = append(cluster, z)
	}

	xsize := cfg.voterSize * cfg.clientSize
	for i := range cfg.clientSize {
		var httpPort, grpcPort, raftGRPCPort uint16

		if i == 0 {
			httpPort = dbaseHTTPPort + xsize
			grpcPort = dbaseGRPCPort + xsize
			raftGRPCPort = dbaseRaftPort + xsize
		} else {
			httpPort = dbaseHTTPPort + xsize + 3*uint16(i)
			grpcPort = dbaseGRPCPort + xsize + 3*uint16(i)
			raftGRPCPort = dbaseRaftPort + xsize + 3*uint16(i)
		}

		config := ClusterInitialConfig{
			Logger:               logger.NewLogger(),
			BindAddress:          defaultBindAddress,
			HostIPAddress:        defaultHostIPAddress,
			HTTPPort:             httpPort,
			GRPCPort:             grpcPort,
			RaftGRPCPort:         raftGRPCPort,
			TestRaftMetricPrefix: fake.CharactersN(50),
			Dev:                  true,
			NodePool:             defaultNodePool,
		}

		z, err := NewCluster(config)
		if err != nil {
			log.Fatal(err)
		}
		z.buildAddressAndID()

		// unset dev requirements
		z.dev = false
		z.isVoter = false
		z.members = voters

		nodePool := defaultNodePool
		client := &Client{
			Enabled: true,
			ClusterJoin: &ClusterJoin{
				InitialMembers: voters,
			},
			Raft: &RaftConfig{
				TimeMultiplier:    1,
				SnapshotInterval:  30 * time.Second,
				SnapshotThreshold: defaultSnapshotThreshold,
			},
			NodePool: &nodePool,
		}
		z.config.Client = client
		cluster = append(cluster, z)
	}

	return
}

func TestMakeSizedCluster(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		voterSize, clientSize, clusterIndex, port uint16
		httpPort, grpcPort, raftPort              uint16
		expected                                  []string
	}{
		{
			voterSize:    3,
			clusterIndex: 0,
			httpPort:     20000,
			grpcPort:     20001,
			raftPort:     20002,
			expected: []string{
				"127.0.0.1:20004",
				"127.0.0.1:20007",
			},
		},
		{
			voterSize:    3,
			clusterIndex: 1,
			httpPort:     20003,
			grpcPort:     20004,
			raftPort:     20005,
			expected: []string{
				"127.0.0.1:20001",
				"127.0.0.1:20007",
			},
		},
		{
			voterSize:    3,
			clusterIndex: 2,
			httpPort:     20006,
			grpcPort:     20007,
			raftPort:     20008,
			expected: []string{
				"127.0.0.1:20001",
				"127.0.0.1:20004",
			},
		},
		{
			voterSize:    5,
			clusterIndex: 3,
			httpPort:     20009,
			grpcPort:     20010,
			raftPort:     20011,
			expected: []string{
				"127.0.0.1:20001",
				"127.0.0.1:20004",
				"127.0.0.1:20007",
				"127.0.0.1:20013",
			},
		},
		{
			voterSize:    3,
			clientSize:   3,
			clusterIndex: 4,
			httpPort:     20012,
			grpcPort:     20013,
			raftPort:     20014,
			expected: []string{
				"127.0.0.1:20001",
				"127.0.0.1:20004",
				"127.0.0.1:20007",
			},
		},
		{
			voterSize:    3,
			clientSize:   3,
			clusterIndex: 5,
			httpPort:     20015,
			grpcPort:     20016,
			raftPort:     20017,
			expected: []string{
				"127.0.0.1:20001",
				"127.0.0.1:20004",
				"127.0.0.1:20007",
			},
		},
	}

	for _, tc := range tests {
		cluster := makeSizedCluster(sizedClusterConfig{voterSize: tc.voterSize, clientSize: tc.clientSize})
		assert.Equal(tc.expected, cluster[tc.clusterIndex].members)
		assert.Equal(tc.httpPort, cluster[tc.clusterIndex].config.HTTPPort)
		assert.Equal(tc.grpcPort, cluster[tc.clusterIndex].config.GRPCPort)
		assert.Equal(tc.raftPort, cluster[tc.clusterIndex].config.RaftGRPCPort)
	}
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

func makeHTTPRequestRecorder(h http.Handler, method, uri string, header map[string]string, payload string) *httptest.ResponseRecorder {
	var req *http.Request
	if payload == "" {
		req, _ = http.NewRequest(method, uri, nil)
	} else {
		req, _ = http.NewRequest(method, uri, bytes.NewBuffer([]byte(payload)))
		req.Header.Add("Content-Length", strconv.Itoa(len(payload)))
	}
	if len(header) > 0 {
		for k, v := range header {
			req.Header.Add(k, v)
		}
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	return w
}
