package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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
	if err := ctx.ShouldBindJSON(obj); err != nil {
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

		if ok, errResponse := bh.parseValidationErrors(err); ok {
			return errResponse, http.StatusBadRequest, err
		}

		var sliceValidationErrors binding.SliceValidationError
		if ok := errors.As(err, &sliceValidationErrors); ok && len(sliceValidationErrors) > 0 {
			if ok, errResponse := bh.parseValidationErrors(sliceValidationErrors[0]); ok {
				return errResponse, http.StatusBadRequest, err
			}
		}

		return nil, http.StatusInternalServerError, err
	}

	return nil, 0, nil
}

func (bh baseHandler) parseValidationErrors(err error) (bool, *models.ErrorResponse) {
	var validationErrors validator.ValidationErrors
	if ok := errors.As(err, &validationErrors); ok && len(validationErrors) > 0 {
		fErr := validationErrors[0]

		var errResponse string
		if fErr.Param() == "" {
			errResponse = fErr.Tag()
		} else {
			errResponse = fmt.Sprintf("%s=%s", fErr.Tag(), fErr.Param())
		}

		return true, &models.ErrorResponse{
			Error: fmt.Sprintf("Field validation for \"%s\" failed on the '%s' tag.", fErr.Field(), errResponse),
		}
	}

	return false, nil
}
