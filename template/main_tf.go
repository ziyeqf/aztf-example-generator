package template

import (
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
)

type MainTF struct {
}

func (MainTF) FileName() string {
	return "main.tf"
}

func (MainTF) Generate(mod *tfconfig.Module) (*hclwrite.File, error) {
	f := hclwrite.NewEmptyFile()
	b := f.Body()

	for _, res := range mod.ManagedResources {
		block, err := ResourceBlock(res)
		if err != nil {
			return nil, err
		}
		b.AppendBlock(block)
		b.AppendNewline()
	}

	for _, ds := range mod.DataResources {
		block, err := DataSourceBlock(ds)
		if err != nil {
			return nil, err
		}
		b.AppendBlock(block)
		b.AppendNewline()
	}

	return f, nil
}

func ResourceBlock(res *tfconfig.Resource) (*hclwrite.Block, error) {
	originalResBlock, err := OriginalResRef(res)
	if err != nil {
		return nil, err
	}

	return originalResBlock, nil
}

func DataSourceBlock(ds *tfconfig.Resource) (*hclwrite.Block, error) {
	originalResBlock, err := OriginalResRef(ds)
	if err != nil {
		return nil, err
	}

	return originalResBlock, nil
}

func OriginalResRef(tfRes *tfconfig.Resource) (*hclwrite.Block, error) {
	f, err := os.OpenFile(tfRes.Pos.Filename, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %+v", err)
	}

	file, diags := hclwrite.ParseConfig(b, tfRes.Pos.Filename, hcl.InitialPos)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse hcl: %+v", diags.Error())
	}

	for _, block := range file.Body().Blocks() {
		if block.Labels()[0] == tfRes.Type && block.Labels()[1] == tfRes.Name {
			return block, nil
		}
	}

	return nil, fmt.Errorf("cannot find  %s definition %s", tfRes.Type, tfRes.Name)
}
