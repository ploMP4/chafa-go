package chafa

import (
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"sync"

	"github.com/ebitengine/purego"
)

const libName = "libchafa"

var libOnce sync.Once

func init() {
	libOnce.Do(func() {
		libchafa, err := loadLibrary()
		if err != nil {
			panic(err)
		}

		// Canvas
		purego.RegisterLibFunc(&CanvasNew, libchafa, "chafa_canvas_new")
		purego.RegisterLibFunc(&CanvasNewSimilar, libchafa, "chafa_canvas_new_similar")
		purego.RegisterLibFunc(&CanvasRef, libchafa, "chafa_canvas_ref")
		purego.RegisterLibFunc(&CanvasUnRef, libchafa, "chafa_canvas_unref")
		purego.RegisterLibFunc(&CanvasPeekConfig, libchafa, "chafa_canvas_peek_config")
		purego.RegisterLibFunc(&CanvasSetPlacement, libchafa, "chafa_canvas_set_placement")
		purego.RegisterLibFunc(&CanvasDrawAllPixels, libchafa, "chafa_canvas_draw_all_pixels")
		purego.RegisterLibFunc(&CanvasPrint, libchafa, "chafa_canvas_print")
		purego.RegisterLibFunc(&CanvasPrintRows, libchafa, "chafa_canvas_print_rows")
		purego.RegisterLibFunc(&CanvasPrintRowsStrv, libchafa, "chafa_canvas_print_rows_strv")
		purego.RegisterLibFunc(&CanvasGetCharAt, libchafa, "chafa_canvas_get_char_at")
		purego.RegisterLibFunc(&CanvasSetCharAt, libchafa, "chafa_canvas_set_char_at")
		purego.RegisterLibFunc(&CanvasGetColorsAt, libchafa, "chafa_canvas_get_colors_at")
		purego.RegisterLibFunc(&CanvasSetColorsAt, libchafa, "chafa_canvas_set_colors_at")
		purego.RegisterLibFunc(&CanvasGetRawColorsAt, libchafa, "chafa_canvas_get_raw_colors_at")
		purego.RegisterLibFunc(&CanvasSetRawColorsAt, libchafa, "chafa_canvas_set_raw_colors_at")

		// Config
		purego.RegisterLibFunc(&CanvasConfigNew, libchafa, "chafa_canvas_config_new")
		purego.RegisterLibFunc(&CanvasConfigCopy, libchafa, "chafa_canvas_config_copy")
		purego.RegisterLibFunc(&CanvasConfigRef, libchafa, "chafa_canvas_config_ref")
		purego.RegisterLibFunc(&CanvasConfigUnref, libchafa, "chafa_canvas_config_unref")
		purego.RegisterLibFunc(
			&CanvasConfigGetGeometry,
			libchafa,
			"chafa_canvas_config_get_geometry",
		)
		purego.RegisterLibFunc(
			&CanvasConfigSetGeometry,
			libchafa,
			"chafa_canvas_config_set_geometry",
		)
		purego.RegisterLibFunc(
			&CanvasConfigGetCellGeometry,
			libchafa,
			"chafa_canvas_config_get_cell_geometry",
		)
		purego.RegisterLibFunc(
			&CanvasConfigSetCellGeometry,
			libchafa,
			"chafa_canvas_config_set_cell_geometry",
		)
		purego.RegisterLibFunc(
			&CanvasConfigGetPixelMode,
			libchafa,
			"chafa_canvas_config_get_pixel_mode",
		)
		purego.RegisterLibFunc(
			&CanvasConfigSetPixelMode,
			libchafa,
			"chafa_canvas_config_set_pixel_mode",
		)
		purego.RegisterLibFunc(
			&CanvasConfigGetCanvasMode,
			libchafa,
			"chafa_canvas_config_get_canvas_mode",
		)
		purego.RegisterLibFunc(
			&CanvasConfigSetCanvasMode,
			libchafa,
			"chafa_canvas_config_set_canvas_mode",
		)
		purego.RegisterLibFunc(
			&CanvasConfigGetColorExtractor,
			libchafa,
			"chafa_canvas_config_get_color_extractor",
		)
		purego.RegisterLibFunc(
			&CanvasConfigSetColorExtractor,
			libchafa,
			"chafa_canvas_config_set_color_extractor",
		)
		purego.RegisterLibFunc(
			&CanvasConfigGetColorSpace,
			libchafa,
			"chafa_canvas_config_get_color_space",
		)
		purego.RegisterLibFunc(
			&CanvasConfigSetColorSpace,
			libchafa,
			"chafa_canvas_config_set_color_space",
		)
		purego.RegisterLibFunc(
			&CanvasConfigGetPreprocessingEnabled,
			libchafa,
			"chafa_canvas_config_get_preprocessing_enabled",
		)
		purego.RegisterLibFunc(
			&CanvasConfigSetPreprocessingEnabled,
			libchafa,
			"chafa_canvas_config_set_preprocessing_enabled",
		)
		purego.RegisterLibFunc(
			&CanvasConfigPeekSymbolMap,
			libchafa,
			"chafa_canvas_config_peek_symbol_map",
		)
		purego.RegisterLibFunc(
			&CanvasConfigSetSymbolMap,
			libchafa,
			"chafa_canvas_config_set_symbol_map",
		)
		purego.RegisterLibFunc(
			&CanvasConfigPeekFillSymbolMap,
			libchafa,
			"chafa_canvas_config_peek_fill_symbol_map",
		)
		purego.RegisterLibFunc(
			&CanvasConfigGetTransparencyThreshold,
			libchafa,
			"chafa_canvas_config_get_transparency_threshold",
		)
		purego.RegisterLibFunc(
			&CanvasConfigSetTransparencyThreshold,
			libchafa,
			"chafa_canvas_config_set_transparency_threshold",
		)
		purego.RegisterLibFunc(
			&CanvasConfigGetFgOnlyEnabled,
			libchafa,
			"chafa_canvas_config_get_fg_only_enabled",
		)
		purego.RegisterLibFunc(
			&CanvasConfigSetFgOnlyEnabled,
			libchafa,
			"chafa_canvas_config_set_fg_only_enabled",
		)
		purego.RegisterLibFunc(
			&CanvasConfigGetFgColor,
			libchafa,
			"chafa_canvas_config_get_fg_color",
		)
		purego.RegisterLibFunc(
			&CanvasConfigSetFgColor,
			libchafa,
			"chafa_canvas_config_set_fg_color",
		)
		purego.RegisterLibFunc(
			&CanvasConfigGetBgColor,
			libchafa,
			"chafa_canvas_config_get_bg_color",
		)
		purego.RegisterLibFunc(
			&CanvasConfigSetBgColor,
			libchafa,
			"chafa_canvas_config_set_bg_color",
		)
		purego.RegisterLibFunc(
			&CanvasConfigGetWorkFactor,
			libchafa,
			"chafa_canvas_config_get_work_factor",
		)
		purego.RegisterLibFunc(
			&CanvasConfigSetWorkFactor,
			libchafa,
			"chafa_canvas_config_set_work_factor",
		)
		purego.RegisterLibFunc(
			&CanvasConfigGetDitherMode,
			libchafa,
			"chafa_canvas_config_get_dither_mode",
		)
		purego.RegisterLibFunc(
			&CanvasConfigSetDitherMode,
			libchafa,
			"chafa_canvas_config_set_dither_mode",
		)
		purego.RegisterLibFunc(
			&CanvasConfigGetDitherGrainSize,
			libchafa,
			"chafa_canvas_config_get_dither_grain_size",
		)
		purego.RegisterLibFunc(
			&CanvasConfigSetDitherGrainSize,
			libchafa,
			"chafa_canvas_config_set_dither_grain_size",
		)
		purego.RegisterLibFunc(
			&CanvasConfigGetDitherIntensity,
			libchafa,
			"chafa_canvas_config_get_dither_intensity",
		)
		purego.RegisterLibFunc(
			&CanvasConfigSetDitherIntensity,
			libchafa,
			"chafa_canvas_config_set_dither_intensity",
		)
		purego.RegisterLibFunc(
			&CanvasConfigGetOptimizations,
			libchafa,
			"chafa_canvas_config_get_optimizations",
		)
		purego.RegisterLibFunc(
			&CanvasConfigSetOptimizations,
			libchafa,
			"chafa_canvas_config_set_optimizations",
		)
		purego.RegisterLibFunc(
			&CanvasConfigGetPassthrough,
			libchafa,
			"chafa_canvas_config_get_passthrough",
		)
		purego.RegisterLibFunc(
			&CanvasConfigSetPassthrough,
			libchafa,
			"chafa_canvas_config_set_passthrough",
		)

		// Placement
		purego.RegisterLibFunc(&PlacementNew, libchafa, "chafa_placement_new")
		purego.RegisterLibFunc(&PlacementRef, libchafa, "chafa_placement_ref")
		purego.RegisterLibFunc(&PlacementUnref, libchafa, "chafa_placement_unref")
		purego.RegisterLibFunc(&PlacementGetTuck, libchafa, "chafa_placement_get_tuck")
		purego.RegisterLibFunc(&PlacementSetTuck, libchafa, "chafa_placement_set_tuck")
		purego.RegisterLibFunc(&PlacementGetHAlign, libchafa, "chafa_placement_get_halign")
		purego.RegisterLibFunc(&PlacementSetHAlign, libchafa, "chafa_placement_set_halign")
		purego.RegisterLibFunc(&PlacementGetVAlign, libchafa, "chafa_placement_get_valign")
		purego.RegisterLibFunc(&PlacementSetVAlign, libchafa, "chafa_placement_set_valign")

		// Image
		purego.RegisterLibFunc(&ImageNew, libchafa, "chafa_image_new")
		purego.RegisterLibFunc(&ImageRef, libchafa, "chafa_image_ref")
		purego.RegisterLibFunc(&ImageUnref, libchafa, "chafa_image_unref")
		purego.RegisterLibFunc(&ImageSetFrame, libchafa, "chafa_image_set_frame")

		// Frame
		purego.RegisterLibFunc(&FrameNew, libchafa, "chafa_frame_new")
		purego.RegisterLibFunc(&FrameNewBorrow, libchafa, "chafa_frame_new_borrow")
		purego.RegisterLibFunc(&FrameNewSteal, libchafa, "chafa_frame_new_steal")
		purego.RegisterLibFunc(&FrameRef, libchafa, "chafa_frame_ref")
		purego.RegisterLibFunc(&FrameUnref, libchafa, "chafa_frame_unref")

		// SymbolMap
		purego.RegisterLibFunc(&SymbolMapNew, libchafa, "chafa_symbol_map_new")
		purego.RegisterLibFunc(&SymbolMapCopy, libchafa, "chafa_symbol_map_copy")
		purego.RegisterLibFunc(&SymbolMapRef, libchafa, "chafa_symbol_map_ref")
		purego.RegisterLibFunc(&SymbolMapUnref, libchafa, "chafa_symbol_map_unref")
		purego.RegisterLibFunc(&SymbolMapAddByTags, libchafa, "chafa_symbol_map_add_by_tags")
		purego.RegisterLibFunc(&SymbolMapAddByRange, libchafa, "chafa_symbol_map_add_by_range")
		purego.RegisterLibFunc(&SymbolMapRemoveByTags, libchafa, "chafa_symbol_map_remove_by_tags")
		purego.RegisterLibFunc(
			&SymbolMapRemoveByRange,
			libchafa,
			"chafa_symbol_map_remove_by_range",
		)
		purego.RegisterLibFunc(
			&SymbolMapApplySelectors,
			libchafa,
			"chafa_symbol_map_apply_selectors",
		)
		purego.RegisterLibFunc(
			&SymbolMapGetAllowBuiltinGlyphs,
			libchafa,
			"chafa_symbol_map_get_allow_builtin_glyphs",
		)
		purego.RegisterLibFunc(
			&SymbolMapSetAllowBuiltinGlyphs,
			libchafa,
			"chafa_symbol_map_set_allow_builtin_glyphs",
		)
		purego.RegisterLibFunc(&SymbolMapGetGlyph, libchafa, "chafa_symbol_map_get_glyph")
		purego.RegisterLibFunc(&SymbolMapAddGlyph, libchafa, "chafa_symbol_map_add_glyph")

		// TermDb
		purego.RegisterLibFunc(&TermDbNew, libchafa, "chafa_term_db_new")
		purego.RegisterLibFunc(&TermDbCopy, libchafa, "chafa_term_db_copy")
		purego.RegisterLibFunc(&TermDbRef, libchafa, "chafa_term_db_ref")
		purego.RegisterLibFunc(&TermDbUnref, libchafa, "chafa_term_db_unref")
		purego.RegisterLibFunc(&TermDbGetDefault, libchafa, "chafa_term_db_get_default")
		purego.RegisterLibFunc(&termDbDetect, libchafa, "chafa_term_db_detect")
		purego.RegisterLibFunc(&TermDbGetFallbackInfo, libchafa, "chafa_term_db_get_fallback_info")

		// TermInfo
		purego.RegisterLibFunc(&TermInfoNew, libchafa, "chafa_term_info_new")
		purego.RegisterLibFunc(&TermInfoCopy, libchafa, "chafa_term_info_copy")
		purego.RegisterLibFunc(&TermInfoRef, libchafa, "chafa_term_info_ref")
		purego.RegisterLibFunc(&TermInfoUnref, libchafa, "chafa_term_info_unref")
		purego.RegisterLibFunc(&TermInfoChain, libchafa, "chafa_term_info_chain")
		purego.RegisterLibFunc(&TermInfoSupplement, libchafa, "chafa_term_info_supplement")
		purego.RegisterLibFunc(&TermInfoGetName, libchafa, "chafa_term_info_get_name")
		purego.RegisterLibFunc(&TermInfoSetName, libchafa, "chafa_term_info_set_name")
		purego.RegisterLibFunc(&TermInfoGetQuirks, libchafa, "chafa_term_info_get_quirks")
		purego.RegisterLibFunc(&TermInfoSetQuirks, libchafa, "chafa_term_info_set_quirks")
		purego.RegisterLibFunc(
			&TermInfoGetSafeSymbolTags,
			libchafa,
			"chafa_term_info_get_safe_symbol_tags",
		)
		purego.RegisterLibFunc(
			&TermInfoSetSafeSymbolTags,
			libchafa,
			"chafa_term_info_set_safe_symbol_tags",
		)
		purego.RegisterLibFunc(&TermInfoGetSeq, libchafa, "chafa_term_info_get_seq")
		purego.RegisterLibFunc(&TermInfoSetSeq, libchafa, "chafa_term_info_set_seq")
		purego.RegisterLibFunc(&TermInfoHaveSeq, libchafa, "chafa_term_info_have_seq")
		purego.RegisterLibFunc(&TermInfoGetInheritSeq, libchafa, "chafa_term_info_get_inherit_seq")
		purego.RegisterLibFunc(&TermInfoSetInheritSeq, libchafa, "chafa_term_info_set_inherit_seq")
		purego.RegisterLibFunc(&TermInfoEmitSeq, libchafa, "chafa_term_info_emit_seq")
		purego.RegisterLibFunc(&TermInfoEmitSeqValist, libchafa, "chafa_term_info_emit_seq_valist")
		purego.RegisterLibFunc(&TermInfoParseSeq, libchafa, "chafa_term_info_parse_seq")
		purego.RegisterLibFunc(
			&TermInfoParseSeqVarargs,
			libchafa,
			"chafa_term_info_parse_seq_varargs",
		)
		purego.RegisterLibFunc(
			&TermInfoIsCanvasModeSupported,
			libchafa,
			"chafa_term_info_is_canvas_mode_supported",
		)
		purego.RegisterLibFunc(
			&TermInfoGetBestCanvasMode,
			libchafa,
			"chafa_term_info_get_best_canvas_mode",
		)
		purego.RegisterLibFunc(
			&TermInfoIsPixelModeSupported,
			libchafa,
			"chafa_term_info_is_pixel_mode_supported",
		)
		purego.RegisterLibFunc(
			&TermInfoGetBestPixelMode,
			libchafa,
			"chafa_term_info_get_best_pixel_mode",
		)
		purego.RegisterLibFunc(
			&TermInfoGetIsPixelPassthroughNeeded,
			libchafa,
			"chafa_term_info_get_is_pixel_passthrough_needed",
		)
		purego.RegisterLibFunc(
			&TermInfoSetIsPixelPassthroughNeeded,
			libchafa,
			"chafa_term_info_set_is_pixel_passthrough_needed",
		)
		purego.RegisterLibFunc(
			&TermInfoGetPassthroughType,
			libchafa,
			"chafa_term_info_get_passthrough_type",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitResetTerminalSoft,
			libchafa,
			"chafa_term_info_emit_reset_terminal_soft",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitResetTerminalHard,
			libchafa,
			"chafa_term_info_emit_reset_terminal_hard",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitResetAttributes,
			libchafa,
			"chafa_term_info_emit_reset_attributes",
		)
		purego.RegisterLibFunc(&TermInfoEmitClear, libchafa, "chafa_term_info_emit_clear")
		purego.RegisterLibFunc(
			&TermInfoEmitCursorToPos,
			libchafa,
			"chafa_term_info_emit_cursor_to_pos",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitCursorToTopLeft,
			libchafa,
			"chafa_term_info_emit_cursor_to_top_left",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitCursorToBottomLeft,
			libchafa,
			"chafa_term_info_emit_cursor_to_bottom_left",
		)
		purego.RegisterLibFunc(&TermInfoEmitCursorUp, libchafa, "chafa_term_info_emit_cursor_up")
		purego.RegisterLibFunc(
			&TermInfoEmitCursorDown,
			libchafa,
			"chafa_term_info_emit_cursor_down",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitCursorLeft,
			libchafa,
			"chafa_term_info_emit_cursor_left",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitCursorRight,
			libchafa,
			"chafa_term_info_emit_cursor_right",
		)
		purego.RegisterLibFunc(&TermInfoEmitCursorUp1, libchafa, "chafa_term_info_emit_cursor_up_1")
		purego.RegisterLibFunc(
			&TermInfoEmitCursorDown1,
			libchafa,
			"chafa_term_info_emit_cursor_down_1",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitCursorLeft1,
			libchafa,
			"chafa_term_info_emit_cursor_left_1",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitCursorRight1,
			libchafa,
			"chafa_term_info_emit_cursor_right_1",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitCursorUpScroll,
			libchafa,
			"chafa_term_info_emit_cursor_up_scroll",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitCursorDownScroll,
			libchafa,
			"chafa_term_info_emit_cursor_down_scroll",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitInsertCells,
			libchafa,
			"chafa_term_info_emit_insert_cells",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitDeleteCells,
			libchafa,
			"chafa_term_info_emit_delete_cells",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitInsertRows,
			libchafa,
			"chafa_term_info_emit_insert_rows",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitDeleteRows,
			libchafa,
			"chafa_term_info_emit_delete_rows",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitEnableCursor,
			libchafa,
			"chafa_term_info_emit_enable_cursor",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitDisableCursor,
			libchafa,
			"chafa_term_info_emit_disable_cursor",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitEnableEcho,
			libchafa,
			"chafa_term_info_emit_enable_echo",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitDisableEcho,
			libchafa,
			"chafa_term_info_emit_disable_echo",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitEnableInsert,
			libchafa,
			"chafa_term_info_emit_enable_insert",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitDisableInsert,
			libchafa,
			"chafa_term_info_emit_disable_insert",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitEnableWrap,
			libchafa,
			"chafa_term_info_emit_enable_wrap",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitDisableWrap,
			libchafa,
			"chafa_term_info_emit_disable_wrap",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitEnableBold,
			libchafa,
			"chafa_term_info_emit_enable_bold",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitInvertColors,
			libchafa,
			"chafa_term_info_emit_invert_colors",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitSetColorBg8,
			libchafa,
			"chafa_term_info_emit_set_color_bg_8",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitSetColorFg8,
			libchafa,
			"chafa_term_info_emit_set_color_fg_8",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitSetColorFgbg8,
			libchafa,
			"chafa_term_info_emit_set_color_fgbg_8",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitSetColorFg16,
			libchafa,
			"chafa_term_info_emit_set_color_fg_16",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitSetColorBg16,
			libchafa,
			"chafa_term_info_emit_set_color_bg_16",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitSetColorFgbg16,
			libchafa,
			"chafa_term_info_emit_set_color_fgbg_16",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitSetColorFg256,
			libchafa,
			"chafa_term_info_emit_set_color_fg_256",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitSetColorBg256,
			libchafa,
			"chafa_term_info_emit_set_color_bg_256",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitSetColorFgbg256,
			libchafa,
			"chafa_term_info_emit_set_color_fgbg_256",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitSetColorFgDirect,
			libchafa,
			"chafa_term_info_emit_set_color_fg_direct",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitSetColorBgDirect,
			libchafa,
			"chafa_term_info_emit_set_color_bg_direct",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitSetColorFgbgDirect,
			libchafa,
			"chafa_term_info_emit_set_color_fgbg_direct",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitResetColorFg,
			libchafa,
			"chafa_term_info_emit_reset_color_fg",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitResetColorBg,
			libchafa,
			"chafa_term_info_emit_reset_color_bg",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitResetColorFgbg,
			libchafa,
			"chafa_term_info_emit_reset_color_fgbg",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitSetDefaultFg,
			libchafa,
			"chafa_term_info_emit_set_default_fg",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitSetDefaultBg,
			libchafa,
			"chafa_term_info_emit_set_default_bg",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitResetDefaultFg,
			libchafa,
			"chafa_term_info_emit_reset_default_fg",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitResetDefaultBg,
			libchafa,
			"chafa_term_info_emit_reset_default_bg",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitQueryDefaultFg,
			libchafa,
			"chafa_term_info_emit_query_default_fg",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitQueryDefaultBg,
			libchafa,
			"chafa_term_info_emit_query_default_bg",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitQueryPrimaryDeviceAttributes,
			libchafa,
			"chafa_term_info_emit_query_primary_device_attributes",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitPrimaryDeviceAttributes,
			libchafa,
			"chafa_term_info_emit_primary_device_attributes",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitQueryCellSizePx,
			libchafa,
			"chafa_term_info_emit_query_cell_size_px",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitCellSizePx,
			libchafa,
			"chafa_term_info_emit_cell_size_px",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitQueryTextAreaSizeCells,
			libchafa,
			"chafa_term_info_emit_query_text_area_size_cells",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitTextAreaSizeCells,
			libchafa,
			"chafa_term_info_emit_text_area_size_cells",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitQueryTextAreaSizePx,
			libchafa,
			"chafa_term_info_emit_query_text_area_size_px",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitTextAreaSizePx,
			libchafa,
			"chafa_term_info_emit_text_area_size_px",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitRepeatChar,
			libchafa,
			"chafa_term_info_emit_repeat_char",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitSetScrollingRows,
			libchafa,
			"chafa_term_info_emit_set_scrolling_rows",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitResetScrollingRows,
			libchafa,
			"chafa_term_info_emit_reset_scrolling_rows",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitSaveCursorPos,
			libchafa,
			"chafa_term_info_emit_save_cursor_pos",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitRestoreCursorPos,
			libchafa,
			"chafa_term_info_emit_restore_cursor_pos",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitBeginSixels,
			libchafa,
			"chafa_term_info_emit_begin_sixels",
		)
		purego.RegisterLibFunc(&TermInfoEmitEndSixels, libchafa, "chafa_term_info_emit_end_sixels")
		purego.RegisterLibFunc(
			&TermInfoEmitEnableSixelScrolling,
			libchafa,
			"chafa_term_info_emit_enable_sixel_scrolling",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitDisableSixelScrolling,
			libchafa,
			"chafa_term_info_emit_disable_sixel_scrolling",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitSetSixelAdvanceDown,
			libchafa,
			"chafa_term_info_emit_set_sixel_advance_down",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitSetSixelAdvanceRight,
			libchafa,
			"chafa_term_info_emit_set_sixel_advance_right",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitBeginKittyImmediateImageV1,
			libchafa,
			"chafa_term_info_emit_begin_kitty_immediate_image_v1",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitBeginKittyImmediateVirtImageV1,
			libchafa,
			"chafa_term_info_emit_begin_kitty_immediate_virt_image_v1",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitEndKittyImage,
			libchafa,
			"chafa_term_info_emit_end_kitty_image",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitBeginKittyImageChunk,
			libchafa,
			"chafa_term_info_emit_begin_kitty_image_chunk",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitEndKittyImageChunk,
			libchafa,
			"chafa_term_info_emit_end_kitty_image_chunk",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitBeginIterm2Image,
			libchafa,
			"chafa_term_info_emit_begin_iterm2_image",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitEndIterm2Image,
			libchafa,
			"chafa_term_info_emit_end_iterm2_image",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitBeginScreenPassthrough,
			libchafa,
			"chafa_term_info_emit_begin_screen_passthrough",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitEndScreenPassthrough,
			libchafa,
			"chafa_term_info_emit_end_screen_passthrough",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitEnableAltScreen,
			libchafa,
			"chafa_term_info_emit_enable_alt_screen",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitDisableAltScreen,
			libchafa,
			"chafa_term_info_emit_disable_alt_screen",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitBeginTmuxPassthrough,
			libchafa,
			"chafa_term_info_emit_begin_tmux_passthrough",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitEndTmuxPassthrough,
			libchafa,
			"chafa_term_info_emit_end_tmux_passthrough",
		)
		purego.RegisterLibFunc(&TermInfoEmitReturnKey, libchafa, "chafa_term_info_emit_return_key")
		purego.RegisterLibFunc(
			&TermInfoEmitBackspaceKey,
			libchafa,
			"chafa_term_info_emit_backspace_key",
		)
		purego.RegisterLibFunc(&TermInfoEmitDeleteKey, libchafa, "chafa_term_info_emit_delete_key")
		purego.RegisterLibFunc(
			&TermInfoEmitDeleteCtrlKey,
			libchafa,
			"chafa_term_info_emit_delete_ctrl_key",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitDeleteShiftKey,
			libchafa,
			"chafa_term_info_emit_delete_shift_key",
		)
		purego.RegisterLibFunc(&TermInfoEmitInsertKey, libchafa, "chafa_term_info_emit_insert_key")
		purego.RegisterLibFunc(
			&TermInfoEmitInsertCtrlKey,
			libchafa,
			"chafa_term_info_emit_insert_ctrl_key",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitInsertShiftKey,
			libchafa,
			"chafa_term_info_emit_insert_shift_key",
		)
		purego.RegisterLibFunc(&TermInfoEmitHomeKey, libchafa, "chafa_term_info_emit_home_key")
		purego.RegisterLibFunc(
			&TermInfoEmitHomeCtrlKey,
			libchafa,
			"chafa_term_info_emit_home_ctrl_key",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitHomeShiftKey,
			libchafa,
			"chafa_term_info_emit_home_shift_key",
		)
		purego.RegisterLibFunc(&TermInfoEmitEndKey, libchafa, "chafa_term_info_emit_end_key")
		purego.RegisterLibFunc(
			&TermInfoEmitEndCtrlKey,
			libchafa,
			"chafa_term_info_emit_end_ctrl_key",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitEndShiftKey,
			libchafa,
			"chafa_term_info_emit_end_shift_key",
		)
		purego.RegisterLibFunc(&TermInfoEmitUpKey, libchafa, "chafa_term_info_emit_up_key")
		purego.RegisterLibFunc(&TermInfoEmitUpCtrlKey, libchafa, "chafa_term_info_emit_up_ctrl_key")
		purego.RegisterLibFunc(
			&TermInfoEmitUpShiftKey,
			libchafa,
			"chafa_term_info_emit_up_shift_key",
		)
		purego.RegisterLibFunc(&TermInfoEmitDownKey, libchafa, "chafa_term_info_emit_down_key")
		purego.RegisterLibFunc(
			&TermInfoEmitDownCtrlKey,
			libchafa,
			"chafa_term_info_emit_down_ctrl_key",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitDownShiftKey,
			libchafa,
			"chafa_term_info_emit_down_shift_key",
		)
		purego.RegisterLibFunc(&TermInfoEmitLeftKey, libchafa, "chafa_term_info_emit_left_key")
		purego.RegisterLibFunc(
			&TermInfoEmitLeftCtrlKey,
			libchafa,
			"chafa_term_info_emit_left_ctrl_key",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitLeftShiftKey,
			libchafa,
			"chafa_term_info_emit_left_shift_key",
		)
		purego.RegisterLibFunc(&TermInfoEmitRightKey, libchafa, "chafa_term_info_emit_right_key")
		purego.RegisterLibFunc(
			&TermInfoEmitRightCtrlKey,
			libchafa,
			"chafa_term_info_emit_right_ctrl_key",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitRightShiftKey,
			libchafa,
			"chafa_term_info_emit_right_shift_key",
		)
		purego.RegisterLibFunc(&TermInfoEmitPageUpKey, libchafa, "chafa_term_info_emit_page_up_key")
		purego.RegisterLibFunc(
			&TermInfoEmitPageUpCtrlKey,
			libchafa,
			"chafa_term_info_emit_page_up_ctrl_key",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitPageUpShiftKey,
			libchafa,
			"chafa_term_info_emit_page_up_shift_key",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitPageDownKey,
			libchafa,
			"chafa_term_info_emit_page_down_key",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitPageDownCtrlKey,
			libchafa,
			"chafa_term_info_emit_page_down_ctrl_key",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitPageDownShiftKey,
			libchafa,
			"chafa_term_info_emit_page_down_shift_key",
		)
		purego.RegisterLibFunc(&TermInfoEmitTabKey, libchafa, "chafa_term_info_emit_tab_key")
		purego.RegisterLibFunc(
			&TermInfoEmitTabShiftKey,
			libchafa,
			"chafa_term_info_emit_tab_shift_key",
		)
		purego.RegisterLibFunc(&TermInfoEmitF1Key, libchafa, "chafa_term_info_emit_f1_key")
		purego.RegisterLibFunc(&TermInfoEmitF1CtrlKey, libchafa, "chafa_term_info_emit_f1_ctrl_key")
		purego.RegisterLibFunc(
			&TermInfoEmitF1ShiftKey,
			libchafa,
			"chafa_term_info_emit_f1_shift_key",
		)
		purego.RegisterLibFunc(&TermInfoEmitF2Key, libchafa, "chafa_term_info_emit_f2_key")
		purego.RegisterLibFunc(&TermInfoEmitF2CtrlKey, libchafa, "chafa_term_info_emit_f2_ctrl_key")
		purego.RegisterLibFunc(
			&TermInfoEmitF2ShiftKey,
			libchafa,
			"chafa_term_info_emit_f2_shift_key",
		)
		purego.RegisterLibFunc(&TermInfoEmitF3Key, libchafa, "chafa_term_info_emit_f3_key")
		purego.RegisterLibFunc(&TermInfoEmitF3CtrlKey, libchafa, "chafa_term_info_emit_f3_ctrl_key")
		purego.RegisterLibFunc(
			&TermInfoEmitF3ShiftKey,
			libchafa,
			"chafa_term_info_emit_f3_shift_key",
		)
		purego.RegisterLibFunc(&TermInfoEmitF4Key, libchafa, "chafa_term_info_emit_f4_key")
		purego.RegisterLibFunc(&TermInfoEmitF4CtrlKey, libchafa, "chafa_term_info_emit_f4_ctrl_key")
		purego.RegisterLibFunc(
			&TermInfoEmitF4ShiftKey,
			libchafa,
			"chafa_term_info_emit_f4_shift_key",
		)
		purego.RegisterLibFunc(&TermInfoEmitF5Key, libchafa, "chafa_term_info_emit_f5_key")
		purego.RegisterLibFunc(&TermInfoEmitF5CtrlKey, libchafa, "chafa_term_info_emit_f5_ctrl_key")
		purego.RegisterLibFunc(
			&TermInfoEmitF5ShiftKey,
			libchafa,
			"chafa_term_info_emit_f5_shift_key",
		)
		purego.RegisterLibFunc(&TermInfoEmitF6Key, libchafa, "chafa_term_info_emit_f6_key")
		purego.RegisterLibFunc(&TermInfoEmitF6CtrlKey, libchafa, "chafa_term_info_emit_f6_ctrl_key")
		purego.RegisterLibFunc(
			&TermInfoEmitF6ShiftKey,
			libchafa,
			"chafa_term_info_emit_f6_shift_key",
		)
		purego.RegisterLibFunc(&TermInfoEmitF7Key, libchafa, "chafa_term_info_emit_f7_key")
		purego.RegisterLibFunc(&TermInfoEmitF7CtrlKey, libchafa, "chafa_term_info_emit_f7_ctrl_key")
		purego.RegisterLibFunc(
			&TermInfoEmitF7ShiftKey,
			libchafa,
			"chafa_term_info_emit_f7_shift_key",
		)
		purego.RegisterLibFunc(&TermInfoEmitF8Key, libchafa, "chafa_term_info_emit_f8_key")
		purego.RegisterLibFunc(&TermInfoEmitF8CtrlKey, libchafa, "chafa_term_info_emit_f8_ctrl_key")
		purego.RegisterLibFunc(
			&TermInfoEmitF8ShiftKey,
			libchafa,
			"chafa_term_info_emit_f8_shift_key",
		)
		purego.RegisterLibFunc(&TermInfoEmitF9Key, libchafa, "chafa_term_info_emit_f9_key")
		purego.RegisterLibFunc(&TermInfoEmitF9CtrlKey, libchafa, "chafa_term_info_emit_f9_ctrl_key")
		purego.RegisterLibFunc(
			&TermInfoEmitF9ShiftKey,
			libchafa,
			"chafa_term_info_emit_f9_shift_key",
		)
		purego.RegisterLibFunc(&TermInfoEmitF10Key, libchafa, "chafa_term_info_emit_f10_key")
		purego.RegisterLibFunc(
			&TermInfoEmitF10CtrlKey,
			libchafa,
			"chafa_term_info_emit_f10_ctrl_key",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitF10ShiftKey,
			libchafa,
			"chafa_term_info_emit_f10_shift_key",
		)
		purego.RegisterLibFunc(&TermInfoEmitF11Key, libchafa, "chafa_term_info_emit_f11_key")
		purego.RegisterLibFunc(
			&TermInfoEmitF11CtrlKey,
			libchafa,
			"chafa_term_info_emit_f11_ctrl_key",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitF11ShiftKey,
			libchafa,
			"chafa_term_info_emit_f11_shift_key",
		)
		purego.RegisterLibFunc(&TermInfoEmitF12Key, libchafa, "chafa_term_info_emit_f12_key")
		purego.RegisterLibFunc(
			&TermInfoEmitF12CtrlKey,
			libchafa,
			"chafa_term_info_emit_f12_ctrl_key",
		)
		purego.RegisterLibFunc(
			&TermInfoEmitF12ShiftKey,
			libchafa,
			"chafa_term_info_emit_f12_shift_key",
		)

		// Features
		purego.RegisterLibFunc(&GetBuiltinFeatures, libchafa, "chafa_get_builtin_features")
		purego.RegisterLibFunc(&GetSupportedFeatures, libchafa, "chafa_get_supported_features")
		purego.RegisterLibFunc(&DescribeFeatures, libchafa, "chafa_describe_features")
		purego.RegisterLibFunc(&GetNThreads, libchafa, "chafa_get_n_threads")
		purego.RegisterLibFunc(&SetNThreads, libchafa, "chafa_set_n_threads")
		purego.RegisterLibFunc(&GetNActualThreads, libchafa, "chafa_get_n_actual_threads")

		// Miscellaneous
		purego.RegisterLibFunc(&CalcCanvasGeometry, libchafa, "chafa_calc_canvas_geometry")
	})
}

type GString struct {
	str string
}

func (gstr *GString) String() string {
	return gstr.str
}

type GError struct {
	Domain  uint32
	Code    int32
	Message string
}

func Load(path string) (pixels []uint8, width, height int32, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, 0, 0, err
	}
	defer file.Close()

	var img image.Image

	switch filepath.Ext(path) {
	case "png":
		img, err = png.Decode(file)
		if err != nil {
			return nil, 0, 0, err
		}
	case "jpg", "jpeg":
		img, err = jpeg.Decode(file)
		if err != nil {
			return nil, 0, 0, err
		}
	case "gif":
		img, err = gif.Decode(file)
		if err != nil {
			return nil, 0, 0, err
		}
	default:
		img, _, err = image.Decode(file)
		if err != nil {
			return nil, 0, 0, err
		}
	}

	bounds := img.Bounds()
	width = int32(bounds.Dx())
	height = int32(bounds.Dy())

	rgbaImg := image.NewRGBA(bounds)
	draw.Draw(rgbaImg, bounds, img, bounds.Min, draw.Src)

	return rgbaImg.Pix, width, height, nil
}
