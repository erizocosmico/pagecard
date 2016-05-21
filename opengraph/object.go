package opengraph

import (
	"errors"
	"strconv"
	"strings"

	"github.com/mvader/pagecard/content"
)

// Object is the representation of a webpage as an object within the graph.
type Object struct {
	Title            string
	Type             string
	URL              string
	Description      string
	Locale           string
	AlternateLocales []string
	Determiners      []string
	SiteName         string
	Images           []*Image
	Videos           []*Video
	Audios           []*Audio
}

// MediaProperties defines the properties of an audio, video or image object.
type MediaProperties struct {
	URL       string
	SecureURL string
	Type      string
}

// Size is the height and width of an media object (video or image).
type Size struct {
	Width  int
	Height int
}

// Image represents an image file to represent the object within the graph.
type Image struct {
	MediaProperties
	Size
}

// Video represents a video file to complement the object.
type Video struct {
	MediaProperties
	Size
}

// Audio represents an audio file to accompany the object.
type Audio struct {
	MediaProperties
}

const (
	ogPrefix    = "og:"
	imagePrefix = "image:"
	audioPrefix = "audio:"
	videoPrefix = "video:"
	title       = "title"
	typ         = "type"
	image       = "image"
	url         = "url"
	audio       = "audio"
	description = "description"
	determiner  = "determiner"
	locale      = "locale"
	altLocale   = "locale:alternate"
	siteName    = "site_name"
	video       = "video"
	secURL      = "secure_url"
	width       = "width"
	height      = "height"
)

var (
	errImgNotInitialized   = errors.New("invalid field: requires og:image declared before")
	errVideoNotInitialized = errors.New("invalid field: requires og:video declared before")
	errAudioNotInitialized = errors.New("invalid field: requires og:audio declared before")
)

// NewObject creates the object representation of the OpenGraph object from
// the metadata on the webpage.
func NewObject(meta []*content.Meta) (*Object, error) {
	var (
		obj = new(Object)
		img *Image
		vid *Video
		aud *Audio
	)

	for _, m := range meta {
		if !strings.HasPrefix(m.Name, ogPrefix) {
			continue
		}

		name := m.Name[3:]
		if strings.HasPrefix(name, imagePrefix) {
			if img == nil {
				return nil, errImgNotInitialized
			}

			name = name[len(imagePrefix):]
			switch name {
			case secURL:
				img.SecureURL = m.Value
			case typ:
				img.Type = m.Value
			case height, width:
				size, err := strconv.Atoi(m.Value)
				if err != nil {
					return nil, err
				}
				if name == height {
					img.Height = size
				} else {
					img.Width = size
				}
			}
			continue
		}

		if strings.HasPrefix(name, videoPrefix) {
			if vid == nil {
				return nil, errVideoNotInitialized
			}

			name = name[len(videoPrefix):]
			switch name {
			case secURL:
				vid.SecureURL = m.Value
			case typ:
				vid.Type = m.Value
			case height, width:
				size, err := strconv.Atoi(m.Value)
				if err != nil {
					return nil, err
				}
				if name == height {
					vid.Height = size
				} else {
					vid.Width = size
				}
			}
			continue
		}

		if strings.HasPrefix(name, audioPrefix) {
			if aud == nil {
				return nil, errAudioNotInitialized
			}

			switch name[len(audioPrefix):] {
			case secURL:
				aud.SecureURL = m.Value
			case typ:
				aud.Type = m.Value
			}
			continue
		}

		switch name {
		case title:
			obj.Title = m.Value
		case typ:
			obj.Type = m.Value
		case image:
			if img != nil {
				obj.Images = append(obj.Images, img)
			}
			img = &Image{}
			img.URL = m.Value
		case audio:
			if aud != nil {
				obj.Audios = append(obj.Audios, aud)
			}
			aud = &Audio{}
			aud.URL = m.Value
		case video:
			if vid != nil {
				obj.Videos = append(obj.Videos, vid)
			}
			vid = &Video{}
			vid.URL = m.Value
		case url:
			obj.URL = m.Value
		case description:
			obj.Description = m.Value
		case determiner:
			obj.Determiners = append(obj.Determiners, m.Value)
		case locale:
			obj.Locale = m.Value
		case altLocale:
			obj.AlternateLocales = append(obj.AlternateLocales, m.Value)
		case siteName:
			obj.SiteName = m.Value
		}
	}

	if img != nil {
		obj.Images = append(obj.Images, img)
	}

	if vid != nil {
		obj.Videos = append(obj.Videos, vid)
	}

	if aud != nil {
		obj.Audios = append(obj.Audios, aud)
	}

	return obj, nil
}
