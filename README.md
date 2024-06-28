![go](https://github.com/NETWAYS/support-collector/actions/workflows/go.yml/badge.svg)
![go](https://github.com/NETWAYS/support-collector/actions/workflows/golangci-lint.yml/badge.svg)

# NETWAYS support collector

The support collector allows to collect relevant information from servers. The resulting ZIP file can be given to second
to get an insight into the system.

> **WARNING:** DO NOT transfer the generated file over insecure connections, it contains potential sensitive
> information!

If you are a customer, you can contact us at [support@netways.de](mailto:support@netways.de) or
[netways.de/en/contact/](https://www.netways.de/en/contact/).

Inspired by [NETWAYS/icinga2-diagnostics](https://github.com/Icinga/icinga2-diagnostics).

## Usage

> **WARNING:** Some passwords or secrets are automatically removed, but this no guarantee, so be careful what you share!

The `--hide` flag can be used multiple times to hide sensitive data, it supports regular expressions.

`# support-collector --hide "Secret:.*" --hide "Password:.*"`

In addition, files and folders that follow a specific pattern are not collected. This affects all files that correspond to the following filters:  
`.*`, `*~`, `*.key`, `*.csr`, `*.crt`, and `*.pem`

By default, we collect all we can find. You can control this by only enabling certain modules, or disabling some.

If you want to see what is collected, add `--verbose`

To collect advanced data for module `Icinga 2`, you can use the Icinga 2 API to collect data from all endpoints provided.
The API requests are performed with a global API user you have to create yourself. Just create that user in a global zone like 'director-global'


| Short | Long                    | Description                                                                                                                                   |
|:-----:|:------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------|
|  -o   | --output                | Output file for the zip content (default: current directory and named like '$HOSTNAME'-netways-support-$TIMESTAMP.zip)                        |
|       | --nodetails             | Disable detailed collection including logs and more                                                                                           |
|       | --enable                | List of enabled modules (default: all)                                                                                                        |
|       | --disable               | List of disabled modules (default: none)                                                                                                      |
|       | --hide                  | List of keywords to obfuscate. Can be used multiple times                                                                                     |
|       | --command-timeout       | Timeout for command execution in modules (default: 1m0s)                                                                                      |
|       | --icinga2-api-user      | Username of global Icinga 2 API user to collect data about Icinga 2 Infrastructure                                                            |
|       | --icinga2-api-pass      | Password for global Icinga 2 API user to collect data about Icinga 2 Infrastructure                                                           |
|       | --icinga2-api-endpoints | List of Icinga 2 API Endpoints (including port) to collect data from. FQDN or IP address must be reachable. (Example: i2-master01.local:5665) |
|  -v   | --verbose               | Enable verbose logging                                                                                                                        |
|  -V   | --version               | Print version and exit                                                                                                                        |

## Modules

A brief overview about the modules, you can check the source code under [modules](modules) for what exactly is
collected.

Most modules check if the component is installed before trying to collect data. If the module is not detected, it will
not be collected.

| Module name    | Description                                                                                                                                              |
|----------------|----------------------------------------------------------------------------------------------------------------------------------------------------------|
| ansible        | Configuration and packages                                                                                                                               |
| base           | Basic information about the system (operating system, kernel, memory, cpu, processes, repositories, firewalls, etc.)                                     |
| corosync       | Includes corosync and pacemaker. Collects configuration, logs, packages and service status                                                               |
| elastic        | Includes elasticsearch, logstash and kibana. Collects configuration, packages and service status                                                         |
| grafana        | Configuration, logs, plugins, packages and service status                                                                                                |
| graphite       | Includes graphite and carbon. Collects configuration, logs, python / pip version and list, packages and service status                                   |
| graylog        | Configuration, packages and service status                                                                                                               |
| icinga2        | Configuration, packages, service status, logs, Icinga 2 objects, Icinga 2 variables, plugins, icinga-installer and data from API endpoints (if provided) |
| icingadb       | Includes IcingaDB and IcingaDB redis. Collects configuration, logs, packages and service status                                                          |
| icingadirector | Packages or git information, logs, Director health status and service status                                                                             |
| icingaweb2     | Configuration, logs, packages, modules, PHP, modules and service status                                                                                  |
| influxdb       | Configuration, logs, packages and service status                                                                                                         |
| keepalived     | Configuration, packages and service status                                                                                                               |
| mongodb        | Configuration, logs, packages and service status                                                                                                         |
| mysql          | Configuration, logs, packages and service status                                                                                                         |
| postgresql     | Configuration, logs, packages and service status                                                                                                         |
| prometheus     | Configuration, packages and service status                                                                                                               |
| puppet         | Configuration, logs, module list, packages and service status                                                                                            |
| webservers     | Includes apache2, httpd and nginx. Collects configuration, logs, packages and service status                                                             |
| foreman        | Configuration, logs, packages and service status                                                                                                         |

## Supported systems

| Distribution    | Tested on                | Supported |
|-----------------|--------------------------|:---------:|
| CentOS / EL     | CentOS 7/8, RHEL 7/8     |     ✅     |
| Debian          | Debian 10/11             |     ✅     |
| Ubuntu          | Ubuntu 18.04/20.04/22.04 |     ✅     |
| SLES / OpenSUSE | openSUSE Leap 15.4       |     ✅     |

## License

Copyright (C) 2021 [NETWAYS GmbH](mailto:info@netways.de)

This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public
License as published by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied
warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not,
see <https://www.gnu.org/licenses/>.