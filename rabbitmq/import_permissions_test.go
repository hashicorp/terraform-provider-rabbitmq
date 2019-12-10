package rabbitmq

import (
	"testing"

	"github.com/michaelklishin/rabbit-hole"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccPermissions_importBasic(t *testing.T) {
	resourceName := "rabbitmq_permissions.test"
	var permissionInfo rabbithole.PermissionInfo

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPermissionsCheckDestroy(&permissionInfo),
		Steps: []resource.TestStep{
			{
				Config: testAccPermissionsConfig_basic,
				Check: testAccPermissionsCheck(
					resourceName, &permissionInfo,
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
