package rabbitmq

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	rabbithole "github.com/michaelklishin/rabbit-hole"
)

func resourceFederationUpstream() *schema.Resource {
	return &schema.Resource{
		Create: CreateFederationUpstream,
		Read:   ReadFederationUpstream,
		Update: UpdateFederationUpstream,
		Delete: DeleteFederationUpstream,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"vhost": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			// "federation-upstream"
			"component": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"definition": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// applicable to both federated exchanges and queues
						"uri": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},

						"prefetch_count": {
							Type:     schema.TypeInt,
							Optional: true,
						},

						"reconnect_delay": {
							Type:     schema.TypeInt,
							Optional: true,
						},

						"ack_mode": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"trust_user_id": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						// applicable to federated exchanges only
						"exchange": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"max_hops": {
							Type:     schema.TypeInt,
							Optional: true,
						},

						"expires": {
							Type:     schema.TypeInt,
							Optional: true,
						},

						"message_ttl": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						// applicable to federated queues only
						"queue": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func CreateFederationUpstream(d *schema.ResourceData, meta interface{}) error {
	rmqc := meta.(*rabbithole.Client)

	name := d.Get("name").(string)
	vhost := d.Get("vhost").(string)
	defList := d.Get("definition").([]interface{})

	defMap, ok := defList[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("Unable to parse federation upstream definition")
	}

	if err := putFederationUpstream(rmqc, vhost, name, defMap); err != nil {
		return err
	}

	id := fmt.Sprintf("%s@%s", name, vhost)
	d.SetId(id)

	return ReadFederationUpstream(d, meta)
}

func ReadFederationUpstream(d *schema.ResourceData, meta interface{}) error {
	rmqc := meta.(*rabbithole.Client)

	id := strings.Split(d.Id(), "@")
	if len(id) < 2 {
		return fmt.Errorf("Unable to determine federation upstream id: %s", d.Id())
	}

	name := id[0]
	vhost := id[1]

	upstream, err := rmqc.GetFederationUpstream(vhost, name)
	if err != nil {
		return checkDeleted(d, err)
	}

	log.Printf("[DEBUG] RabbitMQ: Federation upstream retrieved for %s: %#v", d.Id(), upstream)

	d.Set("name", upstream.Name)
	d.Set("vhost", upstream.Vhost)
	d.Set("component", upstream.Component)

	defMap := map[string]interface{}{
		"uri":             upstream.Definition.Uri,
		"prefetch_count":  upstream.Definition.PrefetchCount,
		"reconnect_delay": upstream.Definition.ReconnectDelay,
		"ack_mode":        upstream.Definition.AckMode,
		"trust_user_id":   upstream.Definition.TrustUserId,
		"exchange":        upstream.Definition.Exchange,
		"max_hops":        upstream.Definition.MaxHops,
		"expires":         upstream.Definition.Expires,
		"message_ttl":     upstream.Definition.MessageTTL,
		"queue":           upstream.Definition.Queue,
	}

	defList := [1]map[string]interface{}{defMap}
	d.Set("definition", defList)

	return nil
}

func UpdateFederationUpstream(d *schema.ResourceData, meta interface{}) error {
	rmqc := meta.(*rabbithole.Client)

	id := strings.Split(d.Id(), "@")
	if len(id) < 2 {
		return fmt.Errorf("Unable to determine federation upstream id: %s", d.Id())
	}

	name := id[0]
	vhost := id[1]

	if d.HasChange("definition") {
		_, newDef := d.GetChange("definition")

		defList := newDef.([]interface{})
		defMap, ok := defList[0].(map[string]interface{})
		if !ok {
			return fmt.Errorf("Unable to parse federation definition")
		}

		if err := putFederationUpstream(rmqc, vhost, name, defMap); err != nil {
			return err
		}
	}

	return ReadFederationUpstream(d, meta)
}

func DeleteFederationUpstream(d *schema.ResourceData, meta interface{}) error {
	rmqc := meta.(*rabbithole.Client)

	id := strings.Split(d.Id(), "@")
	if len(id) < 2 {
		return fmt.Errorf("Unable to determine federation upstream id: %s", d.Id())
	}

	name := id[0]
	vhost := id[1]

	log.Printf("[DEBUG] RabbitMQ: Attempting to delete federation upstream for %s", d.Id())

	resp, err := rmqc.DeleteFederationUpstream(vhost, name)
	log.Printf("[DEBUG] RabbitMQ: Federation upstream delete response: %#v", resp)
	if err != nil {
		return err
	}

	if resp.StatusCode == 404 {
		// the upstream was automatically deleted
		return nil
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("Error deleting RabbitMQ federation upstream: %s", resp.Status)
	}

	return nil
}

func putFederationUpstream(rmqc *rabbithole.Client, vhost string, name string, defMap map[string]interface{}) error {
	definition := rabbithole.FederationDefinition{}

	log.Printf("[DEBUG] RabbitMQ: Attempting to put federation definition for %s@%s: %#v", name, vhost, defMap)

	if v, ok := defMap["uri"].(string); ok {
		definition.Uri = v
	}

	if v, ok := defMap["expires"].(int); ok {
		definition.Expires = v
	}

	if v, ok := defMap["message_ttl"].(int); ok {
		definition.MessageTTL = int32(v)
	}

	if v, ok := defMap["max_hops"].(int); ok {
		definition.MaxHops = v
	}

	if v, ok := defMap["prefetch_count"].(int); ok {
		definition.PrefetchCount = v
	}

	if v, ok := defMap["reconnect_delay"].(int); ok {
		definition.ReconnectDelay = v
	}

	if v, ok := defMap["ack_mode"].(string); ok {
		definition.AckMode = v
	}

	if v, ok := defMap["trust_user_id"].(bool); ok {
		definition.TrustUserId = v
	}

	if v, ok := defMap["exchange"].(string); ok {
		definition.Exchange = v
	}

	if v, ok := defMap["queue"].(string); ok {
		definition.Queue = v
	}

	log.Printf("[DEBUG] RabbitMQ: Attempting to declare federation upstream for %s@%s: %#v", name, vhost, definition)

	resp, err := rmqc.PutFederationUpstream(vhost, name, definition)
	log.Printf("[DEBUG] RabbitMQ: Federation upstream declare response: %#v", resp)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("Error creating RabbitMQ federation upstream: %s", resp.Status)
	}

	return nil
}
