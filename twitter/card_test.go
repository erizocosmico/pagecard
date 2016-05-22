package twitter

import (
	"fmt"
	"testing"

	"github.com/mvader/pagecard/content"
	"github.com/stretchr/testify/assert"
)

func TestNewCard(t *testing.T) {
	cases := []struct {
		meta []*content.Meta
		err  error
		card *Card
	}{
		{makeMeta("twitter:card", "summary"), nil, &Card{Type: SummaryCard}},
		{makeMeta("twitter:card", "app"), nil, &Card{Type: AppCard}},
		{makeMeta("twitter:card", "summary_large_image"), nil, &Card{Type: SummaryBigPictureCard}},
		{makeMeta("twitter:card", "player"), nil, &Card{Type: PlayerCard}},
		{makeMeta("twitter:card", "invented"), fmt.Errorf("invalid card type: %s", "invented"), &Card{}},
		{makeMeta("twitter:site", "site1",
			"twitter:card", "summary",
			"twitter:site:id", "1234",
			"twitter:creator", "creator1",
			"twitter:creator:id", "4321",
			"twitter:title", "title",
			"twitter:description", "description",
			"twitter:image", "image1",
			"twitter:image:alt", "image alt"), nil, &Card{
			Type:        SummaryCard,
			Site:        Site{"1234", "site1"},
			Creator:     Creator{"4321", "creator1"},
			Title:       "title",
			Description: "description",
			Image:       Image{"image1", "image alt"},
		}},
		{makeMeta("twitter:card", "player",
			"twitter:player", "player",
			"twitter:player:width", "5",
			"twitter:player:height", "10",
			"twitter:player:stream", "stream",
			"twitter:player:stream:content_type", "content_type",
		), nil, &Card{
			Type: PlayerCard,
			Player: &Player{
				URL:               "player",
				Width:             5,
				Height:            10,
				Stream:            "stream",
				StreamContentType: "content_type",
			},
		}},
		{makeMeta("twitter:card", "app",
			"twitter:app:id:ipad", "idipad",
			"twitter:app:id:iphone", "idiphone",
			"twitter:app:id:googleplay", "idgp",
			"twitter:app:url:ipad", "urlipad",
			"twitter:app:url:iphone", "urliphone",
			"twitter:app:url:googleplay", "urlgp",
			"twitter:app:name:ipad", "nameipad",
			"twitter:app:name:iphone", "nameiphone",
			"twitter:app:name:googleplay", "namegp",
			"twitter:app:country", "US",
		), nil, &Card{
			Type: AppCard,
			App: &App{
				Country: "US",
				GooglePlay: AppInfo{
					URL:  "urlgp",
					Name: "namegp",
					ID:   "idgp",
				},
				IPhone: AppInfo{
					URL:  "urliphone",
					Name: "nameiphone",
					ID:   "idiphone",
				},
				IPad: AppInfo{
					URL:  "urlipad",
					Name: "nameipad",
					ID:   "idipad",
				},
			},
		}},
	}

	assert := assert.New(t)
	for _, c := range cases {
		card, err := NewCard(c.meta)
		assert.Equal(err, c.err)
		if c.err == nil {
			assert.Equal(card, c.card)
		}
	}
}

func makeMeta(s ...string) []*content.Meta {
	if len(s)%2 != 0 {
		panic("i need k-v pairs")
	}

	var meta []*content.Meta
	for i := 0; i < len(s); i += 2 {
		meta = append(meta, &content.Meta{
			Name:  s[i],
			Value: s[i+1],
		})
	}

	return meta
}
