// discrete simulation of oscillators using complex numbers.
package oscillator

import "math"
import "math/cmplx"


const angleScaler = 1.0/math.Pi // scale phase angle (+pi,-pi) -> (+1,-1)

// an osillator that can be stepped forward by an internally stored amount with an additional drive (out of phase) and also return its current state.
type Stepper interface {
	Step(float64)
	GetState() State
}

type State complex128

func (s State) GetState() State{
	return s
}

func (s State) Amplitude() float64{
	return cmplx.Abs(complex128(s))
}

func (s State) Phase() float64{
	return cmplx.Phase(complex128(s))
}


type Steppers []Stepper

func (ss Steppers) Step(d float64){
	for _,o:=range ss{
		o.Step(d)
	}
}

func (ss Steppers) GetState() (s State){
	for _,o:=range ss{
		s+=o.GetState()
	}
	return
}

type AvgSteppers struct{
	Steppers
}

func (ass AvgSteppers) GetState() State{
	s:=ass.Steppers.GetState()
	s/=State(complex(float64(len(ass.Steppers)),0))
	return s
}
