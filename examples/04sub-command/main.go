package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/podhmo/structflag"
)

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
type App struct {
	b *structflag.Builder
}

func (app *App) Fail(name string) error {
	fmt.Printf("%q is invalid subcommand\n", name)
	os.Exit(2)
	return fmt.Errorf("never")
}

func (app *App) S3(args []string) error {
	app.b.Name = app.b.Name + " s3"
	if len(args) < 1 {
	msg := "available s3 subcommands are {ls, cp, mv}\n"
		fmt.Fprintln(os.Stderr, msg)
		os.Exit(2)
	}
	switch args[0] {
	case "ls":
		return app.S3Ls(args[1:])
	case "cp":
		return app.S3Cp(args[1:])
	case "mv":
		return app.S3Mv(args[1:])
	default:
		return app.Fail(os.Args[1])
	}
}

type S3LsOptions struct {
	Recursive     bool   `flag:"recursive"`
	PageSize      int    `flag:"page-size"`
	HumanReadable bool   `flag:"human-redable"`
	Summarize     bool   `flag:"summarize"`
	RequestPlayer string `flag:"request-player"`
}

func (app *App) S3Ls(args []string) error {
	options := &S3LsOptions{}
	app.b.Name = app.b.Name + " ls"
	fs := app.b.Build(options)
	if err := fs.Parse(args); err != nil {
		return err
	}
	return json.NewEncoder(os.Stdout).Encode(options)

}
func (app *App) S3Cp(args []string) error {
	return nil
}
func (app *App) S3Mv(args []string) error {
	return nil
}
func (app *App) STS(args []string) error {
	switch args[0] {
	case "get-current-identity":
		return app.STSGetCurrentIdentity(args[1:])
	case "assume-role":
		return app.STSAssumeRole(args[1:])
	default:
		return app.Fail(os.Args[1])
	}
}
func (app *App) STSGetCurrentIdentity(args []string) error {
	return nil
}
func (app *App) STSAssumeRole(args []string) error {
	return nil
}

func main() {
	b := structflag.NewBuilder()
	b.Name = "awscli-like"
	app := &App{b: b}
	var err error
	if len(os.Args) <= 1 {
		fmt.Printf("available subcommands are {s3, sts}\n")
		os.Exit(2)
	}
	switch os.Args[1] {
	case "s3":
		err = app.S3(os.Args[2:])
	case "sts":
		err = app.STS(os.Args[2:])
	default:
		err = app.Fail(os.Args[1])
	}
	if err != nil {
		log.Fatalf("!! %+v", err)
	}
}
