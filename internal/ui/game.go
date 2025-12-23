package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/staylor11x/spider-solitaire/internal/game"
)

const (
	CardWidth     = 80
	CardHeight    = 120
	PileSpacing   = 100
	CardStackGap  = 30
	TableauStartX = 50
	TableauStartY = 150
	StatsX        = 20
	StatsY        = 20
)

var (
	CardFaceUpColor   = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	CardFaceDownColor = color.RGBA{R: 50, G: 50, B: 150, A: 255}
	TextColor         = color.RGBA{R: 0, G: 0, B: 0, A: 255}
)

// Game implements the ebiten.Game interface for Spider Solitaire
type Game struct {
	state *game.GameState  // Engine state, mutated only in Update
	view  game.GameViewDTO // Read-only snapshot for rendering
}

// NewGame create a new Ebiten game instance
func NewGame() *Game {
	state, err := game.DealInitialGame()
	if err != nil {
		panic(err) // TODO: Handle this error gracefully
	}

	view := state.View()

	return &Game{
		state: state,
		view:  view,
	}
}

// Update runs game logic at 60 FPS
func (g *Game) Update() error {
	// no logic yet!
	return nil
}

// Draw renders the current frame to the screen
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 0, G: 100, B: 0, A: 255})

	drawTableau(screen, g.view)

	drawStats(screen, g.view)
}

// Layout return the logical screen dimensions
// Ebiten will scale this to fit the actual window size, but you always draw in 1280x720 coordinates
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1280, 720
}