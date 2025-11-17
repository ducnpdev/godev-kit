package translation

import (
	"context"
	"errors"
	"testing"

	"github.com/ducnpdev/godev-kit/internal/entity"
	"github.com/ducnpdev/godev-kit/internal/repo/persistent/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTranslationRepo is a mock implementation of repo.TranslationRepo
type MockTranslationRepo struct {
	mock.Mock
}

func (m *MockTranslationRepo) Store(ctx context.Context, model models.TranslationModel) error {
	args := m.Called(ctx, model)
	return args.Error(0)
}

func (m *MockTranslationRepo) GetHistory(ctx context.Context) ([]entity.Translation, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Translation), args.Error(1)
}

func TestUseCase_History(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockTranslationRepo)
		expectedResult entity.TranslationHistory
		expectedError  error
		validateError  func(*testing.T, error)
	}{
		{
			name: "success - get translation history",
			mockSetup: func(m *MockTranslationRepo) {
				m.On("GetHistory", mock.Anything).Return([]entity.Translation{
					{
						Source:      "auto",
						Destination: "en",
						Original:    "текст для перевода",
						Translation: "text for translation",
					},
					{
						Source:      "en",
						Destination: "ru",
						Original:    "hello world",
						Translation: "привет мир",
					},
				}, nil)
			},
			expectedResult: entity.TranslationHistory{
				History: []entity.Translation{
					{
						Source:      "auto",
						Destination: "en",
						Original:    "текст для перевода",
						Translation: "text for translation",
					},
					{
						Source:      "en",
						Destination: "ru",
						Original:    "hello world",
						Translation: "привет мир",
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "success - empty history",
			mockSetup: func(m *MockTranslationRepo) {
				m.On("GetHistory", mock.Anything).Return([]entity.Translation{}, nil)
			},
			expectedResult: entity.TranslationHistory{
				History: []entity.Translation{},
			},
			expectedError: nil,
		},
		{
			name: "success - single translation in history",
			mockSetup: func(m *MockTranslationRepo) {
				m.On("GetHistory", mock.Anything).Return([]entity.Translation{
					{
						Source:      "auto",
						Destination: "en",
						Original:    "test",
						Translation: "test",
					},
				}, nil)
			},
			expectedResult: entity.TranslationHistory{
				History: []entity.Translation{
					{
						Source:      "auto",
						Destination: "en",
						Original:    "test",
						Translation: "test",
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "error - repository error",
			mockSetup: func(m *MockTranslationRepo) {
				m.On("GetHistory", mock.Anything).Return(nil, errors.New("database connection error"))
			},
			expectedResult: entity.TranslationHistory{},
			expectedError:  errors.New("database connection error"),
			validateError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "TranslationUseCase - History")
				assert.Contains(t, err.Error(), "database connection error")
			},
		},
		{
			name: "error - context deadline exceeded",
			mockSetup: func(m *MockTranslationRepo) {
				m.On("GetHistory", mock.Anything).Return(nil, context.DeadlineExceeded)
			},
			expectedResult: entity.TranslationHistory{},
			expectedError:  context.DeadlineExceeded,
			validateError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "TranslationUseCase - History")
				assert.True(t, errors.Is(err, context.DeadlineExceeded))
			},
		},
		{
			name: "error - context cancelled",
			mockSetup: func(m *MockTranslationRepo) {
				m.On("GetHistory", mock.Anything).Return(nil, context.Canceled)
			},
			expectedResult: entity.TranslationHistory{},
			expectedError:  context.Canceled,
			validateError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "TranslationUseCase - History")
				assert.True(t, errors.Is(err, context.Canceled))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			ctx := context.Background()
			mockRepo := new(MockTranslationRepo)
			tt.mockSetup(mockRepo)

			useCase := New(mockRepo, nil) // webAPI is not used in History function

			// Execute
			result, err := useCase.History(ctx)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				if tt.validateError != nil {
					tt.validateError(t, err)
				} else {
					// Check if error message contains expected error
					assert.Contains(t, err.Error(), tt.expectedError.Error())
				}
				assert.Equal(t, entity.TranslationHistory{}, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
				assert.Len(t, result.History, len(tt.expectedResult.History))

				// Verify all translations match
				for i, expected := range tt.expectedResult.History {
					assert.Equal(t, expected.Source, result.History[i].Source)
					assert.Equal(t, expected.Destination, result.History[i].Destination)
					assert.Equal(t, expected.Original, result.History[i].Original)
					assert.Equal(t, expected.Translation, result.History[i].Translation)
				}
			}

			// Verify mock expectations
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUseCase_History_Context(t *testing.T) {
	t.Run("verify context is passed to repository", func(t *testing.T) {
		// Setup
		ctx := context.WithValue(context.Background(), "test-key", "test-value")
		mockRepo := new(MockTranslationRepo)

		mockRepo.On("GetHistory", mock.MatchedBy(func(c context.Context) bool {
			return c.Value("test-key") == "test-value"
		})).Return([]entity.Translation{}, nil)

		useCase := New(mockRepo, nil)

		// Execute
		_, err := useCase.History(ctx)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}
