package ast

var _ Expr = ExprDict{}

type ExprDict struct {
	Entries []ExprDictEntry

	MetaInfo *MetaExpr
}

func (e ExprDict) Meta() *MetaExpr {
	return e.MetaInfo
}

func MakeExprDict(entries []ExprDictEntry, source *Source) *ExprDict {
	return &ExprDict{
		Entries: entries,
		MetaInfo: &MetaExpr{
			Source: source,
		},
	}
}

func (e ExprDict) EnumerateNestedDecls(enumerate func(interface{}, []Decl)) {
	for _, el := range e.Entries {
		el.Key.EnumerateNestedDecls(enumerate)
		el.Value.EnumerateNestedDecls(enumerate)
	}
}

type ExprDictEntry struct {
	Key      Expr
	Value    Expr
	MetaInfo *MetaExpr
}

func MakeExprDictEntry(key Expr, value Expr, source *Source) *ExprDictEntry {
	return &ExprDictEntry{
		Key:   key,
		Value: value,
		MetaInfo: &MetaExpr{
			Source: source,
		},
	}
}
