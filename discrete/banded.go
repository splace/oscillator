package discrete

// Banded are Triggering oscillators with an additional Threshold, used to separate the value that causes being triggered from the value that causes being untriggered.
// Notice: Banded Triggered have Hysteresis, their state depends on which Threshold was transitioned last.
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
