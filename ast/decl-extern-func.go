package ast

import (
	"fmt"
	"strings"
)

var _ Decl = DeclExternFunc{}
var _ Overviewable = DeclExternFunc{}

type DeclExternFunc struct {
	Name       Identifier
	Parameters []DeclParameter

	Docs     *Docs
	MetaInfo *MetaDecl
}

func (e DeclExternFunc) DeclName() Identifier {
	return e.Name
}

func (e DeclExternFunc) Meta() *MetaDecl {
	return e.MetaInfo
}

func (e DeclExternFunc) IsExportedDecl() bool {
	return true
}

func (e DeclExternFunc) DeclOverview() string {
	if len(e.Parameters) == 0 {
		return fmt.Sprintf("extern %s { => }", e.Name)
	}
	paramNames := make([]string, len(e.Parameters))
	for i, param := range e.Parameters {
		paramNames[i] = string(param.Name)
	}
	return fmt.Sprintf("extern %s { %s => }", e.Name, strings.Join(paramNames, ", "))
}

func MakeDeclExternFunc(name Identifier, params []DeclParameter, source *Source) *DeclExternFunc {
	return &DeclExternFunc{
		Name:       name,
		Parameters: params,
		MetaInfo:   &MetaDecl{source},
	}
}

func (decl DeclExternFunc) ProvidedDocs() *Docs {
	return decl.Docs
}
