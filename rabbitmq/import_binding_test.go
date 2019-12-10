package rabbitmq

import (
	"testing"

	"github.com/michaelklishin/rabbit-hole"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccBinding_importBasic(t *testing.T) {
	resourceName := "rabbitmq_binding.test"
	var bindingInfo rabbithole.BindingInfo

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccBindingCheckDestroy(bindingInfo),
		Steps: []resource.TestStep{
			{
				Config: testAccBindingConfig_basic,
				Check: testAccBindingCheck(
					resourceName, &bindingInfo,
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
