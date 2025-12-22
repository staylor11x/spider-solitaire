package game

// UI-safe value types (primitives to avoid UI depending on internal types)
type SuitDTO int
type RankDTO int

type CardDTO struct {
	Rank   RankDTO
	Suit   SuitDTO
	FaceUp bool
}

type PileDTO struct {
	Cards []CardDTO
}

type GameViewDTO struct {
	Tableau        []PileDTO
	StockCount     int
	CompletedCount int
	Won            bool
	Lost           bool
}

func (g *GameState) View() GameViewDTO {
	tableau := make([]PileDTO, len(g.Tableau.Piles))
	for i := range g.Tableau.Piles {
		tableau[i] = pileToDTO(g.Tableau.Piles[i])
	}

	return GameViewDTO{
		Tableau:        tableau,
		StockCount:     len(g.Stock),
		CompletedCount: len(g.Completed),
		Won:            g.Won,
		Lost:           g.Lost,
	}
}

func pileToDTO(p Pile) PileDTO {
	cards := p.Cards()
	out := make([]CardDTO, len(cards))

	for i, c := range cards {
		out[i] = cardToDTO(c)
	}
	return PileDTO{Cards: out}
}

func cardToDTO(c CardInPile) CardDTO {
	return CardDTO{
		Rank:   RankDTO(c.Card.Rank),
		Suit:   SuitDTO(c.Card.Suit),
		FaceUp: c.FaceUp,
	}
}
