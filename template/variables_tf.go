package template

import (
	"encoding/json"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

// VariablesTF needs to handle already defined variables, and
type VariablesTF struct {
	Prefix string
}

func (VariablesTF) FileName() string {
	return "variables.tf"
}

func (tf VariablesTF) Generate(mod *tfconfig.Module) (*hclwrite.File, error) {
	f := hclwrite.NewEmptyFile()
	b := f.Body()
	// handle the variables defined in input module
	for name, variable := range mod.Variables {
		varBlock, err := VariableBlock(name, variable.Type, variable.Default, variable.Description)
		if err != nil {
			return nil, err
		}
		b.AppendBlock(varBlock)
		b.AppendNewline()
	}

	hasPrefix := false
	hasTags := false
	for _, block := range b.Blocks() {
		if len(block.Labels()) == 1 {
			if block.Labels()[0] == "prefix" {
				hasPrefix = true
			}
			if block.Labels()[0] == "tags" {
				hasTags = true
			}
		}
	}

	if !hasPrefix {
		prefixBlock, err := VariableBlock("prefix", cty.String.FriendlyNameForConstraint(), tf.Prefix, "Prefix of the resource name")
		if err != nil {
			return nil, err
		}
		b.AppendBlock(prefixBlock)
		b.AppendNewline()
	}

	if !hasTags {
		tagsBlock, err := VariableBlock("tags", "map(any)", nil, "Azure Tags for all resources.")
		if err != nil {
			return nil, err
		}
		b.AppendBlock(tagsBlock)
		b.AppendNewline()
	}

	return f, nil
}

func VariableBlock(name string, vType string, defaultValue interface{}, description string) (*hclwrite.Block, error) {
	b := hclwrite.NewBlock("variable", []string{name})
	b.Body().SetAttributeRaw("type", hclwrite.TokensForIdentifier(vType))
	if defaultValue != nil {
		switch vType {
		case cty.String.FriendlyNameForConstraint():
			b.Body().SetAttributeValue("default", cty.StringVal(defaultValue.(string)))
		case cty.Number.FriendlyNameForConstraint():
			b.Body().SetAttributeValue("default", cty.NumberFloatVal(defaultValue.(float64)))
		default:
			defJson, err := json.Marshal(defaultValue)
			if err != nil {
				return nil, err
			}
			defCty, err := ctyjson.Unmarshal(defJson, cty.Map(cty.String))
			b.Body().SetAttributeValue("default", defCty)
		}
	}
	if description != "" {
		b.Body().SetAttributeValue("description", cty.StringVal(description))
	}
	return b, nil
}
