package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidateContentType(t *testing.T) {
	type args struct {
		contentTypeRequest string
		contentTypeNeed    string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Positive text/plain",
			args: args{
				contentTypeRequest: "text/plain",
				contentTypeNeed:    "text/plain",
			},
			want: true,
		},
		{
			name: "Positive json/application",
			args: args{
				contentTypeRequest: "json/application",
				contentTypeNeed:    "json/application",
			},
			want: true,
		},
		{
			name: "Positive without request content-type",
			args: args{
				contentTypeRequest: "",
				contentTypeNeed:    "json/application",
			},
			want: true,
		},
		{
			name: "Negative json/application & text/plain",
			args: args{
				contentTypeRequest: "json/application",
				contentTypeNeed:    "text/plain",
			},
			want: false,
		},
	}

	gin.SetMode(gin.ReleaseMode)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)

			if tt.args.contentTypeRequest != "" {
				ctx.Request.Header.Set("Content-Type", tt.args.contentTypeRequest)
			}
			println(ctx.GetHeader("Content-Type"))

			assert.Equal(t, tt.want, ValidateContentType(ctx, tt.args.contentTypeNeed))
		})
	}
}
