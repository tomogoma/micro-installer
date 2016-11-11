# Consul Installer #

An auto-installer of go-micro api [Hashi Corp](https://www.consul.io) as a service in linux systems


## Install ##

Currently only systemd is supported

```
$ git clone https://github.com/tomogoma/consul-installer
$ cd consul-installer
$ ./systemdInstaller.sh /path/to/downloaded/consul.zip
```


## Install outcome ##


Start

`$ systemctl start micro.service`

Stop

`$ systemctl stop micro.service`

Check status

`$ systemctl status micro.service`


Install Directories

The consul binary is installed into
`/usr/local/bin/micro`

A systemd service unit file is created at
`/etc/systemd/system/micro.service`

## Uninstall ##

`$ cd /path/to/consul-installer`

`$ ./systemdUninstaller.sh`


## Configure a dependent systemd service ##

If your service depends on consul, in the service’s unit file add the following lines:


```
[Unit]

...

After=micro.service

Requires=micro.service

...
```
