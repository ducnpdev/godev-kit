package vietqr

import (
	"context"
	"errors"
	"testing"

	"github.com/ducnpdev/godev-kit/internal/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// contextKey is a type for context keys to avoid collisions
type contextKey string

const testContextKey contextKey = "test-key"

// MockVietQRRepo is a mock implementation of vietqr.VietQRRepo
type MockVietQRRepo struct {
	mock.Mock
}

func (m *MockVietQRRepo) GenerateQR(ctx context.Context, req entity.VietQRGenerateRequest) (string, error) {
	args := m.Called(ctx, req)
	return args.String(0), args.Error(1)
}

// MockVietQRPersistentRepo is a mock implementation of VietQRPersistentRepo
type MockVietQRPersistentRepo struct {
	mock.Mock
}

func (m *MockVietQRPersistentRepo) Store(ctx context.Context, qr entity.VietQR) error {
	args := m.Called(ctx, qr)
	return args.Error(0)
}

func (m *MockVietQRPersistentRepo) FindByID(ctx context.Context, id string) (entity.VietQR, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return entity.VietQR{}, args.Error(1)
	}
	return args.Get(0).(entity.VietQR), args.Error(1)
}

func (m *MockVietQRPersistentRepo) UpdateStatus(ctx context.Context, id string, status entity.VietQRStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func TestVietQRUseCase_GenerateQR(t *testing.T) {
	tests := []struct {
		name           string
		request        entity.VietQRGenerateRequest
		mockSetup      func(*MockVietQRRepo, *MockVietQRPersistentRepo)
		expectedError  error
		validateResult func(*testing.T, *entity.VietQR)
		validateError  func(*testing.T, error)
	}{
		{
			name: "success - generate QR code",
			request: entity.VietQRGenerateRequest{
				AccountNo:    "1234567890",
				Amount:       "100000",
				Description:  "Payment for services",
				MCC:          "5812",
				ReceiverName: "John Doe",
			},
			mockSetup: func(mockRepo *MockVietQRRepo, mockPersistent *MockVietQRPersistentRepo) {
				mockRepo.On("GenerateQR", mock.Anything, mock.MatchedBy(func(req entity.VietQRGenerateRequest) bool {
					return req.AccountNo == "1234567890" && req.Amount == "100000"
				})).Return("00020101021238570010A00000072701270006ACB00000000000000005204000053037045802VN62120808123456790608Payment63040", nil)

				mockPersistent.On("Store", mock.Anything, mock.MatchedBy(func(qr entity.VietQR) bool {
					return qr.Status == entity.VietQRStatusGenerated &&
						qr.Content == "00020101021238570010A00000072701270006ACB00000000000000005204000053037045802VN62120808123456790608Payment63040" &&
						qr.ID != ""
				})).Return(nil)
			},
			expectedError: nil,
			validateResult: func(t *testing.T, qr *entity.VietQR) {
				assert.NotNil(t, qr)
				assert.NotEmpty(t, qr.ID)
				// Validate UUID format
				_, err := uuid.Parse(qr.ID)
				assert.NoError(t, err)
				assert.Equal(t, entity.VietQRStatusGenerated, qr.Status)
				assert.Equal(t, "00020101021238570010A00000072701270006ACB00000000000000005204000053037045802VN62120808123456790608Payment63040", qr.Content)
			},
		},
		{
			name: "success - generate QR with minimal data",
			request: entity.VietQRGenerateRequest{
				AccountNo:    "9876543210",
				Amount:       "50000",
				Description:  "",
				MCC:          "",
				ReceiverName: "",
			},
			mockSetup: func(mockRepo *MockVietQRRepo, mockPersistent *MockVietQRPersistentRepo) {
				mockRepo.On("GenerateQR", mock.Anything, mock.Anything).Return("QR_CONTENT_123", nil)
				mockPersistent.On("Store", mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: nil,
			validateResult: func(t *testing.T, qr *entity.VietQR) {
				assert.NotNil(t, qr)
				assert.NotEmpty(t, qr.ID)
				assert.Equal(t, entity.VietQRStatusGenerated, qr.Status)
				assert.Equal(t, "QR_CONTENT_123", qr.Content)
			},
		},
		{
			name: "error - external API repo error",
			request: entity.VietQRGenerateRequest{
				AccountNo:    "1234567890",
				Amount:       "100000",
				Description:  "Payment",
				MCC:          "5812",
				ReceiverName: "John Doe",
			},
			mockSetup: func(mockRepo *MockVietQRRepo, mockPersistent *MockVietQRPersistentRepo) {
				mockRepo.On("GenerateQR", mock.Anything, mock.Anything).Return("", errors.New("external API error: connection timeout"))
			},
			expectedError: errors.New("external API error: connection timeout"),
			validateError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "external API error")
			},
		},
		{
			name: "error - persistent repo store error",
			request: entity.VietQRGenerateRequest{
				AccountNo:    "1234567890",
				Amount:       "100000",
				Description:  "Payment",
				MCC:          "5812",
				ReceiverName: "John Doe",
			},
			mockSetup: func(mockRepo *MockVietQRRepo, mockPersistent *MockVietQRPersistentRepo) {
				mockRepo.On("GenerateQR", mock.Anything, mock.Anything).Return("QR_CONTENT_456", nil)
				mockPersistent.On("Store", mock.Anything, mock.Anything).Return(errors.New("database error: connection failed"))
			},
			expectedError: errors.New("database error: connection failed"),
			validateError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "database error")
			},
		},
		{
			name: "error - context deadline exceeded in external API",
			request: entity.VietQRGenerateRequest{
				AccountNo:    "1234567890",
				Amount:       "100000",
				Description:  "Payment",
				MCC:          "5812",
				ReceiverName: "John Doe",
			},
			mockSetup: func(mockRepo *MockVietQRRepo, mockPersistent *MockVietQRPersistentRepo) {
				mockRepo.On("GenerateQR", mock.Anything, mock.Anything).Return("", context.DeadlineExceeded)
			},
			expectedError: context.DeadlineExceeded,
			validateError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, context.DeadlineExceeded))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			ctx := context.Background()
			mockRepo := new(MockVietQRRepo)
			mockPersistent := new(MockVietQRPersistentRepo)
			tt.mockSetup(mockRepo, mockPersistent)

			useCase := NewVietQRUseCase(mockRepo, mockPersistent)

			// Execute
			result, err := useCase.GenerateQR(ctx, tt.request)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
				if tt.validateError != nil {
					tt.validateError(t, err)
				} else {
					assert.Equal(t, tt.expectedError.Error(), err.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tt.validateResult != nil {
					tt.validateResult(t, result)
				}
			}

			// Verify mock expectations
			mockRepo.AssertExpectations(t)
			mockPersistent.AssertExpectations(t)
		})
	}
}

func TestVietQRUseCase_InquiryQR(t *testing.T) {
	tests := []struct {
		name           string
		qrID           string
		mockSetup      func(*MockVietQRPersistentRepo)
		expectedError  error
		validateResult func(*testing.T, *entity.VietQR)
		validateError  func(*testing.T, error)
	}{
		{
			name: "success - inquiry QR by ID",
			qrID: "550e8400-e29b-41d4-a716-446655440000",
			mockSetup: func(mockPersistent *MockVietQRPersistentRepo) {
				mockPersistent.On("FindByID", mock.Anything, "550e8400-e29b-41d4-a716-446655440000").Return(
					entity.VietQR{
						ID:      "550e8400-e29b-41d4-a716-446655440000",
						Status:  entity.VietQRStatusGenerated,
						Content: "QR_CONTENT_789",
					},
					nil,
				)
			},
			expectedError: nil,
			validateResult: func(t *testing.T, qr *entity.VietQR) {
				assert.NotNil(t, qr)
				assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", qr.ID)
				assert.Equal(t, entity.VietQRStatusGenerated, qr.Status)
				assert.Equal(t, "QR_CONTENT_789", qr.Content)
			},
		},
		{
			name: "success - inquiry QR with paid status",
			qrID: "test-qr-id-123",
			mockSetup: func(mockPersistent *MockVietQRPersistentRepo) {
				mockPersistent.On("FindByID", mock.Anything, "test-qr-id-123").Return(
					entity.VietQR{
						ID:      "test-qr-id-123",
						Status:  entity.VietQRStatusPaid,
						Content: "PAID_QR_CONTENT",
					},
					nil,
				)
			},
			expectedError: nil,
			validateResult: func(t *testing.T, qr *entity.VietQR) {
				assert.NotNil(t, qr)
				assert.Equal(t, "test-qr-id-123", qr.ID)
				assert.Equal(t, entity.VietQRStatusPaid, qr.Status)
				assert.Equal(t, "PAID_QR_CONTENT", qr.Content)
			},
		},
		{
			name: "error - QR not found",
			qrID: "non-existent-id",
			mockSetup: func(mockPersistent *MockVietQRPersistentRepo) {
				mockPersistent.On("FindByID", mock.Anything, "non-existent-id").Return(
					entity.VietQR{},
					errors.New("record not found"),
				)
			},
			expectedError: errors.New("record not found"),
			validateError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "record not found")
			},
		},
		{
			name: "error - database connection error",
			qrID: "test-id",
			mockSetup: func(mockPersistent *MockVietQRPersistentRepo) {
				mockPersistent.On("FindByID", mock.Anything, "test-id").Return(
					entity.VietQR{},
					errors.New("database connection failed"),
				)
			},
			expectedError: errors.New("database connection failed"),
			validateError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "database connection failed")
			},
		},
		{
			name: "error - context deadline exceeded",
			qrID: "test-id",
			mockSetup: func(mockPersistent *MockVietQRPersistentRepo) {
				mockPersistent.On("FindByID", mock.Anything, "test-id").Return(
					entity.VietQR{},
					context.DeadlineExceeded,
				)
			},
			expectedError: context.DeadlineExceeded,
			validateError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, context.DeadlineExceeded))
			},
		},
		{
			name: "error - context cancelled",
			qrID: "test-id",
			mockSetup: func(mockPersistent *MockVietQRPersistentRepo) {
				mockPersistent.On("FindByID", mock.Anything, "test-id").Return(
					entity.VietQR{},
					context.Canceled,
				)
			},
			expectedError: context.Canceled,
			validateError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, context.Canceled))
			},
		},
		{
			name: "success - inquiry QR with all status types",
			qrID: "status-test-id",
			mockSetup: func(mockPersistent *MockVietQRPersistentRepo) {
				mockPersistent.On("FindByID", mock.Anything, "status-test-id").Return(
					entity.VietQR{
						ID:      "status-test-id",
						Status:  entity.VietQRStatusInProcess,
						Content: "IN_PROCESS_QR",
					},
					nil,
				)
			},
			expectedError: nil,
			validateResult: func(t *testing.T, qr *entity.VietQR) {
				assert.NotNil(t, qr)
				assert.Equal(t, entity.VietQRStatusInProcess, qr.Status)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			ctx := context.Background()
			mockRepo := new(MockVietQRRepo)
			mockPersistent := new(MockVietQRPersistentRepo)
			tt.mockSetup(mockPersistent)

			useCase := NewVietQRUseCase(mockRepo, mockPersistent)

			// Execute
			result, err := useCase.InquiryQR(ctx, tt.qrID)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
				if tt.validateError != nil {
					tt.validateError(t, err)
				} else {
					assert.Contains(t, err.Error(), tt.expectedError.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tt.validateResult != nil {
					tt.validateResult(t, result)
				}
			}

			// Verify mock expectations
			mockPersistent.AssertExpectations(t)
		})
	}
}

func TestVietQRUseCase_GenerateQR_Context(t *testing.T) {
	t.Run("verify context is passed to repositories", func(t *testing.T) {
		// Setup
		ctx := context.WithValue(context.Background(), testContextKey, "test-value")
		mockRepo := new(MockVietQRRepo)
		mockPersistent := new(MockVietQRPersistentRepo)

		request := entity.VietQRGenerateRequest{
			AccountNo:    "1234567890",
			Amount:       "100000",
			Description:  "Test",
			MCC:          "5812",
			ReceiverName: "Test User",
		}

		mockRepo.On("GenerateQR", mock.MatchedBy(func(c context.Context) bool {
			return c.Value(testContextKey) == "test-value"
		}), mock.Anything).Return("QR_CONTENT", nil)

		mockPersistent.On("Store", mock.MatchedBy(func(c context.Context) bool {
			return c.Value(testContextKey) == "test-value"
		}), mock.Anything).Return(nil)

		useCase := NewVietQRUseCase(mockRepo, mockPersistent)

		// Execute
		result, err := useCase.GenerateQR(ctx, request)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		mockRepo.AssertExpectations(t)
		mockPersistent.AssertExpectations(t)
	})
}

func TestVietQRUseCase_InquiryQR_Context(t *testing.T) {
	t.Run("verify context is passed to repository", func(t *testing.T) {
		// Setup
		ctx := context.WithValue(context.Background(), testContextKey, "test-value")
		mockRepo := new(MockVietQRRepo)
		mockPersistent := new(MockVietQRPersistentRepo)

		mockPersistent.On("FindByID", mock.MatchedBy(func(c context.Context) bool {
			return c.Value(testContextKey) == "test-value"
		}), "test-id").Return(entity.VietQR{
			ID:      "test-id",
			Status:  entity.VietQRStatusGenerated,
			Content: "TEST_CONTENT",
		}, nil)

		useCase := NewVietQRUseCase(mockRepo, mockPersistent)

		// Execute
		result, err := useCase.InquiryQR(ctx, "test-id")

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "test-id", result.ID)
		mockPersistent.AssertExpectations(t)
	})
}
