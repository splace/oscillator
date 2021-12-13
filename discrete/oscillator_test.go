package oscillator


import "fmt"
import "os"
//import "time"

import "github.com/splace/signals"

import "image"
import "image/color"
import "image/png"

func ExampleOscillatorImage() {
	const (
		steps = 2000
		cycles = 5
	)

	o:=NewDrivenAndDamped(steps,.5) // damping by half per cycle
	o.Step(127) // initial drive never exceeded, this will always fit byte
	
	var r,a []byte
	for r=make([]byte,0,steps*cycles);len(r)!=cap(r);o.Step(0){
		r=append(r,byte(o.Real()+128))
		a=append(a,byte(o.Abs()+128))
	}

	out, err := os.Create(fmt.Sprintf("%v.png", steps*cycles))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Composing image width:%dpx\n ",steps*cycles)
	var style = Appearance{color.NRGBA{0, 0, 128, 255},color.Gray{128}}
	m := &composable{image.NewRGBA(image.Rect(0, -128, steps*cycles, 128))}
	m.drawOver(
		signals.RGBAImage{
			signals.NewDepiction(
				signals.NewPCM8bit(uint32(steps*cycles), r),
				steps*cycles,
				256,
				// colours swopped to compensate for images origin being at the top.
				style.Background,
				style.Color,
			),
		},
	)
	m.drawOver(
		signals.RGBAImage{
			signals.NewDepiction(
				signals.NewPCM8bit(uint32(steps*cycles), a),
				steps*cycles,
				256,
				// colours swapped to compensate for images origin being at the top.
				color.Transparent,
				style.Background,
				//style.Blend(draw.Over),
			),
		},
	)
	// y-axis line
	m.drawOverAt(image.Rectangle{Max:image.Pt(steps*cycles, 2)},image.Pt(0,128))
	fmt.Printf("Saving image: %s\n ",out.Name())
	png.Encode(out, m)
	// Output:
	//
}
