package chafa

import (
	"bytes"
	"image"
	"image/color"
	"image/gif"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeTestGIF(t *testing.T, animation *gif.GIF) string {
	t.Helper()
	var data bytes.Buffer
	if err := gif.EncodeAll(&data, animation); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "animation.gif")
	if err := os.WriteFile(path, data.Bytes(), 0600); err != nil {
		t.Fatal(err)
	}
	return path
}

func palettedFrame(bounds image.Rectangle, palette color.Palette, pixels ...uint8) *image.Paletted {
	frame := image.NewPaletted(bounds, palette)
	copy(frame.Pix, pixels)
	return frame
}

func TestLoadGIFCompositing(t *testing.T) {
	transparent := color.NRGBA{}
	red := color.NRGBA{R: 255, A: 255}
	blue := color.NRGBA{B: 255, A: 255}
	green := color.NRGBA{G: 255, A: 255}
	palette := color.Palette{transparent, red, blue, green}
	path := writeTestGIF(t, &gif.GIF{
		Image: []*image.Paletted{
			palettedFrame(image.Rect(0, 0, 2, 1), palette, 1, 1),
			palettedFrame(image.Rect(1, 0, 3, 1), palette, 0, 2),
			palettedFrame(image.Rect(1, 0, 2, 1), palette, 2),
			palettedFrame(image.Rect(2, 0, 3, 1), palette, 3),
			palettedFrame(image.Rect(0, 0, 3, 1), palette, 0, 0, 0),
		},
		Delay:     []int{0, 1, 7, 2, 65535},
		Disposal:  []byte{gif.DisposalNone, gif.DisposalPrevious, gif.DisposalBackground, gif.DisposalNone, 0},
		LoopCount: 2,
		Config:    image.Config{ColorModel: palette, Width: 3, Height: 1},
	})
	animation, err := LoadGIF(path)
	if err != nil {
		t.Fatal(err)
	}
	if animation.Width != 3 || animation.Height != 1 || animation.LoopCount != 2 || len(animation.Frames) != 5 {
		t.Fatalf("unexpected animation metadata: %+v", animation)
	}
	want := [][]color.NRGBA{
		{red, red, transparent},
		{red, red, blue},
		{red, blue, transparent},
		{red, transparent, green},
		{red, transparent, green},
	}
	delays := []time.Duration{0, 10 * time.Millisecond, 70 * time.Millisecond, 20 * time.Millisecond, 655350 * time.Millisecond}
	for i, frame := range animation.Frames {
		if frame.Delay != delays[i] {
			t.Errorf("frame %d delay = %v, want %v", i, frame.Delay, delays[i])
		}
		if len(frame.Pixels) != 12 {
			t.Fatalf("frame %d has %d bytes, want 12", i, len(frame.Pixels))
		}
		for x, expected := range want[i] {
			got := color.NRGBA{R: frame.Pixels[x*4], G: frame.Pixels[x*4+1], B: frame.Pixels[x*4+2], A: frame.Pixels[x*4+3]}
			if got != expected {
				t.Errorf("frame %d pixel %d = %v, want %v", i, x, got, expected)
			}
		}
	}
	animation.Frames[0].Pixels[0] = 0
	if animation.Frames[1].Pixels[0] != 255 {
		t.Error("frames share pixel storage")
	}
}

func TestLoadGIFBackground(t *testing.T) {
	yellow := color.NRGBA{R: 255, G: 255, A: 255}
	red := color.NRGBA{R: 255, A: 255}
	blue := color.NRGBA{B: 255, A: 255}
	palette := color.Palette{color.Black, yellow, red, blue}
	path := writeTestGIF(t, &gif.GIF{
		Image: []*image.Paletted{
			palettedFrame(image.Rect(1, 0, 2, 1), palette, 2),
			palettedFrame(image.Rect(2, 0, 3, 1), palette, 3),
		},
		Delay:           []int{1, 2},
		Disposal:        []byte{gif.DisposalBackground, gif.DisposalNone},
		LoopCount:       -1,
		BackgroundIndex: 1,
		Config:          image.Config{ColorModel: palette, Width: 3, Height: 1},
	})
	animation, err := LoadGIF(path)
	if err != nil {
		t.Fatal(err)
	}
	for i, expected := range [][]color.NRGBA{{yellow, red, yellow}, {yellow, yellow, blue}} {
		frame := animation.Frames[i]
		img := &image.NRGBA{Pix: frame.Pixels, Stride: 12, Rect: image.Rect(0, 0, 3, 1)}
		for x, want := range expected {
			if got := img.NRGBAAt(x, 0); got != want {
				t.Errorf("frame %d pixel %d = %v, want %v", i, x, got, want)
			}
		}
	}
}

func TestLoadGIFLoopCount(t *testing.T) {
	palette := color.Palette{color.Black, color.White}
	for _, count := range []int{-1, 0, 3} {
		path := writeTestGIF(t, &gif.GIF{
			Image: []*image.Paletted{
				palettedFrame(image.Rect(0, 0, 1, 1), palette, 0),
				palettedFrame(image.Rect(0, 0, 1, 1), palette, 1),
			},
			Delay: []int{1, 1}, LoopCount: count,
		})
		animation, err := LoadGIF(path)
		if err != nil {
			t.Fatal(err)
		}
		if animation.LoopCount != count {
			t.Errorf("loop count = %d, want %d", animation.LoopCount, count)
		}
	}
}

func TestLoadGIFErrors(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing.gif")
	if animation, err := LoadGIF(missing); err == nil || animation != nil {
		t.Fatalf("missing file: got %v, %v", animation, err)
	}
	path := filepath.Join(t.TempDir(), "invalid.gif")
	if err := os.WriteFile(path, []byte("GIF89a"), 0600); err != nil {
		t.Fatal(err)
	}
	if animation, err := LoadGIF(path); err == nil || animation != nil {
		t.Fatalf("truncated file: got %v, %v", animation, err)
	}
}

func TestCanvasPrintString(t *testing.T) {
	config := CanvasConfigNew()
	defer CanvasConfigUnref(config)
	CanvasConfigSetGeometry(config, 2, 2)
	CanvasConfigSetCanvasMode(config, CHAFA_CANVAS_MODE_TRUECOLOR)
	CanvasConfigSetPixelMode(config, CHAFA_PIXEL_MODE_SYMBOLS)
	canvas := CanvasNew(config)
	defer CanvasUnRef(canvas)

	CanvasDrawAllPixels(canvas, CHAFA_PIXEL_RGBA8_UNASSOCIATED, []uint8{255, 0, 0, 255}, 1, 1, 4)
	red := CanvasPrintString(canvas, nil)
	if red == "" || bytes.Count([]byte(red), []byte("\n")) != 1 {
		t.Fatalf("expected two nonempty rows, got %q", red)
	}
	for range 100 {
		CanvasDrawAllPixels(canvas, CHAFA_PIXEL_RGBA8_UNASSOCIATED, []uint8{0, 0, 255, 255}, 1, 1, 4)
		if blue := CanvasPrintString(canvas, nil); blue == "" || blue == red {
			t.Fatalf("expected a distinct blue frame, got %q", blue)
		}
	}
	CanvasDrawAllPixels(canvas, CHAFA_PIXEL_RGBA8_UNASSOCIATED, []uint8{255, 0, 0, 255}, 1, 1, 4)
	if got := CanvasPrintString(canvas, nil); got != red {
		t.Error("retained Go string changed after freeing or reusing native buffers")
	}
}
