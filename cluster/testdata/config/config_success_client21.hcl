hostname        = "client21"
host_ip_address = "192.168.200.21"

client {
  enabled = true

  raft {
    bootstrap_expected_size = 3
  }

  cluster_join {
    initial_members = [
    "192.168.200.11",
    "192.168.200.12",
    "192.168.200.13",
    ]
  }
}

