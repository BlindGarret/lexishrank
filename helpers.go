package lexishrank

import (
	"fmt"
	"strings"

	"github.com/BlindGarret/lexishrank/base26"
)

// TODO: the old manner of overflow checking doesn't work now that we can do like 43 chars in certain cases. We need to check overflow during difference checks.
// TODO: We need to add a gap to the beginning of rankspace for adds
// TODO: Implement BEFORE
// TODO: Implement NeedsReindex function with Risk ranking
// TODO: Fix up docs
// TODO: Publish

func calculateRankGapSize(objectCount uint64, minGapSize uint64) (gapSize uint64, maxVal uint64, err error) {
	minVal := objectCount*minGapSize + objectCount
	if minVal < objectCount {
		// if we get here, we've overflowed the uint64
		return 0, 0, ErrGapSizeToLargeForObjectCount
	}
	base36MinVal := base26.Encode(minVal)
	minCharLength := len(base36MinVal)
	maxVal = base26.Decode(strings.Repeat("Z", minCharLength))

	// you can get a tighter fit by subtracting one from object count here, but leading zeroes
	// break our growth strategy, so we'll just leave it as is
	return maxVal / objectCount, maxVal, nil
}

// func calculateMiddleRank(first string, second string) string {
// 	paddedFirst, paddedSecond := padLowerValue(first, second)
// 	firstVal := base26.Decode(paddedFirst)
// 	secondVal := base26.Decode(paddedSecond)
// 	difference := secondVal - firstVal
// 	if difference < 2 {
// 		// We don't have a diff so increment the value of the first rank
// 		return paddedFirst + "M" // M is the middle of base26, it gives a bit of room for additional sorting until indexing happens.
// 	} else {
// 		return base26.Encode(firstVal + difference/2)
// 	}
// }

func calculateMiddleRank(first string, second string) string {
	paddedFirst, paddedSecond := padLowerValue(first, second)
	difference := symbolicSubtractionDifference(paddedFirst, paddedSecond)
	if difference < 2 {
		// We don't have a diff so increment the value of the first rank
		return paddedFirst + "M" // M is the middle of base26, it gives a bit of room for additional sorting until indexing happens.
	}

	return moveRankUp(difference/2, paddedFirst)
}

func moveRankUp(distance uint64, rank string) string {
	var sb strings.Builder
	carry := 0
	for i := 0; i < len(rank); i++ {
		diffForCode := distance / iPow(26, uint64(i)) % 26
		pointer := len(rank) - 1 - i
		newCode := rank[pointer] + byte(diffForCode) + byte(carry)
		carry = 0
		if newCode > 'Z' {
			carry++
			newCode -= 26
		}
		sb.WriteByte(newCode)
	}
	return reverse(sb.String())
}

func iPow(a, b uint64) uint64 {
	var result uint64 = 1

	for 0 != b {
		if 0 != (b & 1) {
			result *= a

		}
		b >>= 1
		a *= a
	}

	return result
}

func reverse(s string) string {
	rns := []rune(s) // convert to rune
	for i, j := 0, len(rns)-1; i < j; i, j = i+1, j-1 {

		// swap the letters of the string,
		// like first with last and so on.
		rns[i], rns[j] = rns[j], rns[i]
	}

	// return the reversed string.
	return string(rns)
}

func symbolicSubtractionDifference(first string, second string) uint64 {
	var difference uint64 = 0
	aBytes := []byte(first)
	bBytes := []byte(second)
	for i := len(first) - 1; i >= 0; i-- {
		aByte := aBytes[i]
		bByte := bBytes[i]
		if bByte < aByte {
			bByte += 26
			bBytes[i-1]--
		}
		powerVal := iPow(26, uint64(len(aBytes)-i-1))
		difference += uint64(bByte-aByte) * powerVal
	}
	return difference
}

func padLowerValue(first string, second string) (string, string) {
	for len(first) < len(second) {
		first = first + "A"
	}
	for len(second) < len(first) {
		second = second + "A"
	}
	return first, second
}

func getDirection(nextBucket int) (directionFunction, Direction) {
	if nextBucket == 0 {
		return beginToEndFunc, BeginningToEnd
	}
	return endToBeginFunc, EndToBegining
}

// directionFunction is a function that returns the next rank based on the current rank and the gap size
type directionFunction func(currentID uint64, gapSize uint64) uint64

var beginToEndFunc = func(currentID uint64, gapSize uint64) uint64 {
	return currentID + gapSize
}
var endToBeginFunc = func(currentID uint64, gapSize uint64) uint64 {
	return currentID - gapSize
}

func nextBucket(currentBucket int) int {
	if currentBucket == MaxBucketValue {
		return 0
	}
	return currentBucket + 1
}

func formatRank(bucket int, id string, paddedIDLen int) string {
	return fmt.Sprintf("%d|%s%s", bucket, strings.Repeat("A", max(paddedIDLen-len(id), 0)), id)
}
