# Micro Installer #

An auto-installer of [micro](https://github.com/micro/micro) as a service in
systemd systems.

The scripts `go get` the micro repository, checkout the latest release branch
and build an executable to be installed. SystemD unit files are also created
for the micro commands provided during build.

## Pre-requisite ##

To run the build, you need a fully set up go runtime and working GOPATH.
Documentation here: https://golang.org/doc/install

## Install ##

You may skip The Build step in case you want to use prebuilt versions of the 
unit files and binary - the micro binary may be older than the latest release -
they include `micro web` as `microweb.service` and `micro api` as
`microapi.service`.

1. Clone the repository and cd into it
    ```
    $ git clone https://github.com/tomogoma/micro-installer
    $ cd micro-installer
    ```
2. Build the installer with the micro commands you wish to have unit files for
e.g. to build for `micro api` and `micro web` run:
    ```
    make build commands="api web"
    ```
    A `micro` executable is created at `bin/micro` while the respective unit files
    are created inside the `unit` directory
3. Install `micro` together with the respective unit files:
    ```
    sudo make install
    ```
    
## Install outcome ##


### Start

`$ systemctl start micro[command].service`
e.g.
`$ systemctl start microapi.service`


### Stop

`$ systemctl stop micro[command].service`


### Check status

`$ systemctl status micro[command].service`


### Install Directories

The micro binary is installed into
`/usr/local/bin/micro`

Systemd service unit files are created at
`/etc/systemd/system/micro[command].service`

## Configuring ##

To change config values e.g. change listening address, append the relevant environment
variable to `/etc/systemd/system/micro[command].service`
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

`$ sudo make uninstall`

