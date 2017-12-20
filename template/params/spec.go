package params

type Spec interface {
	Rule() Rule
	Reducers() []Reducer
	Validators() Validators
}

func NewSpec(r Rule, vs ...Validators) Spec {
	return newSpec(r, vs...)
}

func SpecBuilder(r Rule, vs ...Validators) *specBuilder {
	b := &specBuilder{newSpec(r, vs...)}
	b.s.reds = make([]Reducer, 0)
	return b
}

type specBuilder struct {
	s *spec
}

func (b *specBuilder) AddReducer(r reduceFunc, keys ...string) *specBuilder {
	b.s.reds = append(b.s.reds, newReducer(r, keys...))
	return b
}

func (b *specBuilder) Done() Spec {
	return b.s
}

func newSpec(r Rule, vs ...Validators) *spec {
	if len(vs) > 0 {
		return &spec{r: r, v: vs[0]}
	}
	return &spec{r: r}
}

type spec struct {
	r    Rule
	reds []Reducer
	v    Validators
}

func (s *spec) Rule() Rule {
	if s.r == nil {
		return None()
	}
	return s.r
}

func (s *spec) Validators() Validators {
	if s.v == nil {
		return Validators(make(map[string]validatorFunc, 0))
	}
	return s.v
}

func (s *spec) Reducers() []Reducer {
	return s.reds
}
