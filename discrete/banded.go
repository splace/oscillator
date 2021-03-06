package oscillator

// Banded are Triggering oscillators with an additional Threshold, used for transitions to being not-Triggered.
// Notice: Triggered not only depends on the State, and the Thresholds, but also on which Threshold was transitioned last. There is Hysteresis.
type Banded struct {
	*Triggering
	Threshold
}

// returns nil for inappropriate use.
func NewBanded(sc float64, damping float64, upper, lower Threshold) *Banded {
	if lower > upper {
		return nil
	}
	return &Banded{NewTriggering(NewDamped(sc, damping), upper), lower}
}

func (b *Banded) StepChange(d float64) bool {
	b.Step(d)
	if b.Triggered {
		if b.Triggering.Threshold.Exceeded(b.Amplitude()) {
			b.Triggered = true
			return true
		}
	} else {
		if !b.Threshold.Exceeded(b.Amplitude()) {
			b.Triggered = false
			return true
		}
	}
	return false
}
