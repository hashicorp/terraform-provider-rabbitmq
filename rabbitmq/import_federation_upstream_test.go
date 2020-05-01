package rabbitmq

import (
	"testing"

	rabbithole "github.com/michaelklishin/rabbit-hole/v2"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccFederationUpstream_importBasic(t *testing.T) {
	resourceName := "rabbitmq_federation_upstream.foo"
	var upstream rabbithole.FederationUpstream

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccFederationUpstreamCheckDestroy(&upstream),
		Steps: []resource.TestStep{
			{
				Config: testAccFederationUpstream_create,
				Check: testAccFederationUpstreamCheck(
					resourceName, &upstream,
				),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
