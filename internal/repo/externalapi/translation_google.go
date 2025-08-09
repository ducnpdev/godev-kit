package externalapi

import (
	"context"
	"fmt"

	translator "github.com/Conight/go-googletrans"
	"github.com/ducnpdev/godev-kit/internal/entity"
)

// TranslationWebAPI -.
type TranslationWebAPI struct {
	conf translator.Config
}

// New -.
func New() *TranslationWebAPI {
	conf := translator.Config{
		UserAgent:   []string{"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:15.0) Gecko/20100101 Firefox/15.0.1"},
		ServiceUrls: []string{"translate.google.com"},
	}

	return &TranslationWebAPI{
		conf: conf,
	}
}

// Translate - now with context support to prevent memory leaks
func (t *TranslationWebAPI) Translate(ctx context.Context, translation entity.Translation) (entity.Translation, error) {
	// Create a channel to handle the translation result
	resultChan := make(chan struct {
		result entity.Translation
		err    error
	}, 1)

	// Run translation in a goroutine to support cancellation
	go func() {
		trans := translator.New(t.conf)
		result, err := trans.Translate(translation.Original, translation.Source, translation.Destination)

		if err != nil {
			resultChan <- struct {
				result entity.Translation
				err    error
			}{entity.Translation{}, fmt.Errorf("TranslationWebAPI - Translate - trans.Translate: %w", err)}
			return
		}

		translation.Translation = result.Text
		resultChan <- struct {
			result entity.Translation
			err    error
		}{translation, nil}
	}()

	// Wait for either completion or context cancellation
	select {
	case result := <-resultChan:
		return result.result, result.err
	case <-ctx.Done():
		return entity.Translation{}, fmt.Errorf("TranslationWebAPI - Translate - context cancelled: %w", ctx.Err())
	}
}
