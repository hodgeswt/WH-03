package util

func ContainsOnly(s string, allowed map[rune]bool) bool {
    for _, char := range s {
        if !allowed[char] {
            return false
        }
    }

    return true
}
