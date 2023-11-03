package template

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/zclconf/go-cty/cty"
)

type ProviderTF struct {
	TFRequiredVersion string
}

func (ProviderTF) FileName() string {
	return "providers.tf"
}

func Providers(providers map[string]*tfconfig.ProviderRequirement) (*hclwrite.Block, error) {
	b := hclwrite.NewBlock("required_providers", []string{})
	for n, p := range providers {
		pArguments := map[string]cty.Value{}
		defaultValues, ok := defaultProvidersConfig()[n]

		src := p.Source
		if src == "" && ok {
			if defaultSource, ok := defaultValues["source"]; ok {
				src = defaultSource
			}
		}
		if src == "" {
			return nil, fmt.Errorf("%s has not defined `source`, both in input config or default values", n)
		}
		pArguments["source"] = cty.StringVal(src)

		ver := ""
		if len(p.VersionConstraints) != 0 {
			ver = strings.Join(p.VersionConstraints, ",")
		}
		if ver == "" && ok {
			if defaultVersion, ok := defaultValues["version"]; ok {
				ver = defaultVersion
			}
		}
		if ver == "" {
			return nil, fmt.Errorf("%s has not defined `version`, both in input config or default values", n)
		}
		pArguments["version"] = cty.StringVal(ver)

		b.Body().SetAttributeValue(n, cty.MapVal(pArguments))
	}
	return b, nil
}

func (tf ProviderTF) Generate(mod *tfconfig.Module) (*hclwrite.File, error) {
	f := hclwrite.NewEmptyFile()
	tfBlock := hclwrite.NewBlock("terraform", []string{})
	tfBlock.Body().SetAttributeValue("required_version", cty.StringVal(tf.TFRequiredVersion))
	providers, err := Providers(mod.RequiredProviders)
	if err != nil {
		return nil, err
	}
	tfBlock.Body().AppendBlock(providers)
	f.Body().AppendBlock(tfBlock)

	return f, nil
}
