package template

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/zclconf/go-cty/cty"
)

// We can only handle existing outputs defined in input module.

type OutputTF struct {
}

func (OutputTF) FileName() string {
	return "output.tf"
}

func (OutputTF) Generate(mod *tfconfig.Module) (*hclwrite.File, error) {
	f := hclwrite.NewEmptyFile()
	b := f.Body()
	for name, o := range mod.Outputs {
		outputValueLine, err := OutputValueRef(o.Pos)
		if err != nil {
			return nil, err
		}
		outputBlock, err := OutputBlock(name, outputValueLine, o.Sensitive)
		if err != nil {
			return nil, err
		}
		b.AppendBlock(outputBlock)
		b.AppendNewline()
	}

	return f, nil
}

func OutputValueRef(sourcePos tfconfig.SourcePos) (string, error) {
	f, err := os.OpenFile(sourcePos.Filename, os.O_RDONLY, 0666)
	if err != nil {
		return "", err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	line := 0
	start := false
	for scanner.Scan() {
		t := scanner.Text()
		line++
		if line == sourcePos.Line {
			start = true
		}
		if start && strings.Contains(t, "value") {
			re := regexp.MustCompile(`\s+value\s+=\s+(.+)\n?`)
			res := re.FindAllStringSubmatch(t, -1)
			return res[0][1], nil
		}
	}

	return "", fmt.Errorf("cannot find output value")
}

func OutputBlock(name string, valueLine string, sensitive bool) (*hclwrite.Block, error) {
	b := hclwrite.NewBlock("output", []string{name})
	b.Body().SetAttributeRaw("value", hclwrite.TokensForIdentifier(valueLine))
	b.Body().SetAttributeValue("sensitive", cty.BoolVal(sensitive))
	return b, nil
}
