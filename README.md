# The Register CLI

The Register CLI a.k.a. `elreg-cli` is a command-line interface (CLI) tool written in Go that allows you to read articles from The Register (https://theregister.com) directly in your terminal.

## Installation

To install elreg-cli, you need to have Go installed on your system. Once you have Go set up, you can install the tool by running the following command:

```
go install github.com/CADawg/elreg-cli@latest
```

This will install the `elreg-cli` command in your Go bin directory.

## Usage

To start the application, simply run the `elreg-cli` command in your terminal:

```
elreg-cli
```

Upon launching, you will see the homepage of The Register, displaying a list of headlines. You can navigate through the headlines using the following commands:

- `<number>`: Type the number corresponding to the headline you want to read.
- `a`: Show all headlines at once.
- `s`: Show the next 5 headlines.
- `q`: Quit the application.

When reading an article, you can use the following commands:

- `Enter`: Advance to the next paragraph/line of the article.
- `a`: Display the entire article at once.
- `q`: Go back to the main menu.
- `r`: Mark the article as read and go back to the main menu.

## Features

- Read articles from The Register directly in your terminal.
- Navigate through headlines and articles using simple commands.
- Supports pagination for long articles.

## Limitations
- No arrow key support.
- No search functionality.

## Contributing

If you find any issues or have suggestions for improvements, feel free to open an issue or submit a pull request on the [GitHub repository](https://github.com/CADawg/elreg-cli).