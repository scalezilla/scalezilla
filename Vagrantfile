# -*- mode: ruby -*-
# vi: set ft=ruby :

# https://portal.cloud.hashicorp.com/vagrant/discover
DEBIAN_IMAGE = "bento/debian-13"
REDHAT_IMAGE = "cloud-image/almalinux-9"
#REDHAT_IMAGE = "cloud-image/fedora-42"
# network is messy with cloud-image/centos-10-stream
# vagrant is messy with cloud-image/almalinux-9
# for guest additions

MACHINES = {
  'dev' => {
    ip: "192.168.200.10",
    cpus: 4,
    memory: "4096",
  },
  'server' => {
    cpus: 2,
    memory: "2048",
    count: 2, # total 3
  },
  'client_debian' => {
    cpus: 2,
    memory: "2048",
    count: 1, # total 2
  },
  'client_redhat' => {
    cpus: 2,
    memory: "2048",
    count: 0, # total 1
  },
}

Vagrant.configure("2") do |config|
  config.vm.box_check_update = false

  config.vm.synced_folder ".", "/home/vagrant/scalezilla", mount_options: ["ro"]
  config.vm.define "dev" do |srv|
    srv.vm.box = "#{DEBIAN_IMAGE}"
    srv.vm.hostname = "dev"
    srv.vm.network "private_network", ip: MACHINES['dev'][:ip]
    srv.vm.provider "virtualbox" do |vb|
      vb.memory = MACHINES['dev'][:memory]
      vb.cpus   = MACHINES['dev'][:cpus]
      vb.name   = "dev"
        vb.check_guest_additions = false
    end

    srv.vm.provision "os", type: "shell", name: "install_debian_os_package", path: "./scripts/install_debian_os_package.sh"
    srv.vm.provision "golang", type: "shell", name: "install_golang", path: "./scripts/install_golang.sh"
  end

  0.upto(MACHINES['server'][:count]) do |index|
    id = "#{11 + index}"
    server_ssh_port = 2230 + index
    server_name = "server#{id}"
    ip = "192.168.200.#{id}"
    config.vm.define "#{server_name}" do |srv|
      srv.vm.box = "#{DEBIAN_IMAGE}"
      srv.vm.hostname = "#{server_name}"
      srv.vm.network "private_network", ip: "#{ip}"
      srv.vm.network "forwarded_port", guest: 22, host: "#{server_ssh_port}", id: "ssh"
      srv.vm.provider "virtualbox" do |vb|
        vb.memory = MACHINES['server'][:memory]
        vb.cpus   = MACHINES['server'][:cpus]
        vb.name   = "#{server_name}"
        vb.check_guest_additions = false
      end

      srv.vm.provision "os", type: "shell", name: "install_debian_os_package", path: "./scripts/install_debian_os_package.sh"
      srv.vm.provision "golang", type: "shell", name: "install_golang", path: "./scripts/install_golang.sh"
      srv.vm.provision "scalezilla", type: "shell", name: "scalezilla", inline: "/home/vagrant/scalezilla/scripts/scalezilla.sh config_success_server#{id}.hcl"
    end
  end

  0.upto(MACHINES['client_debian'][:count]) do |index|
    id = "#{20 + index}"
    server_ssh_port = 2240 + index
    server_name = "client#{id}"
    ip = "192.168.200.#{id}"
    config.vm.define "#{server_name}" do |srv|
      srv.vm.box = "#{DEBIAN_IMAGE}"
      srv.vm.hostname = "#{server_name}"
      srv.vm.network "private_network", ip: "#{ip}"
      srv.vm.network "forwarded_port", guest: 22, host: "#{server_ssh_port}", id: "ssh"
      srv.vm.provider "virtualbox" do |vb|
        vb.memory = MACHINES['client_debian'][:memory]
        vb.cpus   = MACHINES['client_debian'][:cpus]
        vb.name   = "#{server_name}"
        vb.check_guest_additions = false
      end

      srv.vm.provision "os", type: "shell", name: "install_debian_os_package", path: "./scripts/install_debian_os_package.sh"
      srv.vm.provision "golang", type: "shell", name: "install_golang", path: "./scripts/install_golang.sh"
      srv.vm.provision "scalezilla", type: "shell", name: "scalezilla", inline: "/home/vagrant/scalezilla/scripts/scalezilla.sh config_success_client#{id}.hcl"
    end
  end

#  0.upto(MACHINES['client_redhat'][:count]) do |index|
#    id = "#{25 + index}"
#    server_ssh_port = 2250 + index
#    server_name = "client#{id}"
#    ip = "192.168.200.#{id}"
#    config.vm.define "#{server_name}" do |srv|
#      srv.vm.box = "#{REDHAT_IMAGE}"
#      srv.vm.hostname = "#{server_name}"
#      srv.vm.network "private_network", ip: "#{ip}"
#      srv.vm.network "forwarded_port", guest: 22, host: "#{server_ssh_port}", id: "ssh"
#      srv.vm.provider "virtualbox" do |vb|
#        vb.memory = MACHINES['client_redhat'][:memory]
#        vb.cpus   = MACHINES['client_redhat'][:cpus]
#        vb.name   = "#{server_name}"
#      end
#
#      srv.vm.provision "os", type: "shell", name: "install_redhat_os_package", path: "./scripts/install_redhat_os_package.sh"
#      srv.vm.provision "golang", type: "shell", name: "install_golang", path: "./scripts/install_golang.sh"
#    end
#  end

end

