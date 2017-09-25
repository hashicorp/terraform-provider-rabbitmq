## 0.1.1 (Unreleased)

BACKWARDS INCOMPATIBILITIES / NOTES:

* Due to a bug discovered where bindings were not being correctly stored in state, `rabbitmq_bindings.properties_key` is now a read-only, computed field.

FIXES:

* Fix bindings not being saved to state [GH-8]
* Fix issue in `rabbitmq_user` where tags were removed when a password was changed [GH-7]

## 0.1.0 (June 21, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
