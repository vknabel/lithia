package runtime

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/vknabel/go-lithia/ast"
	"github.com/vknabel/go-lithia/parser"
)

type Interpreter struct {
	ImportRoots         []string
	Parser              *parser.Parser
	Modules             map[ast.ModuleName]*Module
	ExternalDefinitions map[ast.ModuleName]ExternalDefinition
	Prelude             *Environment
}

func defaultImportRootPaths() []string {
	roots := []string{}
	if path, ok := os.LookupEnv("LITHIA_LOCALS"); ok {
		roots = append(roots, path)
	}
	if path, ok := os.LookupEnv("LITHIA_PACKAGES"); ok {
		roots = append(roots, path)
	}
	if path, ok := os.LookupEnv("LITHIA_STDLIB"); ok {
		roots = append(roots, path)
	} else {
		roots = append(roots, "/usr/local/opt/lithia/stdlib")
	}
	return roots
}

func NewInterpreter(importRoots ...string) *Interpreter {
	importRoots = append(importRoots, defaultImportRootPaths()...)
	absoluteImportRoots := make([]string, len(importRoots))
	for i, root := range importRoots {
		absolute, err := filepath.Abs(root)
		if err == nil {
			absoluteImportRoots[i] = absolute
		} else {
			absoluteImportRoots[i] = root
		}
	}
	inter := &Interpreter{
		ImportRoots:         absoluteImportRoots,
		Parser:              parser.NewParser(),
		Modules:             make(map[ast.ModuleName]*Module),
		ExternalDefinitions: make(map[ast.ModuleName]ExternalDefinition),
	}
	// TODO: External definitions
	// inter.ExternalDefinitions["prelude"] = ExternalPrelude{}
	// inter.ExternalDefinitions["os"] = ExternalOS{}
	// inter.ExternalDefinitions["rx"] = ExternalRx{}
	// inter.ExternalDefinitions["docs"] = ExternalDocs{}
	// inter.ExternalDefinitions["fs"] = ExternalFS{}
	return inter
}

func (inter *Interpreter) LoadExternalDefinition(name ast.ModuleName, definition ExternalDefinition) {
	inter.ExternalDefinitions[name] = definition
}

func (inter *Interpreter) Interpret(fileName string, script string) (RuntimeValue, error) {
	moduleName := ast.ModuleName(strings.ReplaceAll(filepath.Base(fileName), ".", "_"))
	module := inter.NewModule(moduleName)
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
	moduleName := ast.ModuleName(strings.ReplaceAll(filepath.Base(fileName), ".", "_"))
	module := inter.Modules[moduleName]
	if module == nil {
		module = inter.NewModule(moduleName)
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

func (inter *Interpreter) LoadFileIntoModule(module *Module, fileName string, script string) (*InterpreterContext, error) {
	fileParser, err := inter.Parser.Parse(module.Name, fileName, script)
	if err != nil {
		return nil, fileParser.SyntaxErrorOrConvert(err)
	}
	sourceFile, errs := fileParser.ParseSourceFile()
	if sourceFile == nil {
		return nil, errs[0] // TODO: Multiple Errors!
	}
	ix := inter.NewInterpreterContext(sourceFile, module, fileParser.Tree.RootNode(), []byte(script), module.Environment.Private())
	module.Files[FileName(fileName)] = ix
	// TODO: Modules?
	// for _, moduleImport := range sourceFile.Imports {
	// 	ix.environment.DeclareUnexported(string(moduleImport), moduleImport)
	// }
	for _, decl := range sourceFile.Declarations {
		ex := &EvaluationContext{
			Environment: ix.environment,
			Interpreter: ix.interpreter,
		}
		declValue := MakeRuntimeValueDecl(ex, decl)
		ix.environment.DeclareExported(string(decl.DeclName()), declValue)
	}
	return ix, nil
}

func (inter *Interpreter) EmbedFileIntoModule(module *Module, fileName string, script string) (*InterpreterContext, error) {
	fileParser, err := inter.Parser.Parse(module.Name, fileName, script)
	if err != nil {
		return nil, fileParser.SyntaxErrorOrConvert(err)
	}
	sourceFile, errs := fileParser.ParseSourceFile()
	if sourceFile == nil {
		return nil, errs[0] // TODO: Multiple Errors!
	}
	ex := inter.NewInterpreterContext(sourceFile, module, fileParser.Tree.RootNode(), []byte(script), module.Environment)
	return ex, nil
}

func (inter *Interpreter) LoadModuleIfNeeded(moduleName ast.ModuleName) (*Module, error) {
	if module, ok := inter.Modules[moduleName]; ok {
		return module, nil
	}
	for _, root := range inter.ImportRoots {
		relativeModulePath := strings.ReplaceAll(string(moduleName), ".", string(filepath.Separator))
		modulePath := filepath.Join(root, relativeModulePath)
		matches, err := filepath.Glob(filepath.Join(modulePath, "*.lithia"))
		if err != nil {
			continue
		}
		if len(matches) == 0 {
			continue
		}

		module := inter.NewModule(moduleName)
		err = inter.LoadFilesIntoModule(module, matches)
		if err != nil {
			return module, err
		}

		return module, nil
	}
	return nil, fmt.Errorf("module %s not found", moduleName)
}

func (inter *Interpreter) LoadFilesIntoModule(module *Module, files []string) error {
	for _, file := range files {
		scriptData, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		_, err = inter.LoadFileIntoModule(module, file, string(scriptData))
		if err != nil {
			return err
		}
	}
	return nil
}
