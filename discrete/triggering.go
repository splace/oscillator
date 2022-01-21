package discrete

// Triggering's are Damped with an amplitude Threshold used to make them StepChanger's.
type Triggering struct {
	*Damped
	Triggered
	Threshold
}

type Threshold float64

func (t Threshold) Exceeded(v float64) bool {
	return float64(t) > v
}

type Triggered bool

func (s Triggered) String() string {
	if s {
		return "<"
	}
	return ">"
}

func NewTriggering(o *Damped, t Threshold) *Triggering {
	return &Triggering{Damped: o, Threshold: t}
}

func (t *Triggering) StepChange(d float64) bool {
	t.Step(d)
	if e := Triggered(t.Exceeded(t.Amplitude())); e != t.Triggered {
		t.Triggered = e
		return true
	}
	return false
}
