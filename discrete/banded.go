package oscillator

// Banded are Triggering oscillators with an additional Threshold, used for transitions to being not-Triggered.
// Notice: Triggered no only depends on the State and the Thresholds but also on which Threshold was transitioned last. There is Hysteresis.
type Banded struct {
	*Triggering
	Threshold
}

// returns nil for inappropriate use.
func NewBanded(sc float64, damping float64, high, low Threshold) *Banded {
	if low>high {return nil}
	return &Banded{NewTriggering(NewDamped(sc,damping),high), low}
}

func (b *Banded) StepChange(d float64) bool {
	b.Step(d)
	if b.Triggered {
		if b.Triggering.Threshold.Exceeded(b.Amplitude()){
			b.Triggered = Triggered(false)
			return true
		}
	} else {
		if !b.Threshold.Exceeded(b.Amplitude()){
			b.Triggered = Triggered(true)
			return true
		}
	}
	return false
}



