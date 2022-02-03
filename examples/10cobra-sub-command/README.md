# sub command example

awscli like command with cobra integration

```go
/*
- s3
	- ls
		-- recursive
		-- page-size <value:integer>
		-- human-redable
		-- summrize
		-- request-player <value:string>
		[s3url]
	- cp
		--dryrun
		--quiet
		--include <value:string>
		--exclude <value:string>
		--acl <valu:stringe
		...
		--metadata <value:map> // not supported yet.
		--metadata-directive <value:string>
		--expected-size <value:integer>
		--recursive
	^ mv
		--dryrun
		--quiet
		--include <value:string>
		--exclude <value:string>
		--acl <valu:stringe
		...
		--metadata <value:map> // not supported yet.
		--metadata-directive <value:string>
		--expected-size <value:integer>
		--recursive
- sts
	- get-current-identity
		--cli-input-json <value:string>
		--generate-cli-skeleton <value:string>
	- assume-role
		--role-arn <value:string> // required
		--role-session-name <value:string> // required
		--policy-arns <value:list> // not supported yet (JSON)
		--policy <value:string>
		--duration-seconds <value:integer>
		--tags <value:list> // not supported yet (JSON)
		--transtive-tag-keys <value:list of string>
		--external-id <value:string>
		--serial-number <value:string>
		--token-code <value:string>
		--source-identity <value:string>
		--cli-input-json <value:string>
		--generate-cli-skeleton <value:string>
		<assume role>
*/
```

## see also

- https://github.com/spf13/cobra/blob/master/cobra/README.md

