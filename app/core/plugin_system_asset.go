package core

import (
	"bytes"
	"fmt"
	"html"
	"io/fs"
	"io/ioutil"
	"net/http"
	"strings"
)

// AssetLocation is a location where assets can be added.
type AssetLocation string

// AssetType is a type of asset.
type AssetType string

// AuthType is a type of authentication.
type AuthType string

// LayoutType is a type of layout.
type LayoutType string

// FieldType is a type of field.
type FieldType string

const (
	// LocationHead is at the bottom of the HTML <head> section.
	LocationHead AssetLocation = "head"
	// LocationHeader is at the top the HTML <header> section.
	LocationHeader AssetLocation = "header"
	// LocationMain is at the bottom of the HTML <main> section.
	LocationMain AssetLocation = "main"
	// LocationFooter is in the HTML <footer> section.
	LocationFooter AssetLocation = "footer"
	// LocationBody is at the bottom of the HTML <body> section.
	LocationBody AssetLocation = "body"

	// AssetStylesheet is a stylesheet element.
	AssetStylesheet AssetType = "stylesheet"
	// AssetJavaScript is a javascript element.
	AssetJavaScript AssetType = "javascript"
	// AssetGeneric is a generic element.
	AssetGeneric AssetType = "generic"

	// AllAuth is both anonymous and authenticated users.
	AllAuth AuthType = "all" // Default.
	// AnonymousOnly is only non-authenticated users.
	AnonymousOnly AuthType = "anonymous"
	// AuthenticatedOnly is only authenticated users.
	AuthenticatedOnly AuthType = "authenticated"

	// Page is a page layout.
	Page LayoutType = "page"
	// Post is a post layout.
	Post LayoutType = "post"
	// Dashboard is a dashboard layout.
	Dashboard LayoutType = "dashboard"
	// Bloglist is a bloglist layout.
	Bloglist LayoutType = "bloglist"
)

// Asset represents an HTML asset like a stylesheet or javascript file.
type Asset struct {
	Filetype   AssetType     `json:"filetype"`
	Location   AssetLocation `json:"location"`
	Auth       AuthType      `json:"auth"`
	Attributes []Attribute   `json:"attributes"`
	LayoutOnly []LayoutType  `json:"layout"`

	TagName    string `json:"tagname"`
	ClosingTag bool   `json:"closingtag"`

	External bool      `json:"external"`
	Inline   bool      `json:"inline"`
	Path     string    `json:"path"`
	Replace  []Replace `json:"replace"`
}

// Replace represents text to find and replace.
type Replace struct {
	Find    string
	Replace string
}

// Attribute represents an HTML attribute.
type Attribute struct {
	Name  string
	Value interface{}
}

// Routable returns true if the file can be served from the embedded filesystem.
func (file Asset) Routable() bool {
	if file.External || file.Inline || file.Filetype == AssetGeneric {
		return false
	}

	return true
}

// SanitizedPath returns an HTML escaped asset path.
func (file Asset) SanitizedPath() string {
	return html.EscapeString(file.Path)
}

// Element returns an HTML element.
func (file *Asset) Element(v IPlugin, assets fs.FS) string {
	// Build the attributes.
	attrs := make([]string, 0)
	for _, attr := range file.Attributes {
		if attr.Value == nil {
			attrs = append(attrs, fmt.Sprintf(`%v`, html.EscapeString(attr.Name)))
		} else {
			attrs = append(attrs, fmt.Sprintf(`%v="%v"`, html.EscapeString(attr.Name), html.EscapeString(fmt.Sprint(attr.Value))))
		}
	}
	attrsJoined := strings.Join(attrs, " ")
	if len(attrsJoined) > 0 {
		// Add a space at the beginning.
		attrsJoined = " " + attrsJoined
	}

	txt := ""
	switch file.Filetype {
	case AssetStylesheet:
		if file.Inline {
			ff, status, err := file.Contents(assets)
			if status != http.StatusOK {
				// FIXME: Do something with these.
				fmt.Println(err.Error())
				return ""
			}
			txt = fmt.Sprintf("<style>%v</style>", string(ff))
		} else {
			if file.External {
				txt = fmt.Sprintf(`<link rel="stylesheet" href="%v"%v>`, file.SanitizedPath(), attrsJoined)

			} else {
				txt = fmt.Sprintf(`<link rel="stylesheet" href="/plugins/%v/%v?v=%v"%v>`, v.PluginName(), file.SanitizedPath(), v.PluginVersion(), attrsJoined)
			}
		}
	case AssetJavaScript:
		if file.Inline {
			ff, status, err := file.Contents(assets)
			if status != http.StatusOK {
				// FIXME: Do something with these.
				fmt.Println(err.Error())
				return ""
			}
			txt = fmt.Sprintf("<script>%v</script>", string(ff))
		} else {
			if file.External {
				txt = fmt.Sprintf(`<script type="application/javascript" src="%v"%v></script>`, file.SanitizedPath(), attrsJoined)
			} else {
				txt = fmt.Sprintf(`<script type="application/javascript" src="/plugins/%v/%v?v=%v"%v></script>`, v.PluginName(), file.SanitizedPath(), v.PluginVersion(), attrsJoined)
			}
		}
	case AssetGeneric:
		if file.Inline {
			ff, status, err := file.Contents(assets)
			if status != http.StatusOK {
				// FIXME: Do something with these.
				fmt.Println(err.Error())
				return ""
			}
			if file.TagName == "" {
				txt = fmt.Sprintf(`%v`, string(ff))
			} else {
				txt = fmt.Sprintf(`<%v%v>%v</%v>`, html.EscapeString(file.TagName), attrsJoined, string(ff), html.EscapeString(file.TagName))
			}
		} else {
			// FIXME: The closing tag could be false but the inline above will still add one.
			if file.ClosingTag {
				txt = fmt.Sprintf(`<%v%v></%v>`, html.EscapeString(file.TagName), attrsJoined, html.EscapeString(file.TagName))
			} else {
				txt = fmt.Sprintf(`<%v%v>`, html.EscapeString(file.TagName), attrsJoined)
			}
		}
	default:
		fmt.Printf("unsupported asset filetype for plugin (%v): %v", v.PluginName(), file.Filetype)
	}

	return txt
}

// Contents returns the text of the file to inline in HTML after doing replace.
func (file *Asset) Contents(assets fs.FS) (ff []byte, status int, err error) {
	// Use the root directory.
	fsys, err := fs.Sub(assets, ".")
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Open the file.
	f, err := fsys.Open(file.Path)
	if err != nil {
		return nil, http.StatusNotFound, nil
	}
	defer f.Close()

	// Get the contents.
	ff, err = ioutil.ReadAll(f)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Loop over the items to replace.
	for _, rep := range file.Replace {
		ff = bytes.ReplaceAll(ff, []byte(rep.Find), []byte(rep.Replace))
	}

	return ff, http.StatusOK, nil
}
