package template

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
)

type TfFile interface {
	FileName() string
	Generate(mod *tfconfig.Module) (*hclwrite.File, error)
}
