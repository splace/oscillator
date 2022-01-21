package discrete

// Stabilising's are Damped StepChanger's whose StepChange returns true when their amplitude first stabilises.(without Triggered being reset no further StepChange's are produced.) 
// they are Steppers whose Step doesn't trigger a change.
type Stabilising struct {
	Stepper
	Triggered
	IsStateTriggering
}

func NewStabilising(s Stepper, m float64) *Stabilising {
	return &Stabilising{Stepper: s, IsStateTriggering: AmplitudeWithin(m)}
}

type IsStateTriggering func(State) bool

// returns a Trigger for when its parameter value changes by less than the presented margin.
func AmplitudeWithin(margin float64) IsStateTriggering{
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
	if s.IsStateTriggering(s.StateGetter()) {
		s.Triggered = true
		return true
	}
	return false
}
