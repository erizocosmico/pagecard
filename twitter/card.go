package twitter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mvader/pagecard/content"
)

// Card contains all the data used to build a twitter card.
type Card struct {
	Type        CardType
	Title       string
	Description string
	Site
	Creator
	Image
	*Player
	*App
}

// Creator is the creator of the content in the card.
type Creator struct {
	ID   string
	User string
}

// Site is the twitter user to publish this card.
type Site struct {
	ID   string
	User string
}

// Image is the representative image of the card.
type Image struct {
	URL string
	Alt string
}

// App contains all the info about an "app" card with all the platforms
// of the application.
type App struct {
	IPhone     AppInfo
	IPad       AppInfo
	GooglePlay AppInfo
	Country    string
}

// AppInfo contains the information of an app for a specific platform.
type AppInfo struct {
	Name string
	ID   string
	URL  string
}

// Player contains all the data for a "player" card.
type Player struct {
	URL               string
	Width             int
	Height            int
	Stream            string
	StreamContentType string
}

// CardType represents the kind of content the card will have.
type CardType byte

const (
	// SummaryCard can be used for many kinds of web content, from blog posts
	// and news articles, to products and restaurants. It is designed to give
	// the reader a preview of the content before clicking through to your
	// website.
	SummaryCard CardType = 1 << iota
	// SummaryBigPictureCard features a large, full-width prominent image
	// alongside a tweet. It is designed to give the reader a rich photo
	// experience, and clicking on the image brings the user to your website.
	SummaryBigPictureCard
	// AppCard is a great way to represent mobile applications on Twitter and
	// to drive installs. We designed the App Card to allow for a name,
	// description and icon, and also to highlight attributes such as the
	// rating and the price.
	AppCard
	// PlayerCard is a great way to deliver audio and video to the audience.
	PlayerCard
)

const (
	twitterPrefix               = "twitter:"
	appPrefix                   = "app:"
	playerPrefix                = "player:"
	cardName                    = "card"
	siteName                    = "site"
	siteIDName                  = "site:id"
	creatorName                 = "creator"
	creatorIDName               = "creator:id"
	titleName                   = "title"
	descriptionName             = "description"
	imageName                   = "image"
	imageAltName                = "image:alt"
	playerName                  = "player"
	playerWidthName             = "player:width"
	playerHeightName            = "player:height"
	playerStreamName            = "player:stream"
	playerStreamContentTypeName = "player:stream:content_type"
	iphoneIDName                = "app:id:iphone"
	ipadIDName                  = "app:id:ipad"
	androidIDName               = "app:id:googleplay"
	androidNameName             = "app:name:googleplay"
	iphoneNameName              = "app:name:iphone"
	ipadNameName                = "app:name:ipad"
	iphoneURLName               = "app:url:iphone"
	ipadURLName                 = "app:url:ipad"
	androidURLName              = "app:url:googleplay"
	appCountryName              = "app:country"
)

// NewCard returns a new twitter card object built with the metatags present
// in the page.
func NewCard(meta []*content.Meta) (*Card, error) {
	var (
		card   = new(Card)
		player *Player
		app    *App
		err    error
	)

	card.Type, meta, err = filterTwitterMeta(meta)
	if err != nil {
		return nil, err
	}

	for _, m := range meta {
		if strings.HasPrefix(m.Name, appPrefix) {
			if card.Type != AppCard {
				continue
			}

			if app == nil {
				app = new(App)
			}

			switch m.Name {
			case iphoneIDName:
				app.IPhone.ID = m.Value
			case ipadIDName:
				app.IPad.ID = m.Value
			case androidIDName:
				app.GooglePlay.ID = m.Value
			case androidNameName:
				app.GooglePlay.Name = m.Value
			case iphoneNameName:
				app.IPhone.Name = m.Value
			case ipadNameName:
				app.IPad.Name = m.Value
			case iphoneURLName:
				app.IPhone.URL = m.Value
			case ipadURLName:
				app.IPad.URL = m.Value
			case androidURLName:
				app.GooglePlay.URL = m.Value
			case appCountryName:
				app.Country = m.Value
			}
		}

		if strings.HasPrefix(m.Name, playerPrefix) || m.Name == playerName {
			if card.Type != PlayerCard {
				continue
			}

			if player == nil {
				player = new(Player)
			}

			switch m.Name {
			case playerName:
				player.URL = m.Value
			case playerHeightName, playerWidthName:
				n, err := strconv.Atoi(m.Value)
				if err != nil {
					return nil, err
				}
				if m.Name == playerHeightName {
					player.Height = n
				} else {
					player.Width = n
				}
			case playerStreamName:
				player.Stream = m.Value
			case playerStreamContentTypeName:
				player.StreamContentType = m.Value
			}
		}

		switch m.Name {
		case siteName:
			card.Site.User = m.Value
		case siteIDName:
			card.Site.ID = m.Value
		case titleName:
			card.Title = m.Value
		case descriptionName:
			card.Description = m.Value
		case creatorName:
			card.Creator.User = m.Value
		case creatorIDName:
			card.Creator.ID = m.Value
		case imageName:
			card.Image.URL = m.Value
		case imageAltName:
			card.Image.Alt = m.Value
		}
	}

	if player != nil {
		card.Player = player
	}

	if app != nil {
		card.App = app
	}

	return card, nil
}

func filterTwitterMeta(meta []*content.Meta) (CardType, []*content.Meta, error) {
	var (
		typ         CardType
		err         error
		twittermeta []*content.Meta
	)

	for _, m := range meta {
		if strings.HasPrefix(m.Name, twitterPrefix) {
			m.Name = m.Name[len(twitterPrefix):]
			if m.Name == cardName {
				typ, err = cardType(m.Value)
				if err != nil {
					return typ, nil, err
				}
			} else {
				twittermeta = append(twittermeta, m)
			}
		}
	}

	return typ, twittermeta, nil
}

const (
	summary             = "summary"
	summaryLargePicture = "summary_large_image"
	app                 = "app"
	player              = "player"
)

func cardType(t string) (typ CardType, err error) {
	switch t {
	case summary:
		typ = SummaryCard
	case summaryLargePicture:
		typ = SummaryBigPictureCard
	case app:
		typ = AppCard
	case player:
		typ = PlayerCard
	default:
		err = fmt.Errorf("invalid card type: %s", t)
	}
	return
}
