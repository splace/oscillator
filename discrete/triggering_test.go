package oscillator


import "fmt"

func ExampleTriggeringImage() {
	const (
		steps = 8000
		cycles = 2
	)

	o:=NewTriggering(NewDamped(steps/2,.85),50)
	o.Step(127) // initial drive, never exceeded, so this will always fit byte
	
	var a,p,c,t []byte
	for a=make([]byte,0,steps*cycles);len(a)!=cap(a);o.Step(0){
		a=append(a,byte(o.Amplitude()+128))
		p=append(p,byte(o.Phase()*angleScaler*128))
		c=append(c,Iff(o.StepChange(0),byte(255),byte(0)).(byte))
		t=append(t,Iff(bool(o.Triggered),byte(255),byte(0)).(byte))
	}
	fmt.Println(Encode(fmt.Sprintf("TriggerDamping(%v).png", len(a)),t))
	// Output:
	// <nil>
}



func Iff(b bool, t,f interface{}) interface{}{
	if b {
		return t
	}
	return f
}
