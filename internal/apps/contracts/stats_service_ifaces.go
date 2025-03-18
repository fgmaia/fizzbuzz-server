package contracts

import "fizzbuzz-server/internal/entities"

type StatsServiceIface interface {
	ParseStatsKey(key string) (entities.StatsKeys, error)
}
