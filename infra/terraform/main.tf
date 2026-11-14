terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 4.2.0"
    }
  }
}

provider "docker" {
  // point to colima socket
  host = "unix://${pathexpand("~/.colima/default/docker.sock")}"
}

resource "docker_image" "machete_img" {
  name = "alpine"
  build {
    context    = "../docker/"
    dockerfile = "Dockerfile"
  }
}

resource "docker_volume" "machete_vol" {
  name = "machete_cache"
}

resource "docker_container" "machete_container" {
  name  = "machete"
  image = docker_image.machete_img.image_id
  volumes {
    volume_name    = docker_volume.machete_vol.name
    container_path = "/home/dev/.cache"
  }
  stdin_open = true # Equivalent to -i
  tty        = true # Equivalent to -t
  ports {
    internal = 80
    external = 8000
  }
}

output "container_id" {
  description = "ID of docker container"
  value       = docker_container.machete_container.id
}

output "image_id" {
  description = "ID of docker image"
  value       = docker_image.machete_img.id
}
