package helper

import "github.com/c-bata/go-prompt"

var createAppSuggest = SmartSuggestGroup{
	{Suggest: prompt.Suggest{Text: "--folder", Description: "Folder to create the project. Ex: . or ./newdir"}, Required: true},
	{Suggest: prompt.Suggest{Text: "--template", Description: "Template project to use. Ex: default or github.com/josephspurrier/template"}, Required: true},
}
