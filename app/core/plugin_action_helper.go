package core

import (
	"embed"
	"io/fs"
)

// fieldArrayEqual tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
func fieldArrayEqual(pa PluginSettings, arrB []Field) bool {
	a := pa.Fields
	if len(pa.Fields) != len(arrB) {
		return false
	}

	// Verify the length of the order field.
	if len(pa.Order) != len(arrB) {
		return false
	}

	b := make(map[string]Field)
	for index, v := range arrB {
		// Verify the order of the order field.
		if pa.Order[index] != v.Name {
			return false
		}

		b[v.Name] = v
	}

	for i, v := range a {
		if v.Name != b[v.Name].Name {
			return false
		}
		if string(v.Type) != string(b[i].Type) {
			return false
		}
		if v.Description.Text != b[i].Description.Text {
			return false
		}
		if v.Description.URL != b[i].Description.URL {
			return false
		}
		if v.Default != b[i].Default {
			return false
		}
	}
	return true
}

// stringArrayEqual tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
func stringArrayEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// fileExists determines if an embedded file exists.
func fileExists(assets *embed.FS, filename string) bool {
	// Use the root directory.
	fsys, err := fs.Sub(assets, ".")
	if err != nil {
		return false
	}

	// Open the file.
	f, err := fsys.Open(filename)
	if err != nil {
		return false
	}
	defer f.Close()

	return true
}

// authAssetAllowed return true if the user has access to the asset.
func authAssetAllowed(loggedIn bool, f Asset) bool {
	switch true {
	case f.Auth == AuthOnly && !loggedIn:
		return false
	case f.Auth == AuthOnly && loggedIn:
		return true
	case f.Auth == AuthAnonymousOnly && !loggedIn:
		return true
	case f.Auth == AuthAnonymousOnly && loggedIn:
		return false
	}

	//f.Auth == All:
	return true
}
