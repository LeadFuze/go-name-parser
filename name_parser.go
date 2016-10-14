package names

import (
	"strings"

	"github.com/blendlabs/go-util"
	"github.com/blendlabs/go-util/collections"
)

var validSuffixes collections.StringArray = []string{
	"I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX", "X", "XI", "XII", "XIII", "XIV", "XV", "XVI", "XVII", "XVIII", "XIX", "XX",
	"Senior", "Junior", "Jr", "Sr",
	"PhD", "APR", "RPh", "PE", "MD", "MA", "DMD", "CME",
}

var compoundLastNames collections.StringArray = []string{
	"vere", "von", "van", "de", "del", "della", "di", "da", "pietro",
	"vanden", "du", "st.", "st", "la", "lo", "ter", "bin", "ibn",
}

// Name is a structured/parsed name.
type Name struct {
	Salutation string
	FirstName  string
	MiddleName string
	LastName   string
	Suffix     string
}

// String returns the string representation of a name.
func (n Name) String() string {
	fullName := util.StringEmpty

	if !util.String.IsEmpty(n.Salutation) {
		fullName = fullName + n.Salutation
	}

	if !util.String.IsEmpty(n.FirstName) {
		if !util.String.IsEmpty(fullName) {
			fullName = fullName + " "
		}
		fullName = fullName + n.FirstName
	}

	if !util.String.IsEmpty(n.MiddleName) {
		if !util.String.IsEmpty(fullName) {
			fullName = fullName + " "
		}
		fullName = fullName + n.MiddleName
	}

	if !util.String.IsEmpty(n.LastName) {
		if !util.String.IsEmpty(fullName) {
			fullName = fullName + " "
		}
		fullName = fullName + n.LastName
	}
	if !util.String.IsEmpty(n.Suffix) {
		if !util.String.IsEmpty(fullName) {
			fullName = fullName + " "
		}
		fullName = fullName + n.Suffix
	}

	return fullName
}

// Parse parses a string into a name.
func Parse(input string) *Name {
	fullName := util.String.TrimWhitespace(input)

	rawNameParts := strings.Split(fullName, " ")

	name := new(Name)

	nameParts := []string{}

	lastName := util.StringEmpty
	firstName := util.StringEmpty
	initials := util.StringEmpty
	for _, part := range rawNameParts {
		if !strings.Contains(part, "(") {
			nameParts = append(nameParts, part)
		}
	}

	numWords := len(nameParts)
	salutation := processSalutation(nameParts[0])
	suffix := processSuffix(nameParts[len(nameParts)-1])

	start := 0
	if !util.String.IsEmpty(salutation) {
		start = 1
	}

	end := numWords
	if !util.String.IsEmpty(suffix) {
		end = numWords - 1
	}

	i := 0
	for i = start; i < (end - 1); i++ {
		word := nameParts[i]
		if isCompoundLastName(word) && i != start {
			break
		}
		if isMiddleName(word) {
			if i == start {
				if isMiddleName(nameParts[i+1]) {
					firstName = firstName + " " + strings.ToUpper(word)
				} else {
					initials = initials + " " + strings.ToUpper(word)
				}
			} else {
				initials = initials + " " + strings.ToUpper(word)
			}
		} else {
			firstName = firstName + " " + fixCase(word)
		}
	}

	if (end - start) > 1 {
		for j := i; j < end; j++ {
			lastName = lastName + " " + fixCase(nameParts[j])
		}
	} else {
		firstName = fixCase(nameParts[i])
	}

	name.Salutation = salutation
	name.FirstName = util.String.TrimWhitespace(firstName)
	name.MiddleName = util.String.TrimWhitespace(initials)
	name.LastName = util.String.TrimWhitespace(lastName)
	name.Suffix = suffix

	return name
}

func processSalutation(input string) string {
	word := cleanString(input)

	switch word {
	case "mr", "master", "mister":
		return "Mr."
	case "mrs", "misses":
		return "Mrs."
	case "ms", "miss":
		return "Ms."
	case "dr":
		return "Dr."
	case "rev":
		return "Rev."
	case "fr":
		return "Fr."
	}

	return util.StringEmpty
}

func processSuffix(input string) string {
	word := cleanString(input)
	return validSuffixes.GetByLower(word)
}

func isCompoundLastName(input string) bool {
	word := cleanString(input)
	exists := compoundLastNames.ContainsLower(word)
	return exists
}

func isMiddleName(input string) bool {
	word := cleanString(input)
	return len(word) == 1
}

func uppercaseFirstAll(input string, seperator string) string {
	words := []string{}
	parts := strings.Split(input, seperator)
	for _, thisWord := range parts {
		toAppend := util.StringEmpty
		if util.String.IsCamelCase(thisWord) {
			toAppend = thisWord
		} else {
			toAppend = strings.ToLower(upperCaseFirst(thisWord))
		}
		words = append(words, toAppend)
	}
	return strings.Join(words, seperator)
}

func upperCaseFirst(input string) string {
	return strings.Title(strings.ToLower(input))
}

func fixCase(input string) string {
	word := uppercaseFirstAll(input, "-")
	word = uppercaseFirstAll(word, ".")
	return word
}

func cleanString(input string) string {
	return strings.ToLower(strings.Replace(input, ".", "", -1))
}
