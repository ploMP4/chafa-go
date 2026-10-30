package chafa

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"os"
	"time"
)

// Animation contains fully composited frames with a shared size.
type Animation struct {
	Frames        []AnimationFrame
	Width, Height int32

	// LoopCount is the number of repetitions after the first playback.
	// Zero means repeat forever, and -1 means play only once.
	LoopCount int
}

// AnimationFrame contains tightly packed, unassociated RGBA8 pixels suitable
// for [CanvasDrawAllPixels] with [CHAFA_PIXEL_RGBA8_UNASSOCIATED]. The rowstride
// is the animation's Width multiplied by four. Each frame owns its pixel data.
type AnimationFrame struct {
	Pixels []uint8
	Delay  time.Duration
}

// LoadGIF loads every frame of a GIF into memory, preserving its delays and loop
// count. Partial frames and transparency are composited on the logical screen,
// applying the GIF's disposal methods, so every returned frame can be drawn
// independently. A zero frame delay is preserved; the caller chooses how to
// schedule such frames.
func LoadGIF(path string) (*Animation, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoded, err := gif.DecodeAll(file)
	if err != nil {
		return nil, err
	}
	if decoded.Config.Width <= 0 || decoded.Config.Height <= 0 {
		return nil, fmt.Errorf("gif: invalid logical screen size %dx%d", decoded.Config.Width, decoded.Config.Height)
	}

	animation := &Animation{
		Frames:    make([]AnimationFrame, len(decoded.Image)),
		Width:     int32(decoded.Config.Width),
		Height:    int32(decoded.Config.Height),
		LoopCount: decoded.LoopCount,
	}
	bounds := image.Rect(0, 0, decoded.Config.Width, decoded.Config.Height)
	screen := image.NewNRGBA(bounds)
	draw.Draw(screen, bounds, image.NewUniform(gifBackground(decoded, decoded.Image[0])), image.Point{}, draw.Src)

	var previous []uint8
	for i, frame := range decoded.Image {
		// RestorePrevious needs the screen as it was before drawing this frame.
		if decoded.Disposal[i] == gif.DisposalPrevious {
			previous = append(previous[:0], screen.Pix...)
		}
		draw.Draw(screen, frame.Bounds(), frame, frame.Bounds().Min, draw.Over)
		animation.Frames[i] = AnimationFrame{
			Pixels: append([]uint8(nil), screen.Pix...),
			Delay:  time.Duration(decoded.Delay[i]) * 10 * time.Millisecond,
		}

		// Disposal applies after a frame has been displayed, before the next.
		switch decoded.Disposal[i] {
		case gif.DisposalBackground:
			draw.Draw(screen, frame.Bounds(), image.NewUniform(gifBackground(decoded, frame)), image.Point{}, draw.Src)
		case gif.DisposalPrevious:
			copy(screen.Pix, previous)
		}
	}
	return animation, nil
}

func gifBackground(decoded *gif.GIF, frame *image.Paletted) color.Color {
	// Transparent GIFs use a transparent screen, including when disposing to
	// the background. The global palette itself does not carry transparency.
	for _, c := range frame.Palette {
		if _, _, _, alpha := c.RGBA(); alpha == 0 {
			return color.Transparent
		}
	}
	if palette, ok := decoded.Config.ColorModel.(color.Palette); ok && int(decoded.BackgroundIndex) < len(palette) {
		return palette[decoded.BackgroundIndex]
	}
	return color.Transparent
}
