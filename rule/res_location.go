package rule

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

type LocationRule struct{}

var _ Rule = LocationRule{}

func (LocationRule) Name() string {
	return "resource_location"
}

func (LocationRule) Process(ctx Context) error {
	// turn locations to resource group's location
	for name, file := range ctx.Files {
		if name == "main.tf" {
			rgName, err := resourceGroupName(file.Body())
			if err != nil {
				return err
			}
			for _, block := range file.Body().Blocks() {
				if block.Type() == "resource" {
					hasLocation := false
					for name, _ := range block.Body().Attributes() {
						if name == "location" {
							hasLocation = true
						}
					}
					if hasLocation {
						block.Body().SetAttributeTraversal("location", hcl.Traversal{
							hcl.TraverseRoot{
								Name: "azurerm_resource_group",
							},
							hcl.TraverseAttr{
								Name: rgName,
							},
							hcl.TraverseAttr{
								Name: "location",
							},
						})
					}

				}
			}
		}
	}
	return nil
}

func resourceGroupName(body *hclwrite.Body) (string, error) {
	for _, block := range body.Blocks() {
		if block.Type() == "resource" &&
			block.Labels()[0] == "azurerm_resource_group" {
			return block.Labels()[1], nil
		}
	}
	return "", fmt.Errorf("no resource group found")
}
