package internal

import "github.com/c-bata/go-prompt"

func filterString(list []string, val string) []string {
	// Don't filter out if it's the only element.
	if len(list) == 1 && list[0] == val {
		return list
	}

	for k, v := range list {
		if v == val {
			list = append(list[:k], list[k+1:]...)
		}
	}

	return list
}

func filterAlreadyUsed(list []prompt.Suggest, args []string) []prompt.Suggest {
	for _, typed := range args {
		// Skip words if they don't start with a -- because they may be just a typed input.
		// Not sure if we need this or not.
		// if !strings.HasPrefix(typed, "-") {
		// 	continue
		// }

		for k, autocomplete := range list {
			if typed == autocomplete.Text {
				list = append(list[:k], list[k+1:]...)
			}
		}
	}

	return list
}
