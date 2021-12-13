package oscillator

// StepChanger's are Stepper's that return if a step caused an internal State change.
type StepChanger interface {
	Stepper
	StepChange(float64) bool
}

// wraps any StepChanger to channel the State after any change.
// uint should be set to identify which StepChanger caused the event.
type Eventing struct {
	StepChanger
	Ident   uint
	Changes chan Change
}

func (e *Eventing) Step(d float64) {
	if e.StepChange(d) {
		e.Changes <- Change{e.Ident, e.StateGetter()}
	}
}

type Change struct {
	uint
	State
}
