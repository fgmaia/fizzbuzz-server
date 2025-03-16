package services

import "strconv"

type FizzBuzzService struct {
}

func NewFizzBuzzService() *FizzBuzzService {
	return &FizzBuzzService{}
}

func (f *FizzBuzzService) GenerateFizzBuzz(int1, int2, limit int, str1, str2 string) []string {
	result := make([]string, limit)
	for i := 1; i <= limit; i++ {
		switch {
		case i%int1 == 0 && i%int2 == 0:
			result[i-1] = str1 + str2
		case i%int1 == 0:
			result[i-1] = str1
		case i%int2 == 0:
			result[i-1] = str2
		default:
			result[i-1] = strconv.Itoa(i)
		}
	}
	return result
}
