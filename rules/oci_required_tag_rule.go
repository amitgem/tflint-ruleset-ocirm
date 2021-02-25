package rules

import (
	"encoding/gob"
	"fmt"
	"log"
	"regexp"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/terraform/configs"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// OciRequiredTagRule checks whether ...
type OciRequiredTagRule struct{}

type tag struct {
	Type   string   `hcl:"type,label"`
	Values []string `hcl:"values"`
}

type ociRequiredTag struct {
	Type string `hcl:"type,label"`
	Name string `hcl:"name,label"`
	Tags []tag  `hcl:"tag,block"`
}

type ociRequiredTagsRuleConfig struct {
	OciRequiredTags []ociRequiredTag `hcl:"required_tags,block"`
}

const (
	tagsAttributeName = "defined_tags"
	tagBlockName      = "tag"
)

// NewOciRequiredTagRule returns a new rule
func NewOciRequiredTagRule() *OciRequiredTagRule {
	gob.Register(map[string]string{})
	return &OciRequiredTagRule{}
}

// Name returns the rule name
func (r *OciRequiredTagRule) Name() string {
	return "oci_required_tags"
}

// Enabled returns whether the rule is enabled by default
func (r *OciRequiredTagRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *OciRequiredTagRule) Severity() string {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *OciRequiredTagRule) Link() string {
	return ""
}

// Check checks resources for missing tags
func (r *OciRequiredTagRule) Check(runner tflint.Runner) error {

	config := ociRequiredTagsRuleConfig{}
	if err := runner.DecodeRuleConfig(r.Name(), &config); err != nil {
		return err
	}

	/*
		For each rule
			find resourceType
			find tagList
				for each tag
					find tagName
					find tagValues
					WalkResources(resourceType, )
	*/

	resourceType := "oci_core_instance"

	err := runner.WalkResources(resourceType, func(resource *configs.Resource) error {
		body, _, diags := resource.Config.PartialContent(&hcl.BodySchema{
			Attributes: []hcl.AttributeSchema{
				{
					Name: tagsAttributeName,
				},
			},
		})
		if diags.HasErrors() {
			return diags
		}

		if attribute, ok := body.Attributes[tagsAttributeName]; ok {
			resourceTags := make(map[string]string)
			err := runner.EvaluateExpr(attribute.Expr, &resourceTags, nil)
			err = runner.EnsureNoError(err, func() error {
				return nil
			})
			if err != nil {
				return err
			}
			found := false
			for key := range resourceTags {
				if matched, _ := regexp.MatchString(".*\\.CostCenter$", key); matched {
					found = true
					break
				}
			}
			if found {
				return nil
			}
			r.emitIssue(runner, resource.DeclRange)
		} else {
			log.Printf("[DEBUG] Walk `%s` Resource", resource.Type+"."+resource.Name)
			r.emitIssue(runner, resource.DeclRange)
		}

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func (r *OciRequiredTagRule) emitIssue(runner tflint.Runner, location hcl.Range) {
	issue := fmt.Sprintf("The resource is missing a required tag")
	runner.EmitIssue(r, issue, location)
}
