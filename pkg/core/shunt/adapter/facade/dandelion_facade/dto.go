package dandelion_facade

type TemplateGenerateDTO struct {
	DBName       string   `json:"dbName" structs:"dbName"`
	ModuleName   string   `json:"moduleName" structs:"moduleName"`
	TemplatePath string   `json:"templatePath" structs:"templatePath"`
	SavePath     string   `json:"savePath" structs:"savePath"`
	MatchSuffix  []string `json:"matchSuffix" structs:"matchSuffix"`
	Tables       []string `json:"tables" structs:"tables"`
}
