package rabbitmq

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	rabbithole "github.com/michaelklishin/rabbit-hole"
)

func TestAccFederationUpstream(t *testing.T) {
	var upstream rabbithole.FederationUpstream
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccFederationUpstreamCheckDestroy(&upstream),
		Steps: []resource.TestStep{
			{
				Config: testAccFederationUpstreamCreate,
				Check: testAccFederationUpstreamCheck(
					"rabbitmq_federation_upstream.foo", &upstream,
				),
			},
			{
				Config: testAccFederationUpstreamUpdate,
				Check: testAccFederationUpstreamCheck(
					"rabbitmq_federation_upstream.foo", &upstream,
				),
			},
		},
	})
}

func testAccFederationUpstreamCheck(rn string, upstream *rabbithole.FederationUpstream) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("resource not found: %s", rn)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("federation upstream id not set")
		}

		id := strings.Split(rs.Primary.ID, "@")
		name := id[0]
		vhost := id[1]

		rmqc := testAccProvider.Meta().(*rabbithole.Client)
		upstreams, err := rmqc.ListFederationUpstreamsIn(vhost)
		if err != nil {
			return fmt.Errorf("Error retrieving federation upstreams: %s", err)
		}

		for _, up := range upstreams {
			if up.Name == name && up.Vhost == vhost {
				upstream = &up
				return nil
			}
		}

		return fmt.Errorf("Unable to find federation upstream %s", rn)
	}
}

func testAccFederationUpstreamCheckDestroy(upstream *rabbithole.FederationUpstream) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rmqc := testAccProvider.Meta().(*rabbithole.Client)

		upstreams, err := rmqc.ListFederationUpstreamsIn(upstream.Vhost)
		if err != nil {
			return fmt.Errorf("Error retrieving federation upstreams: %s", err)
		}

		for _, up := range upstreams {
			if up.Name == upstream.Name && up.Vhost == upstream.Vhost {
				return fmt.Errorf("Federation upstream %s@%s still exists", upstream.Name, upstream.Vhost)
			}
		}

		return nil
	}
}

const testAccFederationUpstreamCreate = `
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

				exchange = ""
				max_hops = 1
				expires = 0
				message_ttl = 0

				queue = ""
		}
}
`

const testAccFederationUpstreamUpdate = `
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
				prefetch_count = 500
				reconnect_delay = 10
				ack_mode = "on-publish"
				trust_user_id = true

				exchange = ""
				max_hops = 2
				expires = 1800000
				message_ttl = 60000

				queue = ""
		}
}
`
