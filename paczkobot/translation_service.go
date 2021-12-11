package paczkobot

import (
	"embed"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
)

//go:embed translations/*.json
var translationFiles embed.FS

type TranslationService struct {
	// first key is language code, second key is translation key
	translations map[string]map[string]string
}

func NewTranslationService() *TranslationService {
	return &TranslationService{
		translations: make(map[string]map[string]string),
	}
}

func (ts *TranslationService) ParseTranslationFiles() error {
	files, err := translationFiles.ReadDir("translations")
	if err != nil {
		return fmt.Errorf("could not read translation files dir: %v", err)
	}
	for _, fPath := range files {
		f, err := translationFiles.Open("translations/" + fPath.Name())
		if err != nil {
			return fmt.Errorf("could not open translation file %v: %v", fPath.Name(), err)
		}
		defer f.Close()
		langCode := strings.TrimSuffix(fPath.Name(), filepath.Ext(fPath.Name()))
		ts.translations[langCode] = make(map[string]string)
		decoder := json.NewDecoder(f)
		dict := ts.translations[langCode]
		err = decoder.Decode(&dict)
		if err != nil {
			return fmt.Errorf("could not decode translation file %v: %v", fPath.Name(), err)
		}
	}
	return nil
}
