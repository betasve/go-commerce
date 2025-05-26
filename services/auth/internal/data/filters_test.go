package data

import (
	"testing"

	"github.com/betasve/go-commerce/services/auth/internal/validator"
)

func TestValidateFilters(t *testing.T) {

	tests := []struct {
		name        string
		filters     Filters
		validator   *validator.Validator
		expectedKey string
		expectedErr map[string]string
	}{
		{"All valid filters", Filters{1, 20, "name", []string{"name"}}, validator.New(), "", map[string]string{}},
		{"Invalidly small page number", Filters{0, 20, "name", []string{"name"}}, validator.New(), "page", map[string]string{"page": "must be greater than zero"}},
		{"Invalidly large page number", Filters{10_000_001, 20, "name", []string{"name"}}, validator.New(), "page", map[string]string{"page": "must be a maximum of 10 million"}},
		{"Invalidly small page size", Filters{1, 0, "name", []string{"name"}}, validator.New(), "page_size", map[string]string{"page_size": "must be greater than zero"}},
		{"Invalidly large page size", Filters{1, 20, "name", []string{"name"}}, validator.New(), "page_size", map[string]string{"page_size": "must be a maximum of 100"}},
		{"Invalid sort value", Filters{1, 20, "email", []string{"name"}}, validator.New(), "sort", map[string]string{"sort": "invalid sort value"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ValidateFilters(tc.validator, tc.filters)
			if tc.validator.Invalid() && tc.validator.Errors[tc.expectedKey] != tc.expectedErr[tc.expectedKey] {
				t.Errorf("Expected '%v', got '%v'", tc.expectedErr[tc.expectedKey], tc.validator.Errors[tc.expectedKey])
			}
		})
	}
}
