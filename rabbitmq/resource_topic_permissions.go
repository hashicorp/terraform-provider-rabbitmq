package rabbitmq

import (
	"fmt"
	"log"
	"strings"

	rabbithole "github.com/michaelklishin/rabbit-hole"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTopicPermissions() *schema.Resource {
	return &schema.Resource{
		Create: CreateTopicPermissions,
		Update: UpdateTopicPermissions,
		Read:   ReadTopicPermissions,
		Delete: DeleteTopicPermissions,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"user": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"vhost": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "/",
				ForceNew: true,
			},

			"permissions": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"exchange": {
							Type:     schema.TypeString,
							Required: true,
						},

						"write": {
							Type:     schema.TypeString,
							Required: true,
						},

						"read": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func CreateTopicPermissions(d *schema.ResourceData, meta interface{}) error {
	rmqc := meta.(*rabbithole.Client)

	user := d.Get("user").(string)
	vhost := d.Get("vhost").(string)
	permsList := d.Get("permissions").([]interface{})

	for _, exchange := range permsList {

		permsMap, ok := exchange.(map[string]interface{})
		if !ok {
			return fmt.Errorf("Unable to parse permissions")
		}

		if err := setTopicPermissionsIn(rmqc, vhost, user, permsMap); err != nil {
			return err
		}
	}

	id := fmt.Sprintf("%s@%s", user, vhost)
	d.SetId(id)

	return ReadTopicPermissions(d, meta)
}

func ReadTopicPermissions(d *schema.ResourceData, meta interface{}) error {
	rmqc := meta.(*rabbithole.Client)

	permissionId := strings.Split(d.Id(), "@")
	if len(permissionId) < 2 {
		return fmt.Errorf("Unable to determine Permission ID")
	}

	user := permissionId[0]
	vhost := permissionId[1]

	userPerms, err := rmqc.GetTopicPermissionsIn(vhost, user)
	if err != nil {
		return checkDeleted(d, err)
	}

	log.Printf("[DEBUG] RabbitMQ: Topic permission retrieved for %s: %#v", d.Id(), userPerms)

	d.Set("user", userPerms[0].User)
	d.Set("vhost", userPerms[0].Vhost)

	perms := make([]map[string]interface{}, len(userPerms))
	for i, perm := range userPerms {
		p := make(map[string]interface{})
		p["exchange"] = perm.Exchange
		p["write"] = perm.Write
		p["read"] = perm.Read
		perms[i] = p
	}
	d.Set("permissions", perms)

	return nil
}

func UpdateTopicPermissions(d *schema.ResourceData, meta interface{}) error {
	rmqc := meta.(*rabbithole.Client)

	permissionId := strings.Split(d.Id(), "@")
	if len(permissionId) < 2 {
		return fmt.Errorf("Unable to determine Permission ID")
	}

	user := permissionId[0]
	vhost := permissionId[1]

	if d.HasChange("permissions") {
		_, newPerms := d.GetChange("permissions")

		newPermsList := newPerms.([]interface{})
		for _, exchange := range newPermsList {
			permsMap, ok := exchange.(map[string]interface{})
			if !ok {
				return fmt.Errorf("Unable to parse permissions")
			}

			if err := setTopicPermissionsIn(rmqc, vhost, user, permsMap); err != nil {
				return err
			}
		}
	}

	return ReadTopicPermissions(d, meta)
}

func DeleteTopicPermissions(d *schema.ResourceData, meta interface{}) error {
	rmqc := meta.(*rabbithole.Client)

	permissionId := strings.Split(d.Id(), "@")
	if len(permissionId) < 2 {
		return fmt.Errorf("Unable to determine Permission ID")
	}

	user := permissionId[0]
	vhost := permissionId[1]

	log.Printf("[DEBUG] RabbitMQ: Attempting to delete topic permission for %s", d.Id())

	resp, err := rmqc.ClearTopicPermissionsIn(vhost, user)
	log.Printf("[DEBUG] RabbitMQ: Topic permission delete response: %#v", resp)
	if err != nil {
		return err
	}

	if resp.StatusCode == 404 {
		// The permissions were already deleted
		return nil
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("Error deleting RabbitMQ topic permission: %s", resp.Status)
	}

	return nil
}

func setTopicPermissionsIn(rmqc *rabbithole.Client, vhost string, user string, permsMap map[string]interface{}) error {
	perms := rabbithole.TopicPermissions{}

	if v, ok := permsMap["exchange"].(string); ok {
		perms.Exchange = v
	}

	if v, ok := permsMap["write"].(string); ok {
		perms.Write = v
	}

	if v, ok := permsMap["read"].(string); ok {
		perms.Read = v
	}

	log.Printf("[DEBUG] RabbitMQ: Attempting to set topic permissions for %s@%s: %#v", user, vhost, perms)

	resp, err := rmqc.UpdateTopicPermissionsIn(vhost, user, perms)
	log.Printf("[DEBUG] RabbitMQ: Permission response: %#v", resp)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("Error setting topic permissions: %s", resp.Status)
	}

	return nil
}
