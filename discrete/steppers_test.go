package oscillator_test


import "fmt"

import osc "../../oscillator/discrete"

func ExampleSteppersImage() {
	const (
		steps = 12800
		cycles = 2
	)
	// build oscillator group, undamped and set with frequency and amplitude of square wave fourier transform
	var os osc.Steppers 
	for n:=1.0;n<10;n+=2{
		o:=osc.NewDamped(steps/n,0)
		o.Step(64/n)
		os=append(os,o)
	}

	var r []byte
	for r=make([]byte,0,steps*cycles);len(r)!=cap(r);os.Step(0){
		r=append(r,byte(real(os.GetState())+128))
//		r=append(r,byte(real(o[0].(*osc.Damped).GetState())/float64(len(o))+128))
//		a=append(a,byte(o[1].(*osc.Damped).Amplitude()+128))
	}
	fmt.Println(osc.Encode(fmt.Sprintf("Steppers(%v).png", len(r)),r))
	// Output:
	// <nil>
}
