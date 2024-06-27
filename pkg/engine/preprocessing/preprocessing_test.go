package preprocessing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var arabicPhrase = NewspecialArabicPhraseNormalizer()
var persianDigit = NewPersianDigitNormalizer()
var unicodeRep = NewUnicodeReplacementPersianNormalizer()


func TestUnicodeReplacementNormalization(t *testing.T) {
	preprocessor := NewPreprocessor([]PreprocessingStep{
		&unicodeRep,
	})

	testCases := map[string]string{
		"ﺁقای فرش": "آقای فرش",
		"ﺁﺑﺎﺩ": "آباد",
	}

	for input, expected := range testCases {
		actual := preprocessor.Process(input)
		assert.Equal(t, expected, actual, "they should be equal")
	}
}

func TestPersianDigitNormalization(t *testing.T) {

	preprocessor := NewPreprocessor([]PreprocessingStep{
		&persianDigit,
	})

	testCases := map[string]string{
		"1234567890": "۱۲۳۴۵۶۷۸۹۰" ,
		"0123" : "۰۱۲۳" ,
	}

	for input, expected := range testCases {
		actual := preprocessor.Process(input)
		assert.Equal(t, expected, actual, "they should be equal")
	}
}

func TestSpecialArabicPhraseNormalization(t *testing.T) {
	preprocessor := NewPreprocessor([]PreprocessingStep{
		&arabicPhrase,
	})

	testCases := map[string]string{
		"\ufdfb": "جل جلاله",
		"پیامبر \ufdfb": "پیامبر جل جلاله" ,
	}

	for input, expected := range testCases {
		actual := preprocessor.Process(input)
		assert.Equal(t, expected, actual, "they should be equal")
	}
}

func TestCombinedNormalization(t *testing.T) {

	preprocessor := NewPreprocessor([]PreprocessingStep{
		&unicodeRep,
		&persianDigit,
		&arabicPhrase,
	})

	testCases := map[string]string{
		"ﺁقای 123 ﷽": "آقای ۱۲۳ بسم الله الرحمن الرحیم",
		"0ﺁ شروع":    "۰آ شروع",
	}

	for input, expected := range testCases {
		actual := preprocessor.Process(input)
		assert.Equal(t, expected, actual, "they should be equal")
	}
}
