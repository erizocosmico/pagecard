package content

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const fixture = `
<!DOCTYPE html>
<html xmlns:og="http://ogp.me/ns#">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="apple-itunes-app" content="foo app">
    <meta property="og:url" content="foo url" />
    <meta name="theme-color" content="#000000" />
    <meta property='og:image' content="foo image" />
    <meta property='og:type' content="foo type" />
    <meta property='og:title' content="Bar title" />
    <meta property='og:site_name' content='Foo' />
    <meta name="title" content="Foo title" />
    <meta name="description" content="Foo baz bar."
    />
    <meta property="og:description" content="Foo bar baz."
    />
</head>
<body>
  <meta property="not showing up" content="nope" />
</body>
<meta property="not showing up" content="nope" />
</html>
`

func TestRead(t *testing.T) {
	assert := assert.New(t)
	f, err := ioutil.TempFile("", "pagecard")
	assert.Nil(err)

	_, err = f.Write([]byte(fixture))
	assert.Nil(err)

	tr := &http.Transport{}
	tr.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
	client = &http.Client{Transport: tr}

	metas, err := Read("file://" + f.Name())
	assert.Nil(err)

	results := []Meta{
		{"apple-itunes-app", "foo app"},
		{"og:url", "foo url"},
		{"theme-color", "#000000"},
		{"og:image", "foo image"},
		{"og:type", "foo type"},
		{"og:title", "Bar title"},
		{"og:site_name", "Foo"},
		{"title", "Foo title"},
		{"description", "Foo baz bar."},
		{"og:description", "Foo bar baz."},
	}

	assert.Equal(len(metas), len(results))
	for i, m := range metas {
		assert.Equal(m.Prop, results[i].Prop)
		assert.Equal(m.Value, results[i].Value)
	}
}
