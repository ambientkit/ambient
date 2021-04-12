package sitemap

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

// Returns a page for web crawlers.
func (p *Plugin) index(w http.ResponseWriter, r *http.Request) (status int, err error) {
	// Resource: https://www.sitemaps.org/protocol.html
	// Resource: https://golang.org/src/encoding/xml/example_test.go

	type URL struct {
		Location     string `xml:"loc"`
		LastModified string `xml:"lastmod"`
	}

	type Sitemap struct {
		XMLName xml.Name `xml:"urlset"`
		XMLNS   string   `xml:"xmlns,attr"`
		XHTML   string   `xml:"xmlns:xhtml,attr"`
		URL     []URL    `xml:"url"`
	}

	m := &Sitemap{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		XHTML: "http://www.w3.org/1999/xhtml",
	}

	siteURL, err := p.Site.FullURL()
	if err != nil {
		return p.Site.Error(err)
	}

	siteUpdated, err := p.Site.Updated()
	if err != nil {
		return p.Site.Error(err)
	}

	// Home page
	m.URL = append(m.URL, URL{
		Location:     siteURL,
		LastModified: siteUpdated.Format("2006-01-02"),
	})

	// Posts and pages
	postsAndPages, err := p.Site.PostsAndPages(true)
	if err != nil {
		return p.Site.Error(err)
	}
	for _, v := range postsAndPages {
		m.URL = append(m.URL, URL{
			Location:     siteURL + "/" + v.URL,
			LastModified: v.Timestamp.Format("2006-01-02"),
		})
	}

	// Tags
	tags, err := p.Site.Tags(true)
	if err != nil {
		return p.Site.Error(err)
	}
	for _, v := range tags {
		m.URL = append(m.URL, URL{
			Location:     siteURL + "/blog?q=" + v.Name,
			LastModified: v.Timestamp.Format("2006-01-02"),
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
