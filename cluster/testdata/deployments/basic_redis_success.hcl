deployment "redis-test" {
  kind = "service"
  namespace = "default"

  metadata = {
    a = "b"
    c = "d"
  }

  pod "redis-pod" {
    container "redis-container" {
      image = "docker.io/library/redis:latest"
      resources {
        cpu    = 128
        memory = 128
      }
    }
  }
}
