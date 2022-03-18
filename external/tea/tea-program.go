package tea

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vknabel/lithia/runtime"
)

var _ runtime.RuntimeValue = TeaProgram{}

type TeaProgram struct {
	programType *TeaProgramType
	program     *tea.Program
}

func MakeTeaProgram(programType *TeaProgramType, model TeaModel) TeaProgram {
	p := tea.NewProgram(model, tea.WithoutCatchPanics())
	return TeaProgram{
		programType: programType,
		program:     p,
	}
}

func (TeaProgram) RuntimeType() runtime.RuntimeTypeRef {
	return TeaProgramTypeRef
}

func (v TeaProgram) String() string {
	return string(v.RuntimeType().Name)
}

func (v TeaProgram) Lookup(member string) (runtime.Evaluatable, *runtime.RuntimeError) {
	switch member {
	case "start":
		return runtime.NewLazyRuntimeValue(func() (runtime.RuntimeValue, *runtime.RuntimeError) {
			anyModel, err := v.program.StartReturningModel()
			if err != nil {
				return nil, runtime.NewRuntimeError(err)
			}
			model, ok := anyModel.(TeaModel)
			if !ok {
				return nil, runtime.NewRuntimeErrorf(fmt.Sprintf("unexpected model type %T", anyModel))
			}
			return model.state.Evaluate()
		}), nil
	default:
		return nil, runtime.NewRuntimeErrorf("tea.Program %s has no member %s", fmt.Sprint(v), member)
	}
}
