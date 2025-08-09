package v1

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ducnpdev/godev-kit/internal/controller/http/v1/request"
	"github.com/ducnpdev/godev-kit/internal/controller/http/v1/response"
	"github.com/ducnpdev/godev-kit/internal/entity"
	"github.com/gin-gonic/gin"
)

// @Summary     Create user
// @Description Create a new user
// @ID          create-user
// @Tags  	    user
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       request body request.CreateUser true "Create user"
// @Success     201 {object} entity.User
// @Failure     400 {object} response.Error
// @Failure     401 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /v1/user [post]
func (r *V1) CreateUser(c *gin.Context) {
	var body request.CreateUser
	if err := c.ShouldBindJSON(&body); err != nil {
		r.l.Error(err, "http - v1 - createUser")
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error(err, "http - v1 - createUser")
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := r.user.Create(
		c.Request.Context(),
		entity.User{
			Email:    body.Email,
			Username: body.Username,
			Password: body.Password,
		},
	)
	if err != nil {
		r.l.Error(err, "http - v1 - createUser")
		errorResponse(c, http.StatusInternalServerError, "user service problems")
		return
	}

	c.JSON(http.StatusCreated, user)
}

// @Summary     Get user
// @Description Get user by ID
// @ID          get-user
// @Tags  	    user
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id path int true "User ID"
// @Success     200 {object} entity.User
// @Failure     400 {object} response.Error
// @Failure     401 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /v1/user/{id} [get]
func (r *V1) GetUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		r.l.Error(err, "http - v1 - getUser")
		errorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := r.user.GetByID(c.Request.Context(), id)
	if err != nil {
		r.l.Error(err, "http - v1 - getUser")
		errorResponse(c, http.StatusInternalServerError, "user service problems")
		return
	}

	if user.ID == 0 {
		errorResponse(c, http.StatusNotFound, "user not found")
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary     List users
// @Description Get all users
// @ID          list-users
// @Tags  	    user
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} entity.UserHistory
// @Failure     401 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /v1/user [get]
func (r *V1) ListUsers(c *gin.Context) {
	// For debugging timeout behavior, uncomment the sleep below:
	// This will trigger the 8-second request timeout defined in middleware
	for i := 0; i < 10; i++ {
		fmt.Printf("[DEBUG] ListUsers processing step %d/10\n", i+1)
		time.Sleep(1 * time.Second)

		// Check if context was cancelled (timeout or client disconnect)
		select {
		case <-c.Request.Context().Done():
			r.l.Warn("ListUsers request was cancelled",
				"reason", c.Request.Context().Err().Error(),
				"client_ip", c.ClientIP(),
				"step", i+1)
			return
		default:
			// Continue processing
		}
	}

	userHistory, err := r.user.List(c.Request.Context())
	if err != nil {
		r.l.Error(err, "http - v1 - listUsers")
		errorResponse(c, http.StatusInternalServerError, "user service problems")
		return
	}

	r.l.Info("ListUsers completed successfully", "client_ip", c.ClientIP())
	c.JSON(http.StatusOK, userHistory)
}

// @Summary     Update user
// @Description Update user by ID
// @ID          update-user
// @Tags  	    user
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id path int true "User ID"
// @Param       request body request.UpdateUser true "Update user"
// @Success     200 {object} response.Success
// @Failure     400 {object} response.Error
// @Failure     401 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /v1/user/{id} [put]
func (r *V1) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		r.l.Error(err, "http - v1 - updateUser")
		errorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}

	var body request.UpdateUser
	if err := c.ShouldBindJSON(&body); err != nil {
		r.l.Error(err, "http - v1 - updateUser")
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error(err, "http - v1 - updateUser")
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	err = r.user.Update(
		c.Request.Context(),
		entity.User{
			ID:       id,
			Email:    body.Email,
			Username: body.Username,
			Password: body.Password,
		},
	)
	if err != nil {
		r.l.Error(err, "http - v1 - updateUser")
		errorResponse(c, http.StatusInternalServerError, "user service problems")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user updated successfully"})
}

// @Summary     Delete user
// @Description Delete user by ID
// @ID          delete-user
// @Tags  	    user
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id path int true "User ID"
// @Success     200 {object} response.Success
// @Failure     400 {object} response.Error
// @Failure     401 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /v1/user/{id} [delete]
func (r *V1) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		r.l.Error(err, "http - v1 - deleteUser")
		errorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}

	err = r.user.Delete(c.Request.Context(), id)
	if err != nil {
		r.l.Error(err, "http - v1 - deleteUser")
		errorResponse(c, http.StatusInternalServerError, "user service problems")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

// @Summary     Login user
// @Description Login user with email and password
// @ID          login-user
// @Tags  	    auth
// @Accept      json
// @Produce     json
// @Param       request body request.LoginUser true "Login user"
// @Success     200 {object} response.LoginResponse
// @Failure     400 {object} response.Error
// @Failure     401 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /v1/auth/login [post]
func (r *V1) LoginUser(c *gin.Context) {
	var body request.LoginUser
	if err := c.ShouldBindJSON(&body); err != nil {
		r.l.Error(err, "http - v1 - loginUser")
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error(err, "http - v1 - loginUser")
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	token, user, err := r.user.Login(c.Request.Context(), body.Email, body.Password)
	if err != nil {
		r.l.Error(err, "http - v1 - loginUser")
		if err.Error() == "UserUseCase - Login - bcrypt.CompareHashAndPassword: crypto/bcrypt: hashedPassword is not the hash of the given password" {
			errorResponse(c, http.StatusUnauthorized, "invalid credentials")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "user service problems")
		return
	}

	c.JSON(http.StatusOK, response.LoginResponse{
		Token: token,
		User: struct {
			ID       int64  `json:"id"`
			Email    string `json:"email"`
			Username string `json:"username"`
		}{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
		},
	})
}
