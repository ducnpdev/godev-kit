package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ducnpdev/godev-kit/internal/controller/http/v1/request"
	"github.com/ducnpdev/godev-kit/internal/entity"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserUseCase is a mock implementation of usecase.User
type MockUserUseCase struct {
	mock.Mock
}

func (m *MockUserUseCase) Create(ctx context.Context, user entity.User) (entity.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return entity.User{}, args.Error(1)
	}
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *MockUserUseCase) GetByID(ctx context.Context, id int64) (entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return entity.User{}, args.Error(1)
	}
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *MockUserUseCase) Update(ctx context.Context, user entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserUseCase) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserUseCase) List(ctx context.Context) (entity.UserHistory, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return entity.UserHistory{}, args.Error(1)
	}
	return args.Get(0).(entity.UserHistory), args.Error(1)
}

func (m *MockUserUseCase) Login(ctx context.Context, email, password string) (string, entity.User, error) {
	args := m.Called(ctx, email, password)
	return args.String(0), args.Get(1).(entity.User), args.Error(2)
}

// MockLogger is a mock implementation of logger.Interface
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debug(message interface{}, args ...interface{}) {
	m.Called(message, args)
}

func (m *MockLogger) Info(message string, args ...interface{}) {
	m.Called(message, args)
}

func (m *MockLogger) Warn(message string, args ...interface{}) {
	m.Called(message, args)
}

func (m *MockLogger) Error(message interface{}, args ...interface{}) {
	m.Called(message, args)
}

func (m *MockLogger) Fatal(message interface{}, args ...interface{}) {
	m.Called(message, args)
}

func (m *MockLogger) Zerolog() interface{} {
	return nil
}

func (m *MockLogger) ZerologPtr() interface{} {
	return nil
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestNewUserRoutes_CreateUser(t *testing.T) {
	tests := []struct {
		name             string
		requestBody      interface{}
		mockSetup        func(*MockUserUseCase)
		expectedStatus   int
		validateResponse func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "success - create user",
			requestBody: request.CreateUser{
				Email:    "test@example.com",
				Username: "testuser",
				Password: "password123",
			},
			mockSetup: func(m *MockUserUseCase) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(u entity.User) bool {
					return u.Email == "test@example.com" && u.Username == "testuser"
				})).Return(entity.User{
					ID:        1,
					Email:     "test@example.com",
					Username:  "testuser",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var user entity.User
				err := json.Unmarshal(w.Body.Bytes(), &user)
				assert.NoError(t, err)
				assert.Equal(t, int64(1), user.ID)
				assert.Equal(t, "test@example.com", user.Email)
				assert.Equal(t, "testuser", user.Username)
			},
		},
		{
			name: "error - invalid request body",
			requestBody: map[string]interface{}{
				"email": "invalid-email",
			},
			mockSetup:      func(m *MockUserUseCase) {},
			expectedStatus: http.StatusBadRequest,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
		{
			name: "error - validation failure - missing email",
			requestBody: request.CreateUser{
				Username: "testuser",
				Password: "password123",
			},
			mockSetup:      func(m *MockUserUseCase) {},
			expectedStatus: http.StatusBadRequest,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
		{
			name: "error - service error",
			requestBody: request.CreateUser{
				Email:    "test@example.com",
				Username: "testuser",
				Password: "password123",
			},
			mockSetup: func(m *MockUserUseCase) {
				m.On("Create", mock.Anything, mock.Anything).Return(entity.User{}, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			router := setupRouter()
			mockUserUseCase := new(MockUserUseCase)
			mockLogger := new(MockLogger)
			// Allow any logger calls during tests
			mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()
			mockLogger.On("Warn", mock.Anything, mock.Anything).Return()
			mockLogger.On("Info", mock.Anything, mock.Anything).Return()
			tt.mockSetup(mockUserUseCase)

			apiV1Group := router.Group("/v1")
			NewUserRoutes(apiV1Group, mockUserUseCase, mockLogger)

			// Create request
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/v1/user", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Execute
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.validateResponse != nil {
				tt.validateResponse(t, w)
			}
			mockUserUseCase.AssertExpectations(t)
		})
	}
}

func TestNewUserRoutes_GetUser(t *testing.T) {
	tests := []struct {
		name             string
		userID           string
		mockSetup        func(*MockUserUseCase)
		expectedStatus   int
		validateResponse func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:   "success - get user",
			userID: "1",
			mockSetup: func(m *MockUserUseCase) {
				m.On("GetByID", mock.Anything, int64(1)).Return(entity.User{
					ID:        1,
					Email:     "test@example.com",
					Username:  "testuser",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil)
			},
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var user entity.User
				err := json.Unmarshal(w.Body.Bytes(), &user)
				assert.NoError(t, err)
				assert.Equal(t, int64(1), user.ID)
				assert.Equal(t, "test@example.com", user.Email)
			},
		},
		{
			name:           "error - invalid user ID",
			userID:         "invalid",
			mockSetup:      func(m *MockUserUseCase) {},
			expectedStatus: http.StatusBadRequest,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
		{
			name:   "error - user not found",
			userID: "999",
			mockSetup: func(m *MockUserUseCase) {
				m.On("GetByID", mock.Anything, int64(999)).Return(entity.User{ID: 0}, nil)
			},
			expectedStatus: http.StatusNotFound,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
		{
			name:   "error - service error",
			userID: "1",
			mockSetup: func(m *MockUserUseCase) {
				m.On("GetByID", mock.Anything, int64(1)).Return(entity.User{}, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			router := setupRouter()
			mockUserUseCase := new(MockUserUseCase)
			mockLogger := new(MockLogger)
			// Allow any logger calls during tests
			mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()
			mockLogger.On("Warn", mock.Anything, mock.Anything).Return()
			mockLogger.On("Info", mock.Anything, mock.Anything).Return()
			tt.mockSetup(mockUserUseCase)

			apiV1Group := router.Group("/v1")
			NewUserRoutes(apiV1Group, mockUserUseCase, mockLogger)

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/v1/user/"+tt.userID, nil)
			w := httptest.NewRecorder()

			// Execute
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.validateResponse != nil {
				tt.validateResponse(t, w)
			}
			mockUserUseCase.AssertExpectations(t)
		})
	}
}

func TestNewUserRoutes_ListUsers(t *testing.T) {
	tests := []struct {
		name             string
		mockSetup        func(*MockUserUseCase)
		expectedStatus   int
		validateResponse func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "success - list users",
			mockSetup: func(m *MockUserUseCase) {
				m.On("List", mock.Anything).Return(entity.UserHistory{
					Users: []entity.User{
						{
							ID:        1,
							Email:     "user1@example.com",
							Username:  "user1",
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						},
						{
							ID:        2,
							Email:     "user2@example.com",
							Username:  "user2",
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						},
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var history entity.UserHistory
				err := json.Unmarshal(w.Body.Bytes(), &history)
				assert.NoError(t, err)
				assert.Len(t, history.Users, 2)
				assert.Equal(t, "user1@example.com", history.Users[0].Email)
			},
		},
		{
			name: "error - service error",
			mockSetup: func(m *MockUserUseCase) {
				m.On("List", mock.Anything).Return(entity.UserHistory{}, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			router := setupRouter()
			mockUserUseCase := new(MockUserUseCase)
			mockLogger := new(MockLogger)
			// Allow any logger calls during tests
			mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()
			mockLogger.On("Warn", mock.Anything, mock.Anything).Return()
			mockLogger.On("Info", mock.Anything, mock.Anything).Return()
			tt.mockSetup(mockUserUseCase)

			apiV1Group := router.Group("/v1")
			NewUserRoutes(apiV1Group, mockUserUseCase, mockLogger)

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/v1/user", nil)
			w := httptest.NewRecorder()

			// Execute
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.validateResponse != nil {
				tt.validateResponse(t, w)
			}
			mockUserUseCase.AssertExpectations(t)
		})
	}
}

func TestNewUserRoutes_UpdateUser(t *testing.T) {
	tests := []struct {
		name             string
		userID           string
		requestBody      interface{}
		mockSetup        func(*MockUserUseCase)
		expectedStatus   int
		validateResponse func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:   "success - update user",
			userID: "1",
			requestBody: request.UpdateUser{
				Email:    "updated@example.com",
				Username: "updateduser",
			},
			mockSetup: func(m *MockUserUseCase) {
				m.On("Update", mock.Anything, mock.MatchedBy(func(u entity.User) bool {
					return u.ID == 1 && u.Email == "updated@example.com"
				})).Return(nil)
			},
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "message")
			},
		},
		{
			name:           "error - invalid user ID",
			userID:         "invalid",
			requestBody:    request.UpdateUser{Email: "test@example.com"},
			mockSetup:      func(m *MockUserUseCase) {},
			expectedStatus: http.StatusBadRequest,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
		{
			name:   "error - invalid request body",
			userID: "1",
			requestBody: map[string]interface{}{
				"email": "invalid-email",
			},
			mockSetup:      func(m *MockUserUseCase) {},
			expectedStatus: http.StatusBadRequest,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
		{
			name:   "error - service error",
			userID: "1",
			requestBody: request.UpdateUser{
				Email:    "updated@example.com",
				Username: "updateduser",
			},
			mockSetup: func(m *MockUserUseCase) {
				m.On("Update", mock.Anything, mock.Anything).Return(errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			router := setupRouter()
			mockUserUseCase := new(MockUserUseCase)
			mockLogger := new(MockLogger)
			// Allow any logger calls during tests
			mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()
			mockLogger.On("Warn", mock.Anything, mock.Anything).Return()
			mockLogger.On("Info", mock.Anything, mock.Anything).Return()
			tt.mockSetup(mockUserUseCase)

			apiV1Group := router.Group("/v1")
			NewUserRoutes(apiV1Group, mockUserUseCase, mockLogger)

			// Create request
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPut, "/v1/user/"+tt.userID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Execute
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.validateResponse != nil {
				tt.validateResponse(t, w)
			}
			mockUserUseCase.AssertExpectations(t)
		})
	}
}

func TestNewUserRoutes_DeleteUser(t *testing.T) {
	tests := []struct {
		name             string
		userID           string
		mockSetup        func(*MockUserUseCase)
		expectedStatus   int
		validateResponse func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:   "success - delete user",
			userID: "1",
			mockSetup: func(m *MockUserUseCase) {
				m.On("Delete", mock.Anything, int64(1)).Return(nil)
			},
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "message")
			},
		},
		{
			name:           "error - invalid user ID",
			userID:         "invalid",
			mockSetup:      func(m *MockUserUseCase) {},
			expectedStatus: http.StatusBadRequest,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
		{
			name:   "error - service error",
			userID: "1",
			mockSetup: func(m *MockUserUseCase) {
				m.On("Delete", mock.Anything, int64(1)).Return(errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			router := setupRouter()
			mockUserUseCase := new(MockUserUseCase)
			mockLogger := new(MockLogger)
			// Allow any logger calls during tests
			mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()
			mockLogger.On("Warn", mock.Anything, mock.Anything).Return()
			mockLogger.On("Info", mock.Anything, mock.Anything).Return()
			tt.mockSetup(mockUserUseCase)

			apiV1Group := router.Group("/v1")
			NewUserRoutes(apiV1Group, mockUserUseCase, mockLogger)

			// Create request
			req := httptest.NewRequest(http.MethodDelete, "/v1/user/"+tt.userID, nil)
			w := httptest.NewRecorder()

			// Execute
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.validateResponse != nil {
				tt.validateResponse(t, w)
			}
			mockUserUseCase.AssertExpectations(t)
		})
	}
}

func TestNewUserRoutes_LoginUser(t *testing.T) {
	tests := []struct {
		name             string
		requestBody      interface{}
		mockSetup        func(*MockUserUseCase)
		expectedStatus   int
		validateResponse func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "success - login user",
			requestBody: request.LoginUser{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(m *MockUserUseCase) {
				m.On("Login", mock.Anything, "test@example.com", "password123").Return(
					"mock-jwt-token",
					entity.User{
						ID:       1,
						Email:    "test@example.com",
						Username: "testuser",
					},
					nil,
				)
			},
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "token")
				assert.Contains(t, response, "user")
				assert.Equal(t, "mock-jwt-token", response["token"])
			},
		},
		{
			name: "error - invalid request body",
			requestBody: map[string]interface{}{
				"email": "invalid-email",
			},
			mockSetup:      func(m *MockUserUseCase) {},
			expectedStatus: http.StatusBadRequest,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
		{
			name: "error - validation failure - missing email",
			requestBody: request.LoginUser{
				Password: "password123",
			},
			mockSetup:      func(m *MockUserUseCase) {},
			expectedStatus: http.StatusBadRequest,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
		{
			name: "error - invalid credentials",
			requestBody: request.LoginUser{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			mockSetup: func(m *MockUserUseCase) {
				m.On("Login", mock.Anything, "test@example.com", "wrongpassword").Return(
					"",
					entity.User{},
					errors.New("UserUseCase - Login - bcrypt.CompareHashAndPassword: crypto/bcrypt: hashedPassword is not the hash of the given password"),
				)
			},
			expectedStatus: http.StatusUnauthorized,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
		{
			name: "error - service error",
			requestBody: request.LoginUser{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(m *MockUserUseCase) {
				m.On("Login", mock.Anything, "test@example.com", "password123").Return(
					"",
					entity.User{},
					errors.New("database error"),
				)
			},
			expectedStatus: http.StatusInternalServerError,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			router := setupRouter()
			mockUserUseCase := new(MockUserUseCase)
			mockLogger := new(MockLogger)
			// Allow any logger calls during tests
			mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()
			mockLogger.On("Warn", mock.Anything, mock.Anything).Return()
			mockLogger.On("Info", mock.Anything, mock.Anything).Return()
			tt.mockSetup(mockUserUseCase)

			apiV1Group := router.Group("/v1")
			NewUserRoutes(apiV1Group, mockUserUseCase, mockLogger)

			// Create request
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Execute
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.validateResponse != nil {
				tt.validateResponse(t, w)
			}
			mockUserUseCase.AssertExpectations(t)
		})
	}
}
