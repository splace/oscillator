package oscillator

import "fmt"

func ExampleImagingLines(){
	const (
		steps  = 8000
		cycles = 2
	)
	// build oscillator group
	// frequency and amplitude set to square wave Fourier transform.
	var oss Steppers
	for n := 1.0; n < 40; n += 2 {
		o := NewDamped(steps/n, 0) 
		o.Step(100 / n)  // 100 won't exceed uint8 (+127 - -127) when all summed
		oss = append(oss, o)
	}

	as := make([][]uint8,len(oss))
	for i := range as{
		as[i]=make([]uint8, steps*cycles)
	}
	for i := range as[0]{
		for j:= range as{
			as[j][i]=uint8(real(oss[j].StateGetter())+127)
		}	
		oss.Step(0) 
	}
	
//	fmt.Println(Encode(fmt.Sprintf("Steppers(%v).png", len(as)),as...))
//	fmt.Println(Encode(fmt.Sprintf("Steppers(%v).png",len(as)), as[0]))
	fmt.Println(EncodeSize(fmt.Sprintf("Steppers(%v).png",len(as)), 2048,1024,as[0]))
	// Output:
	// <nil>

}