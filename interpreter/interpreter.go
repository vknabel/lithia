package interpreter

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/vknabel/go-lithia/parser"
)

type Interpreter struct {
	importRoots []string
	parser      *parser.Parser
	modules     map[ModuleName]*Module
	prelude     *Environment
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
	return &Interpreter{
		importRoots: absoluteImportRoots,
		parser:      parser.NewParser(),
		modules:     make(map[ModuleName]*Module),
	}
}

func (inter *Interpreter) Interpret(fileName string, script string) (RuntimeValue, error) {
	moduleName := ModuleName(strings.ReplaceAll(filepath.Base(fileName), ".", "_"))
	module := inter.NewModule(moduleName)
	ex, err := inter.LoadFileIntoModule(module, fileName, script)
	if err != nil {
		return nil, err
	}
	lazyValue, err := ex.EvaluateNode()
	if err != nil {
		return nil, err
	}
	return lazyValue.Evaluate()
}
func (inter *Interpreter) InterpretEmbed(fileName string, script string) (RuntimeValue, error) {
	moduleName := ModuleName(strings.ReplaceAll(filepath.Base(fileName), ".", "_"))
	module := inter.modules[moduleName]
	if module == nil {
		module = inter.NewModule(moduleName)
	}
	ex, err := inter.EmbedFileIntoModule(module, fileName, script)
	if err != nil {
		return nil, err
	}
	lazyValue, err := ex.EvaluateNode()
	if err != nil {
		return nil, err
	}
	return lazyValue.Evaluate()
}

func (inter *Interpreter) LoadFileIntoModule(module *Module, fileName string, script string) (*EvaluationContext, error) {
	tree, err := inter.parser.Parse(script)
	if err != nil {
		return nil, inter.SyntaxParsingError(fileName, script, tree)
	}
	ex := inter.NewEvaluationContext(fileName, module, tree.RootNode(), []byte(script), module.environment.Private())
	return ex, nil
}

func (inter *Interpreter) EmbedFileIntoModule(module *Module, fileName string, script string) (*EvaluationContext, error) {
	tree, err := inter.parser.Parse(script)
	if err != nil {
		return nil, inter.SyntaxParsingError(fileName, script, tree)
	}
	ex := inter.NewEvaluationContext(fileName, module, tree.RootNode(), []byte(script), module.environment)
	return ex, nil
}

func (ex *EvaluationContext) EvaluateSourceFile() (*LazyRuntimeValue, error) {
	count := ex.node.ChildCount()
	children := make([]*sitter.Node, count)
	for i := uint32(0); i < count; i++ {
		child := ex.node.Child(int(i))
		children[i] = child
	}
	sort.SliceStable(children, func(i, j int) bool {
		lp := priority(children[i].Type())
		rp := priority(children[j].Type())
		return lp > rp
	})

	return NewLazyRuntimeValue(func() (RuntimeValue, error) {
		var lastValue RuntimeValue
		for _, child := range children {
			lazyValue, err := ex.ChildNodeExecutionContext(child).EvaluateNode()
			if err != nil {
				return nil, err
			}
			if lazyValue != nil {
				lastValue, err = lazyValue.Evaluate()
				if err != nil {
					return nil, err
				}
			}
		}
		return lastValue, nil
	}), nil
}

func priority(nodeType string) int {
	switch nodeType {
	case parser.TYPE_NODE_MODULE_DECLARATION:
		return 19
	case parser.TYPE_NODE_IMPORT_DECLARATION:
		return 17
	case parser.TYPE_NODE_DATA_DECLARATION:
		return 15
	case parser.TYPE_NODE_ENUM_DECLARATION:
		return 13
	case parser.TYPE_NODE_FUNCTION_DECLARATION:
		return 7
	case parser.TYPE_NODE_LET_DECLARATION:
		return 3
	default:
		return 0
	}
}

func (ex *EvaluationContext) EvaluateModule() (*LazyRuntimeValue, error) {
	internalName := ex.node.ChildByFieldName("name").Content(ex.source)
	runtimeModule := NewConstantRuntimeValue(RuntimeModule{module: ex.module})
	ex.environment.DeclareUnexported(internalName, runtimeModule)
	return runtimeModule, nil
}

func (inter *Interpreter) LoadModuleIfNeeded(moduleName ModuleName) (*Module, error) {
	if module, ok := inter.modules[moduleName]; ok {
		return module, nil
	}
	for _, root := range inter.importRoots {
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
		childContext, err := inter.LoadFileIntoModule(module, file, string(scriptData))
		if err != nil {
			return err
		}
		source, err := childContext.EvaluateNode()
		if err != nil {
			return err
		}
		_, err = source.Evaluate()
		if err != nil {
			return err
		}
	}
	return nil
}
