package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/staylor11x/spider-solitaire/internal/deck"
	spiderui "github.com/staylor11x/spider-solitaire/internal/ui"
)

var (
	Version   = "dev" // set by -ldflags at build time
	BuildTime = "unknown"
)

func main() {

	log.Printf("Spider Solitaire %s (built %s)", Version, BuildTime)

	// set window properties
	ebiten.SetWindowTitle(fmt.Sprintf("Spider Solitaire %s", Version))
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Spider Solitaire")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// create the game instance
	game := spiderui.NewGame(deck.OneSuit)

	// run the game loop, this blocks until the window closes or an error occurs
	if err := ebiten.RunGame(game); err != nil {
		log.Fatalf("game loop failed: %v", err)
	}
}
