deployment "nginx-test" {
  kind = "service"
  namespace = "default"

  metadata = {
    a = "b"
    c = "d"
  }

  pod "nginx-pod" {
    container "nginx-container" {
      image = "docker.io/library/nginx:latest"
      resources {
        cpu    = 128
        memory = 128
      }
    }
  }
}