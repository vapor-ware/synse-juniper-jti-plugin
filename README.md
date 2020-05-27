[![Build Status](https://build.vio.sh/buildStatus/icon?job=vapor-ware/synse-juniper-jti-plugin/master)](https://build.vio.sh/blue/organizations/jenkins/vapor-ware%2Fsynse-juniper-jti-plugin/activity)

# Synse Juniper JTI Plugin

Synse plugin for Juniper JTI metrics over UDP stream.

This plugin ingests Juniper metrics over a UDP stream, making them available to the
Synse platform. Juniper supports many sensors/resources, not all of which are yet
supported by the plugin. Current plugin capabilities include support for the following
resources:

* `/junos/system/linecard/interface/`
* `/junos/system/linecard/optics/`

The plugin has been tested against Juniper routers running Junos OS 18.2R3. Compatibility
with other Junos versions is not guaranteed.

## Getting Started

### Getting

You can install the plugin via a [release](https://github.com/vapor-ware/synse-juniper-jti-plugin/releases)
binary or via Docker image

```
docker pull vaporio/juniper-jti-plugin
```

If you wish to use a development build, fork/clone the repo and build the plugin
from source.

### Running

A [compose file](docker-compose.yaml) is included in this repo which provides a basic example of
how to run the Juniper JTI plugin in conjunction with Synse Server. You may also run the plugin  on
its own

```
docker run -d \
    --name juniper-jti \
    -p 5010:5010 \
    -p 5566:5566/udp \
    -v ./config.yaml:/etc/synse/plugin/config/config.yaml \
    vaporio/juniper-jti-plugin
```

and use the [Synse CLI](https://github.com/vapor-ware/synse-cli) to query the plugin's gRPC API.

The Juniper JTI plugin will run with minimal configuration (e.g. with config for the UDP server), but
will not provide any data unless JTI data is streamed to the exposed UDP server. As such, the example
deployment will run, but will not provide any useful data.

## Juniper JTI Plugin Configuration

Plugin and device configuration are described in detail in the [SDK Documentation](https://synse.readthedocs.io/en/latest/sdk/intro/).

When deploying, you will need to provide your own plugin configuration (`config.yaml`)
with dynamic configuration defined. This is how the Juniper JTI plugin's UDP server is
configured, allowing Juniper devices to stream telemetry data to it.

As an example:

```yaml
dynamicRegistration:
  config:
  - address: udp://0.0.0.0:5566
```

Note that the IP address in this example is `0.0.0.0`. When running the plugin via a Docker
container, you will want to use this address so it is able to correctly capture the incoming packets.

### Dynamic Registration Options

Below are the fields that are expected in each of the dynamic registration items.
If no default is specified (`-`), the field is required.

| Field   | Description | Default |
| ------- | ----------- | ------- |
| address | The protocol/address/port for the UDP server to listen for incoming telemetry data. The protocol may be one of: [`udp`, `udp4`, `udp6`]. When running in a docker container, the address should be `0.0.0.0`. | `-` |

### Reading Outputs

Outputs are referenced by name. A single device may have more than one instance
of an output type. A value of `-` in the table below indicates that there is no value
set for that field. The *custom* section describes outputs which this plugin defines
while the *built-in* section describes outputs this plugin uses which are [built-in to
the SDK](https://synse.readthedocs.io/en/latest/sdk/concepts/reading_outputs/#built-ins).

**Custom**

| Name               | Description                                                                         | Unit    | Type         | Precision |
| ------------------ | ----------------------------------------------------------------------------------- | :-----: | ------------ | :-------: |
| boolean            | A true/false value.                                                                 | -       | `bool`       | -         |
| bytes              | A count of bytes. This is not associated with any time scale.                       | bytes   | `counter`    | -         |
| bytes-per-second   | The rate of bytes over a second.                                                    | bytes/s | `throughput` | -         |
| decibel-milliwatt  | A measure of absolute power expressed as a ratio between decibels to one milliwatt. | dBm     | `power`      | -         |
| megabit-per-second | The rate of 1,000,000 bits over a second.                                           | Mbit/s  | `throughput` | -         |
| milliampere        | A measure of electric current, in thousandths of an Ampere.                         | mA      | `current`    | -         |
| packets            | A count of packets. This is not associated with any time scale.                     | pkts    | `counter`    | -         |
| packets-per-second | The rate of packets over a second.                                                  | pkts/s  | `throughput` | -         |
| time-ticks         | A measure of time, described in "time ticks".                                       | ticks   | `time`       | -         |

**Built-in**

| Name          | Description                                   | Unit  | Type          | Precision |
| ------------- | --------------------------------------------- | :---: | ------------- | :-------: |
| number        | An arbitrary, unit-less number.               | -     | `number`      | 2         |
| status        | A generic description of status.              | -     | `status`      | -         |
| string        | A generic output for string data.             | -     | `string`      | -         |
| temperature   | A measure of temperature, in degrees Celsius. | C     | `temperature` | 2         |
| timestamp     | A string describing a timestamp.              | -     | `timestamp`   | -         |

### Device Handlers

Device Handlers are referenced by name.

| Name | Description                        | Outputs | Read  | Write | Bulk Read | Listen |
| ---- | ---------------------------------- | ------- | :---: | :---: | :-------: | :----: |
| jti  | A handler for all Juniper devices. | -       | ✓     | ✗     | ✗         | ✗      |

### Write Values

This plugin does not support writing values to devices.

## Compatibility

Below is a table describing the compatibility of plugin versions with Synse platform versions.

|             | Synse v2 | Synse v3 |
| ----------- | -------- | -------- |
| plugin v0.x | ✗        | ✓        |

## Troubleshooting

### Debugging

The plugin can be run in debug mode for additional logging. This is done by:

- Setting the `debug` option  to `true` in the plugin configuration YAML ([config.yml](example/config.yaml))

  ```yaml
  debug: true
  ```

- Passing the `--debug` flag when running the binary/image

  ```
  docker run vaporio/juniper-jti-plugin --debug
  ```

- Running the image with the `PLUGIN_DEBUG` environment variable set to `true`

  ```
  docker run -e PLUGIN_DEBUG=true vaporio/juniper-jti-plugin
  ```

### Developing

A [development/debug Dockerfile](Dockerfile.dev) is provided in the project repository to enable
building image which may be useful when developing or debugging a plugin. Unlike the slim `scratch`-based
production image, the development image uses an ubuntu base, bringing with it all the standard command line
tools one would expect. To build a development image:

```
make docker-dev
```

The built image will be tagged using the format `dev-{COMMIT}`, where `COMMIT` is the short commit for
the repository at the time. This image is not published as part of the CI pipeline, but those with access
to the Docker Hub repo may publish manually.

## Contributing / Reporting

If you experience a bug, would like to ask a question, or request a feature, open a
[new issue](https://github.com/vapor-ware/synse-juniper-jti-plugin/issues) and provide as much
context as possible. All contributions, questions, and feedback are welcomed and appreciated.

## For Developers

### Notes on (re)compiling the .proto source files

I either do not know enough about how `protoc` works, or am just otherwise struggling
to get each source file to compile to a package correctly.

At present, running `$ ./scripts/gen_proto.sh` will generate the compiled Go source
for the proto files. The issue lies with the import paths.  The generated Go source
is put into `pkg/protocol/jti/protos/{name}/{name}.pb.go`. All files (except
`telemetry_top`) import `telemetry_top`, but the import is wrong as generated. As such,
it needs to be updated from

```go
import (
	telemetry_top "protos/telemetry_top"
)
``` 

to

```go
import (
	telemetry_top "github.com/vapor-ware/synse-juniper-jti-plugin/pkg/protocol/jti/protos/telemetry_top"
)
``` 

Note also that the `.proto` files should have a `go_import` option added, e.g. for `port.proto`:

```proto
option go_package = "protos/port";
```

# License

The Synse Juniper JTI Plugin is licensed under GPLv3. See [LICENSE](LICENSE) for more info.
