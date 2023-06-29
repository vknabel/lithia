package worldtest

type MapEnv struct {
	Map      map[string]string
	ExitCode *int
}

func NewMapEnv(m map[string]string) *MapEnv {
	return &MapEnv{
		Map: m,
	}
}

func (e *MapEnv) LookupEnv(key string) (string, bool) {
	val, ok := e.Map[key]
	return val, ok
}

func (e *MapEnv) Exit(code int) {
	e.ExitCode = &code
}

func (e *MapEnv) Environ() []string {
	env := make([]string, 0, len(e.Map))
	for key, val := range e.Map {
		env = append(env, key+"="+val)
	}
	return env
}
