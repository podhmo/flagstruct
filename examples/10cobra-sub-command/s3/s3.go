package s3

type LsOptions struct {
	Recursive     bool   `flag:"recursive"`
	PageSize      int    `flag:"page-size"`
	HumanReadable bool   `flag:"human-redable"`
	Summarize     bool   `flag:"summarize"`
	RequestPlayer string `flag:"request-player"`
}
