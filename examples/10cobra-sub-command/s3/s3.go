package s3

type LsOptions struct {
	Recursive     bool   `flag:"recursive"`
	PageSize      int    `flag:"page-size"`
	HumanReadable bool   `flag:"human-redable"`
	Summarize     bool   `flag:"summarize"`
	RequestPlayer string `flag:"request-player"`
}

type CpOptions struct {
	Dryrun             bool     `flag:"dryrun"`
	Quiet              bool     `flag:"quiet"`
	Include            string   `flag:"include"`
	Exclude            string   `flag:"exclude"`
	Acl                string   `flag:"acl"`
	MetadataDirective  string   `flag:"metadata-directive"`
	ExpectedSize       int      `flag:"expected-size"`
	Recursive          bool     `flag:"recursive"`
}

type MvOptions struct {
	Dryrun             bool     `flag:"dryrun"`
	Quiet              bool     `flag:"quiet"`
	Include            string   `flag:"include"`
	Exclude            string   `flag:"exclude"`
	Acl                string   `flag:"acl"`
	MetadataDirective  string   `flag:"metadata-directive"`
	ExpectedSize       int      `flag:"expected-size"`
	Recursive          bool     `flag:"recursive"`
}
