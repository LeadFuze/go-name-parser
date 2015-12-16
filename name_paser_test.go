package util

import (
	"testing"

	"github.com/blendlabs/go-assert"
)

func TestNames(t *testing.T) {
	assert := assert.New(t)

	np := NameParser{}

	names := map[string]NameParserResult{}
	names["John Doe"] = NameParserResult{"", "John", "", "Doe", ""}
	names["Mr Anthony R Von Fange III"] = NameParserResult{"Mr.", "Anthony", "R", "Von Fange", "III"}
	names["Sara Ann Fraser"] = NameParserResult{"", "Sara Ann", "", "Fraser", ""}
	names["Adam"] = NameParserResult{"", "Adam", "", "", ""}
	names["Jonathan Smith"] = NameParserResult{"", "Jonathan", "", "Smith", ""}
	names["Anthony R Von Fange III"] = NameParserResult{"", "Anthony", "R", "Von Fange", "III"}
	names["Anthony Von Fange III"] = NameParserResult{"", "Anthony", "", "Von Fange", "III"}
	names["Mr John Doe"] = NameParserResult{"Mr.", "John", "", "Doe", ""}
	names["Justin White Phd"] = NameParserResult{"", "Justin", "", "White", "PhD"}
	names["Mark P Williams"] = NameParserResult{"", "Mark", "P", "Williams", ""}
	names["Aaron bin Omar"] = NameParserResult{"", "Aaron", "", "bin Omar", ""}
	names["Aaron ibn Omar"] = NameParserResult{"", "Aaron", "", "ibn Omar", ""}

	for rawName, expectedResult := range names {
		result := np.Parse(rawName)
		assert.Equal(expectedResult, result)
	}
}
