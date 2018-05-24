# Micro Installer #

An auto-installer of [micro](https://github.com/micro/micro) as a service in
systemd systems.

# Installation Instructions #

There are two approaches to install:

1. [Install pre-compiled micro binary](#install-precompiled-micro-binary) (straight forward approach)
    
    Notes:
    
    1. Most services can safely use this approach.
    1. This is not guaranteed to use the latest version of micro
    as it uses a precompiled binary.
    1. The installer installs `micro web` and `micro api`,
    if your service requires other `micro` products then
    it's best to use the [build and install](#build-and-install-the-latest-release-of-micro)
    approach - this gives you more options.
1. [Build and install the latest release of the micro binary](#build-and-install-the-latest-release-of-micro)
(you will need to have the go compiler installed to achieve this approach)

## Install Precompiled micro binary ##

1. Clone the repository and cd into it
    ```
    $ git clone https://github.com/tomogoma/micro-installer
    $ cd micro-installer
    ```
1. Install `micro` together with the respective unit files:
    ```
    sudo make install
    ```
    This installs the `web` and `api` commands.
    See the [Install outcome](#install-outcome) section for next steps.
    
## Build and install the latest release of micro ##

### Prerequisite ###

To run the build, you need a fully set up go runtime and working GOPATH.
Instructions here: https://golang.org/doc/install

This is because the build scripts use `go get` to fetch the micro repository,
checkout the latest release branch and build an executable to be installed.
SystemD unit files are also created for the micro commands provided during build.

### Procedure ###

1. Clone the repository and cd into it
    ```
    $ git clone https://github.com/tomogoma/micro-installer
    $ cd micro-installer
    ```
1. Build the installer.
Build the installer with the micro commands you wish to have unit files for
e.g. to build for `micro api` and `micro web` run:
    ```
    make build commands="api web"
    ```
    A `micro` executable is created at `bin/micro` while the respective unit files
    are created inside the `unit` directory
1. Install the build outcome:
    ```
    sudo make install
    ```
    This installs the commands provided during the build step.
    See the [Install outcome](#install-outcome) section for next steps.
    
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

the command `micro --help` or `micro [command] --help` provides all configuration options
applicable.

To change config values for a command e.g. change listening address, append the relevant environment
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

