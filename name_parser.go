package util

import "strings"

var validSuffixes StringArray = []string{
	"I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX", "X", "XI", "XII", "XIII", "XIV", "XV", "XVI", "XVII", "XVIII", "XIX", "XX",
	"Senior", "Junior", "Jr", "Sr",
	"PhD", "APR", "RPh", "PE", "MD", "MA", "DMD", "CME",
}

var compoundLastNames StringArray = []string{
	"vere", "von", "van", "de", "del", "della", "di", "da", "pietro",
	"vanden", "du", "st.", "st", "la", "lo", "ter", "bin", "ibn",
}

type NameParserResult struct {
	Salutation string
	FirstName  string
	Initials   string
	LastName   string
	Suffix     string
}

func (npr NameParserResult) String() string {
	fullName := ""

	if !IsEmpty(npr.Salutation) {
		fullName = fullName + npr.Salutation
	}

	if !IsEmpty(npr.FirstName) {
		if !IsEmpty(fullName) {
			fullName = fullName + " "
		}
		fullName = fullName + npr.FirstName
	}

	if !IsEmpty(npr.Initials) {
		if !IsEmpty(fullName) {
			fullName = fullName + " "
		}
		fullName = fullName + npr.Initials
	}

	if !IsEmpty(npr.LastName) {
		if !IsEmpty(fullName) {
			fullName = fullName + " "
		}
		fullName = fullName + npr.LastName
	}
	if !IsEmpty(npr.Suffix) {
		if !IsEmpty(fullName) {
			fullName = fullName + " "
		}
		fullName = fullName + npr.Suffix
	}

	return fullName
}

type NameParser struct{}

func (np NameParser) cleanString(input string) string {
	return strings.ToLower(strings.Replace(input, ".", "", -1))
}

func (np NameParser) processSalutation(input string) string {

	word := np.cleanString(input)

	switch word {
	case "mr", "master", "mister":
		return "Mr."
	case "mrs":
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

func (np NameParser) processSuffix(input string) string {
	word := np.cleanString(input)
	return validSuffixes.GetByLower(word)
}

func (np NameParser) isCompoundLastName(input string) bool {
	word := np.cleanString(input)
	exists := compoundLastNames.ContainsLower(word)
	return exists
}

func (np NameParser) isInitial(input string) bool {
	word := np.cleanString(input)
	return len(word) == 1
}

func (np NameParser) safeUppercaseFirst(seperator string, input string) string {
	words := []string{}
	parts := strings.Split(input, seperator)
	for _, thisWord := range parts {
		toAppend := ""
		if IsCamelCase(thisWord) {
			toAppend = thisWord
		} else {
			toAppend = strings.ToLower(UpperCaseFirst(thisWord))
		}
		words = append(words, toAppend)
	}
	return Join(words, seperator)
}

func (np NameParser) fixCase(input string) string {
	word := np.safeUppercaseFirst("-", input)
	word = np.safeUppercaseFirst(".", word)
	return word
}

func (np NameParser) Parse(input string) NameParserResult {
	fulllastName := TrimWhitespace(input)

	unfilteredNameParts := strings.Split(fulllastName, " ")

	name := NameParserResult{}

	nameParts := []string{}

	lastName := ""
	firstName := ""
	initials := ""
	for _, part := range unfilteredNameParts {
		if !strings.Contains(part, "(") {
			nameParts = append(nameParts, part)
		}
	}

	numWords := len(nameParts)
	salutation := np.processSalutation(nameParts[0])
	suffix := np.processSuffix(nameParts[len(nameParts)-1])

	start := 0
	if !IsEmpty(salutation) {
		start = 1
	}

	end := numWords
	if !IsEmpty(suffix) {
		end = numWords - 1
	}

	i := 0
	for i = start; i < (end - 1); i++ {
		word := nameParts[i]
		if np.isCompoundLastName(word) && i != start {
			break
		}
		if np.isInitial(word) {
			if i == start {
				if np.isInitial(nameParts[i+1]) {
					firstName = firstName + " " + strings.ToUpper(word)
				} else {
					initials = initials + " " + strings.ToUpper(word)
				}
			} else {
				initials = initials + " " + strings.ToUpper(word)
			}
		} else {
			firstName = firstName + " " + np.fixCase(word)
		}
	}

	if (end - start) > 1 {
		for j := i; j < end; j++ {
			lastName = lastName + " " + np.fixCase(nameParts[j])
		}
	} else {
		firstName = np.fixCase(nameParts[i])
	}

	name.Salutation = salutation
	name.FirstName = TrimWhitespace(firstName)
	name.Initials = TrimWhitespace(initials)
	name.LastName = TrimWhitespace(lastName)
	name.Suffix = suffix

	return name
}
