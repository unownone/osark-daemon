package models

// OSQuery is an interface for a table in osquery
type OSQuery[T any] interface {
	Query(query string) (T, error)
}