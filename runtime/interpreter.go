package runtime

import (
	"github.com/vknabel/lithia/ast"
	"github.com/vknabel/lithia/parser"
	"github.com/vknabel/lithia/resolution"
	"github.com/vknabel/lithia/world"
)

var interpreter *Interpreter

type Interpreter struct {
	Resolver            resolution.ModuleResolver
	Parser              *parser.Parser
	Modules             map[ast.ModuleName]*RuntimeModule
	ExternalDefinitions map[ast.ModuleName]ExternalDefinition
	Prelude             *Environment
}

func NewIsolatedInterpreter(referenceFile string, importRoots ...string) *Interpreter {
	inter := &Interpreter{
		Resolver:            resolution.NewDefaultModuleResolver(importRoots...),
		Parser:              parser.NewParser(),
		Modules:             make(map[ast.ModuleName]*RuntimeModule),
		ExternalDefinitions: make(map[ast.ModuleName]ExternalDefinition),
	}
	inter.ExternalDefinitions["prelude"] = ExternalPrelude{}
	interpreter = inter
	return inter
}

func (inter *Interpreter) LoadExternalDefinition(name ast.ModuleName, definition ExternalDefinition) {
	inter.ExternalDefinitions[name] = definition
}

func (inter *Interpreter) Interpret(fileName string, script string) (RuntimeValue, error) {
	pkg := inter.Resolver.ResolvePackageForReferenceFile(fileName)
	resolvedModule := inter.Resolver.CreateSingleFileModule(pkg, fileName)
	module := inter.NewModule(resolvedModule)
	ix, err := inter.LoadFileIntoModule(module, fileName, script)
	if err != nil {
		return nil, err
	}
	rval, rerr := ix.Evaluate()
	if rerr != nil {
		return rval, *rerr
	} else {
		return rval, nil
	}
}

func (inter *Interpreter) InterpretEmbed(fileName string, script string) (RuntimeValue, error) {
	pkg := inter.Resolver.ResolvePackageForReferenceFile(fileName)
	resolvedModule := inter.Resolver.CreateSingleFileModule(pkg, fileName)
	moduleName := resolvedModule.AbsoluteModuleName()
	module := inter.Modules[moduleName]
	if module == nil {
		module = inter.NewModule(resolvedModule)
	}
	ex, err := inter.EmbedFileIntoModule(module, fileName, script)
	if err != nil {
		return nil, err
	}
	rval, rerr := ex.Evaluate()
	if rerr != nil {
		return rval, *rerr
	} else {
		return rval, nil
	}
}

func (inter *Interpreter) LoadFileIntoModule(module *RuntimeModule, fileName string, script string) (*InterpreterContext, error) {
	fileParser, errs := inter.Parser.Parse(module.Name, fileName, script)
	if len(errs) > 0 {
		return nil, parser.NewGroupedSyntaxError(errs)
	}
	sourceFile, errs := fileParser.ParseSourceFile()
	if len(errs) > 0 {
		return nil, parser.NewGroupedSyntaxError(errs)
	}
	module.Decl.AddSourceFile(sourceFile)
	ix := inter.NewInterpreterContext(sourceFile, module, fileParser.Tree.RootNode(), []byte(script), module.Environment.Private())
	module.Files[FileName(fileName)] = ix

	for _, decl := range sourceFile.Declarations {
		declValue, err := MakeRuntimeValueDecl(ix, decl)
		if err != nil {
			return nil, err
		}
		if decl.IsExportedDecl() {
			ix.environment.DeclareExported(string(decl.DeclName()), declValue)
		} else {
			ix.environment.DeclareUnexported(string(decl.DeclName()), declValue)
		}
	}
	return ix, nil
}

func (inter *Interpreter) EmbedFileIntoModule(module *RuntimeModule, fileName string, script string) (*InterpreterContext, error) {
	fileParser, errs := inter.Parser.Parse(module.Name, fileName, script)
	if len(errs) > 0 {
		return nil, parser.NewGroupedSyntaxError(errs)
	}
	sourceFile, errs := fileParser.ParseSourceFile()
	if sourceFile == nil {
		return nil, parser.NewGroupedSyntaxError(errs)
	}
	module.Decl.AddSourceFile(sourceFile)
	ex := inter.NewInterpreterContext(sourceFile, module, fileParser.Tree.RootNode(), []byte(script), module.Environment)

	for _, decl := range sourceFile.Declarations {
		declValue, err := MakeRuntimeValueDecl(ex, decl)
		if err != nil {
			return nil, err
		}
		if decl.IsExportedDecl() {
			ex.environment.DeclareExported(string(decl.DeclName()), declValue)
		} else {
			ex.environment.DeclareUnexported(string(decl.DeclName()), declValue)
		}
	}
	return ex, nil
}

func (inter *Interpreter) LoadModuleIfNeeded(queryModuleName ast.ModuleName, fromResolvedModule resolution.ResolvedModule) (*RuntimeModule, error) {
	resolvedModule, err := inter.Resolver.ResolveModuleFromPackage(fromResolvedModule.Package(), queryModuleName)
	if err != nil {
		return nil, err
	}
	moduleName := resolvedModule.AbsoluteModuleName()
	if module, ok := inter.Modules[moduleName]; ok {
		return module, nil
	}
	module := inter.NewModule(resolvedModule)
	contexts, err := inter.LoadFilesIntoModule(module, resolvedModule.Files)
	if err != nil {
		return module, err
	}
	for _, context := range contexts {
		_, err := context.Evaluate()
		if err != nil {
			return module, err
		}
	}
	return module, nil
}

func (inter *Interpreter) LoadFilesIntoModule(module *RuntimeModule, files []string) ([]InterpreterContext, error) {
	var contexts []InterpreterContext
	for _, file := range files {
		scriptData, err := world.Current.FS.ReadFile(file)
		if err != nil {
			return nil, err
		}
		context, err := inter.LoadFileIntoModule(module, file, string(scriptData))
		if err != nil {
			return nil, err
		}
		if context != nil {
			contexts = append(contexts, *context)
		}
	}
	return contexts, nil
}
