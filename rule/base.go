package rule

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
)

type Context struct {
	Files map[string]*hclwrite.File
}

func NewContext() Context {
	return Context{
		Files: make(map[string]*hclwrite.File),
	}
}

type Rule interface {
	Name() string
	Process(ctx Context) error
}
