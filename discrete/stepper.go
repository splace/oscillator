package oscillator

import "math/cmplx"

// an osillator that can be stepped forward by an internally stored amount and also return its current abstract state.
type Stepper interface {
	Step(float64)
	Abs() float64
	Real() float64
}

type State complex128

func (s State) Abs() float64{
	return cmplx.Abs(complex128(s))
}

func (s State) Real() float64{
	return real(complex128(s))
}

