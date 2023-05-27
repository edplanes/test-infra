build {
  sources = ["source.hcloud.base"]
  name = "base-${var.server_image}-${timestamp()}"

  provisioner "ansible" {
    playbook_file = "./ansible/playbook.yml"
    ansible_ssh_extra_args = ["-o IdentitiesOnly=yes -o HostKeyAlgorithms=+ssh-rsa -o PubkeyAcceptedKeyTypes=+ssh-rsa"]
  }
}