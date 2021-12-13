package oscillator

import "math/cmplx"
import "math"

// Damped are simple ossilators with damping. used as base ossilator type for this package.
type Damped struct{
	State
	Stepping
}

type Stepping complex128

// build from discretization, and damping, per cycle.
func NewDamped(stepsPerCycle float64,damping float64) *Damped {
	if !(stepsPerCycle>0) || damping<0 || damping>1 {return nil}
	return &Damped{Stepping:Stepping(cmplx.Rect((1 - damping/stepsPerCycle), math.Pi*2/stepsPerCycle))} 
}

func (g *Damped) Step(d float64) {
	g.State*=State(g.Stepping)
	(&g.State).Step(d)
} 




