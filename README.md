# NETWAYS support collector

<!-- NOTE: Update `Readme` const in `main.go` when changing the text here -->

The support collector allows our customers to collect relevant information from their servers.
A resulting ZIP file can then be provided to our support team for further inspection.

If you are a customer, you can contact us at [support@netways.de](mailto:support@netways.de) or
[netways.de/contact](https://www.netways.de/contact).

**WARNING:** DO NOT transfer the generated file over insecure connections or by
email, it contains potential sensitive information!

Inspired by [icinga2-diagnostics](https://github.com/Icinga/icinga2-diagnostics).

## Usage

**Warning:** Currently no anonymization is implemented, so be careful.

By default, we collect all we can find. You can control this by only enabling certain modules, or disabling some.

If you want to see what is collected, add `--verbose`.

```
Usage of ./support-collector:
  -o, --output string     Output file for the ZIP content (default "netways-support.zip")
      --enable strings    List of enabled module (default [base,icinga2,icingaweb2,icinga-director,mysql])
      --disable strings   List of disabled module
  -v, --verbose           Enable verbose logging
  -V, --version           Print version and exit
```

## Modules

A brief overview about the modules, you can check the source code under [modules](modules) for what exactly is collected.

Most modules check if the component is installed before trying to collect data.

### Base

Will collect basic information about your system:

* Kernel and Operating system versions
* CPU, memory, disk and other hardware and vendor information
* Current process and load status
* Status of AppArmor and SELinux

### Icinga 2

* Configuration from `/etc/icinga2`
* Package information
* Service status
* Config check result
* Log files from `/var/log/icinga2`
* Object list for `Zone` and `Endpoint`
* Variables like `NodeName` and `ZoneName`

### Icinga Web 2

* Configuration from `/etc/icingaweb2`
* Package information
* Log files from `/var/log/icingaweb2`
* Version for icingaweb2 and its modules, including Git information
* Installed PHP packages and php-fpm service status
* Installed webserver packages

### Icinga Director

* Package or Git version information
* Service status
* Health status

### MySQL

Is looking for standard MySQL or MariaDB installations.

* Package versions
* Service status
* Configuration files from `/etc/my*` (depending on the known paths)

## Supported systems

**Note:** Currently, only Linux is supported by the collector.

Distribution    | Supported | Tested | Comment
----------------|-----------|--------|--------
CentOS / EL     | ✅️ | CentOS 7️ | Should work on all similar enterprise platforms
Debian          | ✅ |
Ubuntu          | ✅ |
SLES / OpenSUSE | ✅ |

## License

Copyright (C) 2021 [NETWAYS GmbH](mailto:info@netways.de)

This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public
License as published by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied
warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not,
see <https://www.gnu.org/licenses/>.
