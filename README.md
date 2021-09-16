# mdim - Markdown Images Maintainer

[![Go](https://github.com/bunnier/mdim/actions/workflows/go.yml/badge.svg)](https://github.com/bunnier/mdim/actions/workflows/go.yml)
[![CodeQL](https://github.com/bunnier/mdim/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/bunnier/mdim/actions/workflows/codeql-analysis.yml)

## Function

The tool helps to maintain the images in the markdown files.

- Fix wrong image relative paths after moved document.
- Download the web images in docs into local folder.
- Clean up no reference image.

## Usage

```explain
Usage:
  mdim [flags]

Flags:
  -d, --delete             Set the option to delete no reference images, otherwise print the paths only.
  -m, --docFolder string   Must not be empty. Assign the folder which markdown documents save in, also can be provided by setting env variable named 'mdim_docFolder'
  -h, --help               help for mdim
  -i, --imgFolder string   Must not be empty. Assign the folder which images save in, also can be provided by setting env variable named 'mdim_imgFolder'.
  -s, --save               Set the option to save markdown document changes, otherwise print scan result only.
  -w, --web                Set the option to download web images to imageFolder. This option might be set with the -s option, otherwise although images have been download to imageFolder, the path in document still be url.
```
