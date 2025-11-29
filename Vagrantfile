# -*- mode: ruby -*-
# vi: set ft=ruby :

MACHINES = {
  'dev' => {
    ip: "192.168.200.10"
  }
}

Vagrant.configure("2") do |config|
  config.vm.box_check_update = false

  config.vm.synced_folder ".", "/home/vagrant/scalezilla", mount_options: ["ro"]
  config.vm.define "dev" do |srv|
    srv.vm.box = "bento/debian-13"
    srv.vm.hostname = "dev"
    srv.vm.network "private_network", ip: MACHINES['dev'][:ip]
    srv.vm.provider "virtualbox" do |vb|
      vb.memory = "4096"
      vb.cpus   = 4
    end
    srv.vm.provision "shell", name: "install_os_package", inline: <<-SHELL
      apt-get update
      apt-get install -y vim curl wget git netcat-openbsd \
        rsync tcpdump dnsutils net-tools traceroute jq ripgrep \
        gcc
    SHELL
    srv.vm.provision "shell", name: "install_golang", path: "install_golang.sh"

    srv.vm.provision "shell", name: "check_golang_version", inline: <<-SHELL
      grep -q 'export PATH=$PATH:/usr/local/go/bin' /root/.bashrc || echo 'export PATH=$PATH:/usr/local/go/bin' >> /root/.bashrc
      grep -q 'export PATH=$PATH:/usr/local/go/bin' /home/vagrant/.bashrc || echo 'export PATH=$PATH:/usr/local/go/bin' >> /home/vagrant/.bashrc
      source ~/.bashrc
      go version
    SHELL
  end
end
