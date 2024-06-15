package preprocessing

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var tokenizer, _ = NewWordTokenizer(WORDS_PATH, VERBS_PATH, true, false, false, false, false, true, false, false)
var TEST_DEBUG_MODE = true

func TestSimpleSentence(t *testing.T) {
	// Tokenize a sample text
	text := "یک جمله ساده برای تست کردن"
	actualTokens := tokenizer.Tokenize(text)
	expectedTokens := []string{"یک", "جمله", "ساده", "برای", "تست", "کردن"}
	assert.Equal(t, expectedTokens, actualTokens)

	if TEST_DEBUG_MODE {
		t.Log(actualTokens)
	}
}

func TestSentenceWithIntegerNumbers(t *testing.T) {
	// Tokenize a sample text
	text := "یک جمله با اعداد 1234"
	actualTokens := tokenizer.Tokenize(text)
	expectedTokens := []string{"یک", "جمله", "با", "اعداد", "1234"}
	assert.Equal(t, expectedTokens, actualTokens)

	if TEST_DEBUG_MODE {
		t.Log(actualTokens)
	}
}

func TestSentenceWithEnglishFloatNumbers(t *testing.T) {
	// Tokenize a sample text
	text := "یک جمله با اعداد 12.34"
	actualTokens := tokenizer.Tokenize(text)
	expectedTokens := []string{"یک", "جمله", "با", "اعداد", "12.34"}
	assert.Equal(t, expectedTokens, actualTokens)

	if TEST_DEBUG_MODE {
		t.Log(actualTokens)
	}
}

func TestSentenceWithPersianFloatNumbers(t *testing.T) {
	// Tokenize a sample text
	text := "یک جمله با اعداد ۱۲.۳۴"
	actualTokens := tokenizer.Tokenize(text)
	expectedTokens := []string{"یک", "جمله", "با", "اعداد", "۱۲.۳۴"}
	assert.Equal(t, expectedTokens, actualTokens)

	if TEST_DEBUG_MODE {
		t.Log(actualTokens)
	}
}

func TestSentenceWithEmail(t *testing.T) {
	// A persian sentence containing an email address
	text := "یک جمله حاوی example.somedomain_1402@example.com یک آدرس ایمیل"
	actualTokens := tokenizer.Tokenize(text)
	expectedTokens := []string{"یک", "جمله", "حاوی", "example.somedomain_1402@example.com", "یک", "آدرس", "ایمیل"}
	assert.Equal(t, expectedTokens, actualTokens)

	if TEST_DEBUG_MODE {
		t.Log(actualTokens)
	}
}

func TestSentenceWithId(t *testing.T) {
	// A persian sentence containing an ID
	text := "یک جمله حاوی 1234567890 یک شناسه"
	actualTokens := tokenizer.Tokenize(text)
	expectedTokens := []string{"یک", "جمله", "حاوی", "1234567890", "یک", "شناسه"}
	assert.Equal(t, expectedTokens, actualTokens)

	if TEST_DEBUG_MODE {
		t.Log(actualTokens)
	}
}

func TestSentenceWithSpecialVerbs(t *testing.T) {
	text := "من می‌توانم به مدرس بروم"
	actualTokens := tokenizer.Tokenize(text)
	expectedTokens := []string{"من", "می‌توانم", "به", "مدرس", "بروم"}

	assert.Equal(t, expectedTokens, actualTokens)
	if TEST_DEBUG_MODE {
		t.Log(actualTokens)
	}
}

func TestVerbAbbreviations(t *testing.T) {
	testMap := map[string][]string{
		"او از کشور خواهد رفت":   {"او", "از", "کشور", "خواهد رفت"},
		"این حرف‌ گفته خواهد شد": {"این", "حرف‌", "گفته", "خواهد شد"},
		"او از زندگی خسته شد":    {"او", "از", "زندگی", "خسته شد"},
	}

	for testCase, expected := range testMap {
		actual := tokenizer.Tokenize(testCase)
		assert.Equal(t, expected, actual)
		if TEST_DEBUG_MODE {
			t.Log(actual)
		}
	}
}
