package activated

import oscillator "../../../../oscillator/stepped"

// Histeresis are Triggering oscillators with a second threshold.
// their State depends on their history;
// between the two thresholds the State can vary depending on which threshold was transited most recently. 
type Banded struct {
	*Triggering
	Threshold
}

func NewBanded(sc uint, damping float64, high, low Threshold) *Banded {
	return &Histeresis{NewTriggering(oscillator.NewDrivenAndDamped(sc,damping),high), low}
}

func (b *Banded) StepChange(d float64) bool {
	b.Step(d)
	if b.Triggered {
		if b.Threshold.Exceeded(h.Abs()){
//		if h.Abs2() < h.Threshold {
			b.Triggered = Triggered(false)
			return true
		}
	} else {
		if b.Triggering.Threshold.Exceeded(h.Abs()){
//		if h.Abs2() > h.Triggering.Threshold2 {
			b.Triggered = Triggered(true)
			return true
		}
	}
	return false
}



