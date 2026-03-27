deployment "nginx-test" {
  kind = "service"

  metadata = {
    a = "b"
    c = "d"
  }

  pod "nginx-pod" {
    container "nginx-container" {
      image = "docker.io/library/nginx:latest"
    }
  }
}