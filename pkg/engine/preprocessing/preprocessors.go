package preprocessing

import (
	"math"
	"regexp"
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
}

func NewMostUsedWordRemover() MostUsedWordRemover {
	return MostUsedWordRemover{
		wordFreqMap: make(map[string]int, 50),
	}
}

func minKey(m map[string]int) (string, int) {
	if len(m) == 0 {
		return "", 0 // or handle the empty map case as needed
	}

	minKey := ""
	minValue := math.MaxInt

	for k, v := range m {
		if v < minValue {
			minKey = k
			minValue = v
		}
	}

	return minKey, minValue
}

func (normalizer *MostUsedWordRemover) Process(documentCollection TkDocumentCollection) TkDocumentCollection {
	for _, toknizedDoc := range documentCollection.DocList {
		for _, token := range toknizedDoc.TokenzedDocContent {
			_, contains := normalizer.wordFreqMap[token]
			if contains && len(normalizer.wordFreqMap) >= 50 {
				minKey, _ := minKey(normalizer.wordFreqMap)
				delete(normalizer.wordFreqMap, minKey)
			}
			normalizer.wordFreqMap[token] = 1
		}
	}

	for _, toknizedDoc := range documentCollection.DocList {
		var filteredTokens []string
		for _, token := range toknizedDoc.TokenzedDocContent {
			_, contains := normalizer.wordFreqMap[token]
			if !contains {
				filteredTokens = append(filteredTokens, token)
			}
		}
		toknizedDoc.TokenzedDocContent = filteredTokens
	}

	return documentCollection
}
