# -*- mode: ruby -*-
# vi: set ft=ruby :

GO_VERSION="1.8"

Vagrant.configure("2") do |config|
  config.vm.box = "boxcutter/centos73"

  config.vm.provision "shell", inline: <<-SHELL
    yum install -y git vim

    if [ ! -e /usr/local/go ]; then
        cd  /tmp
        curl -s -LO https://storage.googleapis.com/golang/go#{GO_VERSION}.linux-amd64.tar.gz
        tar xzf go#{GO_VERSION}.linux-amd64.tar.gz
        mv go /usr/local/
        ln -sf /usr/local/go/bin/* /usr/local/bin
    fi
 SHELL
end
