package ast

type SourceFile struct {
	Path         string
	Imports      []ModuleName
	Declarations []Decl
	Statements   []Expr

	*Source
}

func MakeSourceFile(
	Path string,
	Source *Source,
) *SourceFile {
	return &SourceFile{
		Path:         Path,
		Source:       Source,
		Imports:      make([]ModuleName, 0),
		Declarations: make([]Decl, 0),
		Statements:   make([]Expr, 0),
	}
}

func (sf *SourceFile) AddDecl(decl Decl) {
	if decl == nil {
		return
	}
	if importDecl, ok := decl.(DeclImport); ok {
		sf.Imports = append(sf.Imports, importDecl.ModuleName)
	}
	sf.Declarations = append(sf.Declarations, decl)
}

func (sf *SourceFile) AddExpr(expr Expr) {
	if expr == nil {
		return
	}
	sf.Statements = append(sf.Statements, expr)
}