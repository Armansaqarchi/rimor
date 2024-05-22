package preprocessing

import (
	"strings"
)

var PERSIAN_WORD_SUFFIXES = []string {"ات", "ان", "ترین", "تر", "م", "ت", "ش", "یی", "ی", "ها", "ٔ", "‌ا"}

type Stemmer interface {
	Stem(word string) string
}

type persianStemmer struct {
	WordSuffixes []string
}

func NewPersianStemmer() persianStemmer {
	return persianStemmer {
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