package ml

func ToLangMap(ls []string) map[string]bool {
	m := make(map[string]bool)
	for _, l := range ls {
		m[l] = true
	}
	return m
}
