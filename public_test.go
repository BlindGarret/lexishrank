package lexishrank

import (
	"fmt"
	"slices"
	"sort"
	"testing"
)

func TestUnitDissectRank(t *testing.T) {
	var tests = []struct {
		rank           string
		expectedBucket int
		expectedId     string
		expectedErr    error
	}{
		{"0|B", 0, "B", nil},
		{"12|AABCSQCCS", 12, "AABCSQCCS", nil},
		{"Not A Rank", 0, "", ErrRankFormatInvalid},
		{"notanum|AASS", 0, "", ErrRankFormatInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.rank, func(t *testing.T) {
			bucket, id, err := DissectRank(tt.rank)
			if bucket != tt.expectedBucket {
				t.Errorf("Expected bucket to be %d, got %d", tt.expectedBucket, bucket)
			}
			if id != tt.expectedId {
				t.Errorf("Expected id to be %s, got %s", tt.expectedId, id)
			}
			if err != tt.expectedErr {
				t.Errorf("Expected error to be %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestIntegrationIndexerReturnsExpectedDirection(t *testing.T) {
	var tests = []struct {
		currentBucket int
		expected      Direction
	}{
		{0, EndToBegining},
		{MaxBucketValue, BeginningToEnd},
	}

	for _, tt := range tests {
		testName := fmt.Sprintf("%d", tt.currentBucket)
		t.Run(testName, func(t *testing.T) {
			indexer, err := NewIndexer(1, 100, tt.currentBucket)
			if err != nil {
				t.Errorf("Expected error to be nil, got %v", err)
			}

			if indexer.Direction() != tt.expected {
				t.Errorf("Expected direction to be %v, got %v", tt.expected, indexer.Direction())
			}
		})
	}
}

func TestIntegrationIndexerReturnsExpectedNextForKnownSetFullCoverage(t *testing.T) {
	expectedFirst := "1|N" // It is important this is NOT A
	expectedSecond := "1|Z"

	indexer, err := NewIndexer(2, 2, 0)
	if err != nil {
		t.Errorf("Expected error to be nil, got %v", err)
	}

	// pulling these values in reverse for readability
	second := indexer.Next()
	first := indexer.Next()

	if first != expectedFirst {
		t.Errorf("Expected first to be %s, got %s", expectedFirst, first)
	}

	if second != expectedSecond {
		t.Errorf("Expected second to be %s, got %s", expectedSecond, second)
	}
}

func TestIntegrationIndexerReturnsExpectedNextForKnownSetPartialCoverage(t *testing.T) {
	expectedFirst := "0|AAB" // It is important this is NOT AAA
	expectedSecond := "0|FFG"
	expectedThird := "0|KKL"
	expectedFourth := "0|PPQ"
	expectedFifth := "0|UUV"

	indexer, err := NewIndexer(1000, 5, MaxBucketValue)
	if err != nil {
		t.Errorf("Expected error to be nil, got %v", err)
	}

	first := indexer.Next()
	second := indexer.Next()
	third := indexer.Next()
	fourth := indexer.Next()
	fifth := indexer.Next()

	if first != expectedFirst {
		t.Errorf("Expected first to be %s, got %s", expectedFirst, first)
	}

	if second != expectedSecond {
		t.Errorf("Expected second to be %s, got %s", expectedSecond, second)
	}

	if third != expectedThird {
		t.Errorf("Expected third to be %s, got %s", expectedThird, third)
	}

	if fourth != expectedFourth {
		t.Errorf("Expected fourth to be %s, got %s", expectedFourth, fourth)
	}

	if fifth != expectedFifth {
		t.Errorf("Expected fifth to be %s, got %s", expectedFifth, fifth)
	}
}

func TestIntegrationIndexerReturnsErrorOnHugeObjectCount(t *testing.T) {
	_, err := NewIndexer(1, 18446744073709551615, 0)
	if err != ErrGapSizeToLargeForObjectCount {
		t.Errorf("Expected error to be %v, got %v", ErrGapSizeToLargeForObjectCount, err)
	}
}

func TestUnitBetween(t *testing.T) {
	var tests = []struct {
		first, second, expected string
		expectedErr             error
	}{
		{"0|BA", "0|BA", "0|BAM", nil},
		{"0|AQXKKH", "0|AUHHHD", "0|ASPIVS", nil},
		{"0|AXREDZ", "0|BBBBAV", "0|AZJCPK", nil},
		{"0|BA", "0|BAM", "0|BAG", nil},
		{"1|AAB", "1|BZZ", "1|BAA", nil},
		{"2|AAB", "2|ZZW", "2|MZY", nil},
		{"2|BAAAAAAAAAAAA", "2|BAAAAAAAAAAAAB", "", ErrGapSizeToLargeForObjectCount},
	}

	for _, tt := range tests {
		testName := fmt.Sprintf("%s,%s", tt.first, tt.second)
		t.Run(testName, func(t *testing.T) {
			result, err := Between(tt.first, tt.second)
			if result != tt.expected {
				t.Errorf("Expected result to be %s, got %s", tt.expected, result)
			}
			if err != tt.expectedErr {
				t.Errorf("Expected error to be %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestUnitBetweenReturnsErrorOnInvalidRanks(t *testing.T) {
	_, err := Between("Not A Rank", "0|B")
	if err != ErrRankFormatInvalid {
		t.Errorf("Expected error to be %v, got %v", ErrRankFormatInvalid, err)
	}

	_, err = Between("0|B", "Not A Rank")
	if err != ErrRankFormatInvalid {
		t.Errorf("Expected error to be %v, got %v", ErrRankFormatInvalid, err)
	}
}

func TestUnitNext(t *testing.T) {
	var tests = []struct {
		currentMaxRank string
		maxLengthRank  string
		gapSize        uint64
		expected       string
		expectedErr    error
	}{
		{"0|BA", "0|BAA", 1, "0|BAB", nil},
		{"0|BZ", "0|BAA", 1, "0|BZB", nil},
		{"0|BA", "0|BAAAA", 1000, "0|BABMM", nil},
		{"0|BA", "0|BAA", 18446744073709551615, "", ErrGapSizeToLargeForObjectCount},
		{"0|ZZZZZZZZZZZZZ", "0|ZZZZZZZZZZZ", 1, "", ErrGapSizeToLargeForObjectCount},
		{"NotAVal", "0|BAA", 1, "", ErrRankFormatInvalid},
		{"0|B", "NotAVal", 1, "", ErrRankFormatInvalid},
	}
	for _, tt := range tests {
		testName := fmt.Sprintf("%s,%d", tt.currentMaxRank, tt.gapSize)
		t.Run(testName, func(t *testing.T) {
			result, err := After(tt.currentMaxRank, tt.maxLengthRank, tt.gapSize)
			if err != tt.expectedErr {
				t.Errorf("Expected error to be %v, got %v", tt.expectedErr, err)
			}
			if result != tt.expected {
				t.Errorf("Expected result to be %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestIntegrationIndexerNextValuesSortLexicallyAsExpected(t *testing.T) {
	indexer, err := NewIndexer(1, 3, 1)
	if err != nil {
		t.Errorf("Expected error to be nil, got %v", err)
	}

	// Remember these pull backwards since we're indexing for 2
	third := indexer.Next()
	second := indexer.Next()
	first := indexer.Next()

	expected := []string{first, second, third}
	sorted := []string{first, second, third}
	sort.Strings(sorted)

	if !slices.Equal(expected, sorted) {
		t.Errorf("Expected %v, got %v", expected, sorted)
	}
}

func TestIntegrationNextBetweenValuesSortLexically(t *testing.T) {
	const first = "0|BA"
	const third = "0|BB"
	second, err := Between(first, third)
	if err != nil {
		t.Errorf("Expected error to be nil, got %v", err)
	}

	fourth, err := After(third, second, 100)
	if err != nil {
		t.Errorf("Expected error to be nil, got %v", err)
	}
	expected := []string{first, second, third, fourth}
	sorted := []string{first, second, third, fourth}
	sort.Strings(sorted)

	if !slices.Equal(expected, sorted) {
		t.Errorf("Expected %v, got %v", expected, sorted)
	}
}

func TestIntegrationNextValuesSortLexically(t *testing.T) {
	const start = "0|BA"
	first, err := After(start, start, 10000)
	if err != nil {
		t.Errorf("Expected error to be nil, got %v", err)
	}

	second, err := After(first, first, 10000)
	if err != nil {
		t.Errorf("Expected error to be nil, got %v", err)
	}

	third, err := After(second, second, 10000)
	if err != nil {
		t.Errorf("Expected error to be nil, got %v", err)
	}

	expected := []string{start, first, second, third}
	sorted := []string{start, first, second, third}
	sort.Strings(sorted)

	if !slices.Equal(expected, sorted) {
		t.Errorf("Expected %v, got %v", expected, sorted)
	}
}

func TestIntegrationBetweenValuesSortLexically(t *testing.T) {
	const start = "0|AAAAAAB"
	const end = "0|NGHWHHBS"
	between, err := Between(start, end)
	if err != nil {
		t.Errorf("Expected error to be nil, got %v", err)
	}

	expected := []string{start, between, end}
	sorted := []string{start, between, end}
	sort.Strings(sorted)

	if !slices.Equal(expected, sorted) {
		t.Errorf("Expected %v, got %v", expected, sorted)
	}

}

func TestIntegrationComplicatedUsageValuesStillSortLexically(t *testing.T) {
	const initialObjectCount = 200
	const gapSize = 100000
	indexer, err := NewIndexer(gapSize, initialObjectCount, 2)
	if err != nil {
		t.Errorf("Expected error to be nil, got %v", err)
	}

	set := make([]string, 0)

	// create initial 200 objects
	for i := 0; i < initialObjectCount; i++ {
		set = append(set, indexer.Next())
	}

	// swap 10 of them into singular swaps
	set[30], err = Between(set[1], set[2])
	if err != nil {
		t.Errorf("Expected error to be nil, got %v", err)
	}

	set[40], err = Between(set[3], set[4])
	if err != nil {
		t.Errorf("Expected error to be nil, got %v", err)
	}
	set[50], err = Between(set[5], set[6])
	if err != nil {
		t.Errorf("Expected error to be nil, got %v", err)
	}
	set[60], err = Between(set[7], set[8])
	if err != nil {
		t.Errorf("Expected error to be nil, got %v", err)
	}
	set[70], err = Between(set[9], set[10])
	if err != nil {
		t.Errorf("Expected error to be nil, got %v", err)
	}
	set[80], err = Between(set[11], set[12])
	if err != nil {
		t.Errorf("Expected error to be nil, got %v", err)
	}
	set[90], err = Between(set[13], set[14])
	if err != nil {
		t.Errorf("Expected error to be nil, got %v", err)
	}
	set[100], err = Between(set[15], set[16])
	if err != nil {
		t.Errorf("Expected error to be nil, got %v", err)
	}
	set[110], err = Between(set[17], set[18])
	if err != nil {
		t.Errorf("Expected error to be nil, got %v", err)
	}
	set[120], err = Between(set[19], set[20])
	if err != nil {
		t.Errorf("Expected error to be nil, got %v", err)
	}

	// add 20 more objects
	maxRank := set[199]
	newSet := make([]string, 0)
	for i := 0; i < 168; i++ {
		rank, err := After(maxRank, maxRank, gapSize)
		if err != nil {
			t.Errorf("Expected error to be nil, got %v", err)
		}
		newSet = append(newSet, rank)
		maxRank = rank
	}

	// insert all new ranks between two specific values recursively
	betweenEnd := set[141]
	for i := 0; i < len(newSet); i++ {
		newSet[i], err = Between(set[140], betweenEnd)
		if err != nil {
			t.Errorf("Expected error to be nil, got %v", err)
		}
		betweenEnd = newSet[i]
	}

	// Expectations
	slices.Reverse(newSet) // for easier expected
	expected := append(
		[]string{},
		set[0],
		set[1], set[30], set[2],
		set[3], set[40], set[4],
		set[5], set[50], set[6],
		set[7], set[60], set[8],
		set[9], set[70], set[10],
		set[11], set[80], set[12],
		set[13], set[90], set[14],
		set[15], set[100], set[16],
		set[17], set[110], set[18],
		set[19], set[120], set[20])
	expected = append(expected, set[21:30]...)
	expected = append(expected, set[31:40]...)
	expected = append(expected, set[41:50]...)
	expected = append(expected, set[51:60]...)
	expected = append(expected, set[61:70]...)
	expected = append(expected, set[71:80]...)
	expected = append(expected, set[81:90]...)
	expected = append(expected, set[91:100]...)
	expected = append(expected, set[101:110]...)
	expected = append(expected, set[111:120]...)
	expected = append(expected, set[121:141]...)
	expected = append(expected, newSet...)
	expected = append(expected, set[141:200]...)

	actual := append(set, newSet...)
	sort.Strings(actual)

	if len(expected) != len(actual) {
		t.Errorf("Expected %d, got %d", len(expected), len(actual))
	}

	if !slices.Equal(expected, actual) {
		t.Errorf("Expected %v,\n got %v", expected, actual)
	}

	fmt.Printf("Actual: %v\n", actual)
}
