package preprocessing

import (
	"strings"
	"regexp"
)

var persianNormalizationMap = map[string]string{
	"آ": "ﺁ ﺂ",
	"ا": "ﺍ ﺎ ٲ ٵ ﭐ ﭑ ﺃ ﺄ ٳ ﺇ ﺈ إ أ ꙇ",
	"ب": "ٮ ٻ ڀ ݐ ݒ ݔ ݕ ݖ ﭒ ﭓ ﭔ ﭕ ﺏ ﺐ ﺑ ﺒ",
	"پ": "ﭖ ﭗ ﭘ ﭙ ﭚ ﭛ ﭜ ﭝ",
	"ت": "ٹ ٺ ټ ٿ ݓ ﭞ ﭟ ﭠ ﭡ ﭦ ﭧ ﭨ ﭩ ﺕ ﺖ ﺗ ﺘ",
	"ث": "ٽ ݑ ﺙ ﺚ ﺛ ﺜ ﭢ ﭣ ﭤ ﭥ",
	"ج": "ڃ ڄ ﭲ ﭳ ﭴ ﭵ ﭶ ﭷ ﭸ ﭹ ﺝ ﺞ ﺟ ﺠ",
	"چ": "ڇ ڿ ﭺ ﭻ ݘ ﭼ ﭽ ﭾ ﭿ ﮀ ﮁ ݯ",
	"ح": "ځ ڂ څ ݗ ݮ ﺡ ﺢ ﺣ ﺤ",
	"خ": "ﺥ ﺦ ﺧ ﺨ",
	"د": "ڈ ډ ڊ ڋ ڍ ۮ ݙ ݚ ﮂ ﮃ ﮄ ﮈ ﮉ ﺩ ﺪ",
	"ذ": "ڌ ﱛ ﺫ ﺬ ڎ ڏ ڐ ﮅ ﮆ ﮇ",
	"ر": "ڑ ڒ ړ ڔ ڕ ږ ۯ ݛ ﮌ ﮍ ﱜ ﺭ ﺮ",
	"ز": "ڗ ݫ ݬ ﺯ ﺰ",
	"ژ": "ڙ ﮊ ﮋ",
	"س": "ښ ڛ ﺱ ﺲ ﺳ ﺴ",
	"ش": "ڜ ۺ ﺵ ﺶ ﺷ ﺸ ݜ ݭ",
	"ص": "ڝ ڞ ﺹ ﺺ ﺻ ﺼ",
	"ض": "ۻ ﺽ ﺾ ﺿ ﻀ",
	"ط": "ﻁ ﻂ ﻃ ﻄ",
	"ظ": "ﻅ ﻆ ﻇ ﻈ ڟ",
	"ع": "ڠ ݝ ݞ ݟ ﻉ ﻊ ﻋ ﻌ",
	"غ": "ۼ ﻍ ﻎ ﻏ ﻐ",
	"ف": "ڡ ڢ ڣ ڤ ڥ ڦ ݠ ݡ ﭪ ﭫ ﭬ ﭭ ﭮ ﭯ ﭰ ﭱ ﻑ ﻒ ﻓ ﻔ ᓅ",
	"ق": "ٯ ڧ ڨ ﻕ ﻖ ﻗ ﻘ",
	"ک": "ك ػ ؼ ڪ ګ ڬ ڭ ڮ ݢ ݣ ݤ ﮎ ﮏ ﮐ ﮑ ﯓ ﯔ ﯕ ﯖ ﻙ ﻚ ﻛ ﻜ",
	"گ": "ڰ ڱ ڲ ڳ ڴ ﮒ ﮓ ﮔ ﮕ ﮖ ﮗ ﮘ ﮙ ﮚ ﮛ ﮜ ﮝ",
	"ل": "ڵ ڶ ڷ ڸ ݪ ﻝ ﻞ ﻟ ﻠ",
	"م": "۾ ݥ ݦ ﻡ ﻢ ﻣ ﻤ",
	"ن": "ڹ ں ڻ ڼ ڽ ݧ ݨ ݩ ﮞ ﮟ ﮠ ﮡ ﮢ ﮣ ﻥ ﻦ ﻧ ﻨ",
	"و": "ٶ ٷ ﯗ ﯘ ﯙ ﯚ ﯛ ﯜ ﯝ ﯞ ﯟ ﺅ ﺆ ۄ ۅ ۆ ۇ ۈ ۉ ۊ ۋ ۏ ﯠ ﯡ ﯢ ﯣ ﻭ ﻮ ؤ פ",
	"ه": "ھ ۿ ۀ ہ ۂ ۃ ە ﮤ ﮥ ﮦ ﮧ ﮨ ﮩ ﮪ ﮫ ﮬ ﮭ ﺓ ﺔ ﻩ ﻪ ﻫ ﻬ ة",
	"ی": "ؠ ؽ ؾ ؿ ى ي ٸ ۍ ێ ې ۑ ے ۓ ﮮ ﮯ ﮰ ﮱ ﯤ ﯥ ﯦ ﯧ ﯼ ﯽ ﯾ ﯿ ﻯ ﻰ ﻱ ﻲ ﻳ ﻴ ﯨ ﯩ ۦ ﯪ ﯫ ﯬ ﯭ ﯮ ﯯ ﯰ ﯱ ﯲ ﯳ ﯴ ﯵ ﯶ ﯷ ﯸ ﯹ ﯺ ﯻ ﱝ ﺉ ﺊ ﺋ ﺌ ئ",
}


var digitNormalizationMap = map[string] string {
	"۰": "0 ٠ 𝟢 𝟬",
	"۱": "1 ١ 𝟣 𝟭 ⑴ ⒈ ⓵ ① ❶ 𝟙 𝟷 ı",
	"۲": "2 ٢ 𝟤 𝟮 ⑵ ⒉ ⓶ ② ❷ ² 𝟐 𝟸 𝟚 ᒿ շ",
	"۳": "3 ٣ 𝟥 𝟯 ⑶ ⒊ ⓷ ③ ❸ ³ ვ",
	"۴": "4 ٤ 𝟦 𝟰 ⑷ ⒋ ⓸ ④ ❹ ⁴",
	"۵": "5 ٥ 𝟧 𝟱 ⑸ ⒌ ⓹ ⑤ ❺ ⁵",
	"۶": "6 ٦ 𝟨 𝟲 ⑹ ⒍ ⓺ ⑥ ❻ ⁶",
	"۷": "7 ٧ 𝟩 𝟳 ⑺ ⒎ ⓻ ⑦ ❼ ⁷",
	"۸": "8 ٨ 𝟪 𝟴 ⑻ ⒏ ⓼ ⑧ ❽ ⁸",
	"۹": "9 ٩ 𝟫 𝟵 ⑼ ⒐ ⓽ ⑨ ❾ ⁹",
}


var knownPunctuatuions = ". | - ? ! ? ! % / * : ; , > < « » … ▕ ❘ ❙ ❚ ▏│ ㅡ 一 — – ー ̶ ـ \\u00ad ❔ ؟ � ？ ʕ ʔ ❕ ！ ⁉ ‼ ℅ ٪ ÷ × ： ؛ ； ، › ‹ ＜ 《 》•"

 
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