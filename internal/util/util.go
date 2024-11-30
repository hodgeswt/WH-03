package util

var MaxBit int = IntPow(2, 9) - 1

func GetBit(x int, bit int) int {
    return (x >> bit) & 1
}

func IntPow(x int, y int) int {
    if y == 0 {
        return 1
    }

    if y == 1 {
        return x
    }

    o := x
    for i := 2; i <= y; i++ {
        o *= x
    }

    return o
}
