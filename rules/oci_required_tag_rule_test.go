package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_OciRequiredTagRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Config   string
		Expected helper.Issues
	}{
		{
			Name: "Wanted tags: Bar,Foo, found: bar,foo",
			Content: `
resource "oci_core_instance" "t1" {
  defined_tags = {
    "Applications.CostCenter" = "1234"
  }
}`,
			Config: `
rule "oci_required_tags" {
	enabled = true
	required_tags "resource" "oci_core_instance" {
		tag "Applications.CostCenter" {
			values = ["AA", "BB" ]
		}
		tag "HumanResource.Provider" {
			values = ["XX", "YY"]
		}
	}
	required_tags "resource" "oci_core_compartment" {
		tag "Applications.CostCenter"{
			values = ["CC", "DD" ]
		}
	}
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewOciRequiredTagRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content, ".tflint.hcl": tc.Config})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
