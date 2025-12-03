hostname        = "server11"
host_ip_address = "192.168.200.11"

server {
  enabled = true

  raft {
    bootstrap_expected_size = 3
  }

  cluster_join {
    initial_members = [
    "192.168.200.12",
    "192.168.200.13",
    ]
  }
}

client {
  enabled = false
}
