package activated

// StepChanger's are Stepper's that can return if a step causes an internal state change
type StepChanger interface{
	Triggering
	StepChange(float64) bool
}

type Change struct {
	Index uint
	Triggered
}

// wraps any StepChanger to provide any, indexed, state change on a channel.
type Eventing struct{
	StepChanger
	Index uint
	Changes chan Change
}

func (e *Eventing) Step(d float64) {
	e.Index++
	if e.StepChange(d){
		e.Changes <- Change{e.Index,e.StepChanger.Triggered}
	}
}
