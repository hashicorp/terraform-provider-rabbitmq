---
layout: "rabbitmq"
page_title: "RabbitMQ: rabbitmq_shovel"
sidebar_current: "docs-rabbitmq-resource-shovel"
description: |-
  Creates and manages a shovel on a RabbitMQ server.
---

# rabbitmq\_shovel

The ``rabbitmq_shovel`` resource creates and manages a dynamic shovel.

## Example Usage

```hcl
resource "rabbitmq_vhost" "test" {
    name = "test"
}

resource "rabbitmq_exchange" "test" {
    name = "test_exchange"
    vhost = "${rabbitmq_vhost.test.name}"
    settings {
        type = "fanout"
        durable = false
        auto_delete = true
    }
}

resource "rabbitmq_queue" "test" {
	name = "test_queue"
	vhost = "${rabbitmq_vhost.test.name}"
	settings {
		durable = false
		auto_delete = true
	}
}

resource "rabbitmq_shovel" "shovelTest" {
	name = "shovelTest"
	vhost = "${rabbitmq_vhost.test.name}"
	info {
		source_uri = "amqp:///test"
		source_exchange = "${rabbitmq_exchange.test.name}"
		source_exchange_key = "test"
		destination_uri = "amqp:///test"
		destination_queue = "${rabbitmq_queue.test.name}"
	}
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The shovel name.

* `vhost` - (Required) The vhost to create the resource in.

* `info` - (Required) The settings of the dynamic shovel. The structure is
  described below.

The `info` block supports:

### Essential parameters

* `source_uri` - (Required) The amqp uri for the source.

* `source_protocol` - (Optional) The protocol (`amqp091` or `amqp10`) to use when connecting to the source.
Defaults to `amqp091`.

* `source_queue` - (Optional) The queue from which to consume.
Either this or `source_exchange` must be specified but not both.

* `destination_uri` - (Required) The amqp uri for the destination .

* `destination_protocol` - (Optional) The protocol (`amqp091` or `amqp10`) to use when connecting to the destination.
Defaults to `amqp091`.

* `destination_queue` - (Optional) The queue to which messages should be published.
Either this or `destination_exchange` must be specified but not both.

### Optional parameters

* `ack_mode` - (Optional) Determines how the shovel should acknowledge messages. Possible values are: `on-confirm`, `on-publish` and `no-ack`.
Defaults to `on-confirm`.

* `add_forward_headers` - (Optional; **Deprecated**, please use `destination_add_forward_headers`) Whether to add `x-shovelled` headers to shovelled messages.

* `delete_after` - (Optional; **Deprecated**, please use `source_delete_after`) Determines when (if ever) the shovel should delete itself. Possible values are: `never`, `queue-length` or an integer.

* `destination_add_forward_headers` - (Optional) Whether to add `x-shovelled` headers to shovelled messages.

* `destination_add_timestamp_headers` - (Optional) Whether to add `x-shovelled-timestamp` headers to shovelled messages.
Defaults to `false`.

* `destination_exchange` - (Optional) The exchange to which messages should be published.
Either this or `destination_queue` must be specified but not both.

* `destination_exchange_key` - (Optional) The routing key when using `destination_exchange`.

* `destination_publish_properties` - (Optional) A map of properties to overwrite when shovelling messages.

* `prefetch_count` - (Optional; **Deprecated**, please use `source_prefetch_count`) The maximum number of unacknowledged messages copied over a shovel at any one time.

* `reconnect_delay` - (Optional) The duration in seconds to reconnect to a broker after disconnected.
Defaults to `1`.

* `source_delete_after` - (Optional) Determines when (if ever) the shovel should delete itself. Possible values are: `never`, `queue-length` or an integer.

* `source_exchange` - (Optional) The exchange from which to consume.
Either this or `source_queue` must be specified but not both.

* `source_exchange_key` - (Optional) The routing key when using `source_exchange`.

* `source_prefetch_count` - (Optional) The maximum number of unacknowledged messages copied over a shovel at any one time.

### AMQP 1.0 specific parameters

* `source_address` - (Optional) The AMQP 1.0 source link address.

* `destination_address` - (Optional) The AMQP 1.0 destination link address.

* `destination_application_properties` - (Optional) Application properties to set when shovelling messages.

* `destination_properties` - (Optional) Properties to overwrite when shovelling messages.

For more details regarding dynamic shovel parameters please have a look at the official reference documentaion at [RabbitMQ: Configuring Dynamic Shovels](https://www.rabbitmq.com/shovel-dynamic.html).

## Attributes Reference

No further attributes are exported.

## Import

Shovels can be imported using the `name` and `vhost`
E.g.

```
terraform import rabbitmq_shovel.test shovelTest@test
```
