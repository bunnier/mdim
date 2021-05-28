# mdim - Markdown Images Maintainer

[![Go](https://github.com/bunnier/mdim/actions/workflows/go.yml/badge.svg)](https://github.com/bunnier/mdim/actions/workflows/go.yml)
[![CodeQL](https://github.com/bunnier/mdim/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/bunnier/mdim/actions/workflows/codeql-analysis.yml)

## Function

- Fix wrong image relative paths after moved document.
- Download the web images in docs into local folder.
- Clean up no reference image.
- Read the paths in cli-options from enviroment variable.

## Usage

Build:

```bash
cd ./mdim
go build
```

Usage:

```bash
Usage: mdim [-h] [-d] [-w] [-i imageFolder] [-m markdownFolder] 
```

Description: The tool will help to maintain the images in markdown files.

Options:

- `-s` Set the option to save markdown document changes, otherwise print scan result only.
- `-w` Set the option to download web images to imageFolder. This option often be set with the `-s` option, otherwise although images have been download to imageFolder, the path in document still be web paths.
- `-d` Set the option to delete no reference images, otherwise print the paths only.
- `-i imageFolder` Must not be empty. Assign the folder which images save in, also can be provided by setting env variable named `mdim_imgFolder`.
- `-m markdownFolder` Must not be empty. Assign the folder which markdown documents save in, also can be provided by setting env variable named `mdim_docFolder`.
- `-h` Show this help.
