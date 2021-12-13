// package containing oscillators driven by the same
package group

import "../../stepped"

type Steppers []oscillator.Stepper

func (ss Steppers) Step(d float64){
	for _,o:=range ss{
		o.Step(d)
	}
}

func (ss Steppers) Abs() (f float64){
	for _,o:=range ss{
		f+=o.Abs()
	}
	f/=float64(len(ss))
	return
}

func (ss Steppers) Real() (f float64){
	for _,o:=range ss{
		f+=o.Real()
	}
	return
}

