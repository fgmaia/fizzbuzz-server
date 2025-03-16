package handlers

import "github.com/gofiber/fiber/v2"

func DocHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"fizzbuzz_endpoint": fiber.Map{
			"path":        "/fizzbuzz",
			"method":      "GET",
			"params":      "int1(int), int2(int), limit(int), str1(string), str2(string)",
			"description": "Returns a FizzBuzz sequence based on parameters",
		},
		"stats_endpoint": fiber.Map{
			"path":        "/stats",
			"method":      "GET",
			"params":      "none",
			"description": "Returns statistics about most frequent request",
		},
	})
}
