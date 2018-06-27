## 1.0.1 (Unreleased)

FIXES

* Fixed issue preventing policies from updating [GH-18]

## 1.0.0 (April 27, 2018)

IMPROVEMENTS:

* Allow vhost names to contain slashes ([#11](https://github.com/terraform-providers/terraform-provider-rabbitmq/issues/11))

FIXES:

* Allow integer values for policy definitions ([#13](https://github.com/terraform-providers/terraform-provider-rabbitmq/issues/13))

## 0.2.0 (September 26, 2017)

BACKWARDS INCOMPATIBILITIES / NOTES:

* Due to a bug discovered where bindings were not being correctly stored in state, `rabbitmq_bindings.properties_key` is now a read-only, computed field.

IMPROVEMENTS:

* Added `arguments_json` to `rabbitmq_queue`. This argument can accept a nested JSON string which can contain additional settings for the queue. This is useful for queue settings which have non-string values. ([#6](https://github.com/terraform-providers/terraform-provider-rabbitmq/issues/6))

FIXES:

* Fix bindings not being saved to state ([#8](https://github.com/terraform-providers/terraform-provider-rabbitmq/issues/8))
* Fix issue in `rabbitmq_user` where tags were removed when a password was changed ([#7](https://github.com/terraform-providers/terraform-provider-rabbitmq/issues/7))

## 0.1.0 (June 21, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
