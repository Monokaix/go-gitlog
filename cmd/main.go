package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/tsuyoshiwada/go-gitlog"
)

func init() {
	flag.StringVar(&start, "start", "14570c6f6278a9fc3a50202bfdc0e9b8a728a27f", "Start commit id.")
	flag.StringVar(&projectPath, "project-path", "D:\\go\\src\\volcano", "Your project absolute path.")
	flag.StringVar(&version, "version", "v1.8.2", "The release version.")
}

var (
	// flag
	start       string
	projectPath string
	version     string

	// release note template
	tpl = `- {{.Msg}} ([#{{.PRSeq}}](https://github.com/volcano-sh/volcano/pull/{{.PRSeq}}) **@{{.Author}}**)
`
	// regex expression
	reg           = regexp.MustCompile(`Merge pull request (.*) from (.*)`)
	cherryPickReg = regexp.MustCompile(`\[cherry.*\]`)
)

type Git struct {
	Msg    string
	PRSeq  string
	Author string
}

func main() {
	flag.Parse()
	// New gitlog
	git := gitlog.New(&gitlog.Config{
		Bin:  `git`, // default "git"
		Path: projectPath,
	})

	// List git-log
	commits, err := git.Log(nil, nil)
	if err != nil {
		log.Fatalln(err)
	}

	f, err := os.Create(fmt.Sprintf("release note for %s.md", version))
	if err != nil {
		panic(err)
	}
	if _, err = f.Write([]byte(fmt.Sprintf("### release note for %s\n", version))); err != nil {
		panic(err)
	}

	// Output
	for _, commit := range commits {
		if commit.Hash.Long == start {
			break
		}
		if commit.Author.Name != "Volcano Bot" {
			continue
		}

		matchArr := reg.FindStringSubmatch(commit.Subject)
		if len(matchArr) != 3 {
			panic(fmt.Errorf("wrong pr format, commit message: %s", commit.Subject))
		}

		// extract real author of github id.
		parts := strings.Split(matchArr[len(matchArr)-1], "/")
		if len(parts) != 2 {
			panic(fmt.Errorf("failed to extract author: %s", matchArr[len(matchArr)-1]))
		}

		// extract pr sequence
		prSeq := strings.TrimPrefix(matchArr[len(matchArr)-2], "#")

		gitLog := Git{
			Msg:    strings.TrimSpace(cherryPickReg.ReplaceAllString(commit.Body, "")),
			PRSeq:  prSeq,
			Author: parts[0],
		}
		templ, err := template.New("release note").Parse(tpl)
		if err != nil {
			panic(err)
		}
		err = templ.Execute(f, gitLog)
		if err != nil {
			panic(err)
		}
	}
}
