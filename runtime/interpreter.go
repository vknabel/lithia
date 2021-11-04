package runtime

import (
	"os"
	"path/filepath"

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

// func (inter *Interpreter) Interpret(fileName string, script string) (RuntimeValue, error) {
// 	moduleName := ast.ModuleName(strings.ReplaceAll(filepath.Base(fileName), ".", "_"))
// 	module := inter.NewModule(moduleName)
// 	ex, err := inter.LoadFileIntoModule(module, fileName, script)
// 	if err != nil {
// 		return nil, err
// 	}
// 	lazyValue, err := ex.EvaluateNode()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return lazyValue.Evaluate()
// }
// func (inter *Interpreter) InterpretEmbed(fileName string, script string) (RuntimeValue, error) {
// 	moduleName := ast.ModuleName(strings.ReplaceAll(filepath.Base(fileName), ".", "_"))
// 	module := inter.Modules[moduleName]
// 	if module == nil {
// 		module = inter.NewModule(moduleName)
// 	}
// 	ex, err := inter.EmbedFileIntoModule(module, fileName, script)
// 	if err != nil {
// 		return nil, err
// 	}
// 	lazyValue, err := ex.EvaluateNode()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return lazyValue.Evaluate()
// }

// func (inter *Interpreter) LoadFileIntoModule(module *Module, fileName string, script string) (*InterpreterContext, error) {
// 	fileParser, err := inter.Parser.Parse(module.Name, fileName, script)
// 	sourceFile, errs := fileParser.ParseSourceFile()
// 	if err != nil {
// 		return nil, inter.SyntaxParsingError(fileName, script, tree)
// 	}
// 	ex := inter.NewInterpreterContext(fileName, module, tree.RootNode(), []byte(script), module.environment.Private())
// 	return ex, nil
// }

// func (inter *Interpreter) EmbedFileIntoModule(module *Module, fileName string, script string) (*InterpreterContext, error) {
// 	tree, err := inter.Parser.Parse(script)
// 	if err != nil {
// 		return nil, inter.SyntaxParsingError(fileName, script, tree)
// 	}
// 	ex := inter.NewInterpreterContext(fileName, module, tree.Node, []byte(script), module.environment)
// 	return ex, nil
// }

// func (inter *Interpreter) LoadModuleIfNeeded(moduleName ast.ModuleName) (*Module, error) {
// 	if module, ok := inter.Modules[moduleName]; ok {
// 		return module, nil
// 	}
// 	for _, root := range inter.ImportRoots {
// 		relativeModulePath := strings.ReplaceAll(string(moduleName), ".", string(filepath.Separator))
// 		modulePath := filepath.Join(root, relativeModulePath)
// 		matches, err := filepath.Glob(filepath.Join(modulePath, "*.lithia"))
// 		if err != nil {
// 			continue
// 		}
// 		if len(matches) == 0 {
// 			continue
// 		}

// 		module := inter.NewModule(moduleName)
// 		err = inter.LoadFilesIntoModule(module, matches)
// 		if err != nil {
// 			return module, err
// 		}

// 		return module, nil
// 	}
// 	return nil, fmt.Errorf("module %s not found", moduleName)
// }

// func (inter *Interpreter) LoadFilesIntoModule(module *Module, files []string) error {
// 	for _, file := range files {
// 		scriptData, err := os.ReadFile(file)
// 		if err != nil {
// 			return err
// 		}
// 		childContext, err := inter.LoadFileIntoModule(module, file, string(scriptData))
// 		if err != nil {
// 			return err
// 		}
// 		source, err := childContext.EvaluateNode()
// 		if err != nil {
// 			return err
// 		}
// 		_, err = source.Evaluate()
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }
