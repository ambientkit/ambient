package ambient

import (
	"html/template"
	"os"
)

// GlobalFuncMap returns the FuncMaps available in all templates.
func GlobalFuncMap(fm template.FuncMap) template.FuncMap {
	if fm == nil {
		fm = template.FuncMap{}
	}

	fm["URLPrefix"] = func() string {
		return os.Getenv("AMB_URL_PREFIX")
	}
	fm["TrustHTML"] = func(s string) template.HTML {
		return template.HTML(s)
	}
	fm["TrustHTMLAttr"] = func(s string) template.HTMLAttr {
		return template.HTMLAttr(s)
	}
	fm["TrustJS"] = func(s string) template.JS {
		return template.JS(s)
	}
	fm["TrustJSStr"] = func(s string) template.JSStr {
		return template.JSStr(s)
	}
	fm["TrustSrcset"] = func(s string) template.Srcset {
		return template.Srcset(s)
	}

	return fm
}
