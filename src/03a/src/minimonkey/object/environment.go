package object

type Environment struct {
	store map[string]Object
	outer *Environment
}

func (e *Environment) Get(name string) (Object, bool) {
	val, ok := e.store[name]
	if !ok && e.outer != nil {
		val, ok = e.outer.Get(name)
	}
	return val, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: outer}
}
