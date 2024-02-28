package base26

import (
	"math/big"
	"slices"
)

var base26 = []byte{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}
var index = map[byte]int{'A': 0, 'B': 1, 'C': 2, 'D': 3, 'E': 4, 'F': 5, 'G': 6, 'H': 7, 'I': 8, 'J': 9, 'K': 10, 'L': 11, 'M': 12, 'N': 13, 'O': 14, 'P': 15, 'Q': 16, 'R': 17, 'S': 18, 'T': 19, 'U': 20, 'V': 21, 'W': 22, 'X': 23, 'Y': 24, 'Z': 25}

func Encode(val uint64) string {
	var result []byte
	for val > 0 {
		//val--
		result = append(result, base26[val%26])
		val /= 26
	}

	slices.Reverse(result)
	return string(result)
}

func Decode(s string) uint64 {
	var result uint64 = 0
	length := len(s) - 1
	bigInt := big.NewInt(26)
	byteIndex := big.NewInt(0)
	power := big.NewInt(0)
	for i := range s {
		c := s[length-i]
		byteOffset := index[c]
		byteIndex.SetInt64(int64(i))
		result += uint64(byteOffset) * power.Exp(bigInt, byteIndex, nil).Uint64()
	}
	return result
}
