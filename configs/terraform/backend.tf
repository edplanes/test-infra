terraform {
  cloud {
    organization = "KacperMalachowski"

    workspaces {
      name = "prod"
    }
  }
}