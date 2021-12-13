package oscillator_test

import "fmt"

import osc "github.com/splace/oscillator/discrete"

import "github.com/anthonynsimon/bild/transform"

import "os"
import "image/color"
import "image/png"

func ExampleSteppersImage() {
	const (
		steps  = 12800
		cycles = 2
	)
	// build oscillator group
	// frequency and amplitude set to square wave Fourier transform.
	var oss osc.Steppers
	for n := 1.0; n < 40; n += 2 {
		o := osc.NewDamped(steps/n, 1/n) // 1/n damping makes all oscillators damp at the same rate.
		o.Step(100 / n)                  // 100 won't exceed byte (+127 - -127) when all summed
		oss = append(oss, o)
	}

	// step through and store the real part a in byte
	var r []byte
	for r = make([]byte, 0, steps*cycles); len(r) != cap(r); oss.Step(0) {
		r = append(r, byte(real(oss.StateGetter())+128))
	}

	os.Chdir("images") // TODO  simple but not really thread safe
	out, err := os.Create(fmt.Sprintf("Steppers(%v).png", len(r)))
	os.Chdir("")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(
		// save as squared-up (2048 x 1024 px) signal strip (steps x 256 px)
		png.Encode(
			out,
			transform.Resize(
				osc.Image8bit(
					osc.Appearance{
						color.NRGBA{0, 0, 128, 255},
						color.Gray{128},
					},
					r,
				),
				2048,                       // len(r)/2,
				1024,                       // len(r)/2,
				transform.ResampleFilter{}, //  transform.Linear,
			),
		),
	)
	// Output:
	// <nil>
}
