package printer

import (
	"fmt"
	"strings"

	"github.com/staylor11x/spider-solitaire/internal/game"
)

type Options struct {
	UnicodeSuits bool // if false, default to ASCII
}

func Render(view game.GameViewDTO, opts Options) string {
	var b strings.Builder

	// header
	fmt.Fprintf(&b, "Stock: %d | Completed: %d | Won %v | Lost: %v \n",
		view.StockCount, view.CompletedCount, view.Won, view.Lost)

	// Tableau: one line per pile, bottom->top order
	for i, pile := range view.Tableau {
		fmt.Fprintf(&b, "p%d: ", i)
		for j, c := range pile.Cards {
			if j > 0 {
				b.WriteString(", ")
			}
			b.WriteString(formatCard(c, opts))
		}
		b.WriteByte('\n')
	}

	return b.String()
}

func formatCard(c game.CardDTO, opts Options) string {
	if !c.FaceUp {
		return "##"
	}
	return fmt.Sprintf("%s%s", rankStr(int(c.Rank)), suitStr(int(c.Suit), opts))
}

func rankStr(r int) string {
	switch r {
	case 1:
		return "A"
	case 11:
		return "J"
	case 12:
		return "Q"
	case 13:
		return "K"
	default:
		return fmt.Sprintf("%d", r)
	}
}

func suitStr(s int, opts Options) string {
	if opts.UnicodeSuits {
		switch s {
		case 0:
			return "♠"
		case 1:
			return "♥"
		case 2: 
			return "♦"
		case 3:
			return "♣"
		default:
			return "?"
		}
	}
	// ASCI fallback
	switch s {
	case 0:
		return "S"
	case 1:
		return "H"
	case 2: 
		return "D"
	case 3:
		return "C"
	default:
		return "?"
	}
}
