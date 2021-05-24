# mdim - Markdown Images Maintainer

[![Go](https://github.com/bunnier/mdim/actions/workflows/go.yml/badge.svg)](https://github.com/bunnier/mdim/actions/workflows/go.yml)
[![CodeQL](https://github.com/bunnier/mdim/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/bunnier/mdim/actions/workflows/codeql-analysis.yml)

## Function

Now:

- Fixed wrong image relative path after you move document.
- Clean up no reference image.

Next:

- Convert the web images in your docs to local images.
- Read the paths from enviroment variable.

## Usage

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

- `-d` Set the option to delete no reference images, otherwise print the paths only.
- `-f` Set the option to fix image relative paths of markdown documents, otherwise print the paths only.
- `-h` Show this help.
- `-i string` Must be not empty. The folder images save in.
- `-m string` Must be not empty. The folder markdown documents save in.
