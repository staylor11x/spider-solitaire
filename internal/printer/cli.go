package printer

import (
	"fmt"
	"strings"

	"github.com/staylor11x/spider-solitaire/internal/deck"
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
	card := deck.Card{Suit: deck.Suit(c.Suit), Rank: deck.Rank(c.Rank)}
	return fmt.Sprintf("%s%s", card.RankName(), card.SuitSymbol())
}
