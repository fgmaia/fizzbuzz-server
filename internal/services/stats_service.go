package services

import (
	"fizzbuzz-server/internal/entities"
	"fmt"
	"strconv"
	"strings"
)

type StatsService struct {
}

func NewStatsService() *StatsService {
	return &StatsService{}
}

func (s *StatsService) ParseStatsKey(key string) (entities.StatsKeys, error) {
	var statsKeys entities.StatsKeys
	parts := strings.Split(key, ",")
	if len(parts) != 5 {
		return statsKeys, fmt.Errorf("invalid stats key format")
	}

	int1, err := strconv.Atoi(parts[0])
	if err != nil {
		return statsKeys, err
	}
	int2, err := strconv.Atoi(parts[1])
	if err != nil {
		return statsKeys, err
	}
	limit, err := strconv.Atoi(parts[2])
	if err != nil {
		return statsKeys, err
	}

	statsKeys.Int1 = int1
	statsKeys.Int2 = int2
	statsKeys.Limit = limit
	statsKeys.Str1 = parts[3]
	statsKeys.Str2 = parts[4]

	return statsKeys, nil
}
