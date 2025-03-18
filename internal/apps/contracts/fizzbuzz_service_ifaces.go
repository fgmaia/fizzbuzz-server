package contracts

type FizzBuzzServiceIface interface {
	GenerateFizzBuzz(int1, int2, limit int, str1, str2 string) []string
}
