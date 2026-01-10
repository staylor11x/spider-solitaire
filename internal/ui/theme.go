package ui

import "image/color"

type Layout struct {
	CardWidth            int
	CardHeight           int
	PileSpacing          int
	CardStackGap         int
	TableauStartX        int
	TableauStartY        int
	StatsX               int
	StatsY               int
	LogicalWidth         int
	LogicalHeight        int
	ErrorDisplayDuration int
}

type Colors struct {
	Background        color.RGBA
	CardFaceUp        color.RGBA
	CardFaceDown      color.RGBA
	CardText          color.RGBA
	SelectionOverlay  color.RGBA
	ErrorPillBG       color.RGBA
	ErrorPillText     color.RGBA
	HelpOverlayBG     color.RGBA
	HelpOverlayText   color.RGBA
	PlaceholderFill   color.RGBA
	PlaceholderStroke color.RGBA
	HoverOverlay      color.RGBA
}

// Theme combines layout and color definition
type Theme struct {
	Layout Layout
	Colors Colors
}

var DefaultTheme = Theme{
	Layout: Layout{
		CardWidth:            80,
		CardHeight:           120,
		PileSpacing:          100,
		CardStackGap:         30,
		TableauStartX:        50,
		TableauStartY:        150,
		StatsX:               20,
		StatsY:               20,
		LogicalWidth:         1280,
		LogicalHeight:        720,
		ErrorDisplayDuration: 180, // 3 seconds at 60 FPS
	},
	Colors: Colors{
		Background:        color.RGBA{R: 0, G: 100, B: 0, A: 255},
		CardFaceUp:        color.RGBA{R: 255, G: 255, B: 255, A: 255},
		CardFaceDown:      color.RGBA{R: 0, G: 0, B: 139, A: 255},
		CardText:          color.RGBA{R: 0, G: 0, B: 0, A: 255},
		SelectionOverlay:  color.RGBA{R: 255, G: 215, B: 0, A: 100},
		ErrorPillBG:       color.RGBA{R: 200, G: 0, B: 0, A: 220},
		ErrorPillText:     color.RGBA{R: 255, G: 255, B: 255, A: 255},
		HelpOverlayBG:     color.RGBA{R: 0, G: 0, B: 0, A: 200},
		HelpOverlayText:   color.RGBA{R: 255, G: 255, B: 255, A: 255},
		PlaceholderFill:   color.RGBA{R: 255, G: 255, B: 255, A: 20},
		PlaceholderStroke: color.RGBA{R: 255, G: 255, B: 255, A: 50},
		HoverOverlay:      color.RGBA{R: 255, G: 255, B: 255, A: 30},
	},
}
