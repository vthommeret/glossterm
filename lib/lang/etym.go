package lang

var etymMap map[string]string

func init() {
	etymMap = make(map[string]string)
	for _, e := range etyms {
		etymMap[e.Parent] = e.Parent
		for _, c := range e.Codes {
			etymMap[c] = e.Parent
		}
	}
}

func ToParent(l string) string {
	if p, ok := etymMap[l]; ok {
		return p
	}
	return l
}

type Etym struct {
	Canonical string
	Parent    string
	Codes     []string
}

// From https://en.wiktionary.org/wiki/Module:etymology_languages/data
var etyms = []Etym{
	// German varieties
	{
		Canonical: "Austrian German",
		Parent:    "de",
		Codes:     []string{"de-AT", "Austrian German", "de-AT", "AG.", "de-AT"},
	},
	{
		Canonical: "Viennese German",
		Parent:    "de",
		Codes:     []string{"de-AT-vie", "Viennese German", "de-AT-vie", "VG.", "de-AT-vie"},
	},

	// English varieties
	{
		Canonical: "British English",
		Parent:    "en",
		Codes:     []string{"en-GB", "British English", "en-GB", "BE.", "en-GB"},
	},
	{
		Canonical: "American English",
		Parent:    "en",
		Codes:     []string{"en-US", "American English", "en-US", "AE.", "en-US"},
	},

	// French varieties
	{
		Canonical: "Canadian French",
		Parent:    "fr",
		Codes:     []string{"fr-CA", "Canadian French", "fr-CA", "CF.", "fr-CA"},
	},
	{
		Canonical: "Acadian French",
		Parent:    "fr",
		Codes:     []string{"fr-aca", "Acadian French", "fr-aca", "fra-aca", "fr-aca"},
	},
	{
		Canonical: "Cajun French",
		Parent:    "fr",
		Codes:     []string{"frc"},
	},

	// Italian varieties
	{
		Canonical: "Old Italian",
		Parent:    "it",
		Codes:     []string{"roa-oit"},
	},

	// Latin varieties by period
	{
		Canonical: "Late Latin",
		Parent:    "la",
		Codes:     []string{"la-lat", "Late Latin", "la-lat", "LL.", "la-lat", "LL", "la-lat"},
	},
	{
		Canonical: "Vulgar Latin",
		Parent:    "la",
		Codes:     []string{"la-vul", "Vulgar Latin", "la-vul", "VL.", "la-vul"},
	},
	{
		Canonical: "Medieval Latin",
		Parent:    "la",
		Codes:     []string{"la-med", "Medieval Latin", "la-med", "ML.", "la-med", "ML", "la-med"},
	},
	{
		Canonical: "Ecclesiastical Latin",
		Parent:    "la",
		Codes:     []string{"la-ecc", "Ecclesiastical Latin", "la-ecc", "EL.", "la-ecc"},
	},
	{
		Canonical: "Renaissance Latin",
		Parent:    "la",
		Codes:     []string{"la-ren", "Renaissance Latin", "la-ren", "RL.", "la-ren"},
	},
	{
		Canonical: "New Latin",
		Parent:    "la",
		Codes:     []string{"la-new", "New Latin", "la-new", "NL.", "la-new"},
	},

	// Greek varieties
	{
		Canonical: "Koine Greek",
		Parent:    "grc",
		Codes:     []string{"grc-koi", "Koine", "grc-koi"},
	},
	{
		Canonical: "Byzantine Greek",
		Parent:    "grc",
		Codes:     []string{"gkm", "Medieval Greek", "gkm"},
	},
	{
		Canonical: "Doric Greek",
		Parent:    "grc",
		Codes:     []string{"grc-dor"},
	},
	{
		Canonical: "Attic Greek",
		Parent:    "grc",
		Codes:     []string{"grc-att"},
	},
	{
		Canonical: "Ionic Greek",
		Parent:    "grc",
		Codes:     []string{"grc-ion"},
	},
	{
		Canonical: "Pamphylian Greek",
		Parent:    "grc",
		Codes:     []string{"grc-pam"},
	},
	{
		Canonical: "Arcadocypriot Greek",
		Parent:    "grc",
		Codes:     []string{"grc-arp"},
	},
	{
		Canonical: "Aeolic Greek",
		Parent:    "grc",
		Codes:     []string{"grc-aeo"},
	},
	{
		Canonical: "Elean Greek",
		Parent:    "grc",
		Codes:     []string{"grc-ela"},
	},
	{
		Canonical: "Epic Greek",
		Parent:    "grc",
		Codes:     []string{"grc-epc"},
	},
	{
		Canonical: "Homeric Greek",
		Parent:    "grc",
		Codes:     []string{"grc-hmr"},
	},
}
