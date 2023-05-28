data "hcloud_image" "base-ubuntu" {
  with_selector = "name=base-ubuntu-22.04"
  most_recent   = true
}

resource "hcloud_network" "main" {
  name     = "main"
  ip_range = "10.0.0.0/16"
}

resource "hcloud_network_subnet" "main" {
  network_id   = hcloud_network.main.id
  type         = "cloud"
  network_zone = "eu-central"
  ip_range     = "10.0.1.0/24"
}

resource "hcloud_load_balancer" "main" {
  name               = "main"
  load_balancer_type = "lb11"
  network_zone       = "eu-central"
}

resource "hcloud_load_balancer_network" "main" {
  network_id       = hcloud_network.main.id
  load_balancer_id = hcloud_load_balancer.main.id
}

resource "hcloud_server" "main" {
  count       = var.main_server_count
  server_type = var.main_server_type
  name        = format("main-%02d", count.index + 1)

  image = data.hcloud_image.base-ubuntu.id

  labels = {
    "role" = "main"
  }

  network {
    network_id = hcloud_network.main.id
  }

  depends_on = [hcloud_network_subnet.main]
}

resource "hcloud_server" "agent" {
  count       = var.agents_server_count
  server_type = var.agents_server_type
  name        = format("agent-%02d", count.index + 1)

  image = data.hcloud_image.base-ubuntu.id

  labels = {
    "role" = "agent"
  }

  network {
    network_id = hcloud_network.main.id
  }
}

resource "hcloud_server" "data" {
  count       = var.data_server_count
  server_type = var.data_server_type
  name        = format("data-%02d", count.index + 1)

  image = data.hcloud_image.base-ubuntu.id

  labels = {
    "role" = "data"
  }

  network {
    network_id = hcloud_network.main.id
  }
}

data "cloudflare_zone" "kacpermalachowskipl" {
  account_id = var.cloudflare_account_id
  name      = "kacpermalachowski.pl"
}

resource "cloudflare_record" "root" {
  zone_id = data.cloudflare_zone.kacpermalachowskipl.id
  name    = "kacpermalachowski.pl"
  value   = hcloud_load_balancer.main.ipv4
  type    = "A"
  proxied = true
  ttl     = 1
}

resource "cloudflare_record" "wildcard" {
  zone_id = data.cloudflare_zone.kacpermalachowskipl.id
  name    = "*.kacpermalachowski.pl"
  value   = hcloud_load_balancer.main.ipv4
  type    = "A"
  proxied = true
  ttl     = 1
}