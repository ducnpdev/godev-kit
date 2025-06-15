package v1

import (
	"net/http"
	"strconv"

	"github.com/ducnpdev/godev-kit/internal/controller/http/v1/request"
	"github.com/ducnpdev/godev-kit/internal/entity"
	"github.com/gin-gonic/gin"
)

// @Summary     Create user
// @Description Create a new user
// @ID          create-user
// @Tags  	    user
// @Accept      json
// @Produce     json
// @Param       request body request.CreateUser true "Create user"
// @Success     201 {object} entity.User
// @Failure     400 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /user [post]
func (r *V1) createUser(c *gin.Context) {
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
// @Param       id path int true "User ID"
// @Success     200 {object} entity.User
// @Failure     400 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /user/{id} [get]
func (r *V1) getUser(c *gin.Context) {
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
// @Success     200 {object} entity.UserHistory
// @Failure     500 {object} response.Error
// @Router      /user [get]
func (r *V1) listUsers(c *gin.Context) {
	userHistory, err := r.user.List(c.Request.Context())
	if err != nil {
		r.l.Error(err, "http - v1 - listUsers")
		errorResponse(c, http.StatusInternalServerError, "user service problems")
		return
	}

	// for i := 0; i < 5; i++ {
	// 	time.Sleep(time.Second)
	// 	fmt.Println(i)
	// }

	c.JSON(http.StatusOK, userHistory)
}

// @Summary     Update user
// @Description Update user by ID
// @ID          update-user
// @Tags  	    user
// @Accept      json
// @Produce     json
// @Param       id path int true "User ID"
// @Param       request body request.UpdateUser true "Update user"
// @Success     200 {object} response.Success
// @Failure     400 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /user/{id} [put]
func (r *V1) updateUser(c *gin.Context) {
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
// @Param       id path int true "User ID"
// @Success     200 {object} response.Success
// @Failure     400 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /user/{id} [delete]
func (r *V1) deleteUser(c *gin.Context) {
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
