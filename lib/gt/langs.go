package gt

var SourceLangs = map[string]bool{}
var sourceLangs = []string{"en", "fr", "es", "pt"}

func init() {
	for _, l := range sourceLangs {
		SourceLangs[l] = true
	}
}
