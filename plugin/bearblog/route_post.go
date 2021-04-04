package bearblog

import (
	"net/http"
	"strings"

	"github.com/josephspurrier/ambient/app/model"
	"github.com/russross/blackfriday/v2"
	"jaytaylor.com/html2text"
)

func (p *Plugin) postIndex(w http.ResponseWriter, r *http.Request) (status int, err error) {
	vars := make(map[string]interface{})

	tags, err := p.Site.Tags(true)
	if err != nil {
		return p.Site.Error(err)
	}
	vars["tags"] = tags

	// Determine if there is query.
	if q := r.URL.Query().Get("q"); len(q) > 0 {
		vars["query"] = q
		// Don't show tags when there is a filter.
		delete(vars, "tags")

		postsAndPages, err := p.Site.PostsAndPages(true)
		if err != nil {
			return p.Site.Error(err)
		}

		posts := make([]model.PostWithID, 0)
		for _, v := range postsAndPages {
			match := false
			for _, tag := range v.Tags {
				if tag.Name == q {
					match = true
					break
				}
			}

			if match {
				posts = append(posts, v)
			}
		}

		vars["posts"] = posts
	} else {

		pubPosts, err := p.Site.PublishedPosts()
		if err != nil {
			return p.Site.Error(err)
		}

		vars["posts"] = pubPosts
	}

	return p.Render.PluginPage(w, r, assets, "template/bloglist_index", p.FuncMap(r), vars)
}

// func (p *Plugin) postShow(w http.ResponseWriter, r *http.Request) (status int, err error) {
// 	slug := c.Router.Param(r, "slug")
// 	p := c.Storage.Site.PostBySlug(slug)

// 	// Determine if in preview mode.
// 	preview := false
// 	if q := r.URL.Query().Get("preview"); len(q) > 0 && strings.ToLower(q) == "true" {
// 		preview = true
// 	}

// 	// Show 404 if not published and not in preview mode.
// 	if !p.Published && !preview {
// 		return http.StatusNotFound, nil
// 	}

// 	vars := make(map[string]interface{})
// 	// Don't show certain items on pages.
// 	if !p.Page {
// 		vars["title"] = p.Title
// 		vars["pubdate"] = p.Timestamp
// 	}

// 	vars["tags"] = p.Tags
// 	vars["canonical"] = p.Canonical
// 	vars["id"] = p.ID
// 	vars["posturl"] = p.URL
// 	vars["metadescription"] = plaintextBlurb(p.Content)

// 	return c.Render.Post(w, r, p.Post.Content, vars)
// }

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
