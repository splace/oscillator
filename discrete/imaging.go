package oscillator

import "os"
import "fmt"

import "github.com/splace/signals"

import "image"
import "image/draw"
import "image/color"
import "image/png"

const subFolder = "images"

var blueOverGrey = Appearance{
	color.NRGBA{0, 0, 128, 255},
	color.Gray{128},
}

func Encode(name string, bs ...[]byte) error {
	os.Chdir(subFolder) // TODO  not thread safe
	out, err := os.Create(name)
	os.Chdir("")
	if err != nil {
		return err
	}
	if len(bs) < 1 {
		return fmt.Errorf("Nothing to encode.")
	}
	return png.Encode(out, Image8bit(blueOverGrey, bs[0], bs[1:]...))
}

func Image8bit(s Appearance, bs0 []byte, bss ...[]byte) image.Image {
	c := &composable{image.NewRGBA(image.Rect(0, -128, len(bs0), 128))}
	c.drawOver(
		signals.RGBAImage{
			signals.NewDepiction(
				signals.NewPCM8bit(uint32(len(bs0)), bs0),
				len(bs0),
				256,
				// colours reversed, effectively making origin appear at the bottom.
				s.Background,
				s.Color,
			),
		},
	)
	// draw others with transparent background
	for _, bs := range bss {
		c.draw(
			signals.RGBAImage{
				signals.NewDepiction(
					signals.NewPCM8bit(uint32(len(bs)), bs[:len(bs0)]),
					len(bs0),
					256,
					s.Color,
					//s.Background,
					color.Transparent,
				),
			},
		)
	}
	// draw centre line
	c.drawOverAt(image.Rectangle{Max: image.Pt(len(bs0), 2)}, image.Pt(0, 128))
	return c
}

type Appearance struct {
	color.Color
	Background color.Color
}

// optimisation for block colours
// find the blended colour of this style, can then use, the faster, draw instead of drawOver.
func (a Appearance) Blend(d draw.Drawer) color.Color {
	r := image.NewNRGBA(image.Rectangle{Max: image.Pt(1, 1)})
	d.Draw(r, r.Bounds(), image.NewUniform(a.Background), image.ZP)
	d.Draw(r, r.Bounds(), image.NewUniform(a.Color), image.ZP)
	return r.At(0, 0)
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
