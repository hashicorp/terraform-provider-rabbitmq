package rabbitmq

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
)

func checkDeleted(d *schema.ResourceData, err error) error {
	switch e := err.(type) {
	case rabbithole.ErrorResponse:
		if e.StatusCode == 404 {
			d.SetId("")
			return nil
		}
	}

	return err
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
