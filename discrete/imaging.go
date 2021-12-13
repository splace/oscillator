package oscillator


import "os"
//import "log"

import "github.com/splace/signals"

import "image"
import "image/draw"
import "image/color"
import "image/png"


func Encode(name string,bs0 []byte,bss ...[]byte) error{
	out, err := os.Create(name)
	if err!=nil{
		return err
	}
	m := &composable{image.NewRGBA(image.Rect(0, -128, len(bs0), 128))}
	m.drawOver(
		signals.RGBAImage{
			signals.NewDepiction(
				signals.NewPCM8bit(uint32(len(bs0)), bs0),
				len(bs0),
				256,
				// colours swapped to compensate for images origin being at the top.
				style.Background,
				style.Color,
			),
		},
	)
	for _,bs:= range bss{
		m.draw(
			signals.RGBAImage{
				signals.NewDepiction(
					signals.NewPCM8bit(uint32(len(bs)), bs),
					len(bs0),
					256,
					style.Color,
					//style.Background,
					color.Transparent,
				),
			},
		)
	}
	m.drawOverAt(image.Rectangle{Max:image.Pt(len(bs0), 2)},image.Pt(0,128))
//	fmt.Printf("Saving image: %s\n ",out.Name())
	return png.Encode(out, m)
}


var style = Appearance{
	color.NRGBA{0, 0, 128, 255},
	color.Gray{128},
}

type Appearance struct{
	color.Color
	Background color.Color
}

// optimisation for block colours
// find the blended colour of this style, can then use, the faster, draw instead of drawOver.
func (a Appearance) Blend(d draw.Drawer) color.Color{
	r:=image.NewNRGBA(image.Rectangle{Max:image.Pt(1, 1)})
//	log.Println(r.At(0,0))
	d.Draw(r,r.Bounds(),image.NewUniform(a.Background),image.ZP)
//	log.Println(r.At(0,0))
	d.Draw(r,r.Bounds(),image.NewUniform(a.Color),image.ZP)
//	log.Println(r.At(0,0))
	return r.At(0,0)
}


// composable simplifies draw.Draw for incremental composition of images.
// basically draw.Draw's parameter options are moved to separate methods, cutting parameter count and making function clear.
type composable struct {
	draw.Image
}

func (i *composable) draw(isrc image.Image) {
	draw.Draw(i, i.Bounds(), isrc, isrc.Bounds().Min, draw.Src)
}

func (i *composable) drawAt(isrc image.Image, pt image.Point) {
	draw.Draw(i, i.Bounds(), isrc, image.ZP.Sub(pt), draw.Src)
}

func (i *composable) drawOffset(isrc image.Image, pt image.Point) {
	draw.Draw(i, i.Bounds(), isrc, isrc.Bounds().Min.Sub(pt), draw.Src)
}

func (i *composable) drawOver(isrc image.Image) {
	draw.Draw(i, i.Bounds(), isrc, isrc.Bounds().Min, draw.Over)
}

func (i *composable) drawOverAt(isrc image.Image, pt image.Point) {
	draw.Draw(i, i.Bounds(), isrc, image.ZP.Sub(pt), draw.Over)
}

func (i *composable) drawOverOffset(isrc image.Image, pt image.Point) {
	draw.Draw(i, i.Bounds(), isrc, isrc.Bounds().Min.Sub(pt), draw.Over)
} 

