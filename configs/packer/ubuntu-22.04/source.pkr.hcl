source "hcloud" "base" {
  image = var.server_image
  server_type = var.server_type
  location = var.server_location
  ssh_username = "root"
  snapshot_labels = {
    name = "base-${var.server_image}"
  }

  snapshot_name = "base-${var.server_image}-${timestamp()}"
  token = var.hcloud_token
}