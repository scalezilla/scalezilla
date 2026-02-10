host_ip_address = "127.0.0.10"

server {
  enabled = true

  raft {
    bootstrap_expected_size = 3
  }

  cluster_join {
    initial_members = ["127.0.0.11", "127.0.0.12"]
  }
}

client {
  enabled = false
}
