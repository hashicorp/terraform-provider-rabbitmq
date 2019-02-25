package rabbitmq

import (
	"testing"

	rabbithole "github.com/michaelklishin/rabbit-hole"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTopicPermissions_importBasic(t *testing.T) {
	resourceName := "rabbitmq_topic_permissions.test"
	var topicPermissionInfo rabbithole.TopicPermissionInfo

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccTopicPermissionsCheckDestroy(&topicPermissionInfo),
		Steps: []resource.TestStep{
			{
				Config: testAccTopicPermissionsConfig_basic,
				Check: testAccTopicPermissionsCheck(
					resourceName, &topicPermissionInfo,
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
