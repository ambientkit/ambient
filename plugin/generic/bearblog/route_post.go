package bearblog

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/ambientkit/ambient"
	"github.com/microcosm-cc/bluemonday"
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

		posts := make([]ambient.PostWithID, 0)
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

	return p.Render.Page(w, r, assets, "template/content/bloglist_index", p.funcMap(r), vars)
}

func (p *Plugin) postShow(w http.ResponseWriter, r *http.Request) (status int, err error) {
	slug := p.Mux.Param(r, "slug")

	post, err := p.Site.PostBySlug(slug)
	if err != nil {
		return p.Site.Error(err)
	}

	// Determine if in preview mode.
	preview := false
	if q := r.URL.Query().Get("preview"); len(q) > 0 && strings.ToLower(q) == "true" {
		preview = true
	}

	// Show 404 if not published and not in preview mode.
	if !post.Published && !preview {
		return http.StatusNotFound, nil
	}

	vars := make(map[string]interface{})
	// Don't show certain items on pages.
	if !post.Page {
		vars["title"] = post.Title
		vars["pubdate"] = post.Timestamp
	}

	vars["tags"] = post.Tags
	vars["canonical"] = post.Canonical
	vars["id"] = post.ID
	vars["posturl"] = post.URL
	vars["pagetitle"] = post.Title
	vars["pagedescription"] = plaintextBlurb(post.Content)
	vars["postcontent"] = p.sanitized(post.Content)

	return p.Render.Post(w, r, assets, "template/content/post", p.funcMap(r), vars)
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

// sanitized returns a sanitized content block or an error is one occurs.
func (p *Plugin) sanitized(content string) template.HTML {
	b := []byte(content)
	// Ensure unit line endings are used when pulling out of JSON.
	markdownWithUnixLineEndings := strings.Replace(string(b), "\r\n", "\n", -1)
	htmlCode := blackfriday.Run([]byte(markdownWithUnixLineEndings))

	// Determine if raw HTML is allowed.
	allowed, err := p.Site.PluginSettingBool(AllowHTMLinMarkdown)
	if err != nil {
		p.Log.Debug("plugins: error in sanitized() getting plugin field: %v", err)
	}

	// Sanitize by removing HTML if allowed.
	if !allowed {
		htmlCode = bluemonday.UGCPolicy().SanitizeBytes(htmlCode)
	}

	return template.HTML(htmlCode)
}
