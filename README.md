# Uby

🥦 Github scraper that helps you find all dependents of a project.

## The Name

The name Uby is derived from "Used by" in GitHub.

## Introduction

Uby is a command-line application built with Go and the Go-Colly package. It allows you to find repositories that depend on a specific GitHub repository. Uby scrapes & collects the necessary data and writes the results to `deps.json`.

## Installation

You can install the appropriate binary from the [releases page](https://github.com/somnek/uby/releases/tag/v0.1.0).

#### Note:

If you're on macOS, you may need to run xattr -c ./nvim-macos.tar.gz to (to avoid "unknown developer" warning)

## Usage

Use the following command in your terminal:

```sh
$ ./uby
```

1. Uby will ask for the repo dependents url which you can get by clicking `Used by` of the repo you want to search, for example: `https://github.com/spf13/cobra/network/dependents`
2. The application will begin crawling through the repo's dependents network. The search results will be saved in the `deps.json` file in the project directory.
