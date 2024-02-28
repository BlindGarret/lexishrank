package lexishrank

import "fmt"

var ErrGapSizeToLargeForObjectCount = fmt.Errorf("minGapSize is too large for objectCount")
var ErrRankFormatInvalid = fmt.Errorf("rank format is invalid")
