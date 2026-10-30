package chafa

import (
	"strings"
	"unsafe"
)

var (
	// Creates a new canvas with the specified configuration. The canvas makes
	// a private copy of the configuration, so it will not be affected by subsequent changes.
	CanvasNew func(config *CanvasConfig) *Canvas

	// Creates a new canvas configured similarly to orig.
	CanvasNewSimilar func(orig *Canvas) *Canvas

	// Adds a reference to canvas.
	CanvasRef func(canvas *Canvas)

	// Removes a reference from canvas. When remaining references drops to zero,
	// the canvas is freed and can no longer be used.
	CanvasUnRef func(canvas *Canvas)

	// Returns a pointer to the configuration belonging to canvas.
	// This can be inspected using the [CanvasConfig] getter functions, but not changed.
	CanvasPeekConfig func(canvas *Canvas) *CanvasConfig

	// Places placement on canvas, replacing the latter's content. The placement will cover the entire canvas.
	//
	// The canvas will keep a reference to the placement until it is replaced or the canvas itself is freed.
	CanvasSetPlacement func(canvas *Canvas, placement *Placement)

	// Replaces pixel data of canvas with a copy of that found at src_pixels ,
	// which must be in one of the formats supported by [PixelType].
	CanvasDrawAllPixels func(
		canvas *Canvas,
		srcPixelType PixelType,
		srcPixels []uint8,
		srcWidth int32,
		srcHeight int32,
		srcRowstride int32,
	)

	// Builds a UTF-8 string of terminal control sequences and symbols representing the canvas' current contents.
	// This can be printed to a terminal. The exact choice of escape sequences and symbols, dimensions, etc.
	// is determined by the configuration assigned to canvas on its creation.
	//
	// All output lines except for the last one will end in a newline.
	CanvasPrint func(canvas *Canvas, termInfo *TermInfo) *GString

	//Builds an array of UTF-8 strings made up of terminal control sequences and symbols
	// representing the canvas' current contents. These can be printed to a terminal.
	// The exact choice of escape sequences and symbols, dimensions, etc. is determined
	// by the configuration assigned to canvas on its creation.
	//
	// The array will be NULL-terminated. The element count does not include the terminator.
	//
	// When the canvas' pixel mode is [CHAFA_PIXEL_MODE_SYMBOLS], each element will hold
	// the contents of exactly one symbol row. There will be no row separators,
	// newlines or control sequences to reposition the cursor between rows.
	// Row positioning is left to the caller.
	//
	// In other pixel modes, there may be one or more strings, but the splitting
	// criteria should not be relied on. They must be printed in sequence, exactly as they appear.
	CanvasPrintRows func(canvas *Canvas, termInfo *TermInfo, arrayOut *unsafe.Pointer, lenOut *int32)

	// Builds an array of UTF-8 strings made up of terminal control sequences and symbols
	// representing the canvas' current contents. These can be printed to a terminal.
	// The exact choice of escape sequences and symbols, dimensions, etc. is determined by
	// the configuration assigned to canvas on its creation.
	//
	// When the canvas' pixel mode is [CHAFA_PIXEL_MODE_SYMBOLS], each element
	// will hold the contents of exactly one symbol row. There will be no row separators,
	// newlines or control sequences to reposition the cursor between rows.
	// Row positioning is left to the caller.
	//
	// In other pixel modes, there may be one or more strings,
	// but the splitting criteria should not be relied on.
	// They must be printed in sequence, exactly as they appear.
	CanvasPrintRowsStrv func(canvas *Canvas, termInfo *TermInfo) unsafe.Pointer

	// Frees a NULL-terminated array returned by [CanvasPrintRows], including
	// each GString in the array. Passing nil is allowed.
	FreeGStringArray func(array unsafe.Pointer)

	// Returns the character at cell (x, y). The coordinates are zero-indexed.
	// For double-width characters, the leftmost cell will contain the character
	// and the rightmost cell will contain 0.
	CanvasGetCharAt func(canvas *Canvas, x, y int32) rune

	// Sets the character at cell (x, y). The coordinates are zero-indexed.
	// For double-width characters, the leftmost cell must contain the character
	// and the cell to the right of it will automatically be set to 0.
	//
	// If the character is a nonprintable or zero-width, no change will be made.
	CanvasSetCharAt func(canvas *Canvas, x, y int32, c rune) int32

	// Gets the colors at cell (x, y). The coordinates are zero-indexed.
	// For double-width characters, both cells will contain the same colors.
	//
	// The colors will be -1 for transparency, packed 8bpc RGB otherwise, i.e. 0x00RRGGBB hex.
	//
	// If the canvas is in an indexed mode, palette lookups will be made for you.
	CanvasGetColorsAt func(canvas *Canvas, x, y int32, fgOut, bgOut *int32)

	// Sets the colors at cell (x, y). The coordinates are zero-indexed.
	// For double-width characters, both cells will be set to the same color.
	//
	// The colors must be -1 for transparency, packed 8bpc RGB otherwise, i.e. 0x00RRGGBB hex.
	//
	// If the canvas is in an indexed mode, palette lookups will be made for you.
	CanvasSetColorsAt func(canvas *Canvas, x, y, fg, bg int32)

	// Gets the colors at cell (x, y). The coordinates are zero-indexed.
	// For double-width characters, both cells will contain the same colors.
	//
	// The colors will be -1 for transparency, packed 8bpc RGB, i.e. 0x00RRGGBB hex
	// in truecolor mode, or the raw pen value (0-255) in indexed modes.
	//
	// It's the caller's responsibility to handle the color values correctly
	// according to the canvas mode (truecolor or indexed).
	CanvasGetRawColorsAt func(canvas *Canvas, x, y int32, fgOut, bgOut *int32)

	// Sets the colors at cell (x, y). The coordinates are zero-indexed.
	// For double-width characters, both cells will be set to the same color.
	//
	// The colors must be -1 for transparency, packed 8bpc RGB, i.e. 0x00RRGGBB hex
	// in truecolor mode, or the raw pen value (0-255) in indexed modes.
	//
	// It's the caller's responsibility to handle the color values correctly
	// according to the canvas mode (truecolor or indexed).
	CanvasSetRawColorsAt func(canvas *Canvas, x, y, fg, bg int32)
)

// CanvasPrintString renders canvas into a Go-owned string and frees the native
// output buffers. It is useful for repeated rendering, such as animation.
func CanvasPrintString(canvas *Canvas, termInfo *TermInfo) string {
	var rows unsafe.Pointer
	var count int32
	CanvasPrintRows(canvas, termInfo, &rows, &count)
	defer FreeGStringArray(rows)

	var output strings.Builder
	symbols := CanvasConfigGetPixelMode(CanvasPeekConfig(canvas)) == CHAFA_PIXEL_MODE_SYMBOLS
	for i, row := range unsafe.Slice((**GString)(rows), int(count)) {
		if symbols && i > 0 {
			output.WriteByte('\n')
		}
		output.WriteString(row.String())
	}
	return output.String()
}

type Canvas struct {
	Refs int32

	WidthPixels, HeightPixels int32
	Pixels                    *Pixel
	Cells                     *CanvasCell
	HaveAlpha                 bool
	NeedsClear                bool

	// Whether to consider inverted symbols; FALSE if using FG only
	ConsiderInverted bool

	// Whether to extract symbol colors; FALSE if using default colors
	ExtractColors bool

	// Whether to quantize colors before calculating error (slower, but
	// yields better results in palettized modes, especially 16/8)
	UseQuantizedError bool

	DefaultColors ColorPair
	WorkFactorInt uint32

	// Character to use in cells where fg color == bg color. Typically
	// space, but could be something else depending on the symbol map.
	BlankChar uint8 // gunichar

	// Character to use in cells where fg color == bg color and the color
	// is only legal in FG. Typically 0x2588 (solid block), but could be
	// something else depending on the symbol map. Can be zero if there is
	// no good candidate!
	SolidChar uint8 // gunichar

	Config CanvasConfig

	// Used when setting pixel data
	Dither Dither

	// This is NULL in CHAFA_PIXEL_MODE_SYMBOLS, otherwise one of:
	// (ChafaSixelCanvas *), (ChafaKittyCanvas *), (ChafaIterm2Canvas *)
	PixelCanvas unsafe.Pointer

	// It's possible to have a single placement that covers the entire
	// canvas. In this case, it is stored here.
	Placement *Placement

	// Our palettes. Kind of a big structure, so they go last.
	FgPalette Palette
	BgPalette Palette
}

type Pixel struct {
	Col Color
}

type Color struct {
	Ch [4]uint8
}

type CanvasCell struct {
	C uint8

	// Colors can be either packed RGBA or index
	FgColor uint32
	BgColor uint32
}

type ColorPair struct {
	Colors [2]Color
}

type Dither struct {
	Mode             DitherMode
	Intensity        float64
	GrainWidthShift  int32
	GrainHeightShift int32

	TextureSizeShift int32
	TextureSizeMask  uint32
	TextureData      *int32
}

type DitherMode int32

const (
	CHAFA_DITHER_MODE_NONE      DitherMode = 0
	CHAFA_DITHER_MODE_ORDERED   DitherMode = 1
	CHAFA_DITHER_MODE_DIFFUSION DitherMode = 2
	CHAFA_DITHER_MODE_MAX       DitherMode = 3

	// CHAFA_DITHER_MODE_NOISE     DitherMode = 3
)

type PixelType int32

const (
	/* 32 bits per pixel */

	CHAFA_PIXEL_RGBA8_PREMULTIPLIED PixelType = 0
	CHAFA_PIXEL_BGRA8_PREMULTIPLIED PixelType = 1
	CHAFA_PIXEL_ARGB8_PREMULTIPLIED PixelType = 2
	CHAFA_PIXEL_ABGR8_PREMULTIPLIED PixelType = 3

	CHAFA_PIXEL_RGBA8_UNASSOCIATED PixelType = 4
	CHAFA_PIXEL_BGRA8_UNASSOCIATED PixelType = 5
	CHAFA_PIXEL_ARGB8_UNASSOCIATED PixelType = 6
	CHAFA_PIXEL_ABGR8_UNASSOCIATED PixelType = 7

	/* 24 bits per pixel */

	CHAFA_PIXEL_RGB8 PixelType = 8
	CHAFA_PIXEL_BGR8 PixelType = 9

	CHAFA_PIXEL_MAX PixelType = 10
)

const CHAFA_PALETTE_INDEX_MAX = 259

type Palette struct {
	Type             PaletteType
	Colors           [CHAFA_PALETTE_INDEX_MAX]PaletteColor
	Table            [CHAFA_COLOR_SPACE_MAX]ColorTable
	FirstColor       int32
	NColors          int32
	AlphaThreshold   int32
	TransparentIndex int32
}

type PaletteType int32

const (
	CHAFA_PALETTE_TYPE_DYNAMIC_256 PaletteType = 0
	CHAFA_PALETTE_TYPE_FIXED_256   PaletteType = 1
	CHAFA_PALETTE_TYPE_FIXED_240   PaletteType = 2
	CHAFA_PALETTE_TYPE_FIXED_16    PaletteType = 3
	CHAFA_PALETTE_TYPE_FIXED_8     PaletteType = 4
	CHAFA_PALETTE_TYPE_FIXED_FGBG  PaletteType = 5
)

type PaletteColor struct {
	Col [CHAFA_COLOR_SPACE_MAX]Color
}

type ColorSpace int32

const (
	CHAFA_COLOR_SPACE_RGB    ColorSpace = 0
	CHAFA_COLOR_SPACE_DIN99D ColorSpace = 1
	CHAFA_COLOR_SPACE_MAX    ColorSpace = 2
)

const CHAFA_COLOR_TABLE_MAX_ENTRIES = 256

type ColorTableEntry struct {
	V   [2]int32
	Pen int32
}

type ColorTable struct {
	Entries [CHAFA_COLOR_TABLE_MAX_ENTRIES]ColorTableEntry

	// Each pen is 24 bits (B8G8R8) of color information
	Pens [CHAFA_COLOR_TABLE_MAX_ENTRIES]uint32

	NEntries int32
	IsSorted bool

	Eigenvectors [2]Vec3i32
	Average      Vec3i32

	EigenMul [2]uint32
}

type Vec3i32 struct {
	V [3]int32
}

type CanvasMode int32

const (
	CHAFA_CANVAS_MODE_TRUECOLOR    CanvasMode = 0
	CHAFA_CANVAS_MODE_INDEXED_256  CanvasMode = 1
	CHAFA_CANVAS_MODE_INDEXED_240  CanvasMode = 2
	CHAFA_CANVAS_MODE_INDEXED_16   CanvasMode = 3
	CHAFA_CANVAS_MODE_FGBG_BGFG    CanvasMode = 4
	CHAFA_CANVAS_MODE_FGBG         CanvasMode = 5
	CHAFA_CANVAS_MODE_INDEXED_8    CanvasMode = 6
	CHAFA_CANVAS_MODE_INDEXED_16_8 CanvasMode = 7
	CHAFA_CANVAS_MODE_MAX          CanvasMode = 8
)

type ColorExtractor int32

const (
	CHAFA_COLOR_EXTRACTOR_AVERAGE ColorExtractor = 0
	CHAFA_COLOR_EXTRACTOR_MEDIAN  ColorExtractor = 1
	CHAFA_COLOR_EXTRACTOR_MAX     ColorExtractor = 2
)

type PixelMode int32

const (
	CHAFA_PIXEL_MODE_SYMBOLS PixelMode = 0
	CHAFA_PIXEL_MODE_SIXELS  PixelMode = 1
	CHAFA_PIXEL_MODE_KITTY   PixelMode = 2
	CHAFA_PIXEL_MODE_ITERM2  PixelMode = 3
	CHAFA_PIXEL_MODE_MAX     PixelMode = 4
)

type Optimizations int32

const (
	CHAFA_OPTIMIZATION_REUSE_ATTRIBUTES Optimizations = (1 << 0)
	CHAFA_OPTIMIZATION_SKIP_CELLS       Optimizations = (1 << 1)
	CHAFA_OPTIMIZATION_REPEAT_CELLS     Optimizations = (1 << 2)
	CHAFA_OPTIMIZATION_NONE             Optimizations = 0
	CHAFA_OPTIMIZATION_ALL              Optimizations = 0x7fffffff
)

type Passthrough int32

const (
	CHAFA_PASSTHROUGH_NONE   Passthrough = 0
	CHAFA_PASSTHROUGH_SCREEN Passthrough = 1
	CHAFA_PASSTHROUGH_TMUX   Passthrough = 2
	CHAFA_PASSTHROUGH_MAX    Passthrough = 3
)
