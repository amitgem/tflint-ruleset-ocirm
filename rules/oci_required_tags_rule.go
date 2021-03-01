package rules

import (
	"fmt"
	"log"
	"regexp"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/terraform/configs"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// OciRequiredTagsRule checks whether ...
type OciRequiredTagsRule struct{}

type tag struct {
	Type   string   `hcl:"type,label"`
	Values []string `hcl:"values"`
}

type ociRequiredTag struct {
	ResourceType string `hcl:"type,label"`
	Tags         []tag  `hcl:"tag,block"`
}

type ociRequiredTagsRuleConfig struct {
	OciRequiredTags []ociRequiredTag `hcl:"required_tags,block"`
}

var tagAttributeNames = []string{
	"defined_tags",
	"tag",
}

// NewOciRequiredTagsRule returns a new rule
func NewOciRequiredTagsRule() *OciRequiredTagsRule {
	//gob.Register(map[string]string{})
	return &OciRequiredTagsRule{}
}

// Name returns the rule name
func (r *OciRequiredTagsRule) Name() string {
	log.Printf("[INFO] Plugin name retrived.")
	return "oci_required_tagss"
}

// Enabled returns whether the rule is enabled by default
func (r *OciRequiredTagsRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *OciRequiredTagsRule) Severity() string {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *OciRequiredTagsRule) Link() string {
	return ""
}

// Check checks resources for missing tags
func (r *OciRequiredTagsRule) Check(runner tflint.Runner) error {

	config := ociRequiredTagsRuleConfig{}
	if err := runner.DecodeRuleConfig(r.Name(), &config); err != nil {
		return err
	}

	for _, requiredTags := range config.OciRequiredTags {
		for _, tag := range requiredTags.Tags {

			err := runner.WalkResources(requiredTags.ResourceType, func(resource *configs.Resource) error {

				for _, tagsAttributeName := range tagAttributeNames {
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
							//r.emitIssueOnExpr(runner, resourceTags, config, attribute.Expr)
							return nil
						})
						if err != nil {
							return err
						}

						for key := range resourceTags {
							if matched, _ := regexp.MatchString(tag.Type, key); matched {
								return nil
							}
						}
					}
				}
				r.emitIssue(runner, resource.DeclRange, resource.Type+"/"+resource.Name, tag.Type)
				return nil

			})
			if err != nil {
				return err
			}

		}
	}

	return nil
}

func (r *OciRequiredTagsRule) emitIssue(runner tflint.Runner, location hcl.Range, resourceName string, attribute string) {
	issue := fmt.Sprintf("The resource %s is missing required tag %s", resourceName, attribute)
	runner.EmitIssue(r, issue, location)
}

func (r *OciRequiredTagsRule) emitResourceIssue(runner tflint.Runner, location hcl.Range) {
	issue := fmt.Sprintf("The resource is missing required tags")
	runner.EmitIssue(r, issue, location)
}
