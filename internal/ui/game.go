package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// Game implements the ebiten.Game interface for Spider Solitaire
type Game struct {
	// we'll add engine state in in the future
}

// NewGame create a new Ebiten game instance
func NewGame() *Game {
	return &Game{}
}

// Update runs game logic at 60 FPS
func (g *Game) Update() error {
	// no logic yet!
	return nil
}

// Draw renders the current frame to the screen
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 0, G: 100, B: 0, A: 255})
}

// Layout return the logical screen dimensions
// Ebiten will scale this to fit the actual window size, but you always draw in 1280x720 coordinates
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1280, 720
}
