# Jetty 🚢

A simple CLI tool to quickly open web interfaces for running Docker container in your browser.

As a developer, I frequently run multiple Docker containers for my projects, including services like phpMyAdmin and MailHog, each on different ports to avoid conflicts. I built Jetty to simplify accessing these web interfaces.

## Features

- 🔍 Smart container name matching
- 🎯 Auto-selects single matches
- 🌐 Interactive selection for multiple containers/ports
- 🖥️ Cross-platform browser support (macOS, Linux, Windows)
- ⚡️ Command-line flags for automation

## Installation

### Download binary

You can download a binary from the releases page: https://github.com/smitmartijn/jetty/releases

### Using Homebrew

```
brew tap smitmartijn/jetty
brew install jetty-docker
```

### Building from source

Requires Go 1.23.3+:

```bash
git clone https://github.com/username/jetty.git
cd jetty
go build
```

## Usage
Interactive mode:

```bash
./jetty
```

Direct mode:

```bash
./jetty --name container-name
```

### Pro-tip: Add aliases

To gain quicker access to the web interfaces you need, add command aliases to your `.bash_profile`, `.zshrc` or other equivalent shell profile. Here's an example:

```bash
export JETTY_PATH=/Users/martijn/Projects/jetty/jetty # binary path
alias jetty="$JETTY_PATH"
alias phpmyadmin="$JETTY_PATH --name phpmyadmin"
alias mailhog="$JETTY_PATH --name mailhog"
```

With these, you can open the PHPMyAdmin interface by just typing `phpmyadmin` in your terminal!

## How It Works
1. Lists running Docker containers (optionally filtered by name)
2. Presents interactive selection if multiple matches found
3. Shows available port mappings for chosen container
4. Opens localhost URL in your default browser

## Example
```bash
$ ./jetty --name nginx
Enter the container name: nginx
Choose a container:
  > nginx-proxy (Up 2 hours)
    nginx-app (Up 30 minutes)
Choose a port:
  > 80->80/tcp
    443->443/tcp
```

## Requirements
- Go 1.23.3+
- Docker daemon running
- Web browser
- System commands:
  - macOS: open
  - Linux: xdg-open
  - Windows: start

## Contributing

Pull requests welcome! Please open an issue first to discuss major changes.

