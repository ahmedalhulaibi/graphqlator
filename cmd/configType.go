package cmd

/*gqlpackage struct to read in graphqlator-pkg.json*/
type gqlpackage struct {
	ProjectName      string   `json:"project_name"`
	DatabaseType     string   `json:"database_type"`
	ConnectionString string   `json:"connection_string"`
	TableNames       []string `json:"table_names"`
	GitRepo          string   `json:"git_repo"`
	GenMode          string   `json:"gen_mode"`
}
