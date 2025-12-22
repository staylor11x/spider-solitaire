package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/staylor11x/spider-solitaire/internal/game"
	"github.com/staylor11x/spider-solitaire/internal/printer"
)

func main() {
	ascii := flag.Bool("ascii", false, "use ASCII suits (S/H/D/C) instead of Unicode")
	flag.Parse()

	g, err := game.DealInitialGame()
	if err != nil {
		log.Fatalf("deal failed: %v", err)
	}

	view := g.View()
	out := printer.Render(view, printer.Options{
		UnicodeSuits: !*ascii,
	})
	fmt.Print(out)
}