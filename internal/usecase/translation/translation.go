package translation

import (
	"context"
	"fmt"

	"github.com/ducnpdev/godev-kit/internal/entity"
	"github.com/ducnpdev/godev-kit/internal/repo"
	"github.com/ducnpdev/godev-kit/internal/repo/persistent/models"
	"github.com/ducnpdev/godev-kit/pkg/logger"
)

// UseCase -.
type UseCase struct {
	logger logger.Interface
	//
	repo   repo.TranslationRepo
	webAPI repo.TranslationWebAPI
}

// New -.
func New(
	r repo.TranslationRepo,
	w repo.TranslationWebAPI,
) *UseCase {
	return &UseCase{
		repo:   r,
		webAPI: w,
	}
}

// History - getting translate history from store.
func (uc *UseCase) History(ctx context.Context) (entity.TranslationHistory, error) {
	translations, err := uc.repo.GetHistory(ctx)
	if err != nil {
		return entity.TranslationHistory{}, fmt.Errorf("TranslationUseCase - History - s.repo.GetHistory: %w", err)
	}

	return entity.TranslationHistory{History: translations}, nil
}

// Translate -.
func (uc *UseCase) Translate(ctx context.Context, t entity.Translation) (entity.Translation, error) {
	translation, err := uc.webAPI.Translate(t)
	if err != nil {
		return entity.Translation{}, fmt.Errorf("TranslationUseCase - Translate - s.webAPI.Translate: %w", err)
	}

	err = uc.repo.Store(ctx, models.TranslationModel{
		Source:      t.Source,
		Destination: t.Destination,
		Original:    t.Original,
		Translation: t.Translation,
	})
	if err != nil {
		return entity.Translation{}, fmt.Errorf("TranslationUseCase - Translate - s.repo.Store: %w", err)
	}

	return translation, nil
}
