package rssfeed

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/russross/blackfriday/v2"
	"jaytaylor.com/html2text"
)

// Returns a page for web crawlers.
func (p *Plugin) index(w http.ResponseWriter, r *http.Request) (status int, err error) {
	// Resource: https://www.rssboard.org/rss-specification
	// Rsource: https://validator.w3.org/feed/check.cgi

	type Item struct {
		Title       string `xml:"title"`
		Link        string `xml:"link"`
		PubDate     string `xml:"pubDate"`
		GUID        string `xml:"guid"`
		Description string `xml:"description"`
	}

	type AtomLink struct {
		Href string `xml:"href,attr"`
		Rel  string `xml:"rel,attr"`
		Type string `xml:"type,attr"`
	}

	type Sitemap struct {
		XMLName       xml.Name `xml:"rss"`
		Version       string   `xml:"version,attr"`
		Atom          string   `xml:"xmlns:atom,attr"`
		Title         string   `xml:"channel>title"`
		Link          string   `xml:"channel>link"`
		Description   string   `xml:"channel>description"`
		Generator     string   `xml:"channel>generator"`
		Language      string   `xml:"channel>language"`
		LastBuildDate string   `xml:"channel>lastBuildDate"`
		AtomLink      AtomLink `xml:"channel>atom:link"`
		Items         []Item   `xml:"channel>item"`
	}

	title, err := p.Site.Title()
	if err != nil {
		return p.Site.Error(err)
	}

	siteURL, err := p.Site.FullURL()
	if err != nil {
		return p.Site.Error(err)
	}

	description, err := p.Site.PluginSettingString(Description)
	if err != nil {
		return p.Site.Error(err)
	}

	feedURL, err := p.Site.PluginSettingString(FeedURL)
	if err != nil {
		return p.Site.Error(err)
	}

	m := &Sitemap{
		Version:       "2.0",
		Atom:          "http://www.w3.org/2005/Atom",
		Title:         title,
		Link:          siteURL,
		Description:   description,
		Generator:     "Ambient",
		Language:      "en-us",
		LastBuildDate: time.Now().Format(time.RFC1123Z),
		AtomLink: AtomLink{
			Href: path.Join(siteURL, feedURL),
			Rel:  "self",
			Type: "application/rss+xml",
		},
	}

	postAndPages, err := p.Site.PostsAndPages(true)
	if err != nil {
		return p.Site.Error(err)
	}

	for _, v := range postAndPages {
		plaintext := plaintextBlurb(v.Post.Content)
		m.Items = append(m.Items, Item{
			Title:       v.Title,
			Link:        siteURL + "/" + v.URL,
			PubDate:     v.Timestamp.Format(time.RFC1123Z),
			GUID:        siteURL + "/" + v.URL,
			Description: plaintext,
		})
	}

	output, err := xml.MarshalIndent(m, "  ", "    ")
	if err != nil {
		return http.StatusInternalServerError, err
	}

	header := []byte(xml.Header)
	output = append(header[:], output[:]...)

	w.Header().Set("Content-Type", "application/xml")
	fmt.Fprint(w, string(output))
	return
}

// plaintextBlurb returns a plaintext blurb from markdown content.
func plaintextBlurb(s string) string {
	unsafeHTML := blackfriday.Run([]byte(s))
	plaintext, err := html2text.FromString(string(unsafeHTML))
	if err != nil {
		plaintext = s
	}
	period := strings.Index(plaintext, ". ")
	if period > 0 {
		plaintext = plaintext[:period+1]
	}

	return plaintext
}
