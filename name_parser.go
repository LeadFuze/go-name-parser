package names

import (
	"strings"
	"unicode"
)

const (
	EMPTY = ""
)

var validSuffixes stringArray = []string{
	"I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX", "X", "XI", "XII", "XIII", "XIV", "XV", "XVI", "XVII", "XVIII", "XIX", "XX",
	"Senior", "Junior", "Jr", "Sr",
	"PhD", "APR", "RPh", "PE", "MD", "MA", "DMD", "CME",
}

var compoundLastNames stringArray = []string{
	"vere", "von", "van", "de", "del", "della", "di", "da", "pietro",
	"vanden", "du", "st.", "st", "la", "lo", "ter", "bin", "ibn",
}

type Name struct {
	Salutation string
	FirstName  string
	MiddleName string
	LastName   string
	Suffix     string
}

func (n Name) String() string {
	fullName := EMPTY

	if !isEmpty(n.Salutation) {
		fullName = fullName + n.Salutation
	}

	if !isEmpty(n.FirstName) {
		if !isEmpty(fullName) {
			fullName = fullName + " "
		}
		fullName = fullName + n.FirstName
	}

	if !isEmpty(n.MiddleName) {
		if !isEmpty(fullName) {
			fullName = fullName + " "
		}
		fullName = fullName + n.MiddleName
	}

	if !isEmpty(n.LastName) {
		if !isEmpty(fullName) {
			fullName = fullName + " "
		}
		fullName = fullName + n.LastName
	}
	if !isEmpty(n.Suffix) {
		if !isEmpty(fullName) {
			fullName = fullName + " "
		}
		fullName = fullName + n.Suffix
	}

	return fullName
}

func Parse(input string) *Name {
	fullName := trimWhitespace(input)

	rawNameParts := strings.Split(fullName, " ")

	name := new(Name)

	nameParts := []string{}

	lastName := EMPTY
	firstName := EMPTY
	initials := EMPTY
	for _, part := range rawNameParts {
		if !strings.Contains(part, "(") {
			nameParts = append(nameParts, part)
		}
	}

	numWords := len(nameParts)
	salutation := processSalutation(nameParts[0])
	suffix := processSuffix(nameParts[len(nameParts)-1])

	start := 0
	if !isEmpty(salutation) {
		start = 1
	}

	end := numWords
	if !isEmpty(suffix) {
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
	name.FirstName = trimWhitespace(firstName)
	name.MiddleName = trimWhitespace(initials)
	name.LastName = trimWhitespace(lastName)
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

	return EMPTY
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
		toAppend := EMPTY
		if isCamelCase(thisWord) {
			toAppend = thisWord
		} else {
			toAppend = strings.ToLower(upperCaseFirst(thisWord))
		}
		words = append(words, toAppend)
	}
	return strings.Join(words, seperator)
}

func fixCase(input string) string {
	word := uppercaseFirstAll(input, "-")
	word = uppercaseFirstAll(word, ".")
	return word
}

func cleanString(input string) string {
	return strings.ToLower(strings.Replace(input, ".", "", -1))
}

func trimWhitespace(input string) string {
	return strings.Trim(input, " \t")
}

func upperCaseFirst(input string) string {
	return strings.Title(strings.ToLower(input))
}

func isCamelCase(input string) bool {
	hasLowers := false
	hasUppers := false

	for _, c := range input {
		if unicode.IsUpper(c) {
			hasUppers = true
		}
		if unicode.IsLower(c) {
			hasLowers = true
		}
	}

	return hasLowers && hasUppers
}

func isEmpty(input string) bool {
	return len(input) == 0
}

type stringArray []string

func (sa stringArray) Contains(elem string) bool {
	for _, arrayElem := range sa {
		if arrayElem == elem {
			return true
		}
	}
	return false
}

func (sa stringArray) ContainsLower(elem string) bool {
	for _, arrayElem := range sa {
		if strings.ToLower(arrayElem) == elem {
			return true
		}
	}
	return false
}

func (sa stringArray) GetByLower(elem string) string {
	for _, arrayElem := range sa {
		if strings.ToLower(arrayElem) == elem {
			return arrayElem
		}
	}
	return EMPTY
}
