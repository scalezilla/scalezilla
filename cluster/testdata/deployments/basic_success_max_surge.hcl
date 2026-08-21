deployment "nginx-test" {
  kind = "service"
  namespace = "default"

  metadata = {
    a = "b"
    c = "d"
  }

  max_surge = 15

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