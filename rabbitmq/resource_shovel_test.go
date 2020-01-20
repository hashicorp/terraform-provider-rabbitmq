package rabbitmq

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	rabbithole "github.com/michaelklishin/rabbit-hole"
)

func TestAccShovel(t *testing.T) {
	var shovel string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccShovelCheckDestroy(shovel),
		Steps: []resource.TestStep{
			{
				Config: testAccShovelConfig_basic,
				Check: testAccShovelCheck(
					"rabbitmq_shovel.shovelTest", &shovel,
				),
			},
		},
	})
}

func testAccShovelCheck(rn string, name *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("resource not found: %s", rn)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("shovel id not set")
		}

		rmqc := testAccProvider.Meta().(*rabbithole.Client)
		shovelInfos, err := rmqc.ListShovels()
		if err != nil {
			return fmt.Errorf("Error retrieving shovels: %s", err)
		}

		for _, info := range shovelInfos {
			if info.Name == rs.Primary.ID {
				*name = rs.Primary.ID
				return nil
			}
		}

		return fmt.Errorf("Unable to find shovel %s", rn)
	}
}

func testAccShovelCheckDestroy(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rmqc := testAccProvider.Meta().(*rabbithole.Client)
		shovelInfos, err := rmqc.ListShovels()
		if err != nil {
			return fmt.Errorf("Error retrieving shovels: %s", err)
		}

		for _, info := range shovelInfos {
			if info.Name == name {
				return fmt.Errorf("shovel still exists: %v", info)
			}
		}

		return nil
	}
}

const testAccShovelConfig_basic = `
resource "rabbitmq_vhost" "test" {
    name = "test"
}

resource "rabbitmq_exchange" "test" {
    name = "test_exchange"
    vhost = "${rabbitmq_vhost.test.name}"
    settings {
        type = "fanout"
        durable = false
        auto_delete = true
    }
}

resource "rabbitmq_queue" "test" {
	name = "test_queue"
	vhost = "${rabbitmq_vhost.test.name}"
	settings {
		durable = false
		auto_delete = true
	}
}

resource "rabbitmq_shovel" "shovelTest" {
	name = "shovelTest"
	vhost = "${rabbitmq_vhost.test.name}"
	info {
		source_uri = "amqp:///test"
		source_exchange = "${rabbitmq_exchange.test.name}"
		source_exchange_key = "test"
		destination_uri = "amqp:///test"
		destination_queue = "${rabbitmq_queue.test.name}"
	}
}`
