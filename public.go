package lexishrank

import (
	"strconv"
	"strings"

	"github.com/BlindGarret/lexishrank/base26"
)

// MaxBucketValue is the maximum value a bucket can be before it wraps around to 0, default is 2
var MaxBucketValue = 2

// Direction is the direction of movement for reindexing through the set. Whether it's going from begining to end or end to begining will matter based on which bucket we are moving to
type Direction string

const (
	BeginningToEnd = "beginingToEnd"
	EndToBegining  = "endToBegining"
)

// Indexer is a lexorank index generator for reindexing objects to ensure they are spread throughout the lexorank space
type Indexer struct {
	minGapSize        uint64
	currentId         uint64
	gapSize           uint64
	newBucket         int
	directionFunction directionFunction
	direction         Direction
	idLength          int
}

// NewIndexer creates a new indexer for reindexing objects in the lexorank space
func NewIndexer(minGapSize uint64, objectCount uint64, currentBucket int) (*Indexer, error) {
	gapSize, maxVal, err := calculateRankGapSize(objectCount, minGapSize)
	if err != nil {
		return nil, err
	}
	newBucket := nextBucket(currentBucket)
	directionFunc, direction := getDirection(newBucket)
	var currentId uint64 = 1
	if direction == EndToBegining {
		currentId = maxVal
	}

	maxId := base26.Encode(maxVal)

	return &Indexer{
		minGapSize:        minGapSize,
		gapSize:           gapSize,
		newBucket:         newBucket,
		directionFunction: directionFunc,
		direction:         direction,
		currentId:         currentId,
		idLength:          len(maxId),
	}, err
}

// Next returns the next rank in the lexorank space
func (i *Indexer) Next() string {
	nextRank := formatRank(i.newBucket, base26.Encode(i.currentId), i.idLength)
	i.currentId = i.directionFunction(i.currentId, i.gapSize)
	return nextRank
}

// Direction returns the direction you should travel through your set during reindexing
func (i *Indexer) Direction() Direction {
	return i.direction
}

// DissectRank is a helper function which takes a rank and returns the bucket and id of the rank
func DissectRank(rank string) (bucket int, id string, err error) {
	parts := strings.Split(rank, "|")
	if len(parts) != 2 || parts[1] == "" {
		return 0, "", ErrRankFormatInvalid
	}
	bucket, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", ErrRankFormatInvalid // this hides the ATOI error but it unifies the interface
	}
	id = parts[1]
	return
}

// Between returns the rank that is between the two provided ranks
func Between(first string, second string) (string, error) {
	firstBucket, firstID, err := DissectRank(first)
	if err != nil {
		return "", err
	}

	_, secondID, err := DissectRank(second)
	if err != nil {
		return "", err
	}

	rank := calculateMiddleRank(firstID, secondID)
	if len(rank) > 100 {
		// This is most likely an overflow. There exists some set of values where this isn't but testing for them is costly,
		// as normal overflow tests don't work. For example a number starting with Z which is 14 digits long has overflows 5 times or more,
		// meaning it could be greater than or less than the original value when it's done.
		return "", ErrGapSizeToLargeForObjectCount
	}
	return formatRank(
		firstBucket,
		calculateMiddleRank(firstID, secondID),
		max(len(firstID), len(secondID)),
	), nil
}

func After(currentMaxRank string, widestRank string, stepSize uint64) (string, error) {
	bucket, id, err := DissectRank(currentMaxRank)
	if err != nil {
		return "", err
	}
	_, wideId, err := DissectRank(widestRank)
	if err != nil {
		return "", err
	}

	maxRank, _ := padLowerValue(id, wideId)
	if len(maxRank) > 100 {
		// This is most likely an overflow. There exists some set of values where this isn't but testing for them is costly,
		// as normal overflow tests don't work. For example a number starting with Z which is 14 digits long has overflows 5 times or more,
		// meaning it could be greater than or less than the original value when it's done.
		return "", ErrGapSizeToLargeForObjectCount
	}
	maxRankValue := base26.Decode(maxRank)

	newRankValue := maxRankValue + stepSize
	newRank := base26.Encode(newRankValue)
	if newRankValue < stepSize || newRankValue < maxRankValue {
		// We rolled over
		return "", ErrGapSizeToLargeForObjectCount
	}
	if len(newRank) > len(maxRank) {
		// We rolled over a digit which will mess up the lexical sorting,
		// run this again but with a wider widestRank.
		// This should only ever recurse once.
		return After(currentMaxRank, formatRank(bucket, newRank, len(id)), stepSize)

	}

	return formatRank(bucket, newRank, len(id)), nil
}
