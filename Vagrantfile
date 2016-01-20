Vagrant.configure(2) do |config|
  config.vm.box = "ubuntu/wily64"
  config.vm.synced_folder ".", "/opt/gopath/src/github.com/SebastianCzoch/lxc-exporter"
  config.vm.network "private_network", ip: "192.168.100.100"
  config.vm.provider :virtualbox do |vb|
    vb.customize ["modifyvm", :id, "--memory", 2048]
    vb.customize ["modifyvm", :id, "--cpus", "3"]
  end
end
