package oscillator


import "fmt"

func ExampleDampingImage() {
	const (
		steps = 2000
		cycles = 5
	)

	o:=NewDamped(steps,.5) // damping by half per cycle
	o.Step(127) // initial drive, never exceeded so this will always fit byte
	
	// build byte arrays for parameters of interest
	var a,p,r []byte
	for a=make([]byte,0,steps*cycles);len(a)!=cap(a);o.Step(0){
		a=append(a,byte(o.Amplitude()+128))
		p=append(p,byte(o.Phase()*angleScaler*128))
		r=append(r,byte(real(o.GetState())+128))
	}
	
	fmt.Println(Encode(fmt.Sprintf("Damping(%v).png", len(r)),a))
	// Output:
	// <nil>
}

func ExampleDamping2Image() {
	const (
		steps = 1000
		cycles = 15
	)

	d:=NewDamped(steps,0) // driver ossilator, no damping
	d.Step(127.0/steps)

	o:=NewDamped(steps*2,.5) // damping by half per cycle
	
	var a,p,r []byte
	for a=make([]byte,0,steps*cycles);len(a)!=cap(a);o.Step(real(d.GetState())){
		a=append(a,byte(o.Amplitude()+128))
		p=append(p,byte(o.Phase()*angleScaler*128))
		r=append(r,byte(real(o.GetState())+128))
		d.Step(0)
	}
	
	fmt.Println(Encode(fmt.Sprintf("Damping2(%v).png", len(r)),r))
	// Output:
	// <nil>
}
