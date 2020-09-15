package gt

var SourceLangs = map[string]bool{}
var sourceLangs = []string{"en", "fr", "es"}

func init() {
	for _, l := range sourceLangs {
		SourceLangs[l] = true
	}
}
