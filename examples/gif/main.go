// Run from the repository root:
//
//	go run ./examples/gif path/to/animation.gif
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ploMP4/chafa-go"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	width := flag.Int("width", 80, "maximum width in terminal cells")
	height := flag.Int("height", 24, "maximum height in terminal cells")
	once := flag.Bool("once", false, "play once instead of using the GIF's loop count")
	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), "Usage: gif [flags] path/to/animation.gif")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		return fmt.Errorf("provide one GIF file")
	}
	if *width <= 0 || *height <= 0 || *width > math.MaxInt32 || *height > math.MaxInt32 {
		return fmt.Errorf("width and height must be positive 32-bit integers")
	}
	animation, err := chafa.LoadGIF(flag.Arg(0))
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return play(ctx, os.Stdout, animation, int32(*width), int32(*height), *once)
}

func play(
	ctx context.Context,
	output io.Writer,
	animation *chafa.Animation,
	width, height int32,
	once bool,
) error {
	termInfo := chafa.TermDbDetect(chafa.TermDbGetDefault(), os.Environ())
	defer chafa.TermInfoUnref(termInfo)

	config := chafa.CanvasConfigNew()
	defer chafa.CanvasConfigUnref(config)
	// Symbol output works in ordinary ANSI terminals without a graphics protocol.
	chafa.CanvasConfigSetPixelMode(config, chafa.CHAFA_PIXEL_MODE_SYMBOLS)
	chafa.CanvasConfigSetCanvasMode(config, chafa.TermInfoGetBestCanvasMode(termInfo))
	chafa.CanvasConfigSetCellGeometry(config, 8, 16)
	chafa.CalcCanvasGeometry(animation.Width, animation.Height, &width, &height, 0.5, false, false)
	chafa.CanvasConfigSetGeometry(config, width, height)
	canvas := chafa.CanvasNew(config)
	defer chafa.CanvasUnRef(canvas)

	// Ask Chafa for the detected terminal's commands. These sequences take no
	// arguments, so they can be used directly without formatting.
	enterAlternateScreen := chafa.TermInfoGetSeq(termInfo, chafa.CHAFA_TERM_SEQ_ENABLE_ALT_SCREEN)
	leaveAlternateScreen := chafa.TermInfoGetSeq(termInfo, chafa.CHAFA_TERM_SEQ_DISABLE_ALT_SCREEN)
	hideCursor := chafa.TermInfoGetSeq(termInfo, chafa.CHAFA_TERM_SEQ_DISABLE_CURSOR)
	showCursor := chafa.TermInfoGetSeq(termInfo, chafa.CHAFA_TERM_SEQ_ENABLE_CURSOR)
	clearScreen := chafa.TermInfoGetSeq(termInfo, chafa.CHAFA_TERM_SEQ_CLEAR)
	cursorHome := chafa.TermInfoGetSeq(termInfo, chafa.CHAFA_TERM_SEQ_CURSOR_TO_TOP_LEFT)
	resetAttributes := chafa.TermInfoGetSeq(termInfo, chafa.CHAFA_TERM_SEQ_RESET_ATTRIBUTES)

	// Use the alternate screen and hide the cursor; restore both on every exit,
	// including Ctrl+C and output errors.
	defer io.WriteString(output, resetAttributes+showCursor+leaveAlternateScreen)
	if _, err := io.WriteString(output, enterAlternateScreen+hideCursor+clearScreen); err != nil {
		return err
	}
	for repetition := 0; ; repetition++ {
		for _, frame := range animation.Frames {
			if ctx.Err() != nil {
				return nil
			}
			started := time.Now()
			chafa.CanvasDrawAllPixels(canvas, chafa.CHAFA_PIXEL_RGBA8_UNASSOCIATED,
				frame.Pixels, animation.Width, animation.Height, animation.Width*4)
			if _, err := fmt.Fprint(
				output,
				cursorHome,
				chafa.CanvasPrintString(canvas, termInfo),
			); err != nil {
				return err
			}

			delay := frame.Delay
			if delay <= 0 {
				delay = 50 * time.Millisecond // Match Chafa's default for untimed GIFs.
			}
			// Include rendering time in the frame budget so it doesn't slow playback.
			timer := time.NewTimer(delay - time.Since(started))
			select {
			case <-ctx.Done():
				timer.Stop()
				return nil
			case <-timer.C:
			}
		}
		if once || len(animation.Frames) == 1 || animation.LoopCount < 0 ||
			(animation.LoopCount > 0 && repetition >= animation.LoopCount) {
			return nil
		}
	}
}
