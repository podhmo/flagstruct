package sts

type GetCurrentIdentityOptions struct {
	CliInputJson        string `flag:"cli-input-json"`
	GenerateCliSkeleton string `flag:"generate-cli-skeleton"`
}

type AssumeRoleOptions struct {
	RoleArn             string   `flag:"role-arn" required:"true"`
	RoleSessionName     string   `flag:"role-session-name" required:"true"`
	Policy              string   `flag:"policy"`
	DurationSeconds     int      `flag:"duration-seconds"`
	TransitiveTagKeys   []string `flag:"transitive-tag-keys"`
	ExternalId          string   `flag:"external-id"`
	SerialNumber        string   `flag:"serial-number"`
	TokenCode           string   `flag:"token-code"`
	SourceIdentity      string   `flag:"source-identity"`
	CliInputJson        string   `flag:"cli-input-json"`
	GenerateCliSkeleton string   `flag:"generate-cli-skeleton"`
}
