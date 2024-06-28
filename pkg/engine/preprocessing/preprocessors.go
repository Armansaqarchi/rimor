package preprocessing

import (
	"log"
	"regexp"
	"sort"
	"strings"
)

func applyNormalizingMap(normalizer map[string]string, text string) string {
	var processedContent []byte = []byte(text)

	for key, value := range normalizer {
		normalizingCandidateMatcher := regexp.MustCompile(strings.Join(strings.Split(value, " "), "|"))

		processedContent = normalizingCandidateMatcher.ReplaceAll(processedContent, []byte(key))
	}

	return string(processedContent)
}

type UnicodeReplacementPersianNormalizer struct {
	persianCharsNormalizationMap map[string]string
}

func NewUnicodeReplacementPersianNormalizer() UnicodeReplacementPersianNormalizer {
	return UnicodeReplacementPersianNormalizer{
		persianCharsNormalizationMap: persianNormalizationMap,
	}
}

func (normalizer *UnicodeReplacementPersianNormalizer) Process(text string) string {
	return applyNormalizingMap(normalizer.persianCharsNormalizationMap, text)
}

type PersianDigitNormalizer struct {
	persianDigitsNormalizationMap map[string]string
}

func NewPersianDigitNormalizer() PersianDigitNormalizer {
	return PersianDigitNormalizer{
		persianDigitsNormalizationMap: digitNormalizationMap,
	}
}

func (normalizer *PersianDigitNormalizer) Process(text string) string {
	return applyNormalizingMap(normalizer.persianDigitsNormalizationMap, text)
}

type PunctuationRemover struct {
	knownPunctuations string
}

func NewPunctuationRemover() PunctuationRemover {
	return PunctuationRemover{
		knownPunctuations: knownPunctuatuions,
	}
}

func (normalizer *PunctuationRemover) Process(text string) string {
	normalizingCandidateMatcher := regexp.MustCompile(strings.Join(strings.Split(normalizer.knownPunctuations, " "), "|"))

	processedContent := normalizingCandidateMatcher.ReplaceAll([]byte(text), []byte(""))

	return string(processedContent)

}

type SpecialArabicPhraseNormalizer struct {
	specialArabicPhraseMap map[string]string
}

func NewSpecialArabicPhraseNormalizer() SpecialArabicPhraseNormalizer {
	return SpecialArabicPhraseNormalizer{
		specialArabicPhrases,
	}
}

func (normalizer *SpecialArabicPhraseNormalizer) Process(text string) string {
	return applyNormalizingMap(normalizer.specialArabicPhraseMap, text)
}

type MostUsedWordRemover struct {
	wordFreqMap map[string]int
	top50Words  map[string]struct{}
}

func NewMostUsedWordRemover() MostUsedWordRemover {
	return MostUsedWordRemover{
		wordFreqMap: make(map[string]int, 50),
	}
}

func (normalizer *MostUsedWordRemover) ProcessQuery(query []string) []string {
	filteredQuery := make([]string, 0)
	for _, token := range query {
		if _, isInTop50 := normalizer.top50Words[token]; !isInTop50 {
			filteredQuery = append(filteredQuery, token)
		}
	}

	return filteredQuery
}

func (normalizer *MostUsedWordRemover) ProcessDocCollection(documentCollection TkDocumentCollection) TkDocumentCollection {
	normalizer.createWordFreqMap(documentCollection)
	normalizer.top50Words = normalizer.createTop50WordsMap()

	for i, document := range documentCollection.DocList {
		var newContent []string
		for _, word := range document.TokenzedDocContent {
			if _, found := normalizer.top50Words[word]; !found {
				newContent = append(newContent, word)
			}
		}
		documentCollection.DocList[i].TokenzedDocContent = newContent
	}
	log.Printf("Top 50 words: %v", normalizer.top50Words)
	return documentCollection
}

func (normalizer *MostUsedWordRemover) createWordFreqMap(collection TkDocumentCollection) {
	for _, document := range collection.DocList {
		for _, word := range document.TokenzedDocContent {
			if _, ok := normalizer.wordFreqMap[word]; ok {
				normalizer.wordFreqMap[word]++
			} else {
				normalizer.wordFreqMap[word] = 1
			}
		}
	}
}

func (normalizer *MostUsedWordRemover) createTop50WordsMap() map[string]struct{} {
	// Step 2: Identify the top 50 most frequently used words
	type wordFreqPair struct {
		word  string
		count int
	}
	var wordFreqPairs []wordFreqPair
	for word, count := range normalizer.wordFreqMap {
		wordFreqPairs = append(wordFreqPairs, wordFreqPair{word, count})
	}

	sort.Slice(wordFreqPairs, func(i, j int) bool {
		return wordFreqPairs[i].count > wordFreqPairs[j].count
	})

	var top50Words map[string]struct{}
	if len(wordFreqPairs) > 50 {
		top50Words = make(map[string]struct{}, 50)
		for i := 0; i < 50; i++ {
			top50Words[wordFreqPairs[i].word] = struct{}{}
		}
	} else {
		top50Words = make(map[string]struct{}, len(wordFreqPairs))
		for _, pair := range wordFreqPairs {
			top50Words[pair.word] = struct{}{}
		}
	}

	return top50Words
}
