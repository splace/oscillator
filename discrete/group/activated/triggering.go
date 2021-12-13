package activated

import oscillator "../../../../oscillator/stepped"

// Triggering's are driven, damped, Threshold triggered oscillator
type Triggering struct {
	*oscillator.DrivenAndDamped
	Triggered
	Threshold
}

type Threshold float64

func (t Threshold) Exceeded(v float64) bool{
	return float64(t)>v
}


type Triggered bool

func (s Triggered) String()string{
	if s {return "<"}
	return ">"
}

func NewTriggering(o *oscillator.DrivenAndDamped,t Threshold) *Triggering{
	return &Triggering{DrivenAndDamped:o,Threshold:t}
}

func (t *Triggering) Step(d float64) {
	t.StepChange(d)
}

func (t *Triggering) StepChange(d float64) bool {
	t.DrivenAndDamped.Step(d)
	if isAbove :=  Triggered(t.Exceeded(t.Abs())); isAbove != t.Triggered {
		t.Triggered = isAbove
		return true
	}
	return false
}

