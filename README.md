# naos-esp

**The [Networked Artifacts Operating System](https://github.com/shiftr-io/naos) component for the ESP platform.**

## Requirements

- The [esp-mqtt](https://github.com/256dpi/esp-mqtt) components must be installed alongside `naos-esp`.
- Bluetooth must be enabled via `menuconfig`.
- At least two OTA partitions must be configured via `menuconfig`.

## Device Management Protocol

The following device management protocol is based on an "action down, data up" principle and tries to be as slim as possible to ease implementation. 

### Discovery

All devices will subscribe to the global `naos/collect` topic and publish to the global `naos/announcement` topic if requested. The published data has the following format: `device_type,firmware_version,device_name,base_topic`.

### Heartbeat

All devices will periodically publish to the local `naos/heartbeat` topic. The published data has the following format: `device_type,firmware_version,device_name,free_heap,up_time,running_partition`.

### Configuration

All devices will subscribe to the local `naos/get/+`, `naos/set/+` and `naos/unset/+` topics to handle configuration read, write and delete requests. The device will acknowledge read and write actions by responding with the saved value on the local topic `naos/value/+`.

### Remote Logging

All devices will subscribe to the local `naos/record` topic and enable or disable remote logging when `on` or `off` is supplied
 respectively. The device will send log messages to the local `naos/log` topic if remote logging is activated.

### Firmware Update

All devices will subscribe to the local `naos/update/begin` topic and wait for an update request of the format `total_size`. If one is received the device will publish a request for the next chunk of data by publishing a `max_size` message to the local `naos/update/next` topic. The other end will then publish the next chunk of data to the local `naos/update/write` topic and wait for the next chunk request. When all data has been written the updater publishes a message to the local `naos/update/finish` topic to make the device commit the update and restart to the newly received firmware.
