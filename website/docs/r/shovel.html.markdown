---
layout: "rabbitmq"
page_title: "RabbitMQ: rabbitmq_shovel"
sidebar_current: "docs-rabbitmq-resource-shovel"
description: |-
  Creates and manages a shovel on a RabbitMQ server.
---

# rabbitmq\_shovel

The ``rabbitmq_shovel`` resource creates and manages a shovel.

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

* `info` - (Required) The settings of the shovel. The structure is
  described below.

The `info` block supports:

* `source_uri` - (Required) The amqp uri for the source.

* `source_exchange` - (Optional) The exchange from which to consume.
Either this or source_queue must be specified but not both.

* `source_exchange_key` - (Optional) The routing key when using source_exchange.

* `source_queue` - (Optional) The queue from which to consume.
Either this or source_exchange must be specified but not both.

* `destination_uri` - (Required) The amqp uri for the destination .

* `destination_exchange` - (Optional) The exchange to which messages should be published.
Either this or destination_queue must be specified but not both.

* `destination_exchange_key` - (Optional) The routing key when using destination_exchange.

* `destination_queue` - (Optional) The queue to which messages should be published.
Either this or destination_exchange must be specified but not both.

* `prefetch_count` - (Optional) The maximum number of unacknowledged messages copied over a shovel at any one time.
Defaults to `1000`.

* `reconnect_delay` - (Optional) The duration in seconds to reconnect to a broker after disconnected.
Defaults to `1`.

* `add_forward_headers` - (Optional) Whether to amqp shovel headers.
Defaults to `false`.

* `ack_mode` - (Optional) Determines how the shovel should acknowledge messages.
Defaults to `on-confirm`.

* `delete_after` - (Optional) Determines when (if ever) the shovel should delete itself .
Defaults to `never`.

## Attributes Reference

No further attributes are exported.

## Import

Shovels can be imported using the `name` and `vhost`
E.g.

```
terraform import rabbitmq_shovel.test shovelTest@test
```
