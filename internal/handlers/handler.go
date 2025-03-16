package handlers

import "fizzbuzz-server/internal/entities"

var stats = entities.RequestStats{
	Counts: make(map[string]int),
}

// StatsResponse represents the stats endpoint response
type StatsResponse struct {
	MostFrequentRequest struct {
		Int1  int    `json:"int1"`
		Int2  int    `json:"int2"`
		Limit int    `json:"limit"`
		Str1  string `json:"str1"`
		Str2  string `json:"str2"`
		Hits  int    `json:"hits"`
	} `json:"most_frequent_request"`
}

// ErrorResponse standardizes error responses
type ErrorResponse struct {
	Error string `json:"error"`
}
