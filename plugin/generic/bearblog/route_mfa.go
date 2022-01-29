package bearblog

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/ambientkit/ambient/plugin/generic/bearblog/lib/totp"
	qrcode "github.com/skip2/go-qrcode"
)

func (p *Plugin) mfa(w http.ResponseWriter, r *http.Request) (status int, err error) {
	vars := make(map[string]interface{})
	vars["title"] = "MFA Generate"
	return p.Render.Page(w, r, assets, "template/content/mfa", p.funcMap(r), vars)
}

func (p *Plugin) mfaPost(w http.ResponseWriter, r *http.Request) (status int, err error) {
	vars := make(map[string]interface{})
	vars["title"] = "MFA Generate"

	username := r.FormValue("username")
	issuer := r.FormValue(("issuer"))

	// Generate a MFA.
	URI, secret, err := totp.GenerateURL(username, issuer)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Generate PDF.
	var png []byte
	png, err = qrcode.Encode(URI, qrcode.Medium, 400)
	if err != nil {
		log.Fatalln(err.Error())
	}

	vars["mfa"] = fmt.Sprintf("The secret you can paste into the settings screen is: %v. The URI is: %v", secret, URI)
	vars["qrcode"] = base64.StdEncoding.EncodeToString(png)
	return p.Render.Page(w, r, assets, "template/content/mfa", p.funcMap(r), vars)
}
