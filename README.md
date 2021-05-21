# mdim - Markdown Images Maintainer

[![Go](https://github.com/bunnier/mdim/actions/workflows/go.yml/badge.svg)](https://github.com/bunnier/mdim/actions/workflows/go.yml)
[![CodeQL](https://github.com/bunnier/mdim/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/bunnier/mdim/actions/workflows/codeql-analysis.yml)

The tool will help to maintain the image relative paths of markdown files and cleanup no reference images.

Install:

```bash
cd ./mdim
go install
```

Usage:

```bash
mdim [-h] [-d] [-f] [-i imageFolder] [-m markdownFolder] 
```

Options:

- `-d` Set the option to delete no reference images.
- `-f` Set the option to fix image relative paths of markdown documents.
- `-h` Show this help.
- `-i string` Must be not empty. The folder images save in.
- `-m string` Must be not empty. The folder markdown documents save in.
