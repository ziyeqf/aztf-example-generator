package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/ziyeqf/aztf-example-generator/rule"
	"github.com/ziyeqf/aztf-example-generator/template"
)

func main() {
	inputModule := flag.String("path", "", "the example tf module to process")
	tfRequiredVersion := flag.String("tf-version", "", "the terraform version to use")
	defaultPrefix := flag.String("prefix", "", "the prefix of resource name")
	outputPath := flag.String("output", "", "the output path of the generated tf files")

	flag.Parse()

	tfVersion := ">=1.2"
	if tfRequiredVersion != nil {
		tfVersion = *tfRequiredVersion
	}

	inputAbsPath, err := filepath.Abs(*inputModule)
	if err != nil {
		log.Println(err)
		return
	}

	module, diag := tfconfig.LoadModule(inputAbsPath)
	if diag.HasErrors() {
		log.Println(err.Error())
		return
	}

	files := []template.TfFile{
		template.ProviderTF{
			TFRequiredVersion: tfVersion,
		},
		template.VariablesTF{
			Prefix: *defaultPrefix,
		},
		template.MainTF{},
		template.OutputTF{},
	}
	rules := []rule.Rule{
		rule.LocationRule{},
	}
	ctx := rule.NewContext()

	for _, file := range files {
		output, err := file.Generate(module)
		if err != nil {
			log.Println(err.Error())
			return
		}
		ctx.Files[file.FileName()] = output
	}

	for _, rule := range rules {
		err := rule.Process(ctx)
		if err != nil {
			log.Println(fmt.Errorf("rule %s process err: %+v", rule.Name(), err))
			return
		}
	}

	for fileName, file := range ctx.Files {
		outputAbsPath, err := filepath.Abs(*outputPath + "/" + fileName)
		if err != nil {
			log.Println(err)
			return
		}

		f, err := os.OpenFile(outputAbsPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0660)
		if err != nil {
			log.Println(err.Error())
			return
		}
		defer f.Close()

		_, err = f.Write(file.Bytes())
		if err != nil {
			log.Println(err.Error())
			return
		}
	}
}
