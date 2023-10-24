package mixture

func main() {
	op := ManagerStartOp{
		ModuleName:   "hub_op",
		TemplatePath: "/Users/yapi/go/src/github.com/dandelion/test/goTemplate.zip",
		MatchSuffix:  []string{"go", "mod", "tmpl", "xml"},
		SavePath:     "/Users/yapi/WorkSpace/VscodeWorkSpace/templateTest/test.zip",
		Tables:       []string{"article", "project", "repo", "git_account"},
	}
	manager := NewManager(&op)

	manager.Start()
}
