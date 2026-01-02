package assets

import "embed"

// Files embeds card faces and backs for use at runtime.
//
// Directory layout (relative to this package):
// images/cards/*.png
// images/card-back/*.png
//
// Ensure filenames are lowercase and match rank/suit mapping.

//go:embed images/cards/*.png images/card-back/*.png
var Files embed.FS
