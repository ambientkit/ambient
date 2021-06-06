package ambient

import (
	"bytes"
	"fmt"
	"html"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
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

	// AuthAll is both anonymous and authenticated users.
	AuthAll AuthType = "all" // Default.
	// AuthAnonymousOnly is only non-authenticated users.
	AuthAnonymousOnly AuthType = "anonymous"
	// AuthOnly is only authenticated users.
	AuthOnly AuthType = "authenticated"

	// LayoutPage is a page layout.
	LayoutPage LayoutType = "page"
	// LayoutPost is a post layout.
	LayoutPost LayoutType = "post"
)

// Asset represents an HTML asset like a stylesheet or javascript file.
type Asset struct {
	// Filetype is the type of asset: generic, stylesheet, or javascript. (required)
	Filetype AssetType `json:"filetype"`
	// Location is the location on the HTML page where the asset will be
	// added. (required)
	Location AssetLocation `json:"location"`
	// Auth determines whether to show the asset to all users, only authenticated
	// users, or only non-authenticated users. Will display to all users if
	// not specified. (optional)
	Auth AuthType `json:"auth"`
	// Attributes are a list of HTML attributes on all filetypes except on
	// generic with no TagName. (optional)
	Attributes []Attribute `json:"attributes"`
	// LayoutOnly are a list of layout types where the element will be added.
	// Supports page and post. Will display on all layouts if not specified.
	// (optional)
	LayoutOnly []LayoutType `json:"layout"`

	// TagName is only for generic assets when Inline is true. Will specify the
	// type of element to create. If empty, then the asset will be written to
	// the page without a surrounding HTML element.
	TagName string `json:"tagname"`
	// ClosingTag, if true, will add a closing tag. It's only for generic assets
	// when inline is false.
	ClosingTag bool `json:"closingtag"`

	// External, if true, will just use the path as the source of the element.
	// It is only for stylesheet and javascript filetypes.
	External bool `json:"external"`
	// Inline if true, will output the contents from an embedded file (Path) or
	// the contents (Content) after doing a find/replace (Replace).
	Inline bool `json:"inline"`
	// SkipExistCheck if true, will not check for the file existing because it's
	// managed by a route.
	SkipExistCheck bool `json:"skipexist"`
	// Path is relative path to the embedded file or the full path to the
	// external asset. (optional)
	Path string `json:"path"`
	// Content is the content that will output on the page. Path must be empty
	// for content to be used and content is only used when Inline is true.
	Content string `json:"content"`
	// Replace is a list of find and replace strings that are run on the Path
	// or Content when Inline is true.
	Replace []Replace `json:"replace"`
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
func (file *Asset) Element(logger AppLogger, v Plugin, assets fs.FS, debug bool) string {
	// Build the attributes.
	attrs := make([]string, 0)
	for _, attr := range file.Attributes {
		if attr.Value == nil {
			attrs = append(attrs, fmt.Sprintf(`%v`, html.EscapeString(attr.Name)))
		} else {
			attrs = append(attrs, fmt.Sprintf(`%v="%v"`, html.EscapeString(attr.Name), html.EscapeString(fmt.Sprint(attr.Value))))
		}
	}

	if debug {
		attrs = append(attrs, fmt.Sprintf(`%v="%v"`, html.EscapeString("data-ambplugin"), html.EscapeString(fmt.Sprint(v.PluginName()))))
	}

	attrsJoined := strings.Join(attrs, " ")
	if len(attrsJoined) > 0 {
		// Add a space at the beginning.
		attrsJoined = " " + attrsJoined
	}

	// Get the URL prefix for assets.
	urlprefix := os.Getenv("AMB_URL_PREFIX")
	if len(urlprefix) == 0 {
		urlprefix = ""
	}

	txt := ""
	switch file.Filetype {
	case AssetStylesheet:
		if file.Inline {
			ff, status, err := file.Contents(assets)
			if status != http.StatusOK {
				logger.Error("plugin injector: error getting file contents: %v", err.Error())
				return ""
			}
			txt = fmt.Sprintf("<style%v>%v</style>", attrsJoined, string(ff))
		} else {
			if file.External {
				txt = fmt.Sprintf(`<link rel="stylesheet" href="%v"%v>`, file.SanitizedPath(), attrsJoined)
			} else {
				txt = fmt.Sprintf(`<link rel="stylesheet" href="%v/plugins/%v/%v?v=%v"%v>`, urlprefix, v.PluginName(), file.SanitizedPath(), v.PluginVersion(), attrsJoined)
			}
		}
	case AssetJavaScript:
		if file.Inline {
			ff, status, err := file.Contents(assets)
			if status != http.StatusOK {
				logger.Error("plugin injector: error getting file contents: %v", err.Error())
				return ""
			}
			txt = fmt.Sprintf("<script%v>%v</script>", attrsJoined, string(ff))
		} else {
			if file.External {
				txt = fmt.Sprintf(`<script type="application/javascript" src="%v"%v></script>`, file.SanitizedPath(), attrsJoined)
			} else {
				txt = fmt.Sprintf(`<script type="application/javascript" src="%v/plugins/%v/%v?v=%v"%v></script>`, urlprefix, v.PluginName(), file.SanitizedPath(), v.PluginVersion(), attrsJoined)
			}
		}
	case AssetGeneric:
		if file.Inline {
			ff, status, err := file.Contents(assets)
			if status != http.StatusOK {
				if err != nil {
					logger.Error("plugin injector: error getting file contents: %v %v", status, err.Error())
				} else {
					logger.Error("plugin injector: error getting file contents: %v", status)
				}
				return ""
			}

			if file.TagName == "" {
				if debug {
					txt = fmt.Sprintf(`<span%v data-amblocation="start"></span>%v<span%v data-amblocation="end"></span>`, attrsJoined, string(ff), attrsJoined)
				} else {
					txt = fmt.Sprintf(`%v`, string(ff))
				}
			} else {
				txt = fmt.Sprintf(`<%v%v>%v</%v>`, html.EscapeString(file.TagName), attrsJoined, string(ff), html.EscapeString(file.TagName))
			}
		} else {
			if file.ClosingTag {
				txt = fmt.Sprintf(`<%v%v></%v>`, html.EscapeString(file.TagName), attrsJoined, html.EscapeString(file.TagName))
			} else {
				txt = fmt.Sprintf(`<%v%v>`, html.EscapeString(file.TagName), attrsJoined)
			}
		}
	default:
		logger.Error("plugin injector: unsupported asset filetype for plugin (%v): %v", v.PluginName(), file.Filetype)
	}

	return txt
}

// Contents returns the text of the file to inline in HTML after doing replace.
func (file *Asset) Contents(assets fs.FS) (ff []byte, status int, err error) {
	// Get the contents from the path if the content field is not filled in.
	if len(file.Path) > 0 {
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
	} else {
		ff = []byte(file.Content)
	}

	// Loop over the items to replace.
	for _, rep := range file.Replace {
		ff = bytes.ReplaceAll(ff, []byte(rep.Find), []byte(rep.Replace))
	}

	return ff, http.StatusOK, nil
}
