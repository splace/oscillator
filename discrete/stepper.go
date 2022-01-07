// discrete simulation of oscillators using complex numbers.
package oscillator

import "math"
import "math/cmplx"

const angleScaler = 1.0 / math.Pi // scale phase angle (+pi,-pi) -> (+1,-1)

// Steppers can Step an amount.
// In this package; they will be an oscillator that Step's forward by an internally stored amount with an additional amount to drive them. (out of phase)
// also they can return their current State.
type Stepper interface {
	Step(float64)
	StateGetter() State
}

// State is a Stepper whose drive steps simply add to the imaginary part of itself.
// with drive imaginary only, orthoganal, oscillators (embedding this) are driven out of phase and as such change their amplitude but not their phase.
type State complex128

func (s *State) Step(d float64) {
	*s += State(complex(0, d))
}

func (s State) StateGetter() State {
	return s
}

func (s State) Amplitude() float64 {
	return cmplx.Abs(complex128(s))
}

func (s State) Phase() float64 {
	return cmplx.Phase(complex128(s))
}

// Steppers are group of Stepper's that together are a Stepper, simultaneously Stepping all contained Stepper's and whose StateGetter returns the sum of all the contained Steppers States.
type Steppers []Stepper

func (ss Steppers) Step(d float64) {
	for _, o := range ss {
		o.Step(d)
	}
}

func (ss Steppers) StateGetter() (s State) {
	for _, o := range ss {
		s += o.StateGetter()
	}
	return
}

// AvgSteppers wraps a Steppers whose State is the average of the contained Steppers States.
type AvgSteppers struct {
	Steppers
}

func (ass AvgSteppers) StateGetter() State {
	s := ass.Steppers.StateGetter()
	s /= State(complex(float64(len(ass.Steppers)), 0))
	return s
}
