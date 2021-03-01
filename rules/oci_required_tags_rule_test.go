package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_OciRequiredTagsRule(t *testing.T) {
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
    "HumanResource.Provider" = "1234"
    "Applications.CostCenter" = "1234"
  }
}`,
			Config: `
rule "oci_required_tags" {
	enabled = true
	required_tags "oci_core_instance" {
		tag "Applications.CostCenter" {
			values = ["AA", "BB" ]
		}
		tag "HumanResource.Provider" {
			values = ["XX", "YY"]
	    }
    }
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewOciRequiredTagsRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content, ".tflint.hcl": tc.Config})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
