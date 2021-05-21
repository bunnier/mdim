# mdic - Markdown Images Cleaner

[![Go](https://github.com/bunnier/mdic/actions/workflows/go.yml/badge.svg)](https://github.com/bunnier/mdic/actions/workflows/go.yml)

The tool will help you maintain the image relative paths of markdown files and cleanup no reference images.

Install:

```bash
cd ./mdic
go install
```

Usage:

```bash
mdic [-h] [-d] [-f] [-i imageFolder] [-m markdownFolder] 
```

Options:

- `-d` Set the option to delete no reference images.
- `-f` Set the option to fix image relative paths of markdown documents.
- `-h` Show this help.
- `-i string` Must be not empty. The folder images save in.
- `-m string` Must be not empty. The folder markdown documents save in.
