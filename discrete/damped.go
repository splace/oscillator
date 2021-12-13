package oscillator

import "math/cmplx"
import "math"


type DrivenAndDamped struct{
	State
	Stepping
}

type Stepping complex128

// build from discretization, and damping, per cycle.
func NewDrivenAndDamped(stepsPerCycle uint,damping float64) *DrivenAndDamped {
	return &DrivenAndDamped{Stepping:Stepping(cmplx.Rect((1 - damping/float64(stepsPerCycle)), math.Pi*2/float64(stepsPerCycle)))} 
}

func (g DrivenAndDamped) Abs2() float64 {
	return real(g.State)*real(g.State)+imag(g.State)*imag(g.State)
}

func (g *DrivenAndDamped) Step(d float64) {
	g.State*=State(g.Stepping)
	g.State+=State(complex(0,d)) // drive imaginary, at right angles, no effect on phase, FUNDAMENTAL to algorithm
} 




