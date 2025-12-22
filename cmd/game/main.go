package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	spiderui "github.com/staylor11x/spider-solitaire/internal/ui"
)

func main() {

	// set window properties
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Spider Solitaire")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// create the game instance
	game := spiderui.NewGame()

	// run the game loop, this blocks until the window closes or an error occurs
	if err := ebiten.RunGame(game); err != nil {
		log.Fatalf("game loop failed: %v", err)
	}
}
