package rabbitmq

import (
	"fmt"
	"log"
	"strings"

	rabbithole "github.com/michaelklishin/rabbit-hole/v2"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceShovel() *schema.Resource {
	return &schema.Resource{
		Create: CreateShovel,
		Read:   ReadShovel,
		Delete: DeleteShovel,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
			"info": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_uri": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
						"source_exchange": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  nil,
						},
						"source_exchange_key": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  nil,
						},
						"source_queue": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  nil,
						},
						"destination_uri": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
						"destination_exchange": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  nil,
						},
						"destination_exchange_key": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  nil,
						},
						"destination_queue": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  nil,
						},
						"prefetch_count": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1000,
						},
						"reconnect_delay": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
						},
						"add_forward_headers": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"ack_mode": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "on-confirm",
						},
						"delete_after": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "never",
						},
					},
				},
			},
		},
	}
}

func CreateShovel(d *schema.ResourceData, meta interface{}) error {
	rmqc := meta.(*rabbithole.Client)

	vhost := d.Get("vhost").(string)
	shovelName := d.Get("name").(string)
	shovelInfo := d.Get("info").([]interface{})

	shovelMap, ok := shovelInfo[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("Unable to parse shovel info")
	}

	shovelDefinition := setShovelDefinition(shovelMap).(rabbithole.ShovelDefinition)

	log.Printf("[DEBUG] RabbitMQ: Attempting to declare shovel %s in vhost %s", shovelName, vhost)
	resp, err := rmqc.DeclareShovel(vhost, shovelName, shovelDefinition)
	log.Printf("[DEBUG] RabbitMQ: shovel declartion response: %#v", resp)
	if err != nil {
		return err
	}

	shovelId := fmt.Sprintf("%s@%s", shovelName, vhost)

	d.SetId(shovelId)

	return ReadShovel(d, meta)
}

func ReadShovel(d *schema.ResourceData, meta interface{}) error {
	rmqc := meta.(*rabbithole.Client)

	shovelId := strings.Split(d.Id(), "@")

	name := shovelId[0]
	vhost := shovelId[1]

	shovelInfo, err := rmqc.GetShovel(vhost, name)
	if err != nil {
		return checkDeleted(d, err)
	}

	log.Printf("[DEBUG] RabbitMQ: Shovel retrieved: Vhost: %#v, Name: %#v", vhost, name)

	d.Set("name", shovelInfo.Name)
	d.Set("vhost", shovelInfo.Vhost)

	return nil
}

func DeleteShovel(d *schema.ResourceData, meta interface{}) error {
	rmqc := meta.(*rabbithole.Client)

	shovelId := strings.Split(d.Id(), "@")

	name := shovelId[0]
	vhost := shovelId[1]

	log.Printf("[DEBUG] RabbitMQ: Attempting to delete shovel %s", d.Id())

	resp, err := rmqc.DeleteShovel(vhost, name)
	log.Printf("[DEBUG] RabbitMQ: shovel deletion response: %#v", resp)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("Error deleting RabbitMQ shovel: %s", resp.Status)
	}

	return nil
}

func setShovelDefinition(shovelMap map[string]interface{}) interface{} {
	shovelDefinition := &rabbithole.ShovelDefinition{}

	if v, ok := shovelMap["source_uri"].(string); ok {
		shovelDefinition.SourceURI = v
	}

	if v, ok := shovelMap["source_exchange"].(string); ok {
		shovelDefinition.SourceExchange = v
	}

	if v, ok := shovelMap["source_exchange_key"].(string); ok {
		shovelDefinition.SourceExchangeKey = v
	}

	if v, ok := shovelMap["source_queue"].(string); ok {
		shovelDefinition.SourceQueue = v
	}

	if v, ok := shovelMap["destination_uri"].(string); ok {
		shovelDefinition.DestinationURI = v
	}

	if v, ok := shovelMap["destination_exchange"].(string); ok {
		shovelDefinition.DestinationExchange = v
	}

	if v, ok := shovelMap["destination_exchange_key"].(string); ok {
		shovelDefinition.DestinationExchangeKey = v
	}

	if v, ok := shovelMap["destination_queue"].(string); ok {
		shovelDefinition.DestinationQueue = v
	}

	if v, ok := shovelMap["prefetch_count"].(int); ok {
		shovelDefinition.PrefetchCount = v
	}

	if v, ok := shovelMap["reconnect_delay"].(int); ok {
		shovelDefinition.ReconnectDelay = v
	}

	if v, ok := shovelMap["add_forward_headers"].(bool); ok {
		shovelDefinition.AddForwardHeaders = v
	}

	if v, ok := shovelMap["ack_mode"].(string); ok {
		shovelDefinition.AckMode = v
	}

	if v, ok := shovelMap["delete_after"].(string); ok {
		shovelDefinition.DeleteAfter = v
	}

	return *shovelDefinition
}
