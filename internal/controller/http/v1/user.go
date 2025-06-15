package v1

import (
	"net/http"
	"strconv"

	"github.com/ducnpdev/godev-kit/internal/controller/http/v1/request"
	"github.com/ducnpdev/godev-kit/internal/controller/http/v1/response"
	"github.com/ducnpdev/godev-kit/internal/entity"
	"github.com/gofiber/fiber/v2"
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
func (r *V1) createUser(ctx *fiber.Ctx) error {
	var body request.CreateUser
	if err := ctx.BodyParser(&body); err != nil {
		r.l.Error(err, "http - v1 - createUser")
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error(err, "http - v1 - createUser")
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	user, err := r.user.Create(
		ctx.UserContext(),
		entity.User{
			Email:    body.Email,
			Username: body.Username,
			Password: body.Password,
		},
	)
	if err != nil {
		r.l.Error(err, "http - v1 - createUser")
		return errorResponse(ctx, http.StatusInternalServerError, "user service problems")
	}

	return ctx.Status(http.StatusCreated).JSON(user)
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
func (r *V1) getUser(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		r.l.Error(err, "http - v1 - getUser")
		return errorResponse(ctx, http.StatusBadRequest, "invalid user id")
	}

	user, err := r.user.GetByID(ctx.UserContext(), id)
	if err != nil {
		r.l.Error(err, "http - v1 - getUser")
		return errorResponse(ctx, http.StatusInternalServerError, "user service problems")
	}

	return ctx.Status(http.StatusOK).JSON(user)
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
func (r *V1) updateUser(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		r.l.Error(err, "http - v1 - updateUser")
		return errorResponse(ctx, http.StatusBadRequest, "invalid user id")
	}

	var body request.UpdateUser
	if err := ctx.BodyParser(&body); err != nil {
		r.l.Error(err, "http - v1 - updateUser")
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error(err, "http - v1 - updateUser")
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	user := entity.User{
		ID:       id,
		Email:    body.Email,
		Username: body.Username,
		Password: body.Password,
	}

	if err := r.user.Update(ctx.UserContext(), user); err != nil {
		r.l.Error(err, "http - v1 - updateUser")
		return errorResponse(ctx, http.StatusInternalServerError, "user service problems")
	}

	return ctx.Status(http.StatusOK).JSON(response.Success{Message: "user updated successfully"})
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
func (r *V1) deleteUser(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		r.l.Error(err, "http - v1 - deleteUser")
		return errorResponse(ctx, http.StatusBadRequest, "invalid user id")
	}

	if err := r.user.Delete(ctx.UserContext(), id); err != nil {
		r.l.Error(err, "http - v1 - deleteUser")
		return errorResponse(ctx, http.StatusInternalServerError, "user service problems")
	}

	return ctx.Status(http.StatusOK).JSON(response.Success{Message: "user deleted successfully"})
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
func (r *V1) listUsers(ctx *fiber.Ctx) error {
	users, err := r.user.List(ctx.UserContext())
	if err != nil {
		r.l.Error(err, "http - v1 - listUsers")
		return errorResponse(ctx, http.StatusInternalServerError, "user service problems")
	}

	return ctx.Status(http.StatusOK).JSON(users)
}
