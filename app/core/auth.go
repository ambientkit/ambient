package core

// authAssetAllowed return true if the user has access to the asset.
func authAssetAllowed(loggedIn bool, f Asset) bool {
	switch true {
	case f.Auth == AuthenticatedOnly && !loggedIn:
		return false
	case f.Auth == AuthenticatedOnly && loggedIn:
		return true
	case f.Auth == AnonymousOnly && !loggedIn:
		return true
	case f.Auth == AnonymousOnly && loggedIn:
		return false
	}

	//f.Auth == All:
	return true
}
