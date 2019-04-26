package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	rabbithole "github.com/michaelklishin/rabbit-hole"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceQueue() *schema.Resource {
	return &schema.Resource{
		Create: CreateQueue,
		Read:   ReadQueue,
		Delete: DeleteQueue,
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
				Optional: true,
				Default:  "/",
				ForceNew: true,
			},

			"settings": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"durable": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},

						"auto_delete": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},

						"arguments": {
							Type:          schema.TypeMap,
							Optional:      true,
							ConflictsWith: []string{"settings.0.arguments_json"},
						},

						"arguments_json": {
							Type:          schema.TypeString,
							Optional:      true,
							ValidateFunc:  validateJsonString,
							ConflictsWith: []string{"settings.0.arguments"},
						},
					},
				},
			},
		},
	}
}

func CreateQueue(d *schema.ResourceData, meta interface{}) error {
	rmqc := meta.(*rabbithole.Client)

	name := d.Get("name").(string)
	vhost := d.Get("vhost").(string)
	settingsList := d.Get("settings").([]interface{})

	settingsMap, ok := settingsList[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("Unable to parse settings")
	}

	// If arguments_json is used, unmarshal it into a generic interface
	// and use it as the "arguments" key for the queue.
	if v, ok := settingsMap["arguments_json"].(string); ok && v != "" {
		var arguments interface{}
		err := json.Unmarshal([]byte(v), &arguments)
		if err != nil {
			return err
		}

		delete(settingsMap, "arguments_json")
		settingsMap["arguments"] = arguments
	}

	if err := declareQueue(rmqc, vhost, name, settingsMap); err != nil {
		return err
	}

	id := fmt.Sprintf("%s@%s", name, vhost)
	d.SetId(id)

	return ReadQueue(d, meta)
}

func ReadQueue(d *schema.ResourceData, meta interface{}) error {
	rmqc := meta.(*rabbithole.Client)

	queueId := strings.Split(d.Id(), "@")
	if len(queueId) < 2 {
		return fmt.Errorf("Unable to determine Queue ID")
	}

	user := queueId[0]
	vhost := queueId[1]

	queueSettings, err := rmqc.GetQueue(vhost, user)
	if err != nil {
		return checkDeleted(d, err)
	}

	log.Printf("[DEBUG] RabbitMQ: Queue retrieved for %s: %#v", d.Id(), queueSettings)

	d.Set("name", queueSettings.Name)
	d.Set("vhost", queueSettings.Vhost)

	queue := make([]map[string]interface{}, 1)
	e := make(map[string]interface{})
	e["durable"] = queueSettings.Durable
	e["auto_delete"] = queueSettings.AutoDelete
	e["arguments"] = queueSettings.Arguments
	queue[0] = e

	return d.Set("settings", queue)
}

func DeleteQueue(d *schema.ResourceData, meta interface{}) error {
	rmqc := meta.(*rabbithole.Client)

	queueId := strings.Split(d.Id(), "@")
	if len(queueId) < 2 {
		return fmt.Errorf("Unable to determine Queue ID")
	}

	user := queueId[0]
	vhost := queueId[1]

	log.Printf("[DEBUG] RabbitMQ: Attempting to delete queue for %s", d.Id())

	resp, err := rmqc.DeleteQueue(vhost, user)
	log.Printf("[DEBUG] RabbitMQ: Queue delete response: %#v", resp)
	if err != nil {
		return err
	}

	if resp.StatusCode == 404 {
		// the queue was automatically deleted
		return nil
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("Error deleting RabbitMQ queue: %s", resp.Status)
	}

	return nil
}

func declareQueue(rmqc *rabbithole.Client, vhost string, name string, settingsMap map[string]interface{}) error {
	queueSettings := rabbithole.QueueSettings{}

	if v, ok := settingsMap["durable"].(bool); ok {
		queueSettings.Durable = v
	}

	if v, ok := settingsMap["auto_delete"].(bool); ok {
		queueSettings.AutoDelete = v
	}

	if v, ok := settingsMap["arguments"].(map[string]interface{}); ok {
		queueSettings.Arguments = v
	}

	log.Printf("[DEBUG] RabbitMQ: Attempting to declare queue for %s@%s: %#v", name, vhost, queueSettings)

	resp, err := rmqc.DeclareQueue(vhost, name, queueSettings)
	log.Printf("[DEBUG] RabbitMQ: Queue declare response: %#v", resp)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("Error declaring RabbitMQ queue: %s", resp.Status)
	}

	return nil
}
