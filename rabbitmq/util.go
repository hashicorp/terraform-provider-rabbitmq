package rabbitmq

import (
	"fmt"
	"strings"

	"github.com/appilon/terraform-plugin-sdk/structure"
	"github.com/hashicorp/terraform/helper/schema"
)

func checkDeleted(d *schema.ResourceData, err error) error {
	if err.Error() == "not found" {
		d.SetId("")
		return nil
	}

	return err
}

// pulled from terraform-provider-aws/aws/validators.go
func validateJsonString(v interface{}, k string) (ws []string, errors []error) {
	if _, err := structure.NormalizeJsonString(v); err != nil {
		errors = append(errors, fmt.Errorf("%q contains invalid JSON: %s", k, err))
	}
	return
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
