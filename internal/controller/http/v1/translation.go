package v1

import (
	"net/http"

	"github.com/ducnpdev/godev-kit/internal/controller/http/v1/request"
	"github.com/ducnpdev/godev-kit/internal/entity"
	"github.com/gin-gonic/gin"
)

// @Summary     Show history
// @Description Show all translation history
// @ID          history
// @Tags  	    translation
// @Accept      json
// @Produce     json
// @Success     200 {object} entity.TranslationHistory
// @Failure     500 {object} response.Error
// @Router      /translation/history [get]
func (r *V1) history(c *gin.Context) {
	translationHistory, err := r.t.History(c.Request.Context())
	if err != nil {
		r.l.Error(err, "http - v1 - history")
		errorResponse(c, http.StatusInternalServerError, "database problems")
		return
	}

	c.JSON(http.StatusOK, translationHistory)
}

// @Summary     Translate
// @Description Translate a text
// @ID          do-translate
// @Tags  	    translation
// @Accept      json
// @Produce     json
// @Param       request body request.Translate true "Set up translation"
// @Success     200 {object} entity.Translation
// @Failure     400 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /translation/do-translate [post]
func (r *V1) doTranslate(c *gin.Context) {
	var body request.Translate

	if err := c.ShouldBindJSON(&body); err != nil {
		r.l.Error(err, "http - v1 - doTranslate")
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error(err, "http - v1 - doTranslate")
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	translation, err := r.t.Translate(
		c.Request.Context(),
		entity.Translation{
			Source:      body.Source,
			Destination: body.Destination,
			Original:    body.Original,
		},
	)
	if err != nil {
		r.l.Error(err, "http - v1 - doTranslate")
		errorResponse(c, http.StatusInternalServerError, "translation service problems")
		return
	}

	c.JSON(http.StatusOK, translation)
}

// errorResponse is a helper function to send error responses
func errorResponse(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}
