# Torrent Downloader

A simple command-line torrent downloader written in Go. This application allows users to download files from torrent files or magnet links, displaying progress updates, and automatically exits when the download completes.

## Features
- Supports both `.torrent` files and magnet links
- Downloads files to a `downloads/` directory
- Displays real-time progress (percentage and bytes downloaded)
- Automatically closes when the download is complete
- Cross-platform, with specific build instructions for Windows x64

## Prerequisites
- [Go](https://golang.org/dl/) (version 1.17 or later recommended)
- Git (optional, for cloning the repository)

## Installation

1. **Clone the Repository** (if using Git):
   ```bash
   git clone https://github.com/perakit/torrent-downloader.git
   cd torrent-downloader
   ```
   Alternatively, download the source code manually.

2. **Initialize the Go Module** (if not already done):
   ```bash
   go mod init torrent-downloader
   ```

3. **Install Dependencies**:
   ```bash
   go get github.com/anacrolix/torrent
   ```

## Building

### For Windows x64
To compile the project for Windows 64-bit:
```bash
go build -o torrent-downloader.exe
```
This creates a `torrent-downloader.exe` executable in the project directory.

### Cross-Compiling from Linux/Mac
If you're on Linux or macOS and want to target Windows x64:
```bash
GOOS=windows GOARCH=amd64 go build -o torrent-downloader.exe
```

### Optional Build Flags
- **Strip Debug Info** (smaller executable):
  ```bash
  GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o torrent-downloader.exe
  ```

## Usage
1. Run the compiled executable:
   ```bash
   torrent-downloader.exe  # On Windows
   ./torrent-downloader    # If built for your current OS
   ```

2. When prompted, enter either:
   - A path to a `.torrent` file (e.g., `C:\path\to\file.torrent`)
   - A magnet link (e.g., `magnet:?xt=urn:btih:...`)
   Example:
   ```
   Enter torrent file path or magnet link: C:\torrents\example.torrent
   ```

3. The program will:
   - Create a `downloads/` directory if it doesn’t exist
   - Start downloading and show progress
   - Exit automatically when the download completes

   Sample output:
   ```
   Downloading: example-file
   Progress: 75.50% - 75500000/100000000 bytes
   Download completed!
   Exiting application...
   ```

## Project Structure
```
torrent-downloader/
├── main.go          # Main application code
├── go.mod          # Module definition
├── go.sum          # Dependency checksums
├── .gitignore      # Git ignore file
└── downloads/      # Downloaded files (created at runtime)
```

## Notes
- Ensure you have write permissions in the directory where you run the program.
- For Windows paths with backslashes, use double backslashes (`\\`) or forward slashes (`/`).
- The `.gitignore` excludes the `downloads/` directory and compiled binaries.

## Legal Notice
This is a technical demonstration. Ensure you have the right to download and distribute any content you use with this tool. The author is not responsible for any misuse.

## Contributing
Feel free to fork this repository, submit issues, or send pull requests with improvements!

## License
This project is unlicensed (public domain) unless specified otherwise. Use at your own risk.