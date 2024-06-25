package preprocessing

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

const (
	WORDS_PATH = "./pkg/engine/preprocessing/tokenizer/words.dat"
	VERBS_PATH = "./pkg/engine/preprocessing/tokenizer/verbs.dat"
)

// WordTokenizer defines the tokenizer struct.
type WordTokenizer struct {
	joinVerbParts     bool
	joinAbbreviations bool
	separateEmoji     bool
	replaceLinks      bool
	replaceIDs        bool
	replaceEmails     bool
	replaceNumbers    bool
	replaceHashtags   bool

	emojiPattern       *regexp.Regexp
	linkPattern        *regexp.Regexp
	emailPattern       *regexp.Regexp
	numberIntPattern   *regexp.Regexp
	numberFloatPattern *regexp.Regexp
	hashtagPattern     *regexp.Regexp
	abbreviations      map[string]string

	words       map[string][2]string
	afterVerbs  map[string]bool
	beforeVerbs map[string]bool
	verbs       []string
	bons        map[string]bool
	verbe       map[string]bool
}

// NewWordTokenizer initializes a new WordTokenizer.
func NewWordTokenizer(
	wordsFile, verbsFile string,
	joinVerbParts, joinAbbreviations, separateEmoji, replaceLinks,
	replaceIDs, replaceEmails, replaceNumbers, replaceHashtags bool,
) (*WordTokenizer, error) {

	

	tokenizer := &WordTokenizer{
		joinVerbParts:      joinVerbParts,
		joinAbbreviations:  joinAbbreviations,
		separateEmoji:      separateEmoji,
		replaceLinks:       replaceLinks,
		replaceIDs:         replaceIDs,
		replaceEmails:      replaceEmails,
		replaceNumbers:     replaceNumbers,
		replaceHashtags:    replaceHashtags,
		emojiPattern:       regexp.MustCompile("[\U0001f600-\U0001f64f\U0001f300-\U0001f5ff\U0001f4cc\U0001f4cd]"),
		linkPattern:        regexp.MustCompile(`((https?|ftp)://)?(([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,})[-\w@:%_.+/~#?=&]*`),
		emailPattern:       regexp.MustCompile(`[a-zA-Z0-9._+-]+@([a-zA-Z0-9-]+\.)+[A-Za-z]{2,}`),
		numberIntPattern:   regexp.MustCompile(`[\d۰-۹]+`),
		numberFloatPattern: regexp.MustCompile(`[\d۰-۹]+\\[\d۰-۹]+`),
		hashtagPattern:     regexp.MustCompile(`#(\S+)`),
		words:              make(map[string][2]string),
		afterVerbs:         make(map[string]bool),
		beforeVerbs:        make(map[string]bool),
		bons:               make(map[string]bool),
		verbe:              make(map[string]bool),
		abbreviations:      make(map[string]string),
	}



	tokenizer.beforeVerbs = map[string]bool{
		"خواهم":   true,
		"خواهی":   true,
		"خواهد":   true,
		"خواهیم":  true,
		"خواهید":  true,
		"خواهند":  true,
		"نخواهم":  true,
		"نخواهی":  true,
		"نخواهد":  true,
		"نخواهیم": true,
		"نخواهید": true,
		"نخواهند": true,
	}

	tokenizer.afterVerbs = map[string]bool{
		"ام":         true,
		"ای":         true,
		"است":        true,
		"ایم":        true,
		"اید":        true,
		"اند":        true,
		"بودم":       true,
		"بودی":       true,
		"بود":        true,
		"بودیم":      true,
		"بودید":      true,
		"بودند":      true,
		"باشم":       true,
		"باشی":       true,
		"باشد":       true,
		"باشیم":      true,
		"باشید":      true,
		"باشند":      true,
		"شده_ام":     true,
		"شده_ای":     true,
		"شده_است":    true,
		"شده_ایم":    true,
		"شده_اید":    true,
		"شده_اند":    true,
		"شده_بودم":   true,
		"شده_بودی":   true,
		"شده_بود":    true,
		"شده_بودیم":  true,
		"شده_بودید":  true,
		"شده_بودند":  true,
		"شده_باشم":   true,
		"شده_باشی":   true,
		"شده_باشد":   true,
		"شده_باشیم":  true,
		"شده_باشید":  true,
		"شده_باشند":  true,
		"نشده_ام":    true,
		"نشده_ای":    true,
		"نشده_است":   true,
		"نشده_ایم":   true,
		"نشده_اید":   true,
		"نشده_اند":   true,
		"نشده_بودم":  true,
		"نشده_بودی":  true,
		"نشده_بود":   true,
		"نشده_بودیم": true,
		"نشده_بودید": true,
		"نشده_بودند": true,
		"نشده_باشم":  true,
		"نشده_باشی":  true,
		"نشده_باشد":  true,
		"نشده_باشیم": true,
		"نشده_باشید": true,
		"نشده_باشند": true,
		"شوم":        true,
		"شوی":        true,
		"شود":        true,
		"شویم":       true,
		"شوید":       true,
		"شوند":       true,
		"شدم":        true,
		"شدی":        true,
		"شد":         true,
		"شدیم":       true,
		"شدید":       true,
		"شدند":       true,
		"نشوم":       true,
		"نشوی":       true,
		"نشود":       true,
		"نشویم":      true,
		"نشوید":      true,
		"نشوند":      true,
		"نشدم":       true,
		"نشدی":       true,
		"نشد":        true,
		"نشدیم":      true,
		"نشدید":      true,
		"نشدند":      true,
		"می‌شوم":     true,
		"می‌شوی":     true,
		"می‌شود":     true,
		"می‌شویم":    true,
		"می‌شوید":    true,
		"می‌شوند":    true,
		"می‌شدم":     true,
		"می‌شدی":     true,
		"می‌شد":      true,
		"می‌شدیم":    true,
		"می‌شدید":    true,
		"می‌شدند":    true,
		"نمی‌شوم":    true,
		"نمی‌شوی":    true,
		"نمی‌شود":    true,
		"نمی‌شویم":   true,
		"نمی‌شوید":   true,
		"نمی‌شوند":   true,
		"نمی‌شدم":    true,
		"نمی‌شدی":    true,
		"نمی‌شد":     true,
		"نمی‌شدیم":   true,
		"نمی‌شدید":   true,
		"نمی‌شدند":   true,
		"خواهم_شد":   true,
		"خواهی_شد":   true,
		"خواهد_شد":   true,
		"خواهیم_شد":  true,
		"خواهید_شد":  true,
		"خواهند_شد":  true,
		"نخواهم_شد":  true,
		"نخواهی_شد":  true,
		"نخواهد_شد":  true,
		"نخواهیم_شد": true,
		"نخواهید_شد": true,
		"نخواهند_شد": true,
	}

	err := tokenizer.loadWords(wordsFile)
	if err != nil {
		return nil, err
	}

	err = tokenizer.loadVerbs(verbsFile)
	if err != nil {
		return nil, err
	}

	if joinAbbreviations {
		err = tokenizer.loadAbbreviations()
		if err != nil {
			return nil, err
		}
	}

	return tokenizer, nil
}

// loadWords loads the words from the given file.
func (wt *WordTokenizer) loadWords(wordsFile string) error {
	file, err := os.Open(wordsFile)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ",")
		if len(parts) == 3 {
			wt.words[parts[0]] = [2]string{parts[1], parts[2]}
		}
	}

	return scanner.Err()
}

// loadVerbs loads the verbs from the given file.
func (wt *WordTokenizer) loadVerbs(verbsFile string) error {
	file, err := os.Open(verbsFile)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		verb := scanner.Text()
		wt.verbs = append(wt.verbs, verb)
		parts := strings.Split(verb, "#")
		if len(parts) > 0 {
			wt.bons[parts[0]] = true
		}
	}

	for bon := range wt.bons {
		wt.verbe[bon+"ه"] = true
		wt.verbe["ن"+bon+"ه"] = true
	}

	return scanner.Err()
}

// loadAbbreviations loads the abbreviations from the default file.
func (wt *WordTokenizer) loadAbbreviations() error {
	abbreviationsFile := filepath.Join("path/to/default", "abbreviations.txt")
	file, err := os.Open(abbreviationsFile)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		abbr := scanner.Text()
		wt.abbreviations[abbr] = ""
	}

	return scanner.Err()
}

// Tokenize tokenizes the given text.
func (wt *WordTokenizer) Tokenize(text string) []string {
	// TODO: This tokenizer is replacing info such as email or link with a placeholder. This is not the intended behaviour and should be fixed
	if wt.joinAbbreviations {
		text = wt.replaceAbbreviations(text)
	}
	if wt.separateEmoji {
		text = wt.emojiPattern.ReplaceAllString(text, " $0 ")
	}
	if wt.replaceEmails {
		text = wt.emailPattern.ReplaceAllString(text, " EMAIL ")
	}
	if wt.replaceLinks {
		text = wt.linkPattern.ReplaceAllString(text, " LINK ")
	}
	if wt.replaceHashtags {
		text = wt.hashtagPattern.ReplaceAllString(text, "TAG ")
	}
	if wt.replaceNumbers {
		text = wt.numberIntPattern.ReplaceAllString(text, " NUM ")
	}

	text = regexp.MustCompile(`[؟!?]+|[\d۰-۹.:]+|[:.،؛»\])}"«\[({/\\]`).ReplaceAllString(text, " $0 ")
	tokens := strings.Fields(text)

	if wt.joinVerbParts {
		tokens = wt.joinVerbPartsFunc(tokens)
	}

	if wt.joinAbbreviations {
		tokens = wt.revertAbbreviations(tokens)
	}

	return tokens
}

// joinVerbPartsFunc joins the verb parts in the tokens.
func (wt *WordTokenizer) joinVerbPartsFunc(tokens []string) []string {
	if len(tokens) == 1 {
		return tokens
	}
	result := []string{""}

	for i := len(tokens) - 1; i >= 0; i-- {
		token := tokens[i]

		isInBeforeVerbs, _ := wt.beforeVerbs[token]
		isLastTokenInAfterVerbs, _ := wt.afterVerbs[result[len(result)-1]]
		isInVerbe, _ := wt.verbe[token]

		if isInBeforeVerbs || (isLastTokenInAfterVerbs && isInVerbe) {
			result[len(result)-1] = token + " " + result[len(result)-1]
		} else {
			result = append(result, token)
		}
	}
	result = result[1:]
	slices.Reverse(result)
	return result
}

// replaceAbbreviations replaces abbreviations in the text with placeholders.
func (wt *WordTokenizer) replaceAbbreviations(text string) string {
	// Generate a unique placeholder
	placeholder := "313"
	for strings.Contains(text, placeholder) {
		conv, _ := strconv.Atoi(placeholder)
		placeholder = fmt.Sprint(conv + 1)
	}

	for abbr := range wt.abbreviations {
		text = strings.Replace(text, " "+abbr+" ", " "+placeholder+" ", -1)
		wt.abbreviations[abbr] = placeholder
	}

	return text
}

// revertAbbreviations reverts the placeholders back to abbreviations.
func (wt *WordTokenizer) revertAbbreviations(tokens []string) []string {
	for i, token := range tokens {
		for abbr, placeholder := range wt.abbreviations {
			if token == placeholder {
				tokens[i] = abbr
				break
			}
		}
	}
	return tokens
}
