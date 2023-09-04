# Uby

🥦 Github scraper that helps you find all the other repositories that depend on a particular GitHub repository!

## The Name
The name Uby is derived from "Used by" in GitHub.

## Introduction
Uby is a command-line application built with Go and the Go-Colly package. It allows you to find repositories that depend on a specific GitHub repository. Uby scapes & collects the necessary data and writes the results to `deps.json`.

## How It Works
Uby employs Go-Colly, a powerful web scraping framework, to extract the required information from GitHub. The application prompts the user to enter the URL or name of a repository. It then crawls through the GitHub ecosystem to identify other repositories that depend on the specified repository. The search results are saved in a structured JSON format in the `deps.json` file.

## Installation
To use the Uby, follow these steps:
1. Clone the repository to your local machine.
2. Navigate to the project directory.
3. Install the necessary dependencies by running the following command:
   ```
   go mod download
   ```

## Usage
To run the dependent seach, use the following command in your terminal:
```sh
$ cd ./scraper
$ go run .
```

The application will begin crawling through the GitHub ecosystem to find repositories that depend on the specified repository. The search results will be saved in the `deps.json` file in the project directory.
