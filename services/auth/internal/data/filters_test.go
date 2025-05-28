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

func TestSortColumn(t *testing.T) {
	tests := []struct {
		name            string
		filters         Filters
		expectedSortVal string
		expectedPanic   bool
	}{
		{"A valid sort ascending column", Filters{1, 20, "name", []string{"name"}}, "name", false},
		{"A valid sort descending column", Filters{1, 20, "-name", []string{"-name"}}, "name", false},
		{"An invalid sort column", Filters{1, 20, "name", []string{"email"}}, "", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				r := recover()

				if !tc.expectedPanic && r != nil {
					t.Error("Didn't expect to panic but it did")
				}
			}()

			result := tc.filters.sortColumn()
			if result != tc.expectedSortVal {
				t.Errorf("Expected '%v', got '%v'", tc.expectedSortVal, result)
			}
		})
	}
}

func TestSortDirection(t *testing.T) {
	tests := []struct {
		name                  string
		filters               Filters
		expectedSortDirection string
	}{
		{"An ascending sort direction", Filters{1, 20, "name", []string{"name"}}, "ASC"},
		{"A descending sort direction", Filters{1, 20, "-name", []string{"-name"}}, "DESC"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.filters.sortDirection()
			if result != tc.expectedSortDirection {
				t.Errorf("Expected '%v', got '%v'", tc.expectedSortDirection, result)
			}
		})
	}
}

func TestLimit(t *testing.T) {
	f := Filters{
		Page:         2,
		PageSize:     5,
		Sort:         "name",
		SortSafeList: []string{"name"},
	}

	result := f.limit()
	if result != 5 {
		t.Errorf("Expected '%v', got '%v'", 5, result)
	}
}

func TestOffset(t *testing.T) {
	tests := []struct {
		name           string
		filters        Filters
		expectedOffser int
	}{
		{"Small offset", Filters{2, 20, "name", []string{"name"}}, 20},
		{"Big offset", Filters{5, 20, "name", []string{"name"}}, 80},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.filters.offset()
			if result != tc.expectedOffser {
				t.Errorf("Expected '%v', got '%v'", tc.expectedOffser, result)
			}
		})
	}
}

func TestCalculateMetadata(t *testing.T) {
	result := calculateMetadata(183, 2, 20)

	if result.CurrentPage != 2 ||
		result.FirstPage != 1 ||
		result.PageSize != 20 ||
		result.LastPage != 10 ||
		result.TotalRecords != 183 {
		t.Errorf("Wrong results '%v'", result)
	}
}
