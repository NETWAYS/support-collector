![go](https://github.com/NETWAYS/support-collector/actions/workflows/go.yml/badge.svg)
![go](https://github.com/NETWAYS/support-collector/actions/workflows/golangci-lint.yml/badge.svg)

# NETWAYS support collector

The support collector allows to collect relevant information from servers. The resulting ZIP file can be given to second to get an insight into the system.

> **WARNING:** Do not transfer the generated file over insecure connections, it contains potential sensitive
> information!

If you are a customer, you can contact us at [support@netways.de](mailto:support@netways.de) or
[netways.de/en/contact/](https://www.netways.de/en/contact/).

The initial idea and inspiration came from [NETWAYS/icinga2-diagnostics](https://github.com/Icinga/icinga2-diagnostics).

## Available Modules

A brief overview about the modules, you can check the source code under [modules](modules) for what exactly is collected.

Most modules check if the component is installed before trying to collect data. If the module is not detected, it will not be collected.

| Module name    | Description                                                                                                                                              |
|----------------|----------------------------------------------------------------------------------------------------------------------------------------------------------|
| ansible        | Configuration and packages                                                                                                                               |
| base           | Basic information about the system (operating system, kernel, memory, cpu, processes, repositories, firewalls, etc.)                                     |
| corosync       | Includes corosync and pacemaker. Collects configuration, logs, packages and service status                                                               |
| elastic        | Includes elasticsearch, logstash and kibana. Collects configuration, packages and service status                                                         |
| foreman        | Configuration, logs, packages and service status                                                                                                         |
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
| redis          | Configuration, logs, packages and service status                                                                                                         |
| webservers     | Includes apache2, httpd and nginx. Collects configuration, logs, packages and service status                                                             |


## Usage

`$ support-collector`

The CLI wizard will guide you through the possible arguments after calling the command. If you prefer to skip the wizard, you can use `--disable-wizard` and use the default control values.  
A more detailed control is possible through the use of an answer-file.

**Available arguments:**

| Short | Long                   | Description                                               |
|-------|------------------------|-----------------------------------------------------------|
| -f    | --answer-file          | Provide an answer-file to control the collection          |
|       | --disable-wizard       | Disable interactive wizard and use default control values |
|       | --generate-answer-file | Generate an example answer-file with default values       |
| -V    | --verbose              | Enable verbose logging                                    |
| -v    | --version              | Print version and exit                                    |

## Obfuscation

> **WARNING:** Some passwords or secrets are automatically removed, but this no guarantee, so be careful what you share!

With using an answer-file, you are able to add multiple custom obfuscators.  
As these obfuscators are based on regex, you must add a valid regex pattern that meets your requirements.

For example, `Secret:\s*(.*)` will find `Secret: DummyValue` and set it to `Secret: <hidden>`.

In addition, files and folders that follow a specific pattern are not collected. This affects all files that correspond to the following filters:  
`.*`, `*~`, `*.key`, `*.csr`, `*.crt`, and `*.pem`

## Answer File

By providing an answer-file you can customize the data collection.  
In addition to some general control values that customize the collection, further details for modules - that are not included by default - can be collected and configured.

The answer-file has to be in YAML format.  
To generate a default answer-file, you can use `--generate-answer-file`.    

To provide an answer-file, just use `--answer-file <file>.yml` or `-f <path>.yml`. With using this, the wizard will be skipped.

You can find an example answer-file [here](doc/answer-file.yml.example).

### General

Inside the general section we can configure some general behavior for the support-collector.
````yaml
general:
    outputFile: data.zip      # Name of the resulting zip file
    enabledModules: []        # List of enabled modules (Can also be 'all')
    disabledModules: []       # List of disabled modules
    extraObfuscators: []      # Custom obfuscators that should be applied
    detailedCollection: true  # Enable detailed collection
    commandTimeout: []        # Command timeout for exec commands (Default 1m0s)
````

### Icinga 2

For the module `icinga2` it is possible do define some API endpoints to collect data from.  
There is no limit of endpoints that can be defined.

````yaml
icinga2:
    endpoints:                # List of Icinga 2 API endpoint to collect data from
        - address: 127.0.0.1  # Address of endpoint
          port: 5665          # Icinga 2 port
          username: icinga    # Icinga 2 API user
          password: icinga    # Icinga 2 API password
````

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