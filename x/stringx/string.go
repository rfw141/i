package stringx

import (
	"math/rand"
	"strings"
	"time"
)

// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var (
	strRandSrc = rand.NewSource(time.Now().UnixNano())
)

func Random(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	for i, cache, remain := n-1, strRandSrc.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = strRandSrc.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}

func CamelCase(name string) string {
	var buf []byte
	toggleUpper := true
	for i := 0; i < len(name); i++ {
		if name[i] == '_' {
			toggleUpper = true
		} else {
			c := name[i]
			if toggleUpper {
				toggleUpper = false
				if c >= 'a' && c <= 'z' {
					c = c - 'a' + 'A'
				}
			}
			if c >= '0' && c <= '9' {
				toggleUpper = true
			}
			buf = append(buf, c)
		}
	}
	return string(buf)
}

func SnakeCase(name string) string {
	var posList []int
	i := 1
	for i < len(name) {
		if name[i] >= 'A' && name[i] <= 'Z' {
			posList = append(posList, i)
			i++
			for i < len(name) && name[i] >= 'A' && name[i] <= 'Z' {
				i++
			}
		} else {
			i++
		}
	}
	lower := strings.ToLower(name)
	if len(posList) == 0 {
		return lower
	}
	b := strings.Builder{}
	left := 0
	for _, right := range posList {
		b.WriteString(lower[left:right])
		b.WriteByte('_')
		left = right
	}
	b.WriteString(lower[left:])
	return b.String()
}
