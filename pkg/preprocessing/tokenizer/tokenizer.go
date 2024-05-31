package preprocessing

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	WORDS_PATH = "words.dat"
	VERBS = "verbs.dat"
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
	afterVerbs  map[string]struct{}
	beforeVerbs map[string]struct{}
	verbs       []string
	bons        map[string]struct{}
	verbe       map[string]struct{}
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
		afterVerbs:         make(map[string]struct{}),
		beforeVerbs:        make(map[string]struct{}),
		bons:               make(map[string]struct{}),
		verbe:              make(map[string]struct{}),
		abbreviations:      make(map[string]string),
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
			wt.bons[parts[0]] = struct{}{}
		}
	}

	for bon := range wt.bons {
		wt.verbe[bon+"ه"] = struct{}{}
		wt.verbe["ن"+bon+"ه"] = struct{}{}
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

	text = regexp.MustCompile(`[؟!?]+|[\d.:]+|[:.،؛»\])}"«\[({/\\]`).ReplaceAllString(text, " $0 ")
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

	var result []string
	for i := len(tokens) - 1; i >= 0; i-- {
		token := tokens[i]
		if _, beforeExists := wt.beforeVerbs[token]; beforeExists {
			// If token exists in beforeVerbs
			if len(result) > 0 {
				result[len(result)-1] = token + "_" + result[len(result)-1]
			} else {
				result = append(result, token)
			}
		} else if len(result) > 0 {
			// Check if the last token in result exists in afterVerbs and current token exists in verbe
			lastToken := result[len(result)-1]
			if _, afterExists := wt.afterVerbs[lastToken]; afterExists {
				if _, verbeExists := wt.verbe[token]; verbeExists {
					result[len(result)-1] = token + "_" + result[len(result)-1]
				} else {
					result = append(result, token)
				}
			} else {
				result = append(result, token)
			}
		} else {
			result = append(result, token)
		}
	}

	// Reverse the result slice
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

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

/*
Example usage:

	tokenizer, err := NewWordTokenizer("words.dat", "verbs.dat", true, false, false, false, false, false, false, false)
		if err != nil {
			fmt.Println("Error initializing tokenizer:", err)
			return
		}

		// Tokenize a sample text
		text := "این جمله (خیلی) پیچیده نیست!!!"
		tokens := tokenizer.Tokenize(text)
		fmt.Println(tokens)

*/
