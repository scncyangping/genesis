package dandelion

type MatchOptions func(string) bool

type PaddingOptions func() map[string]any

var SuffixReplaceMap = map[string]string{
	"tmpl": "go",
}
