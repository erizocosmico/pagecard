package opengraph

import (
	"testing"

	"github.com/mvader/pagecard/content"
	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	cases := []struct {
		meta []*content.Meta
		err  error
		obj  *Object
	}{
		// General properties
		{makeMeta("og:title", "title"), nil, &Object{Title: "title"}},
		{makeMeta("og:type", "video.movie"), nil, &Object{Type: "video.movie"}},
		{makeMeta("og:url", "http://foo.bar"), nil, &Object{URL: "http://foo.bar"}},
		{makeMeta("og:description", "Foo is bar."), nil, &Object{Description: "Foo is bar."}},
		{makeMeta("og:determiner", "a"), nil, &Object{Determiners: []string{"a"}}},
		{makeMeta(
			"og:determiner", "a",
			"og:determiner", "an",
		), nil, &Object{Determiners: []string{"a", "an"}}},
		{makeMeta("og:locale", "en"), nil, &Object{Locale: "en"}},
		{makeMeta("og:locale:alternate", "es"), nil, &Object{AlternateLocales: []string{"es"}}},
		{makeMeta(
			"og:locale:alternate", "es",
			"og:locale:alternate", "es-ca",
		), nil, &Object{AlternateLocales: []string{"es", "es-ca"}}},
		{makeMeta("og:site_name", "Foo Site"), nil, &Object{SiteName: "Foo Site"}},

		// Image
		{makeMeta(
			"og:image", "image1",
		), nil, &Object{
			Images: []*Image{
				&Image{
					MediaProperties: MediaProperties{URL: "image1"},
				},
			},
		}},
		{makeMeta(
			"og:image", "image1",
			"og:image", "image2",
		), nil, &Object{
			Images: []*Image{
				&Image{
					MediaProperties: MediaProperties{URL: "image1"},
				},
				&Image{
					MediaProperties: MediaProperties{URL: "image2"},
				},
			},
		}},
		{makeMeta(
			"og:image", "image1",
			"og:image:secure_url", "securl1",
			"og:image:height", "100",
			"og:image:width", "200",
			"og:image", "image2",
			"og:image:type", "type2",
		), nil, &Object{
			Images: []*Image{
				&Image{
					MediaProperties: MediaProperties{URL: "image1", SecureURL: "securl1"},
					Size:            Size{200, 100},
				},
				&Image{
					MediaProperties: MediaProperties{URL: "image2", Type: "type2"},
				},
			},
		}},
		{makeMeta(
			"og:image:secure_url", "securl1",
			"og:image", "image2",
		), errImgNotInitialized, nil},

		// Video
		{makeMeta(
			"og:video", "video1",
		), nil, &Object{
			Videos: []*Video{
				&Video{
					MediaProperties: MediaProperties{URL: "video1"},
				},
			},
		}},
		{makeMeta(
			"og:video", "video1",
			"og:video", "video2",
		), nil, &Object{
			Videos: []*Video{
				&Video{
					MediaProperties: MediaProperties{URL: "video1"},
				},
				&Video{
					MediaProperties: MediaProperties{URL: "video2"},
				},
			},
		}},
		{makeMeta(
			"og:video", "video1",
			"og:video:secure_url", "securl1",
			"og:video:height", "100",
			"og:video:width", "200",
			"og:video", "video2",
			"og:video:type", "type2",
		), nil, &Object{
			Videos: []*Video{
				&Video{
					MediaProperties: MediaProperties{URL: "video1", SecureURL: "securl1"},
					Size:            Size{200, 100},
				},
				&Video{
					MediaProperties: MediaProperties{URL: "video2", Type: "type2"},
				},
			},
		}},
		{makeMeta(
			"og:video:secure_url", "securl1",
			"og:video", "video2",
		), errVideoNotInitialized, nil},

		// Audio
		{makeMeta(
			"og:audio", "audio1",
		), nil, &Object{
			Audios: []*Audio{
				&Audio{
					MediaProperties: MediaProperties{URL: "audio1"},
				},
			},
		}},
		{makeMeta(
			"og:audio", "audio1",
			"og:audio", "audio2",
		), nil, &Object{
			Audios: []*Audio{
				&Audio{
					MediaProperties: MediaProperties{URL: "audio1"},
				},
				&Audio{
					MediaProperties: MediaProperties{URL: "audio2"},
				},
			},
		}},
		{makeMeta(
			"og:audio", "audio1",
			"og:audio:secure_url", "securl1",
			"og:audio", "audio2",
			"og:audio:type", "type2",
		), nil, &Object{
			Audios: []*Audio{
				&Audio{
					MediaProperties: MediaProperties{URL: "audio1", SecureURL: "securl1"},
				},
				&Audio{
					MediaProperties: MediaProperties{URL: "audio2", Type: "type2"},
				},
			},
		}},
		{makeMeta(
			"og:audio:secure_url", "securl1",
			"og:audio", "audio2",
		), errAudioNotInitialized, nil},
	}

	assert := assert.New(t)
	for _, c := range cases {
		obj, err := NewObject(c.meta)
		assert.Equal(err, c.err)
		if c.err == nil {
			assert.Equal(obj, c.obj)
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
