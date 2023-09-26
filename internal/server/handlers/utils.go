package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/server/models"
	"io"
	"net/http"
)

func (bh baseHandler) validateContentType(ctx *gin.Context, contentType string, withoutContentType bool) bool {
	requestContentType := ctx.GetHeader("Content-Type")
	if requestContentType == "" && withoutContentType {
		return true
	}

	return requestContentType == contentType
}

func (bh baseHandler) validateAndShouldBindJSON(ctx *gin.Context, obj any) (*models.ErrorResponse, int, error) {
	if err := ctx.ShouldBindJSON(&obj); err != nil {
		if errors.Is(err, io.EOF) {
			return &models.ErrorResponse{Error: "Request body not provided."}, http.StatusBadRequest, err
		}

		var jsonTypeError *json.UnmarshalTypeError
		if ok := errors.As(err, &jsonTypeError); ok {
			return &models.ErrorResponse{
				Error: fmt.Sprintf("Field value \"%s\" must be %s.", jsonTypeError.Field, jsonTypeError.Type),
			}, http.StatusBadRequest, err
		}

		var jsonError *json.SyntaxError
		if ok := errors.As(err, &jsonError); ok {
			return &models.ErrorResponse{
				Error: fmt.Sprintf("JSON error: %s", jsonError.Error()),
			}, http.StatusBadRequest, err
		}

		var validationErrors validator.ValidationErrors
		if ok := errors.As(err, &validationErrors); ok {
			for _, fErr := range validationErrors {
				var errResponse string
				if fErr.Param() == "" {
					errResponse = fErr.Tag()
				} else {
					errResponse = fmt.Sprintf("%s=%s", fErr.Tag(), fErr.Param())
				}

				return &models.ErrorResponse{
					Error: fmt.Sprintf("Field validation for \"%s\" failed on the '%s' tag.", fErr.Field(), errResponse),
				}, http.StatusBadRequest, err
			}
		}

		return nil, http.StatusInternalServerError, err
	}

	return nil, 0, nil
}
