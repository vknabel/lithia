package tea

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vknabel/lithia/ast"
	"github.com/vknabel/lithia/external/rx"
	"github.com/vknabel/lithia/runtime"
)

var _ tea.Model = TeaModel{}
var _ runtime.RuntimeValue = TeaModel{}

type TeaModel struct {
	env *runtime.Environment

	state  runtime.Evaluatable
	init   runtime.Evaluatable
	update runtime.RuntimeValue
	view   runtime.RuntimeValue
}

func MakeTeaModel(env *runtime.Environment, init runtime.Evaluatable, update runtime.RuntimeValue, view runtime.RuntimeValue) TeaModel {
	model := TeaModel{env, nil, init, update, view}
	state := runtime.NewLazyRuntimeValue(func() (runtime.RuntimeValue, *runtime.RuntimeError) {
		init, err := model.init.Evaluate()
		if err != nil {
			return nil, err
		}
		state, _ := model.applyUpdate(init)
		return state, nil
	})
	model.state = state
	return model
}

// Lookup implements runtime.RuntimeValue
func (m TeaModel) Lookup(name string) (runtime.Evaluatable, *runtime.RuntimeError) {
	switch name {
	case "init":
		return m.init, nil
	case "update":
		return runtime.NewConstantRuntimeValue(m.update), nil
	case "view":
		return runtime.NewConstantRuntimeValue(m.view), nil
	default:
		return nil, runtime.NewRuntimeErrorf("tea.Model %s has no member %s", fmt.Sprint(m), name)
	}
}

// RuntimeType implements runtime.RuntimeValue
func (TeaModel) RuntimeType() runtime.RuntimeTypeRef {
	return TeaModelTypeRef
}

// String implements runtime.RuntimeValue
func (TeaModel) String() string {
	return "tea.Model"
}

func (m TeaModel) applyUpdate(runtimeUpdate runtime.RuntimeValue) (runtime.RuntimeValue, tea.Cmd) {
	data, ok := runtimeUpdate.(runtime.DataRuntimeValue)
	if !ok {
		panic("Model.init and Model.update must return tea.Update")
	}
	lazyState, err := data.Lookup("state")
	if err != nil {
		panic(err)
	}
	state, err := lazyState.Evaluate()
	if err != nil {
		panic(err)
	}

	lazyCmd, err := data.Lookup("cmd")
	if err != nil {
		panic(err)
	}
	cmd, err := lazyCmd.Evaluate()
	if err != nil {
		panic(err)
	}
	noneType := runtime.MakeRuntimeTypeRef("None", "prelude")
	isNone, err := noneType.HasInstance(cmd)
	if err != nil {
		panic(err)
	}
	if isNone {
		return state, nil
	} else if callableCmd, ok := cmd.(runtime.CallableRuntimeValue); ok {
		return state, func() tea.Msg {
			fmt.Printf("calling %T, %T\n", callableCmd, lazyCmd)
			value, err := runtime.Call(callableCmd, []runtime.Evaluatable{lazyState}, nil)
			if err != nil {
				panic(err)
			}
			return toTeaMsgIfNeeded(value)
		}
	} else if futureCmd, ok := cmd.(rx.RxFuture); ok {
		return state, func() tea.Msg {
			value, err := futureCmd.Await()
			if err != nil {
				panic(err)
			}
			return toTeaMsgIfNeeded(value)
		}
	} else {
		return state, func() tea.Msg {
			return toTeaMsgIfNeeded(cmd)
		}
	}
}

func toTeaMsgIfNeeded(value runtime.RuntimeValue) tea.Msg {
	dataValue, ok := value.(runtime.DataRuntimeValue)
	if !ok {
		return value
	}
	mappedData := map[ast.Identifier]func() tea.Msg{
		"ClearScrollArea":       tea.ClearScrollArea,
		"DisableMouse":          tea.DisableMouse,
		"EnableMouseAllMotion":  tea.EnableMouseAllMotion,
		"EnableMouseCellMotion": tea.EnableMouseCellMotion,
		"EnterAltScreen":        tea.EnterAltScreen,
		"ExitAltScreen":         tea.ExitAltScreen,
		"HideCursor":            tea.HideCursor,
		"Quit":                  tea.Quit,
	}
	for name, fn := range mappedData {
		typeRef := runtime.MakeRuntimeTypeRef(name, "tea")
		has, err := typeRef.HasInstance(dataValue)
		if err != nil {
			panic(err)
		}
		if has {
			return fn()
		}
	}
	return value
}

func (m TeaModel) Init() tea.Cmd {
	init, err := m.init.Evaluate()
	if err != nil {
		panic(fmt.Sprintf("init: %s", err))
	}
	_, cmd := m.applyUpdate(init)
	return cmd
}

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (m TeaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var evaluatableMsg runtime.Evaluatable
	switch msg := msg.(type) {
	case runtime.Evaluatable:
		evaluatableMsg = msg
	case runtime.RuntimeValue:
		evaluatableMsg = runtime.NewConstantRuntimeValue(msg)
	case tea.KeyMsg:
		runtimeMsg, err := m.env.MakeDataRuntimeValue("Key", map[string]runtime.Evaluatable{
			// "type":   runtime.NewConstantRuntimeValue(msg.Type),
			// "alt":    m.env.BoolToRuntimeValue(msg.Alt),
			"string": runtime.NewConstantRuntimeValue(runtime.PreludeString(msg.String())),
		})
		if err != nil {
			panic(err)
		}
		evaluatableMsg = runtime.NewConstantRuntimeValue(runtimeMsg)
	case tea.MouseMsg:
		runtimeMsg, err := m.env.MakeDataRuntimeValue("Key", map[string]runtime.Evaluatable{
			// "type":   runtime.NewConstantRuntimeValue(msg.Type),
			// "alt":    m.env.BoolToRuntimeValue(msg.Alt),
			// "string": runtime.NewConstantRuntimeValue(runtime.PreludeString(msg.String())),
		})
		if err != nil {
			panic(err)
		}
		evaluatableMsg = runtime.NewConstantRuntimeValue(runtimeMsg)
	case tea.WindowSizeMsg:
		runtimeMsg, err := m.env.MakeDataRuntimeValue("WindowSize", map[string]runtime.Evaluatable{
			"width":  runtime.NewConstantRuntimeValue(runtime.PreludeInt(msg.Width)),
			"height": runtime.NewConstantRuntimeValue(runtime.PreludeInt(msg.Height)),
		})
		if err != nil {
			panic(err)
		}
		evaluatableMsg = runtime.NewConstantRuntimeValue(runtimeMsg)
	default:
		return m, nil
	}

	updateValue, err := runtime.Call(m.update, []runtime.Evaluatable{
		m.state,
		evaluatableMsg,
	}, nil)
	if err != nil {
		panic(err)
	}
	state, cmd := m.applyUpdate(updateValue)
	return TeaModel{m.env, runtime.NewConstantRuntimeValue(state), m.init, m.update, m.view}, cmd
}

// View renders the program's UI, which is just a string. The view is
// rendered after every Update.
func (m TeaModel) View() string {
	viewValue, err := runtime.Call(m.view, []runtime.Evaluatable{
		m.state,
	}, nil)
	if err != nil {
		panic(err)
	}
	view, ok := viewValue.(runtime.PreludeString)
	if !ok {
		panic(fmt.Sprintf("view must return a string, got %s", viewValue))
	}
	return view.String()
}
