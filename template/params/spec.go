package params

type Spec interface {
	Rule() Rule
	Validators() Validators
}

func NewSpec(r Rule, vs ...Validators) Spec {
	if len(vs) > 0 {
		return &spec{r: r, v: vs[0]}
	}
	return &spec{r: r}
}

type spec struct {
	r Rule
	v Validators
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
