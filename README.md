# gorigctl

[![Build Status](https://travis-ci.org/dh1tw/gorigctl.svg?branch=master)](https://travis-ci.org/dh1tw/gorigctl)

![Alt text](http://i.imgur.com/V8z68Pm.png "Screenshot gorigctl's cli based GUI")

This application is used to control a (ham) radio locally or remotely. Gorigctl
comes with both, a local and a remote command line interface and (cli base) GUI.
gorigctl is written in [Go](https://golang.org) and using [hamlib](http://www.hamlib.log)
via the [goHamlib](https://github.com/dh1tw/goHamlib) bindings.

gorigctl can be considered as a dropin replacement for [hamlib's](http://www.hamlib.log)
rigctl(d). Instead of sending text commands over a TCP channel, gorigctl is
implementing the [Shackbus](https://shackbus.org) standard, using
[Protocol Buffers](https://developers.google.com/protocol-buffers/) for
(de) serialization and [MQTT](http://mqtt.org) for transportation. Thanks to
MQTT's Publish/Subscribe architecture, serveral clients can access the radio
simultaneously.

gorigctl tries to maintain compatibility with hamlib's rigctl cli commands.

**ADVICE**: This project is **under development**. This application is still
under development and not considered ready for production.

**ADVICE**: The user experience depends heavily on the backend implementation
of the selected radio. Radios with a _stable_ backend provide the best
experience. [Hamlib's rig matrix](http://hamlib.sourceforge.net/sup-info/rigmatrix.html)
provides an overview over the available backends and their status.

## Supported Platforms

gorigctl has been tested on the following platforms:

- AMD64
- i386
- ARMv6
- ARMv8

and the following operating Systems:

- Linux (Ubuntu, Raspian, Armbian)
- MacOS (Sierra)

Windows should be supported in the future.
## Download

You can download a tarball / zip archive with the compiled binary for MacOS,
Linux (ARM/AMD64) and Windows from the
[releases](https://github.com/dh1tw/gorigctl/releases) page. gorigctl is
just a single exectuable.

## Installation / Dependencies

gorigctl depends on hamlib as a 3rd party library. You can either install hamlib
on Linux and MacOS through their packet managers or
[build hamlib from source]() for
the latest updates and new rigs.

### Linux (Ubuntu >= 14.04)

```bash
$ sudo apt-get install -y libhamlib2 libhamlib-dev
```

### MacOS

```bash
$ brew update
$ brew install hamlib
```

## Requirements

In order to operate your radio remotely through gorigctl, you need to either
run your own MQTT Broker ([Mosquitto](4) is a good choice) or connect to a
public broker, like `iot.eclipse.org` or `test.mosquitto.org`. The load of
these brokers and the ping to your place will influence the latency. These public
brokers are good for inital tests, however they are sometimes overloaded.

# Getting started

### Configuration

Both, the server and the client provide extensive configuration possibilities,
either through a configuration file (TOML|YAML|JSON), typically located in
your home directory `/home/your_user/.gorigctl.toml`. or through pflags.

An example configuration file named ```gorigctl.toml```is included in the
repository.

All parameters can be set through pflags. The following *example* shows the
options for

```bash
$ gorigctl server mqtt --help
```

```
MQTT server which makes a local radio available on the network

The MQTT Topics follow the Shackbus convention and must match on the
Server and the Client.

The parameters in "<>" can be set through flags or in the config file:
<station>/radios/<radio>/cat

Usage:
  gorigctl server mqtt [flags]

Flags:
  -b, --baudrate int                Baudrate (default 38400)
  -p, --broker-port int             Broker Port (default 1883)
  -u, --broker-url string           Broker URL (default "localhost")
  -d, --databits int                Databits (default 8)
  -a, --handshake string            Handshake (default "none")
  -r, --parity string               Parity (default "none")
  -t, --polling_interval duration   Timer for polling the rig's meter values [ms] (0 = disabled) (default 100ms)
  -o, --portname string             Portname (e.g. COM1) (default "/dev/mhux/cat")
  -Y, --radio string                Radio ID (default "myradio")
  -m, --rig-model int               Hamlib Rig Model ID (default 1)
  -X, --station string              Your station callsign (default "mystation")
  -s, --stopbits int                Stopbits (default 1)
  -k, --sync_interval duration      Timer for syncing all values with the rig [s] (0 = disabled) (default 3s)

Global Flags:
      --config string   config file (default is $HOME/.gorigctl.[yaml|toml|json])
```

## Start a radio server

```bash
$ gorigctl server mqtt
```

## Start the GUI for connecting to a remote radio

```bash
$ gorigctl gui mqtt
```

## Start a CLI interface for a local radio

```bash
$ gorigctl cli local
```

## How build gorigctl

The [Wiki](https://github.com/dh1tw/gorigctl/wiki) contains detailed
instructions on how to build remoteAudio from source code on Linux, MacOS and Windows.


## Known issues

- Running "gorigctl gui local" on an embedded device like the raspberry is lacking
  performance. An alternative is to launch a server and a client is seperate
  terminals / sessions.

## Troubleshooting

Feel free to open an [issue](https://github.com/dh1tw/gorigctl/issues) if you
encounter problems.