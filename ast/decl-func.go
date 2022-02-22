package ast

import (
	"fmt"
	"strings"
)

var _ Decl = DeclFunc{}
var _ Overviewable = DeclFunc{}

type DeclFunc struct {
	Name Identifier
	Impl *ExprFunc

	Docs     *Docs
	MetaInfo *MetaDecl
}

func (e DeclFunc) DeclName() Identifier {
	return e.Name
}

func (e DeclFunc) DeclOverview() string {
	if len(e.Impl.Parameters) == 0 {
		return fmt.Sprintf("func %s { => }", e.Name)
	}
	paramNames := make([]string, len(e.Impl.Parameters))
	for i, param := range e.Impl.Parameters {
		paramNames[i] = string(param.Name)
	}
	return fmt.Sprintf("func %s { %s => }", e.Name, strings.Join(paramNames, ", "))
}

func (e DeclFunc) Meta() *MetaDecl {
	return e.MetaInfo
}

func (e DeclFunc) IsExportedDecl() bool {
	return true
}

func MakeDeclFunc(name Identifier, impl *ExprFunc, source *Source) *DeclFunc {
	return &DeclFunc{
		Name: name,
		Impl: impl,
		MetaInfo: &MetaDecl{
			Source: source,
		},
	}
}

func (decl DeclFunc) ProvidedDocs() *Docs {
	return decl.Docs
}
