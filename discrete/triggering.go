package oscillator

// Triggering's are Damped (implementing Stepper) with amplitude Threshold StepChanger.
type Triggering struct {
	*Damped
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

func NewTriggering(o *Damped,t Threshold) *Triggering{
	return &Triggering{Damped:o,Threshold:t}
}

func (t *Triggering) Step(d float64) {
	t.StepChange(d)
}

func (t *Triggering) StepChange(d float64) bool {
	t.Damped.Step(d)
	if e :=  Triggered(t.Exceeded(t.Amplitude())); e != t.Triggered {
		t.Triggered = e
		return true
	}
	return false
}

