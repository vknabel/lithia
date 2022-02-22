package ast

type Documented interface {
	ProvidedDocs() *Docs
}

type Overviewable interface {
	DeclOverview() string
}
