server {
  enabled = true

  raft {
    bootstrap_expected_size = 3
  }

  cluster_join {
    initial_members = ["127.0.0.2", "127.0.0.3"]
  }
}

client {
  enabled = false
}
