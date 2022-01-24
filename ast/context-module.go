package ast

type ModuleName string

type ContextModule struct {
	Name ModuleName

	Files []*SourceFile
}

func MakeContextModule(name ModuleName) *ContextModule {
	return &ContextModule{
		Name:  name,
		Files: []*SourceFile{},
	}
}

func (m *ContextModule) AddSourceFile(sourceFile *SourceFile) {
	m.Files = append(m.Files, sourceFile)
}
