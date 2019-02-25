package rabbitmq

import (
	"fmt"
	"strings"
	"testing"

	rabbithole "github.com/michaelklishin/rabbit-hole"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTopicPermissions(t *testing.T) {
	var topicPermissionInfo rabbithole.TopicPermissionInfo
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccTopicPermissionsCheckDestroy(&topicPermissionInfo),
		Steps: []resource.TestStep{
			{
				Config: testAccTopicPermissionsConfig_basic,
				Check: testAccTopicPermissionsCheck(
					"rabbitmq_topic_permissions.test", &topicPermissionInfo,
				),
			},
			{
				Config: testAccTopicPermissionsConfig_update,
				Check: testAccTopicPermissionsCheck(
					"rabbitmq_topic_permissions.test", &topicPermissionInfo,
				),
			},
		},
	})
}

func testAccTopicPermissionsCheck(rn string, topicPermissionInfo *rabbithole.TopicPermissionInfo) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("resource not found: %s", rn)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("permission id not set")
		}

		rmqc := testAccProvider.Meta().(*rabbithole.Client)
		perms, err := rmqc.ListTopicPermissions()
		if err != nil {
			return fmt.Errorf("Error retrieving topic permissions: %s", err)
		}

		userParts := strings.Split(rs.Primary.ID, "@")
		for _, perm := range perms {
			if perm.User == userParts[0] && perm.Vhost == userParts[1] {
				topicPermissionInfo = &perm
				return nil
			}
		}

		return fmt.Errorf("Unable to find topic permissions for user %s", rn)
	}
}

func testAccTopicPermissionsCheckDestroy(topicPermissionInfo *rabbithole.TopicPermissionInfo) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rmqc := testAccProvider.Meta().(*rabbithole.Client)
		perms, err := rmqc.ListTopicPermissions()
		if err != nil {
			return fmt.Errorf("Error retrieving topic permissions: %s", err)
		}

		for _, perm := range perms {
			if perm.User == topicPermissionInfo.User && perm.Vhost == topicPermissionInfo.Vhost {
				return fmt.Errorf("Topic permissions still exist for user %s@%s", topicPermissionInfo.User, topicPermissionInfo.Vhost)
			}
		}

		return nil
	}
}

const testAccTopicPermissionsConfig_basic = `
resource "rabbitmq_vhost" "test" {
    name = "test"
}

resource "rabbitmq_user" "test" {
    name = "mctest"
    password = "foobar"
    tags = ["administrator"]
}

resource "rabbitmq_topic_permissions" "test" {
    user = "${rabbitmq_user.test.name}"
    vhost = "${rabbitmq_vhost.test.name}"
    permissions {
        exchange = ".*"
        write = ".*"
        read = ".*"
    }
}`

const testAccTopicPermissionsConfig_update = `
resource "rabbitmq_vhost" "test" {
    name = "test"
}

resource "rabbitmq_user" "test" {
    name = "mctest"
    password = "foobar"
    tags = ["administrator"]
}

resource "rabbitmq_topic_permissions" "test" {
    user = "${rabbitmq_user.test.name}"
    vhost = "${rabbitmq_vhost.test.name}"
    permissions {
        exchange = ".*"
        write = ".*"
        read = ""
    }
}`
