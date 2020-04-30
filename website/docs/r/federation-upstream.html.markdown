---
layout: "rabbitmq"
page_title: "RabbitMQ: rabbitmq_federation_upstream"
sidebar_current: "docs-rabbitmq-resource-federation-upstream"
description: |-
  Creates and manages a federation upstream on a RabbitMQ server.
---

# rabbitmq\_federation\_upstream

The ``rabbitmq_federation_upstream`` resource creates and manages a federation upstream parameter.

## Example Usage

```hcl
resource "rabbitmq_vhost" "test" {
  name = "test"
}

resource "rabbitmq_permissions" "guest" {
  user = "guest"
  vhost = rabbitmq_vhost.test.name
  permissions {
    configure = ".*"
    write = ".*"
    read = ".*"
  }
}

resource "rabbitmq_federation_upstream" "foo" {
  name = "foo"
  vhost = rabbitmq_permissions.guest.vhost

  definition {
    uri = "amqp://server-name"
    prefetch_count = 1000
    reconnect_delay = 1
    ack_mode = "on-confirm"
    trust_user_id = false
    max_hops = 1
  }
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the federation upstream.

* `vhost` - (Required) The vhost to create the resource in.

* `component` - (Computed) Set to `federation-upstream` by the underlying RabbitMQ provider. You do not set this attribute but will see it in state and plan output.

* `definition` - (Required) The configuration of the federation upstream. The structure is described below.

The `definition` block supports the following arguments:

Applicable to Both Federated Exchanges and Queues

* `uri` - (Required) The AMQP URI(s) for the upstream. Note that the URI may contain sensitive information, such as a password.
* `prefetch_count` - (Optional) Maximum number of unacknowledged messages that may be in flight over a federation link at one time.
* `reconnect_delay` - (Optional) Time in seconds to wait after a network link goes down before attempting reconnection.
* `ack_mode` - (Optional) Determines how the link should acknowledge messages. Valid values are `on-confirm`, `on-publish`, and `no-ack`.
* `trust_user_id` - (Optional) Determines how federation should interact with the validated user-id feature.

Applicable to Federated Exchanges Only

* `exchange` - (Optional)  The name of the upstream exchange.
* `max_hops` - (Optional) Maximum number of federation links that messages can traverse before being dropped.
* `expires` - (Optional) The expiry time (in milliseconds) after which an upstream queue for a federated exchange may be deleted if a connection to the upstream is lost.
* `message_ttl` - (Optional) The expiry time (in milliseconds) for messages in the upstream queue for a federated exchange (see expires).

Applicable to Federated Queues Only

* `queue` - (Optional) The name of the upstream queue.

Consult the RabbitMQ [Federation Reference](https://www.rabbitmq.com/federation-reference.html) documentation for detailed information and guidance on setting these values.

## Attributes Reference

No further attributes are exported.

## Import

A Federation upstream can be imported using the `id` which is composed of `name@vhost`.
E.g.

```sh
terraform import rabbitmq_federation_upstream.foo foo@test
```
