package discrete

import "os"
import "fmt"

import "github.com/splace/signals"

import "image"
import "image/draw"
import "image/color"
import "image/png"

import "github.com/anthonynsimon/bild/transform"


const subFolder = "images"

var blueOverGrey = Appearance{
	color.NRGBA{0, 0, 128, 255},
	color.Gray{128},
}

func Encode(name string, bs ...[]uint8) error {
	if len(bs) < 1 {
		return fmt.Errorf("Nothing to encode.")
	}
	return encode(name,Image8bit(blueOverGrey, bs[0], bs[1:]...))
}

func encode(name string,i draw.Image) error{
	os.Chdir(subFolder) // TODO  not thread safe
	out, err := os.Create(name)
	os.Chdir("")
	if err != nil {
		return err
	}
	// draw centre line
	b:=i.Bounds()
	(&composable{i}).drawOverAt(image.Rectangle{Max: image.Pt(b.Min.X, 2)}, image.Pt(0,b.Max.Y/2))
	return png.Encode(out,i)
} 

func EncodeSize(name string, dx,dy int,comp func(s Appearance, bs0 []byte, bss ...[]byte) *composable, bs ...[]uint8) error {
	if len(bs) < 1 {
		return fmt.Errorf("Nothing to encode.")
	}
	if dx<2 || dy<2 || dx>100000 || dy >100000 {
		return fmt.Errorf("Size unreasonable (%dx%d)",dx,dy)
	}
	img:=transform.Resize(
		comp(blueOverGrey, bs[0], bs[1:]...),
		dx,
		dy,
		transform.Linear,
	)
	// draw a centre line
	(&composable{img}).drawOverAt(image.Rectangle{Max: image.Pt(dx, 1)}, image.Pt(0, dy/2))
	return encode(name, img)
}

// use signals.Depiction to make an image of a number of []int8 values
func Image8bit(s Appearance, bs0 []byte, bss ...[]byte) *composable {
	return compose(
		func(r uint32,d []byte) signals.LimitedSignal{
			return signals.LimitedSignal(signals.NewPCM8bit(r,d))
		},
		1<<8,
		s,
		bs0,
		bss,
	)
}

// use signals.Depiction to make an image of a number of []int8 values
func Image16bit(s Appearance, bs0 []byte, bss ...[]byte) *composable {
	return compose(
		func(r uint32,d []byte) signals.LimitedSignal{
			return signals.LimitedSignal(signals.NewPCM16bit(r,d))
		},
		1<<16, // this pixel height would require very large amounts of memory if we weren't using signals' precedural images and generally resizing down or encoding x-context.
		s,
		bs0,
		bss,
	)
}

func compose(lsigfn func(uint32,[]byte) signals.LimitedSignal,height int,s Appearance, bs0 []byte, bss [][]byte) *composable {
	c := &composable{image.NewRGBA(image.Rect(0, -128, len(bs0), 128))}
	c.drawOver(
		Convolute(
			signals.RGBAImage{
				signals.NewDepiction(
					lsigfn(uint32(len(bs0)), bs0),
					len(bs0),
					height,
					// colours reversed, effectively flipping vertically and making more conventional origin at bottom.
					s.Background,
					s.Color,
				),
			},
			[][]int{
				{-1,0,-1},
				{0,5,0},
				{-1,0,-1},
			},
		),
	)
	// draw others with transparent background
	for _, bs := range bss {
		c.draw(
			Convolute(
				signals.RGBAImage{
					signals.NewDepiction(
						lsigfn(uint32(len(bs)), bs[:len(bs0)]),
						len(bs0),
						height,
						s.Blend(draw.Over),
						//s.Background,
						color.Transparent,
					),
				},
				[][]int{
					{-1,0,-1},
					{0,5,0},
					{-1,0,-1},
				},
			),
		)
	}
	// draw centre line
	c.drawOverAt(image.Rectangle{Max: image.Pt(len(bs0), 1)}, image.Pt(0, 127))
	return c
}


//// use signals.Depiction to make an image of a number of []int8 values
//func Image16it(s Appearance, bs0 []int16, bss ...[]int16) *composable {
//	Image8bit(s Appearance, bs0 []uint8, bss ...[]uint8)
//}


type Appearance struct {
	color.Color
	Background color.Color
}

// optimisation for block colours
// find the blended colour of this style, can then use, the faster, draw instead of drawOver.
func (a Appearance) Blend(d draw.Drawer) color.NRGBA {
	r := image.NewNRGBA(image.Rectangle{Max: image.Pt(1, 1)})
	d.Draw(r, r.Bounds(), image.NewUniform(a.Background), image.ZP)
	d.Draw(r, r.Bounds(), image.NewUniform(a.Color), image.ZP)
	return r.NRGBAAt(0, 0)
}

// composable simplifies draw.Draw for incremental composition of images.
// in essence draw.Draw's parameter options are moved to separate methods, cutting parameter count and making function clearer.
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

func Resize(img image.Image, dx,dy uint) image.Image {
	return resized{img,float32(img.Bounds().Dx())/float32(dx),float32(img.Bounds().Dy())/float32(dy)}
}

// procedural resized intermediate image
type resized struct{
	image.Image
	xd,yd float32
}

func (r resized) At(x,y int) color.Color{
	return r.Image.At(int(float32(x)*r.xd),int(float32(y)*r.yd))
}

// procedural over sized intermediate image
type oversized struct{
	image.Image
	xd,yd int
}

func (o oversized) At(x,y int) color.Color{
	return o.Image.At(x*o.xd,y*o.yd)
}




// Convolution returns an new image.Image that is the convolution of the provided image.Image, with a number of provided matrices, these are applied to the images components colours depending on the number of matrices provided.
// TODO target image different size
func Convolute(img image.Image, matrices ...[][]int) *image.NRGBA {
	ir := img.Bounds()
	//xd, yd := ir.Dx(), ir.Dy()

	var r, g, b, a uint32

	var normp,normn uint32
	for y := range matrices[0] {
		for x := range matrices[0][y] {
			if v:=matrices[0][y][x];v>0{
				normp += uint32(v)
			}else{
				normn += uint32(-v)
			}
		}
	}

	if normn>normp {
		panic("Matrix unsupported, normalisation was negative.")
	}
	norm:=int32(normp-normn)
	if norm>>16 != 0 {
		panic("Matrix unsupported, normalisation factor overflows 16 bits (65536).")
	}
	conv := image.NewNRGBA(ir)
	var x, y, dx, dy, mx, my int
	var mf int
	for y = ir.Min.Y; y < ir.Max.Y; y++ {
		for x = ir.Min.X; x < ir.Max.X; x++ {
			halfHeight := len(matrices[0]) / 2
			var tr, tg, tb, ta int32
			//var trr, tgg, tbb uint32
			switch len(matrices){
			case 1:
				for dy, my = -halfHeight, 0; my < len(matrices[0]); dy, my = dy+1, my+1 {
					halfWidth := len(matrices[0][my]) / 2
					for dx, mx = -halfWidth, 0; mx < len(matrices[0][my]); dx, mx = dx+1, mx+1 {
						mf = matrices[0][my][mx]
						if mf==0 {continue}
						r, g, b, _ = img.At(x+dx, y+dy).RGBA()
						if mf < 0 {
							f := uint32(-mf)
							tr, tg, tb = tr-int32(r*f), tg-int32(g*f), tb-int32(b*f)
						} else {
							f := uint32(mf)
							tr, tg, tb = tr+int32(r*f), tg+int32(g*f), tb+int32(b*f)
						}
					}
				}
				ta=0xffff*norm
			case 2:
				for dy, my = -halfHeight, 0; my < len(matrices[0]); dy, my = dy+1, my+1 {
					halfWidth := len(matrices[0][my]) / 2
					for dx, mx = -halfWidth, 0; mx < len(matrices[0][my]); dx, mx = dx+1, mx+1 {
						mf = matrices[0][my][mx]
						if mf==0 {continue}
						r, g, b, _ = img.At(x+dx, y+dy).RGBA()
						if mf < 0 {
							f := uint32(-mf)
							tr, tg, tb = tr-int32(r*f), tg-int32(g*f), tb-int32(b*f)
						} else {
							f := uint32(mf)
							tr, tg, tb = tr+int32(r*f), tg+int32(g*f), tb+int32(b*f)
						}
					}
				}
				for dy, my = -halfHeight, 0; my < len(matrices[1]); dy, my = dy+1, my+1 {
					halfWidth := len(matrices[1][my]) / 2
					for dx, mx = -halfWidth, 0; mx < len(matrices[1][my]); dx, mx = dx+1, mx+1 {
						mf = matrices[1][my][mx]
						if mf==0 {continue}
						_,_,_, a = img.At(x+dx, y+dy).RGBA()
						if mf < 0 {
							f := uint32(-mf)
							ta = ta-int32(a*f)
						} else {
							f := uint32(mf)
							ta = ta+int32(a*f)
						}
					}
				}
			case 3:
				for dy, my = -halfHeight, 0; my < len(matrices[0]); dy, my = dy+1, my+1 {
					halfWidth := len(matrices[0][my]) / 2
					for dx, mx = -halfWidth, 0; mx < len(matrices[0][my]); dx, mx = dx+1, mx+1 {
						mf = matrices[0][my][mx]
						if mf==0 {continue}
						r, _, _, _ = img.At(x+dx, y+dy).RGBA()
						if mf < 0 {
							f := uint32(-mf)
							tr = tr-int32(r*f)
						} else {
							f := uint32(mf)
							tr = tr+int32(r*f)
						}
					}
				}
				for dy, my = -halfHeight, 0; my < len(matrices[1]); dy, my = dy+1, my+1 {
					halfWidth := len(matrices[1][my]) / 2
					for dx, mx = -halfWidth, 0; mx < len(matrices[1][my]); dx, mx = dx+1, mx+1 {
						mf = matrices[1][my][mx]
						if mf==0 {continue}
						_, g, _, _ = img.At(x+dx, y+dy).RGBA()
						if mf < 0 {
							f := uint32(-mf)
							tg = tg-int32(g*f)
						} else {
							f := uint32(mf)
							tg = tg+int32(g*f)
						}
					}
				}
				for dy, my = -halfHeight, 0; my < len(matrices[2]); dy, my = dy+1, my+1 {
					halfWidth := len(matrices[2][my]) / 2
					for dx, mx = -halfWidth, 0; mx < len(matrices[2][my]); dx, mx = dx+1, mx+1 {
						mf = matrices[2][my][mx]
						if mf==0 {continue}
						_, _, b, _ = img.At(x+dx, y+dy).RGBA()
						if mf < 0 {
							f := uint32(-mf)
							tb = tb-int32(b*f)
						} else {
							f := uint32(mf)
							tb = tb+int32(b*f)
						}
					}
				}
				ta=0xffff*norm
			case 4:
				for dy, my = -halfHeight, 0; my < len(matrices[0]); dy, my = dy+1, my+1 {
					halfWidth := len(matrices[0][my]) / 2
					for dx, mx = -halfWidth, 0; mx < len(matrices[0][my]); dx, mx = dx+1, mx+1 {
						mf = matrices[0][my][mx]
						if mf==0 {continue}
						r, g, b, a = img.At(x+dx, y+dy).RGBA()
						if mf < 0 {
							f := uint32(-mf)
							tr, tg, tb, ta = tr-int32(r*f), tg-int32(g*f), tb-int32(b*f), ta-int32(a*f)
						} else {
							f := uint32(mf)
							tr, tg, tb, ta = tr+int32(r*f), tg+int32(g*f), tb+int32(b*f), ta+int32(a*f)
						}
					}
				}
				for dy, my = -halfHeight, 0; my < len(matrices[1]); dy, my = dy+1, my+1 {
					halfWidth := len(matrices[1][my]) / 2
					for dx, mx = -halfWidth, 0; mx < len(matrices[1][my]); dx, mx = dx+1, mx+1 {
						mf = matrices[1][my][mx]
						if mf==0 {continue}
						_, g, _, _ = img.At(x+dx, y+dy).RGBA()
						if mf < 0 {
							f := uint32(-mf)
							tg = tg-int32(g*f)
						} else {
							f := uint32(mf)
							tg = tg+int32(g*f)
						}
					}
				}
				for dy, my = -halfHeight, 0; my < len(matrices[2]); dy, my = dy+1, my+1 {
					halfWidth := len(matrices[2][my]) / 2
					for dx, mx = -halfWidth, 0; mx < len(matrices[2][my]); dx, mx = dx+1, mx+1 {
						mf = matrices[2][my][mx]
						if mf==0 {continue}
						_, _, b, _ = img.At(x+dx, y+dy).RGBA()
						if mf < 0 {
							f := uint32(-mf)
							tb = tb-int32(b*f)
						} else {
							f := uint32(mf)
							tb = tb+int32(b*f)
						}
					}
				}
				for dy, my = -halfHeight, 0; my < len(matrices[3]); dy, my = dy+1, my+1 {
					halfWidth := len(matrices[3][my]) / 2
					for dx, mx = -halfWidth, 0; mx < len(matrices[3][my]); dx, mx = dx+1, mx+1 {
						mf = matrices[3][my][mx]
						if mf==0 {continue}
						_,_,_, a = img.At(x+dx, y+dy).RGBA()
						if mf < 0 {
							f := uint32(-mf)
							ta = ta-int32(a*f)
						} else {
							f := uint32(mf)
							ta = ta+int32(a*f)
						}
					}
				}
			}
			if norm<2{
				conv.Set(x, y, color.NRGBA{
					clamp(tr),
					clamp(tg),
					clamp(tb),
					clamp(ta),
				})
			}else{
				conv.Set(x, y, color.NRGBA{
					clamp(tr / norm),
					clamp(tg / norm),
					clamp(tb / norm),
					clamp(ta / norm),
				})
			}
		}
	}
	return conv
}

func clamp(i int32) uint8{
	if i<0 {return 0}
	if v:=i>>8;v>>8>0{
		return 0xff
	}else{
		return uint8(v)
	}
}
