package data

import (
	"github.com/denim-bluu/movie-db-app/internal/validator"
)

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafelist []string
}

func ValidateFilters(v *validator.Validator, filters *Filters) {
	v.Check(filters.Page <= 10_000_000, "page", "Max 10 million")
	v.Check(filters.Page > 0, "page", "must be a positive integer")
	v.Check(filters.PageSize <= 100, "page_size", "Max 100")
	v.Check(filters.PageSize > 0, "page_size", "must be a positive integer")
	v.Check(validator.PermitteValues(filters.Sort, filters.SortSafelist...), "sort", "invalid sort value")
}
