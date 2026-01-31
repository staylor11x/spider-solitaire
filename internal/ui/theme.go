package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/basicfont"
)

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
	SelectionBorder   color.RGBA
	HoverOverlay      color.RGBA
	ErrorPillBG       color.RGBA
	ErrorPillText     color.RGBA
	HelpOverlayBG     color.RGBA
	HelpOverlayText   color.RGBA
	PlaceholderBG     color.RGBA
	PlaceholderBorder color.RGBA
}

// Theme combines layout and color definition
type Theme struct {
	Layout Layout
	Colors Colors
	Font   *text.GoXFace
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
		SelectionOverlay:  color.RGBA{R: 0, G: 0, B: 0, A: 100},
		SelectionBorder:   color.RGBA{R: 0, G: 0, B: 0, A: 255},
		HoverOverlay:      color.RGBA{R: 0, G: 0, B: 0, A: 50},
		ErrorPillBG:       color.RGBA{R: 80, G: 70, B: 90, A: 180},
		ErrorPillText:     color.RGBA{R: 230, G: 230, B: 240, A: 255},
		HelpOverlayBG:     color.RGBA{R: 0, G: 0, B: 0, A: 200},
		HelpOverlayText:   color.RGBA{R: 255, G: 255, B: 255, A: 255},
		PlaceholderBG:     color.RGBA{R: 0, G: 100, B: 0, A: 255},
		PlaceholderBorder: color.RGBA{R: 255, G: 255, B: 255, A: 50},
	},
	Font: text.NewGoXFace(basicfont.Face7x13),
}
