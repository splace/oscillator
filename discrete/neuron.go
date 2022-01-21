package discrete

// Neuron's are Triggering oscillators with a Threshold that changes.
// Threshold changed by different Adjusters when either Triggered or not.
type Neuron struct {
	*Triggering
	whenTriggered, whenNotTriggered Adjuster
}

// change Threshold
type Adjuster interface {
	Adjust(*Threshold)
}

type Geometric float64

func (a Geometric) Adjust(t *Threshold) {
	(*t) *= Threshold(a)
}

// Geometric but scaled from an offset value.
type OffsetGeometric struct {
	Geometric
	Offset Threshold
}

func (a OffsetGeometric) Adjust(t *Threshold) {
	(*t) -= a.Offset
	a.Geometric.Adjust(t)
	(*t) += a.Offset
}

// makes a Neuron with the same/different Adjusters.
// returns nil for inappropriate use.
func NewNeuron(sc float64, damping float64, t Threshold, f float64, fs ...float64) *Neuron {
	nt := NewTriggering(NewDamped(sc, damping), t)
	if f > 1.0 {
		return nil
	}
	if len(fs) == 0 {
		return &Neuron{nt, Geometric(1 - f/float64(sc)), Geometric(1 + f/float64(sc))}
	}
	if fs[0] > 1.0 {
		return nil
	}
	return &Neuron{nt, Geometric(1 - f/float64(sc)), Geometric(1 + fs[0]/float64(sc))}
}

func (n *Neuron) StepChange(d float64) bool {
	n.Step(d)
	if n.Triggered {
		n.whenTriggered.Adjust(&n.Threshold)
		if !n.Threshold.Exceeded(n.Amplitude()) {
			n.Triggered = Triggered(false)
			return true
		}
	} else {
		n.whenNotTriggered.Adjust(&n.Threshold)
		if n.Threshold.Exceeded(n.Amplitude()) {
			n.Triggered = Triggered(true)
			return true
		}
	}
	return false
}
