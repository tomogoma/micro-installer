# Micro Installer #

An auto-installer of [micro](https://github.com/micro/micro) as a service in linux systems


## Install ##

Currently only systemd is supported

```
$ git clone https://github.com/tomogoma/micro-installer
$ cd micro-installer
$ ./systemdInstaller.sh
```


## Install outcome ##


Start

`$ systemctl start micro@web.service`

`$ systemctl start micro@api.service`

`$ systemctl start micro@[name].service`

Stop

`$ systemctl stop micro@[name].service`

Check status

`$ systemctl status micro@[name].service`


Install Directories

The micro binary is installed into
`/usr/local/bin/micro`

A systemd service unit file is created at
`/etc/systemd/system/micro@.service`

## Configuring ##

To change config values e.g. change listening address, append the relevant environment
variable to `/etc/systemd/system/micro@.service`
e.g.
```
[Service]
Environment=MICRO_API_ADDRESS=0.0.0.0:8089 "ANOTHER_ENV_VAR=some value"
...
```
Once done, perform a daemon reload for changes to take effect:
```
systemctl daemon-reload
```


## Uninstall ##

`$ cd /path/to/micro-installer`

`$ ./systemdUninstaller.sh`


## Configure a dependent systemd service ##

If your service depends on micro, in the service’s unit file add the following lines:


```
[Unit]

...

After=micro@[name].service

Requires=micro@[name].service

...
```
