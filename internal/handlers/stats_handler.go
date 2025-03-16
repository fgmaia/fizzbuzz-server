package handlers

import (
	"fizzbuzz-server/internal/apps"

	"github.com/gofiber/fiber/v2"
)

func Stats(c *fiber.Ctx) error {
	stats.Mutex.RLock()
	defer stats.Mutex.RUnlock()

	if len(stats.Counts) == 0 {
		return c.JSON(fiber.Map{
			"message": "No requests have been made yet",
		})
	}

	// Find most frequent request
	var maxKey string
	maxCount := 0
	for key, count := range stats.Counts {
		if count > maxCount {
			maxCount = count
			maxKey = key
		}
	}

	// Parse the key back into components
	parts, err := apps.App().StatsService.ParseStatsKey(maxKey)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Internal stats error",
		})
	}

	resp := StatsResponse{}
	resp.MostFrequentRequest.Int1 = parts.Int1
	resp.MostFrequentRequest.Int2 = parts.Int2
	resp.MostFrequentRequest.Limit = parts.Limit
	resp.MostFrequentRequest.Str1 = parts.Str1
	resp.MostFrequentRequest.Str2 = parts.Str2
	resp.MostFrequentRequest.Hits = maxCount

	return c.JSON(resp)
}
