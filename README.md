# Rat (Go)
<p align="center">
<img width="400" src="doc/images/Rat.png" alt="Rat" title="Rat" />
</p>

## Why Rat Exists

Rat is a small terminal file viewer built to explore how syntax highlighting, buffered output, and CLI ergonomics work under the hood.

It is intentionally simpler than `bat`, but it solves the same everyday problem: quickly inspect source files in the terminal with readable colour and predictable output.

The project focuses on practical implementation details such as file handling, token scanning, language-specific highlighting rules, and keeping output fast enough for larger files.

## 🛠️ Prerequisites

- [Go](https://golang.org/dl/) installed to compile the code

## 📁 Setup

1. Clone this repository:

   ```bash
   git clone https://github.com/jonathon-chew/Rat.git
   cd Rat
   ```

2. Compile the binary:

   ```bash
   go build .
   ```
    or

   ```bash
   go build -o $HOME/go/bin .
   ```
## Usage

View a single file:

    ```bash
    Rat main.go
    ```

View multiple files:
    ```bash
    Rat main.go parser.go config.json
    ```
    or

    ```bash
    Rat *
    ```
## 📂 Output

Pass in one or more files and the output will pretty-print them.

## Current Status

Rat currently supports syntax highlighting for common source formats and is usable as a lightweight file viewer. It is still intentionally smaller than `bat`.

## Supported Languages

Rat currently supports syntax highlighting for:
- **Python** (.py) - Keywords, comments, strings, and numbers
- **JavaScript** (.js) - Keywords, comments, strings, and numbers
- **Go** (.go) - Keywords, comments, strings, and numbers
- **PowerShell** (.ps1) - Keywords, comments, strings, variables, and numbers
- **JSON** (.json) - Basic structure highlighting

Performance improvements help the tool handle large files efficiently through buffered output and optimized token scanning.

## 📜 License

This project is licensed under the MIT License. See the LICENSE file for details.

## 🖌️ Attribution

The Go Gopher was originally designed by [Renee French](https://reneefrench.blogspot.com/).  
Used under the [Creative Commons Attribution 4.0 License](https://creativecommons.org/licenses/by/4.0/).  
