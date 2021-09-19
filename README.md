# mdim - Markdown Images Maintainer

[![Go](https://github.com/bunnier/mdim/actions/workflows/go.yml/badge.svg)](https://github.com/bunnier/mdim/actions/workflows/go.yml) [![CodeQL](https://github.com/bunnier/mdim/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/bunnier/mdim/actions/workflows/codeql-analysis.yml)

> The tool helps to maintain the images in the markdown files.

## Usage

### Fix wrong image relative paths after moved document

```bash
./mdim --doc {doc_path} --imgFolder {img_folder} --relfix # for single document
./mdim --docFolder {doc_folder} --imgFolder {img_folder} --relfix # for folder
```

### Download the web images in docs into local folder

```bash
./mdim --doc {doc_path} --imgFolder {img_folder} --web # for single document
./mdim --docFolder {doc_folder} --imgFolder {img_folder} --web # for folder
```

### Clean up no reference images

```bash
./mdim --doc {doc_path} --imgFolder {img_folder} --delete # for single document
./mdim --docFolder {doc_folder} --imgFolder {img_folder} --delete # for folder
```

### Convert local markdown images to web images

```bash
./mdim qiniu --doc {doc_path} --imgFolder {img_folder} --ak {qiniu_ak} --sk {qiniu_sk} --bucket {qiniu_bucket} # for single document
./mdim qiniu --docFolder {doc_folder} --imgFolder {img_folder} --ak {qiniu_ak} --sk {qiniu_sk} --bucket {qiniu_bucket} # for folder
```

## Document

### mdim

```explain
The tool helps to maintain the images in the markdown files.
Github: https://github.com/bunnier/mdim

Usage:
  mdim [flags]
  mdim [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  help        Help about any command
  qiniu       Uploading the local image files in specific markdown files to Qiniu cloud space.

Flags:
  -d, --delete             Set this option to delete no reference images, otherwise print the paths only.
  -m, --doc string         Assign the target markdown document. There must be at least one of '--doc' and '--docFolder'.
  -f, --docFolder string   Assign the folder which markdown documents save in, also can be provided by setting env variable named 'mdim_docFolder'
  -h, --help               help for mdim
  -i, --imgFolder string   Must not be empty. Assign the folder which images save in, also can be provided by setting env variable named 'mdim_imgFolder'.
  -r, --relfix             Set this option to fix wrong local relative path of images after moved document.
  -v, --version            version for mdim
  -w, --web                Set this option to download reference web images to 'imgFolder'.

Use "mdim [command] --help" for more information about a command.
```

### mdim qiniu

```explain
The tool helps to upload the local image files in specific markdown files to Qiniu cloud space.
Github: https://github.com/bunnier/mdim

Usage:
  mdim qiniu [flags]

Flags:
      --ak string          Must not be empty. Assign the AK(Access Key) of Qiniu SDK, also can be provided by setting env variable named 'mdim_qiniu_ak'.
  -b, --bucket string      Must not be empty. Assign the Bucket of Qiniu SDK, also can be provided by setting env variable named 'mdim_qiniu_bucket'.
  -m, --doc string         Assign the target markdown document. There must be at least one of '--doc' and '--docFolder'.
  -f, --docFolder string   Assign the folder which markdown documents save in, also can be provided by setting env variable named 'mdim_docFolder'
  -d, --domain string      The domain of uploaded image url, also can be provided by setting env variable named 'mdim_qiniu_domain'. If do not assign the option, will use first domain in specific bucket
  -u, --folder string      After uploaded, you images can access in this url '{protocal}://{domain}/{path}/{img_name}'.
  -h, --help               help for qiniu
  -s, --https              If assign this option, will use https instead of http.
  -i, --imgFolder string   Must not be empty. Assign the folder which images save in, also can be provided by setting env variable named 'mdim_imgFolder'.
      --sk string          Must not be empty. Assign the SK(Secret Key) of Qiniu SDK, also can be provided by setting env variable named 'mdim_qiniu_sk'.
```
