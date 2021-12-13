package activated


import "fmt"
//import "log"

import oscillator "../../../../oscillator/stepped"
// import "../../stepped" // actually same as above, just not so obviously.
import "../../../../oscillator/stepped/group"

func ExampleTriggeringImage() {
	const (
		steps = 800
		cycles = 2
	)

	o:=NewTriggering(oscillator.NewDrivenAndDamped(steps/2,.85),50)
	o.Step(127) // initial drive, never exceeded, so this will always fit byte
	
	var r,a,c,t []byte
	for r=make([]byte,0,steps*cycles);len(r)!=cap(r);o.Step(0){
		r=append(r,byte(o.Real()+128))
		a=append(a,byte(o.Abs()+128))
		c=append(c,Iff(o.StepChange(0),byte(255),byte(0)).(byte))
		t=append(t,Iff(bool(o.Triggered),byte(255),byte(0)).(byte))
	}
	fmt.Println(oscillator.Encode(fmt.Sprintf("%v.png", len(r)),t))
	// Output:
	// <nil>
}

func ExampleTriggeringGroupImage() {
	const (
		steps = 800
		cycles = 2
	)

	o:=group.Steppers{
		NewTriggering(oscillator.NewDrivenAndDamped(steps/2,.85),50),	
	}
	o.Step(127) // initial drive, never exceeded, so this will always fit byte
	
	var r,a,t []byte
	for r=make([]byte,0,steps*cycles);len(r)!=cap(r);o.Step(0){
		r=append(r,byte(o.Real()+128))
		a=append(a,byte(o.Abs()+128))
		t=append(t,Iff(bool(o[0].(*Triggering).Triggered),byte(255),byte(0)).(byte))
	}
	fmt.Println(oscillator.Encode(fmt.Sprintf("g%v.png", len(r)),t))
	// Output:
	// <nil>
}


func Iff(b bool, t,f interface{}) interface{}{
	if b {
		return t
	}
	return f
}
