package preprocessing

import (
	"strings"
	"regexp"
	"math"
)


func applyNormalizingMap(normalizer map[string] string, doc Document) Document {
	var processedContent []byte = []byte(doc.DocContent)

	for key, value := range normalizer {
		normalizingCandidateMatcher := regexp.MustCompile(strings.Join(strings.Split(value, " "), "|"))
		
		processedContent = normalizingCandidateMatcher.ReplaceAll(processedContent, []byte(key))
	}

	return Document {
		Id: doc.Id,
		DocUrl: doc.DocUrl,
		DocContent: string(processedContent),
	}
}

type unicodeReplacementPersianNormalizer struct {
	persianCharsNormalizationMap map[string] string
}

func NewUnicodeReplacementPersianNormalizer() unicodeReplacementPersianNormalizer {
	return unicodeReplacementPersianNormalizer {
		persianCharsNormalizationMap :persianNormalizationMap,
	}
}

func (normalizer *unicodeReplacementPersianNormalizer) Process(document Document) Document{
	return applyNormalizingMap(normalizer.persianCharsNormalizationMap, document)
}

type persianDigitNormalizer struct {
	persianDigitsNormalizationMap map[string] string
}

func NewPersianDigitNormalizer() persianDigitNormalizer {
	return persianDigitNormalizer {
		persianDigitsNormalizationMap :digitNormalizationMap,
	}
}

func (normalizer *persianDigitNormalizer) Process(document Document) Document{
	return applyNormalizingMap(normalizer.persianDigitsNormalizationMap, document)
}

type puncutationRemover struct {
	knownPunctuatuions	string
}

func NewPunctuationRemover() puncutationRemover {
	return puncutationRemover {
		knownPunctuatuions: knownPunctuatuions,
	}
}

func (normalizer *puncutationRemover) Process(document Document) Document {
	normalizingCandidateMatcher := regexp.MustCompile(strings.Join(strings.Split(normalizer.knownPunctuatuions, " "), "|"))
		
	processedContent := normalizingCandidateMatcher.ReplaceAll([]byte(document.DocContent), []byte(""))
	
	return Document {
		Id: document.Id,
		DocUrl: document.DocUrl,
		DocContent: string(processedContent),
	} 
}

type specialArabicPhraseNormalizer struct {
	specialArabicPhraseMap	map[string] string
}

func NewspecialArabicPhraseNormalizer() specialArabicPhraseNormalizer {
	return specialArabicPhraseNormalizer {
		specialArabicPhrases,
	}
}

func (normalizer *specialArabicPhraseNormalizer) Process(document Document) Document {
	return applyNormalizingMap(normalizer.specialArabicPhraseMap, document)
}

type mostUsedWordRemover struct {
	wordFreqMap map[string] int
}

func NewMostUsedWordRemover() mostUsedWordRemover {
	return mostUsedWordRemover {
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

func (normalizer *mostUsedWordRemover) Process(documentCollection TkDocumentCollection) TkDocumentCollection {
	for _, toknizedDoc := range documentCollection.DocList {
		for _, token := range toknizedDoc.TokenzedDocContent {
			_, contains := normalizer.wordFreqMap[token]
			if contains && len(normalizer.wordFreqMap) >= 50 {
				minKey,_ := minKey(normalizer.wordFreqMap)
				delete(normalizer.wordFreqMap, minKey)	
			}
			normalizer.wordFreqMap[token] = 1
		}
	}

	i := 0 // output index
	for _, toknizedDoc := range documentCollection.DocList {
		for _, token := range toknizedDoc.TokenzedDocContent {
			if _, isInMostUsed := normalizer.wordFreqMap[token]; !isInMostUsed {
				toknizedDoc.TokenzedDocContent[i] = token
				i++
			}
		}
	}

	return documentCollection
}