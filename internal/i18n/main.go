package i18n

import (
	_ "embed"
	"encoding/json"
	"log"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locales/en.json
var enJSON []byte

var bundle *i18n.Bundle

func Init() {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	if _, err := bundle.ParseMessageFileBytes(enJSON, "en.json"); err != nil {
		log.Fatalf("failed to load en.json: %v", err)
	}
}

func getLocalizer(lang string) *i18n.Localizer {
	return i18n.NewLocalizer(bundle, lang, "en")
}

func GetText(id string) string {
	localizer := getLocalizer("en")
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    id,
		TemplateData: map[string]interface{}{},
	})

	if err != nil {
		return id
	}

	return msg
}
