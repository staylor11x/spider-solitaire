package game

// UI-safe value types (primitives to avoid UI depending on internal types).
// These map 1:1 to internal deck enums but keep the UI decoupled.
type SuitDTO int
type RankDTO int

// CardDTO is a single card snapshot for rendering.
// Rank/Suit values mirror internal deck enums; FaceUp indicates visibility.
type CardDTO struct {
	Rank   RankDTO
	Suit   SuitDTO
	FaceUp bool
}

// PileDTO is a rendering-friendly pile snapshot.
// Cards are ordered from bottom (index 0) to top (last index).
type PileDTO struct {
	Cards []CardDTO
}

// GameViewDTO is the full UI snapshot.
// - Tableau: leftmost pile is index 0, rightmost is index 9.
// - StockCount: cards remaining in stock.
// - CompletedCount: completed runs removed from tableau.
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
