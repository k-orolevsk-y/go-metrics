package handlers

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidateHTTPMethod(t *testing.T) {
	type args struct {
		methodRequest string
		methodNeed    string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Positive GET",
			args: args{
				methodRequest: http.MethodGet,
				methodNeed:    http.MethodGet,
			},
			want: true,
		},
		{
			name: "Positive POST",
			args: args{
				methodRequest: http.MethodPost,
				methodNeed:    http.MethodPost,
			},
			want: true,
		},
		{
			name: "Negative GET/POST",
			args: args{
				methodRequest: http.MethodGet,
				methodNeed:    http.MethodPost,
			},
			want: false,
		},
		{
			name: "Negative PUT/DELETE",
			args: args{
				methodRequest: http.MethodPut,
				methodNeed:    http.MethodDelete,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.args.methodRequest, "/", nil)
			assert.Equal(t, tt.want, ValidateHTTPMethod(request, tt.args.methodNeed))
		})
	}
}

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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.args.contentTypeRequest != "" {
				request.Header.Set("Content-Type", tt.args.contentTypeRequest)
			}

			assert.Equal(t, tt.want, ValidateContentType(request, tt.args.contentTypeNeed))
		})
	}
}

func TestParseURLParams(t *testing.T) {
	type args struct {
		pathRequest string
		path        string
	}
	tests := []struct {
		name string
		args args
		want ParsedURLParams
	}{
		{
			name: "Positive with params",
			args: args{
				pathRequest: "/test/Key/Value",
				path:        "/test/",
			},
			want: ParsedURLParams{"Key", "Value"},
		},
		{
			name: "Positive without params",
			args: args{
				pathRequest: "/test/",
				path:        "/test/",
			},
			want: ParsedURLParams{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.args.pathRequest, nil)
			assert.Equal(t, tt.want, ParseURLParams(tt.args.path, request.URL))
		})
	}
}
