deployment "deployment-name" {
  kind = "service"
  namespace = "default"

  metadata = {
    a = "b"
    c = "d"
  }

  pod "pod-name" {
    container "container_name" {
      image = "nginx:1.29-trixie"
      resources {
        cpu    = 128
        memory = 128
      }
    }
  }
}