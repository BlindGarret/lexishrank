package lexishrank

import (
	"fmt"
	"testing"
)

func TestUnitPadLowerValue(t *testing.T) {
	var tests = []struct {
		a, b                 string
		expectedA, expectedB string
	}{
		{"B", "BB", "BA", "BB"},
		{"C", "BB", "CA", "BB"},
		{"F", "ZAW", "FAA", "ZAW"},
		{"ZAW", "FS", "ZAW", "FSA"},
	}

	for _, tt := range tests {
		testName := fmt.Sprintf("%s,%s", tt.a, tt.b)
		t.Run(testName, func(t *testing.T) {
			paddedFirst, paddedSecond := padLowerValue(tt.a, tt.b)
			if paddedFirst != tt.expectedA {
				t.Errorf("Expected paddedFirst to be %s, got %s", tt.expectedA, paddedFirst)
			}
			if paddedSecond != tt.expectedB {
				t.Errorf("Expected paddedSecond to be %s, got %s", tt.expectedB, paddedSecond)
			}
		})
	}
}

func TestUnitCalculateMiddleRank(t *testing.T) {
	var tests = []struct {
		a, b, expected string
	}{
		{"BA", "BC", "BB"},
		{"BA", "BB", "BAM"},
		{"B", "BM", "BG"},
		{"BHRE", "ZZZZZZ", "NQVOZZ"},
	}

	for _, tt := range tests {
		testName := fmt.Sprintf("%s,%s", tt.a, tt.b)
		t.Run(testName, func(t *testing.T) {
			result := calculateMiddleRank(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Expected result to be %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestUnitCalculateRankGapSize(t *testing.T) {
	var tests = []struct {
		objectCount, minGapSize, expectedVal, expectedMaxVal uint64
		expectedErr                                          error
	}{
		{26, 2, 25, 675, nil},
		{2, 1, 12, 25, nil},
		{3, 1, 8, 25, nil},
		{10, 2, 67, 675, nil},
		{10, 200, 1757, 17575, nil},
		{18446744073709551615, 2, 0, 0, ErrGapSizeToLargeForObjectCount},
	}

	for _, tt := range tests {
		testName := fmt.Sprintf("%d,%d", tt.objectCount, tt.minGapSize)
		t.Run(testName, func(t *testing.T) {
			result, maxVal, err := calculateRankGapSize(tt.objectCount, tt.minGapSize)
			if result != tt.expectedVal {
				t.Errorf("Expected result to be %d, got %d", tt.expectedVal, result)
			}
			if maxVal != tt.expectedMaxVal {
				t.Errorf("Expected maxVal to be %d, got %d", tt.expectedMaxVal, maxVal)
			}
			if err != tt.expectedErr {
				t.Errorf("Expected error to be %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestUnitGetDirection(t *testing.T) {
	var tests = []struct {
		bucket            int
		expectedDirection Direction
	}{
		{1, EndToBegining},
		{0, BeginningToEnd},
	}

	for _, tt := range tests {
		testName := fmt.Sprintf("%d", tt.bucket)
		t.Run(testName, func(t *testing.T) {
			_, direction := getDirection(tt.bucket)
			if direction != tt.expectedDirection {
				t.Errorf("Expected direction to be %v, got %v", tt.expectedDirection, direction)
			}
		})
	}
}

// This is litterally just an addition function
func TestUnitBeginToEndFunc(t *testing.T) {
	result := beginToEndFunc(1, 2)
	if result != 3 {
		t.Errorf("Expected result to be 3, got %d", result)
	}
}

// This is litterally just a subtraction function
func TestUnitEndToBeginFunc(t *testing.T) {
	result := endToBeginFunc(3, 2)
	if result != 1 {
		t.Errorf("Expected result to be 1, got %d", result)
	}
}

func TestUnitNextBucket(t *testing.T) {
	var tests = []struct {
		currentBucket, expected int
	}{
		{0, 1},
		{MaxBucketValue, 0},
	}

	for _, tt := range tests {
		testName := fmt.Sprintf("%d", tt.currentBucket)
		t.Run(testName, func(t *testing.T) {
			result := nextBucket(tt.currentBucket)
			if result != tt.expected {
				t.Errorf("Expected result to be %d, got %d", tt.expected, result)
			}
		})
	}
}
