# ToTxt

ToTxt is a command-line tool that converts a directory with code files into a single text file. This tool is useful for sharing code snippets or entire project structures in a plain text format, such as in a forum or chat application like [ChatGPT](https://chat.openai.com/).

## Features

- Concatenate multiple code files into a single text file.
- Specify a preamble file to include at the beginning of the output.
- Use a `.totxtignore` file to exclude specific files or patterns from the output.

## Installation

1. Clone this repository

2. Change into the project directory and build the executable: `go build`

## Usage

```
./totxt /path/to/directory [-p /path/to/preamble.txt] [-o /path/to/output.txt]
```

- `/path/to/directory`: The directory containing the code files you want to convert.
- `-p /path/to/preamble.txt`: (Optional) The path to a preamble file to include at the beginning of the output. If not specified, a default preamble will be used.
- `-o /path/to/output.txt`: (Optional) The path to the output file. If not specified, the default output will be `output.txt`.

### .totxtignore

You can create a `.totxtignore` file in the root directory of the directory you want to convert. This file should contain a list of file patterns to exclude from the output, one pattern per line.
