## Introduction

This program generates release note for volcano automatically, it mostly use `git log` to extract git message and output them as a Markdown format file.

## Usage

```shell
git clone https://github.com/Monokaix/go-gitlog.git

cd go-gitlog && go run cmd/main.go --start $start commit id --version $release version --project-path $volcano project path
```

### example

```go
go run cmd/main.go --start 14570c6f6278a9fc3a50202bfdc0e9b8a728a27f --version v1.8.2 --project-path "D:\go\src\volcano"
```

The `--start` param indicates from which commit id begin to count, the `--version` param indicates release version, the `--project-path` param indicates the volcano project absolute path, and it will generate a markdown file in current dir named `release note for v1.8.2.md` after execute it.