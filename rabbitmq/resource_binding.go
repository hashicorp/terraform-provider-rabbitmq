package rabbitmq

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/michaelklishin/rabbit-hole"
)

func resourceBinding() *schema.Resource {
	return &schema.Resource{
		Create: CreateBinding,
		Read:   ReadBinding,
		Delete: DeleteBinding,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"source": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"vhost": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"destination": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"destination_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"properties_key": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"routing_key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"arguments": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

// Because slashes are used to separate different components when constructing binding IDs,
// we need a way to ensure any components that include slashes can survive the round trip.
// Percent-encoding is a straightforward way of doing so.
// (reference: https://developer.mozilla.org/en-US/docs/Glossary/percent-encoding)

func percentEncodeSlashes(s string) string {
	// Encode any percent signs, then encode any forward slashes.
	return strings.Replace(strings.Replace(s, "%", "%25", -1), "/", "%2F", -1)
}

func percentDecodeSlashes(s string) string {
	// Decode any forward slashes, then decode any percent signs.
	return strings.Replace(strings.Replace(s, "%2F", "/", -1), "%25", "%", -1)
}

func CreateBinding(d *schema.ResourceData, meta interface{}) error {
	rmqc := meta.(*rabbithole.Client)

	vhost := d.Get("vhost").(string)
	bindingInfo := rabbithole.BindingInfo{
		Source:          d.Get("source").(string),
		Destination:     d.Get("destination").(string),
		DestinationType: d.Get("destination_type").(string),
		RoutingKey:      d.Get("routing_key").(string),
		Arguments:       d.Get("arguments").(map[string]interface{}),
	}

	propertiesKey, err := declareBinding(rmqc, vhost, bindingInfo)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] RabbitMQ: Binding properties key: %s", propertiesKey)
	bindingInfo.PropertiesKey = propertiesKey
	name := fmt.Sprintf("%s/%s/%s/%s/%s", percentEncodeSlashes(vhost), bindingInfo.Source, bindingInfo.Destination, bindingInfo.DestinationType, bindingInfo.PropertiesKey)
	d.SetId(name)

	return ReadBinding(d, meta)
}

func ReadBinding(d *schema.ResourceData, meta interface{}) error {
	rmqc := meta.(*rabbithole.Client)

	log.Printf("[TRACE] RabbitMQ: read binding resource ID (pre-split): %s", d.Id())
	bindingId := strings.Split(d.Id(), "/")
	log.Printf("[DEBUG] RabbitMQ: binding ID: %#v", bindingId)
	if len(bindingId) < 5 {
		return fmt.Errorf("Unable to determine binding ID")
	}

	vhost := percentDecodeSlashes(bindingId[0])
	source := bindingId[1]
	destination := bindingId[2]
	destinationType := bindingId[3]
	propertiesKey := bindingId[4]
	log.Printf("[DEBUG] RabbitMQ: Attempting to find binding for: vhost=%s source=%s destination=%s destinationType=%s propertiesKey=%s",
		vhost, source, destination, destinationType, propertiesKey)

	bindings, err := rmqc.ListBindingsIn(vhost)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] RabbitMQ: Bindings retrieved: %#v", bindings)
	bindingFound := false
	for _, binding := range bindings {
		log.Printf("[TRACE] RabbitMQ: Assessing binding: %#v", binding)
		if binding.Source == source && binding.Destination == destination && binding.DestinationType == destinationType && binding.PropertiesKey == propertiesKey {
			log.Printf("[DEBUG] RabbitMQ: Found Binding: %#v", binding)
			bindingFound = true

			d.Set("vhost", binding.Vhost)
			d.Set("source", binding.Source)
			d.Set("destination", binding.Destination)
			d.Set("destination_type", binding.DestinationType)
			d.Set("routing_key", binding.RoutingKey)
			d.Set("properties_key", binding.PropertiesKey)
			d.Set("arguments", binding.Arguments)
		}
	}

	// The binding could not be found,
	// so consider it deleted and remove from state
	if !bindingFound {
		d.SetId("")
	}

	return nil
}

func DeleteBinding(d *schema.ResourceData, meta interface{}) error {
	rmqc := meta.(*rabbithole.Client)

	bindingId := strings.Split(d.Id(), "/")
	if len(bindingId) < 5 {
		return fmt.Errorf("Unable to determine binding ID")
	}

	vhost := percentDecodeSlashes(bindingId[0])
	source := bindingId[1]
	destination := bindingId[2]
	destinationType := bindingId[3]
	propertiesKey := bindingId[4]

	bindingInfo := rabbithole.BindingInfo{
		Vhost:           vhost,
		Source:          source,
		Destination:     destination,
		DestinationType: destinationType,
		PropertiesKey:   propertiesKey,
	}

	log.Printf("[DEBUG] RabbitMQ: Attempting to delete binding for: vhost=%s source=%s destination=%s destinationType=%s propertiesKey=%s",
		vhost, source, destination, destinationType, propertiesKey)

	resp, err := rmqc.DeleteBinding(vhost, bindingInfo)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] RabbitMQ: Binding delete response: %#v", resp)

	if resp.StatusCode == 404 {
		// The binding was already deleted
		return nil
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("Error deleting RabbitMQ binding: %s", resp.Status)
	}

	return nil
}

func declareBinding(rmqc *rabbithole.Client, vhost string, bindingInfo rabbithole.BindingInfo) (string, error) {
	log.Printf("[DEBUG] RabbitMQ: Attempting to declare binding for: vhost=%s source=%s destination=%s destinationType=%s",
		vhost, bindingInfo.Source, bindingInfo.Destination, bindingInfo.DestinationType)

	resp, err := rmqc.DeclareBinding(vhost, bindingInfo)
	log.Printf("[DEBUG] RabbitMQ: Binding declare response: %#v", resp)
	if err != nil {
		return "", err
	}

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("Error declaring RabbitMQ binding: %s", resp.Status)
	}

	location := strings.Split(resp.Header.Get("Location"), "/")
	propertiesKey, err := url.PathUnescape(location[len(location)-1])

	if err != nil {
		return "", err
	}

	return propertiesKey, nil
}
