package base26

import "testing"

func TestUnitEncode(t *testing.T) {
	var tests = []struct {
		input    uint64
		expected string
	}{
		{1, "B"},
		{25, "Z"},
		{26, "BA"},
		{27, "BB"},
		{51, "BZ"},
		{52, "CA"},
		{701, "BAZ"},
		{702, "BBA"},
		{703, "BBB"},
		{3139, "EQT"},
		{4000, "FXW"},
		{475255, "BBBBB"},
		{18446744073709551615, "HLHXCZMXSYUMQP"},
	}
	for _, test := range tests {
		if output := Encode(test.input); output != test.expected {
			t.Errorf("Test failed: input %v, expected %v, got %v", test.input, test.expected, output)
		}
	}
}

func TestUnitDecode(t *testing.T) {
	var tests = []struct {
		input    string
		expected uint64
	}{
		{"B", 1},
		{"Z", 25},
		{"BA", 26},
		{"BB", 27},
		{"BZ", 51},
		{"CA", 52},
		{"BAZ", 701},
		{"BBA", 702},
		{"BBB", 703},
		{"EQT", 3139},
		{"FXW", 4000},
		{"BBBBB", 475255},
		{"HLHXCZMXSYUMQP", 18446744073709551615},
	}
	for _, test := range tests {
		if output := Decode(test.input); output != test.expected {
			t.Errorf("Test failed: input %v, expected %v, got %v", test.input, test.expected, output)
		}
	}
}
