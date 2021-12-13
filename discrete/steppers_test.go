package oscillator_test


import "fmt"

import osc "../../oscillator/discrete"

func ExampleSteppersImage() {
	const (
		steps = 12800
		cycles = 2
	)
	// build oscillator group
	// set frequency and amplitude to that off a square wave fourier transform.
	var os osc.Steppers 
	for n:=1.0;n<40;n+=2{
		o:=osc.NewDamped(steps/n,1/n)  // 1/n damping makes all ossillators damp at the same rate.
		o.Step(100/n)  // 100 won't exceed byte (+127 - -127) when all summed
		os=append(os,o)
	}

	var r []byte
	for r=make([]byte,0,steps*cycles);len(r)!=cap(r);os.Step(0){
		r=append(r,byte(real(os.StateGetter())+128))
//		r=append(r,byte(real(o[0].(*osc.Damped).GetState())/float64(len(o))+128))
//		a=append(a,byte(o[1].(*osc.Damped).Amplitude()+128))
	}
	fmt.Println(osc.Encode(fmt.Sprintf("Steppers(%v).png", len(r)),r))
	// Output:
	// <nil>
}
