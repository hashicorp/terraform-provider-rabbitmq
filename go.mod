module github.com/terraform-providers/terraform-provider-rabbitmq

require (
	github.com/hashicorp/terraform-plugin-sdk v1.0.0
	github.com/michaelklishin/rabbit-hole/v2 v2.2.1-0.20200601180354-b5a90e068691
)

replace github.com/michaelklishin/rabbit-hole/v2 => github.com/niclic/rabbit-hole/v2 v2.1.1-0.20200426194252-fccfa4cf97f4

go 1.13
