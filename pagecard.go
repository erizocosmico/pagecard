package pagecard

import (
	"github.com/mvader/pagecard/content"
	"github.com/mvader/pagecard/opengraph"
	"github.com/mvader/pagecard/twitter"
)

// Info contains all the data retrieved from the opengraph and twitter cards
// metatags in a webpage.
type Info struct {
	OpenGraph *opengraph.Object
	Twitter   *twitter.Card
}

// Get retrieves the Info of a webpage with the given URL.
func Get(url string) (*Info, error) {
	meta, err := content.Read(url)
	if err != nil {
		return nil, err
	}

	obj, err := opengraph.NewObject(meta)
	if err != nil {
		return nil, err
	}

	card, err := twitter.NewCard(meta)
	if err != nil {
		return nil, err
	}

	return &Info{obj, card}, nil
}
