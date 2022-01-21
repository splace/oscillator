package discrete

// Stabilising's are Damped StepChanger's whose StepChange returns true when their amplitude first stabilises.(without Triggered being reset no further StepChange's are produced.) 
// they are Steppers whose Step doesn't trigger a change.
type Stabilising struct {
	Stepper
	Triggered
	IsStateTrigger
}

func NewStabilising(s Stepper, m float64) *Stabilising {
	return &Stabilising{Stepper: s, IsStateTrigger: WithinMargin(m)}
}

type IsStateTrigger func(State) bool

// returns a Trigger for when its parameter value changes by less than the presented margin.
func WithinMargin(margin float64) IsStateTrigger{
	var lv,dv,v float64
	return func(s State) bool{
		v=s.StateGetter().Amplitude()
		dv,lv=v-lv,v
		return dv<margin && dv>-margin
	}
} 

func (s Stabilising) StepChange(d float64) bool {
	s.Step(d)
	if s.Triggered {
		return false
	}
	if s.IsStateTrigger(s.StateGetter()) {
		s.Triggered = true
		return true
	}
	return false
}
