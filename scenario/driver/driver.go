package driver

type Driver interface {
	Lookup(...Token) DriverFn
}

type DriverFn func(map[Token]interface{}) error

type Scenario struct {
	Lines []*Line
}

type Line struct {
	Action   Token
	Resource Token
	Params   map[Token]interface{}
}

type Runner struct {
	Driver
}

func (r *Runner) Run(s *Scenario) error {
	for _, l := range s.Lines {
		if err := r.processLine(l); err != nil {
			return err
		}
	}
	return nil
}

func (r *Runner) processLine(l *Line) error {
	driverFn := r.Lookup(l.Action, l.Resource)
	return driverFn(l.Params)
}
