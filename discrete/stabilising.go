package oscillator

// Stabilising's are Damped StepChanger's whose StepChange returns true when their amplitude first stabilises.(can only do this more thasn once if Triggered is reset) 
type Stabilising struct {
	*Damped
	Triggered
	Margin
}

type Margin func(float64) bool

func NewMargin(m float64) Margin{
	var lv,dv float64
	return func(v float64) bool{
		dv=v-lv
		lv=v
		return dv<m && dv>-m
	}
} 


func NewStabilising(o *Damped, m float64) *Stabilising {
	return &Stabilising{Damped: o, Margin: NewMargin(m)}
}

func (s *Stabilising) Step(d float64) {
	s.StepChange(d)
}

func (s *Stabilising) StepChange(d float64) bool {
	s.Damped.Step(d)
	if s.Triggered {
		return false
	}
	if s.Margin(s.Amplitude()) {
		s.Triggered = true
		return true
	}
	return false
}
