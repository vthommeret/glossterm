package ml

var DefaultLangs = []string{"en", "es", "fr", "la", "LL"}
var DefaultLangMap map[string]bool

func init() {
	DefaultLangMap = make(map[string]bool)
	for _, l := range DefaultLangs {
		DefaultLangMap[l] = true
	}
}

type Lang struct {
	Code      string
	Canonical string
	Other     []string
}

// From https://en.wiktionary.org/wiki/Module:languages/data2
// Only supporting two-letter language codes for now.
var Langs = map[string]Lang{
	"aa": {
		Code:      "aa",
		Canonical: "Afar",
		Other:     []string{"Qafar"},
	},
	"ab": {
		Code:      "ab",
		Canonical: "Abkhaz",
		Other:     []string{"Abkhazian", "Abxazo"},
	},
	"ae": {
		Code:      "ae",
		Canonical: "Avestan",
		Other:     []string{"Zend", "Old Bactrian"},
	},
	"af": {
		Code:      "af",
		Canonical: "Afrikaans",
	},
	"ak": {
		Code:      "ak",
		Canonical: "Akan",
		Other:     []string{"Twi-Fante", "Twi", "Fante", "Fanti", "Asante", "Akuapem"},
	},
	"am": {
		Code:      "am",
		Canonical: "Amharic",
	},
	"an": {
		Code:      "an",
		Canonical: "Aragonese",
	},
	"ar": {
		Code:      "ar",
		Canonical: "Arabic",
		Other:     []string{"Modern Standard Arabic", "Standard Arabic", "Literary Arabic", "Classical Arabic"},
	},
	"as": {
		Code:      "as",
		Canonical: "Assamese",
	},
	"av": {
		Code:      "av",
		Canonical: "Avar",
		Other:     []string{"Avaric"},
	},
	"ay": {
		Code:      "ay",
		Canonical: "Aymara",
		Other:     []string{"Southern Aymara", "Central Aymara"},
	},
	"az": {
		Code:      "az",
		Canonical: "Azeri",
		Other:     []string{"Azerbaijani", "Azari", "Azeri Turkic", "Azerbaijani Turkic", "North Azerbaijani", "South Azerbaijani", "Afshar", "Afshari", "Afshar Azerbaijani", "Afchar", "Qashqa'i", "Qashqai", "Kashkay", "Sonqor"},
	},
	"ba": {
		Code:      "ba",
		Canonical: "Bashkir",
	},
	"be": {
		Code:      "be",
		Canonical: "Belarusian",
		Other:     []string{"Belorussian", "Belarusan", "Bielorussian", "Byelorussian", "Belarussian", "White Russian"},
	},
	"bg": {
		Code:      "bg",
		Canonical: "Bulgarian",
	},
	"bh": {
		Code:      "bh",
		Canonical: "Bihari",
	},
	"bi": {
		Code:      "bi",
		Canonical: "Bislama",
	},
	"bm": {
		Code:      "bm",
		Canonical: "Bambara",
		Other:     []string{"Bamanankan"},
	},
	"bn": {
		Code:      "bn",
		Canonical: "Bengali",
		Other:     []string{"Bangla"},
	},
	"bo": {
		Code:      "bo",
		Canonical: "Tibetan",
		Other:     []string{"Ü", "Dbus", "Lhasa", "Lhasa Tibetan", "Amdo Tibetan", "Amdo", "Panang", "Khams", "Khams Tibetan", "Khamba", "Tseku", "Dolpo", "Humla", "Limi", "Lhomi", "Shing Saapa", "Mugom", "Mugu", "Nubri", "Walungge", "Gola", "Thudam", "Lowa", "Loke", "Mustang", "Tichurong"},
	},
	"br": {
		Code:      "br",
		Canonical: "Breton",
	},
	"ca": {
		Code:      "ca",
		Canonical: "Catalan",
		Other:     []string{"Valencian"},
	},
	"ce": {
		Code:      "ce",
		Canonical: "Chechen",
	},
	"ch": {
		Code:      "ch",
		Canonical: "Chamorro",
		Other:     []string{"Chamoru"},
	},
	"co": {
		Code:      "co",
		Canonical: "Corsican",
		Other:     []string{"Corsu"},
	},
	"cr": {
		Code:      "cr",
		Canonical: "Cree",
	},
	"cs": {
		Code:      "cs",
		Canonical: "Czech",
	},
	"cu": {
		Code:      "cu",
		Canonical: "Old Church Slavonic",
		Other:     []string{"Old Church Slavic"},
	},
	"cv": {
		Code:      "cv",
		Canonical: "Chuvash",
	},
	"cy": {
		Code:      "cy",
		Canonical: "Welsh",
	},
	"da": {
		Code:      "da",
		Canonical: "Danish",
	},
	"de": {
		Code:      "de",
		Canonical: "German",
		Other:     []string{"High German", "New High German", "Deutsch"},
	},
	"dv": {
		Code:      "dv",
		Canonical: "Dhivehi",
		Other:     []string{"Divehi", "Mahal", "Mahl", "Maldivian"},
	},
	"dz": {
		Code:      "dz",
		Canonical: "Dzongkha",
	},
	"ee": {
		Code:      "ee",
		Canonical: "Ewe",
	},
	"el": {
		Code:      "el",
		Canonical: "Greek",
		Other:     []string{"Modern Greek", "Neo-Hellenic"},
	},
	"en": {
		Code:      "en",
		Canonical: "English",
		Other:     []string{"Modern English", "New English", "Hawaiian Creole English", "Hawai'ian Creole English", "Hawaiian Creole", "Hawai'ian Creole", "Polari", "Yinglish"},
	},
	"eo": {
		Code:      "eo",
		Canonical: "Esperanto",
	},
	"es": {
		Code:      "es",
		Canonical: "Spanish",
		Other:     []string{"Castilian", "Amazonian Spanish", "Amazonic Spanish", "Loreto-Ucayali Spanish"},
	},
	"et": {
		Code:      "et",
		Canonical: "Estonian",
	},
	"eu": {
		Code:      "eu",
		Canonical: "Basque",
		Other:     []string{"Euskara"},
	},
	"fa": {
		Code:      "fa",
		Canonical: "Persian",
		Other:     []string{"Farsi", "New Persian", "Modern Persian", "Western Persian", "Iranian Persian", "Eastern Persian", "Dari", "Aimaq", "Aimak", "Aymaq", "Eimak"},
	},
	"ff": {
		Code:      "ff",
		Canonical: "Fula",
		Other:     []string{"Adamawa Fulfulde", "Bagirmi Fulfulde", "Borgu Fulfulde", "Central-Eastern Niger Fulfulde", "Fulani", "Fulfulde", "Maasina Fulfulde", "Nigerian Fulfulde", "Pular", "Pulaar", "Western Niger Fulfulde"},
	},
	"fi": {
		Code:      "fi",
		Canonical: "Finnish",
		Other:     []string{"Suomi", "Botnian"},
	},
	"fj": {
		Code:      "fj",
		Canonical: "Fijian",
	},
	"fo": {
		Code:      "fo",
		Canonical: "Faroese",
	},
	"fr": {
		Code:      "fr",
		Canonical: "French",
		Other:     []string{"Modern French"},
	},
	"fy": {
		Code:      "fy",
		Canonical: "West Frisian",
		Other:     []string{"Western Frisian", "Frisian"},
	},
	"ga": {
		Code:      "ga",
		Canonical: "Irish",
		Other:     []string{"Irish Gaelic"},
	},
	"gd": {
		Code:      "gd",
		Canonical: "Scottish Gaelic",
		Other:     []string{"Gàidhlig", "Highland Gaelic", "Scots Gaelic", "Scottish"},
	},
	"gl": {
		Code:      "gl",
		Canonical: "Galician",
	},
	"gn": {
		Code:      "gn",
		Canonical: "Guaraní",
	},
	"gu": {
		Code:      "gu",
		Canonical: "Gujarati",
	},
	"gv": {
		Code:      "gv",
		Canonical: "Manx",
		Other:     []string{"Manx Gaelic"},
	},
	"ha": {
		Code:      "ha",
		Canonical: "Hausa",
	},
	"he": {
		Code:      "he",
		Canonical: "Hebrew",
		Other:     []string{"Ivrit"},
	},
	"hi": {
		Code:      "hi",
		Canonical: "Hindi",
	},
	"ho": {
		Code:      "ho",
		Canonical: "Hiri Motu",
		Other:     []string{"Pidgin Motu", "Police Motu"},
	},
	"ht": {
		Code:      "ht",
		Canonical: "Haitian Creole",
		Other:     []string{"Creole", "Haitian", "Kreyòl"},
	},
	"hu": {
		Code:      "hu",
		Canonical: "Hungarian",
		Other:     []string{"Magyar"},
	},
	"hy": {
		Code:      "hy",
		Canonical: "Armenian",
		Other:     []string{"Modern Armenian", "Eastern Armenian", "Western Armenian"},
	},
	"hz": {
		Code:      "hz",
		Canonical: "Herero",
	},
	"ia": {
		Code:      "ia",
		Canonical: "Interlingua",
	},
	"id": {
		Code:      "id",
		Canonical: "Indonesian",
	},
	"ie": {
		Code:      "ie",
		Canonical: "Interlingue",
		Other:     []string{"Occidental"},
	},
	"ig": {
		Code:      "ig",
		Canonical: "Igbo",
	},
	"ii": {
		Code:      "ii",
		Canonical: "Sichuan Yi",
		Other:     []string{"Nuosu", "Nosu", "Northern Yi", "Liangshan Yi"},
	},
	"ik": {
		Code:      "ik",
		Canonical: "Inupiak",
		Other:     []string{"Inupiaq", "Iñupiaq", "Inupiatun"},
	},
	"io": {
		Code:      "io",
		Canonical: "Ido",
	},
	"is": {
		Code:      "is",
		Canonical: "Icelandic",
	},
	"it": {
		Code:      "it",
		Canonical: "Italian",
	},
	"iu": {
		Code:      "iu",
		Canonical: "Inuktitut",
		Other:     []string{"Eastern Canadian Inuktitut", "Eastern Canadian Inuit", "Western Canadian Inuktitut", "Western Canadian Inuit", "Western Canadian Inuktun", "Inuinnaq", "Inuinnaqtun", "Inuvialuk", "Inuvialuktun", "Nunavimmiutit", "Nunatsiavummiut", "Aivilimmiut", "Natsilingmiut", "Kivallirmiut", "Siglit", "Siglitun"},
	},
	"ja": {
		Code:      "ja",
		Canonical: "Japanese",
		Other:     []string{"Modern Japanese", "Nipponese", "Nihongo"},
	},
	"jv": {
		Code:      "jv",
		Canonical: "Javanese",
	},
	"ka": {
		Code:      "ka",
		Canonical: "Georgian",
		Other:     []string{"Kartvelian", "Judeo-Georgian", "Kivruli", "Gruzinic"},
	},
	"kg": {
		Code:      "kg",
		Canonical: "Kongo",
		Other:     []string{"Kikongo", "Koongo", "Laari", "San Salvador Kongo", "Yombe"},
	},
	"ki": {
		Code:      "ki",
		Canonical: "Kikuyu",
		Other:     []string{"Gikuyu", "Gĩkũyũ"},
	},
	"kj": {
		Code:      "kj",
		Canonical: "Kwanyama",
		Other:     []string{"Kuanyama", "Oshikwanyama"},
	},
	"kk": {
		Code:      "kk",
		Canonical: "Kazakh",
	},
	"kl": {
		Code:      "kl",
		Canonical: "Greenlandic",
		Other:     []string{"Kalaallisut"},
	},
	"km": {
		Code:      "km",
		Canonical: "Khmer",
		Other:     []string{"Cambodian"},
	},
	"kn": {
		Code:      "kn",
		Canonical: "Kannada",
	},
	"ko": {
		Code:      "ko",
		Canonical: "Korean",
		Other:     []string{"Modern Korean"},
	},
	"kr": {
		Code:      "kr",
		Canonical: "Kanuri",
		Other:     []string{"Kanembu", "Bilma Kanuri", "Central Kanuri", "Manga Kanuri", "Tumari Kanuri"},
	},
	"ks": {
		Code:      "ks",
		Canonical: "Kashmiri",
	},
	"ku": {
		Code:      "ku",
		Canonical: "Kurdish",
	},
	"kw": {
		Code:      "kw",
		Canonical: "Cornish",
	},
	"ky": {
		Code:      "ky",
		Canonical: "Kyrgyz",
		Other:     []string{"Kirghiz", "Kirgiz"},
	},
	"la": {
		Code:      "la",
		Canonical: "Latin",
	},
	"lb": {
		Code:      "lb",
		Canonical: "Luxembourgish",
	},
	"lg": {
		Code:      "lg",
		Canonical: "Luganda",
		Other:     []string{"Ganda"},
	},
	"li": {
		Code:      "li",
		Canonical: "Limburgish",
		Other:     []string{"Limburgan", "Limburgian", "Limburgic"},
	},
	"ln": {
		Code:      "ln",
		Canonical: "Lingala",
	},
	"lo": {
		Code:      "lo",
		Canonical: "Lao",
		Other:     []string{"Laotian"},
	},
	"lt": {
		Code:      "lt",
		Canonical: "Lithuanian",
	},
	"lu": {
		Code:      "lu",
		Canonical: "Luba-Katanga",
	},
	"lv": {
		Code:      "lv",
		Canonical: "Latvian",
		Other:     []string{"Lettish", "Lett"},
	},
	"mg": {
		Code:      "mg",
		Canonical: "Malagasy",
		Other:     []string{"Betsimisaraka Malagasy", "Betsimisaraka", "Northern Betsimisaraka Malagasy", "Northern Betsimisaraka", "Southern Betsimisaraka Malagasy", "Southern Betsimisaraka", "Bara Malagasy", "Bara", "Masikoro Malagasy", "Masikoro", "Antankarana", "Antankarana Malagasy", "Plateau Malagasy", "Sakalava", "Tandroy Malagasy", "Tandroy", "Tanosy", "Tanosy Malagasy", "Tesaka", "Tsimihety", "Tsimihety Malagasy", "Bushi", "Shibushi", "Kibushi", "Sakalava"},
	},
	"mh": {
		Code:      "mh",
		Canonical: "Marshallese",
	},
	"mi": {
		Code:      "mi",
		Canonical: "Maori",
		Other:     []string{"Māori"},
	},
	"mk": {
		Code:      "mk",
		Canonical: "Macedonian",
	},
	"ml": {
		Code:      "ml",
		Canonical: "Malayalam",
	},
	"mn": {
		Code:      "mn",
		Canonical: "Mongolian",
		Other:     []string{"Khalkha Mongolian"},
	},
	"mr": {
		Code:      "mr",
		Canonical: "Marathi",
	},
	"ms": {
		Code:      "ms",
		Canonical: "Malay",
		Other:     []string{"Malaysian", "Standard Malay", "Orang Seletar", "Orang Kanaq", "Jakun", "Temuan"},
	},
	"mt": {
		Code:      "mt",
		Canonical: "Maltese",
	},
	"my": {
		Code:      "my",
		Canonical: "Burmese",
		Other:     []string{"Myanmar"},
	},
	"na": {
		Code:      "na",
		Canonical: "Nauruan",
		Other:     []string{"Nauru"},
	},
	"nb": {
		Code:      "nb",
		Canonical: "Norwegian Bokmål",
		Other:     []string{"Bokmål"},
	},
	"nd": {
		Code:      "nd",
		Canonical: "Northern Ndebele",
		Other:     []string{"North Ndebele"},
	},
	"ne": {
		Code:      "ne",
		Canonical: "Nepali",
		Other:     []string{"Nepalese"},
	},
	"ng": {
		Code:      "ng",
		Canonical: "Ndonga",
	},
	"nl": {
		Code:      "nl",
		Canonical: "Dutch",
		Other:     []string{"Netherlandic", "Flemish"},
	},
	"nn": {
		Code:      "nn",
		Canonical: "Norwegian Nynorsk",
		Other:     []string{"New Norwegian", "Nynorsk"},
	},
	"no": {
		Code:      "no",
		Canonical: "Norwegian",
	},
	"nr": {
		Code:      "nr",
		Canonical: "Southern Ndebele",
		Other:     []string{"South Ndebele"},
	},
	"nv": {
		Code:      "nv",
		Canonical: "Navajo",
		Other:     []string{"Navaho", "Diné bizaad"},
	},
	"ny": {
		Code:      "ny",
		Canonical: "Chichewa",
		Other:     []string{"Chicheŵa", "Chinyanja", "Nyanja", "Chewa", "Cicewa", "Cewa", "Cinyanja"},
	},
	"oc": {
		Code:      "oc",
		Canonical: "Occitan",
		Other:     []string{"Provençal", "Auvergnat", "Auvernhat", "Gascon", "Languedocien", "Lengadocian", "Shuadit", "Chouhadite", "Chouhadit", "Chouadite", "Chouadit", "Shuhadit", "Judeo-Provençal", "Judeo-Provencal", "Judeo-Comtadin"},
	},
	"oj": {
		Code:      "oj",
		Canonical: "Ojibwe",
		Other:     []string{"Chippewa", "Ojibway", "Ojibwemowin", "Southwestern Ojibwa"},
	},
	"om": {
		Code:      "om",
		Canonical: "Oromo",
		Other:     []string{"Orma", "Borana-Arsi-Guji Oromo", "West Central Oromo"},
	},
	"or": {
		Code:      "or",
		Canonical: "Oriya",
		Other:     []string{"Odia", "Oorya"},
	},
	"os": {
		Code:      "os",
		Canonical: "Ossetian",
		Other:     []string{"Ossete", "Ossetic", "Digor", "Iron"},
	},
	"pa": {
		Code:      "pa",
		Canonical: "Punjabi",
		Other:     []string{"Panjabi"},
	},
	"pi": {
		Code:      "pi",
		Canonical: "Pali",
	},
	"pl": {
		Code:      "pl",
		Canonical: "Polish",
	},
	"ps": {
		Code:      "ps",
		Canonical: "Pashto",
		Other:     []string{"Pashtun", "Pushto", "Pashtu", "Central Pashto", "Northern Pashto", "Southern Pashto", "Pukhto", "Pakhto", "Pakkhto", "Afghani"},
	},
	"pt": {
		Code:      "pt",
		Canonical: "Portuguese",
		Other:     []string{"Modern Portuguese"},
	},
	"qu": {
		Code:      "qu",
		Canonical: "Quechua",
	},
	"rm": {
		Code:      "rm",
		Canonical: "Romansch",
		Other:     []string{"Romansh", "Rumantsch", "Romanche"},
	},
	"ro": {
		Code:      "ro",
		Canonical: "Romanian",
		Other:     []string{"Daco-Romanian", "Roumanian", "Rumanian"},
	},
	"ru": {
		Code:      "ru",
		Canonical: "Russian",
	},
	"rw": {
		Code:      "rw",
		Canonical: "Rwanda-Rundi",
		Other:     []string{"Rwanda", "Kinyarwanda", "Rundi", "Kirundi", "Ha", "Giha", "Hangaza", "Vinza", "Shubi", "Subi"},
	},
	"sa": {
		Code:      "sa",
		Canonical: "Sanskrit",
	},
	"sc": {
		Code:      "sc",
		Canonical: "Sardinian",
		Other:     []string{"Campidanese", "Campidanese Sardinian", "Logudorese", "Logudorese Sardinian", "Nuorese", "Nuorese Sardinian"},
	},
	"sd": {
		Code:      "sd",
		Canonical: "Sindhi",
	},
	"se": {
		Code:      "se",
		Canonical: "Northern Sami",
		Other:     []string{"North Sami", "Northern Saami", "North Saami"},
	},
	"sg": {
		Code:      "sg",
		Canonical: "Sango",
	},
	"sh": {
		Code:      "sh",
		Canonical: "Serbo-Croatian",
		Other:     []string{"BCS", "Croato-Serbian", "Serbocroatian", "Bosnian", "Croatian", "Montenegrin", "Serbian"},
	},
	"si": {
		Code:      "si",
		Canonical: "Sinhalese",
		Other:     []string{"Singhalese", "Sinhala"},
	},
	"sk": {
		Code:      "sk",
		Canonical: "Slovak",
	},
	"sl": {
		Code:      "sl",
		Canonical: "Slovene",
		Other:     []string{"Slovenian"},
	},
	"sm": {
		Code:      "sm",
		Canonical: "Samoan",
	},
	"sn": {
		Code:      "sn",
		Canonical: "Shona",
	},
	"so": {
		Code:      "so",
		Canonical: "Somali",
	},
	"sq": {
		Code:      "sq",
		Canonical: "Albanian",
	},
	"ss": {
		Code:      "ss",
		Canonical: "Swazi",
		Other:     []string{"Swati"},
	},
	"st": {
		Code:      "st",
		Canonical: "Sotho",
		Other:     []string{"Sesotho", "Southern Sesotho", "Southern Sotho"},
	},
	"su": {
		Code:      "su",
		Canonical: "Sundanese",
	},
	"sv": {
		Code:      "sv",
		Canonical: "Swedish",
	},
	"sw": {
		Code:      "sw",
		Canonical: "Swahili",
		Other:     []string{"Settler Swahili", "KiSetla", "KiSettla", "Setla", "Settla", "Kitchen Swahili", "Kihindi", "Indian Swahili", "KiShamba", "Kishamba", "Field Swahili", "Kibabu", "Asian Swahili", "Kimanga", "Arab Swahili", "Kitvita", "Army Swahili"},
	},
	"ta": {
		Code:      "ta",
		Canonical: "Tamil",
	},
	"te": {
		Code:      "te",
		Canonical: "Telugu",
	},
	"tg": {
		Code:      "tg",
		Canonical: "Tajik",
		Other:     []string{"Tadjik", "Tadzhik", "Tajiki", "Tajik Persian"},
	},
	"th": {
		Code:      "th",
		Canonical: "Thai",
	},
	"ti": {
		Code:      "ti",
		Canonical: "Tigrinya",
	},
	"tk": {
		Code:      "tk",
		Canonical: "Turkmen",
	},
	"tl": {
		Code:      "tl",
		Canonical: "Tagalog",
	},
	"tn": {
		Code:      "tn",
		Canonical: "Tswana",
		Other:     []string{"Setswana"},
	},
	"to": {
		Code:      "to",
		Canonical: "Tongan",
	},
	"tr": {
		Code:      "tr",
		Canonical: "Turkish",
	},
	"ts": {
		Code:      "ts",
		Canonical: "Tsonga",
	},
	"tt": {
		Code:      "tt",
		Canonical: "Tatar",
	},
	"ty": {
		Code:      "ty",
		Canonical: "Tahitian",
	},
	"ug": {
		Code:      "ug",
		Canonical: "Uyghur",
		Other:     []string{"Uigur", "Uighur", "Uygur"},
	},
	"uk": {
		Code:      "uk",
		Canonical: "Ukrainian",
	},
	"ur": {
		Code:      "ur",
		Canonical: "Urdu",
	},
	"uz": {
		Code:      "uz",
		Canonical: "Uzbek",
		Other:     []string{"Northern Uzbek", "Southern Uzbek"},
	},
	"ve": {
		Code:      "ve",
		Canonical: "Venda",
	},
	"vi": {
		Code:      "vi",
		Canonical: "Vietnamese",
		Other:     []string{"Annamese", "Annamite"},
	},
	"vo": {
		Code:      "vo",
		Canonical: "Volapük",
	},
	"wa": {
		Code:      "wa",
		Canonical: "Walloon",
	},
	"wo": {
		Code:      "wo",
		Canonical: "Wolof",
		Other:     []string{"Gambian Wolof"},
	},
	"xh": {
		Code:      "xh",
		Canonical: "Xhosa",
	},
	"yi": {
		Code:      "yi",
		Canonical: "Yiddish",
	},
	"yo": {
		Code:      "yo",
		Canonical: "Yoruba",
	},
	"za": {
		Code:      "za",
		Canonical: "Zhuang",
	},
	"zh": {
		Code:      "zh",
		Canonical: "Chinese",
	},
	"zu": {
		Code:      "zu",
		Canonical: "Zulu",
		Other:     []string{"isiZulu"},
	},
}

var CanonicalLangs map[string]Lang

func init() {
	CanonicalLangs = make(map[string]Lang)
	for _, l := range Langs {
		CanonicalLangs[l.Canonical] = l
	}
}
