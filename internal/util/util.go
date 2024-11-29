package util

var MaxBit int = (2 ^ 9) - 1

func GetBit(x int, bit int) int {
    return (x >> bit) & 1
}
