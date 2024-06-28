package preprocessing

import (
	"regexp"
	"strings"
)

var PERSIAN_WORD_SUFFIXES = []string{"ات", "ان", "ترین", "تر", "م", "ت", "ش", "یی", "ی", "ها", "ٔ", "‌ا"}

type Stemmer interface {
	Stem(word string) string
}

type persianStemmer struct {
	WordSuffixes []string
}

func NewSimplePersianStemmer() persianStemmer {
	return persianStemmer{
		WordSuffixes: PERSIAN_WORD_SUFFIXES,
	}
}

func (stemmer *persianStemmer) Stem(word string) string {
	iteration := len(stemmer.WordSuffixes)
	for iteration > 0 {
		for _, end := range stemmer.WordSuffixes {
			if strings.HasSuffix(word, end) {
				word = word[:len(word)-len(end)]
				break
			}
		}
		iteration--
	}
	return word
}

type ComplexStemmer struct {
	patterns  []*regexp.Regexp
	stemIndex []int
}

func NewComplexPersianStemmer() *ComplexStemmer {
	s := &ComplexStemmer{
		patterns:  []*regexp.Regexp{},
		stemIndex: []int{},
	}

	// Mazie Eltezami va Mazie Eltezamie Manfi
	s.addPattern("ن?(.*)ه (باشم|باشی|باشد|باشیم|باشید|باشند)$", 1)

	// Mazie Baeid va Mazie Baeide Manfi
	s.addPattern("ن?(.*)ه (بودم|بودی|بود|بودیم|بودید|بودند)$", 1)

	// Mazie Estemrari, Mazie Estemrarie Manfi, Mozare'e Ekhbari va Mozare'e Ekhbarie Manfi
	s.addPattern("(می|نمی)(.*)(یم|ید|ند)$", 2)
	s.addPattern("(می|نمی)(.*)(م|ی|د)$", 2)
	s.addPattern("(می|نمی)(.*)$", 2)

	// Mazie Naqli va Mazie Naqlie Manfi
	s.addPattern("ن?(.*)ه (ام|ای|است|ایم|اید|اند)$", 1)

	//// Mazie Sade va Mazie Sadeye Manfi
	s.addPattern("ن?(.*)(یم|ید|ند)$", 1)
	s.addPattern("ن?(.*)(م|ی)$", 1)

	return s
}

func (s *ComplexStemmer) addPattern(pattern string, index int) {
	re := regexp.MustCompile(pattern)
	s.patterns = append(s.patterns, re)
	s.stemIndex = append(s.stemIndex, index)
}

func (s *ComplexStemmer) Stem(word string) string {
	for i, pattern := range s.patterns {
		matches := pattern.FindStringSubmatch(word)
		if len(matches) > 0 {
			return matches[s.stemIndex[i]]
		}
	}
	return word
}
