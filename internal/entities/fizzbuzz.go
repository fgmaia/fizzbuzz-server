package entities

import "sync"

// RequestStats holds request statistics with thread-safe access
type RequestStats struct {
	Counts map[string]int
	Mutex  sync.RWMutex
}

// FizzBuzzRequest represents the expected query parameters
type FizzBuzzRequest struct {
	Int1  int    `query:"int1" validate:"required,gt=0"`
	Int2  int    `query:"int2" validate:"required,gt=0"`
	Limit int    `query:"limit" validate:"required,gt=0,lte=10000"`
	Str1  string `query:"str1" validate:"required" default:"fizz"`
	Str2  string `query:"str2" validate:"required" default:"buzz"`
}

type StatsKeys struct {
	Int1  int
	Int2  int
	Limit int
	Str1  string
	Str2  string
}
