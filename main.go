package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/anacrolix/torrent"
)

func main() {
	// Get torrent path from user
	torrentPath := getTorrentPath()
	downloadDir := "./downloads" // Directory where files will be downloaded

	// Create download directory if it doesn't exist
	err := os.MkdirAll(downloadDir, 0755)
	if err != nil {
		log.Fatalf("Failed to create download directory: %v", err)
	}

	// Configure torrent client
	cfg := torrent.NewDefaultClientConfig()
	cfg.DataDir = downloadDir
	cfg.ListenPort = 42069 // You can change the port if needed

	// Create new torrent client
	client, err := torrent.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating torrent client: %v", err)
	}
	defer client.Close()

	// Add torrent
	var t *torrent.Torrent
	if filepath.Ext(torrentPath) == ".torrent" {
		// Add from torrent file
		t, err = client.AddTorrentFromFile(torrentPath)
	} else {
		// Assume it's a magnet link
		t, err = client.AddMagnet(torrentPath)
	}
	if err != nil {
		log.Fatalf("Error adding torrent: %v", err)
	}

	// Wait for torrent info
	<-t.GotInfo()
	fmt.Printf("Downloading: %s\n", t.Name())

	// Start download
	t.DownloadAll()

	// Create a channel to signal completion
	done := make(chan struct{})

	// Monitor progress and signal completion
	go printProgress(t, done)

	// Wait for download to complete, then exit
	<-done
	fmt.Println("Exiting application...")
}

func getTorrentPath() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter torrent file path or magnet link: ")
	path, _ := reader.ReadString('\n')
	return strings.TrimSpace(path)
}

func printProgress(t *torrent.Torrent, done chan<- struct{}) {
	lastBytes := t.BytesCompleted()
	lastTime := time.Now()
	var avgSpeed float64 // exponential moving average of download speed
	firstIteration := true

	for {
		bs := t.BytesCompleted()
		total := t.Length()
		progress := float64(bs) / float64(total) * 100

		// Calculate download speed and estimated time remaining
		var eta string
		currentTime := time.Now()
		elapsedTime := currentTime.Sub(lastTime).Seconds()
		
		// Skip ETA calculation on first iteration to avoid near-zero elapsed time
		if firstIteration {
			eta = "calculating..."
			firstIteration = false
		} else if elapsedTime > 0 {
			downloadedBytes := bs - lastBytes
			instantSpeed := float64(downloadedBytes) / elapsedTime // bytes per second
			
			// Use exponential moving average for smoother speed calculation
			// alpha = 0.3 gives more weight to recent measurements while smoothing volatility
			if avgSpeed == 0 {
				avgSpeed = instantSpeed
			} else {
				avgSpeed = 0.3*instantSpeed + 0.7*avgSpeed
			}
			
			if avgSpeed > 0 {
				remainingBytes := total - bs
				secondsRemaining := float64(remainingBytes) / avgSpeed
				
				// Cap maximum ETA at 99 hours to prevent overflow and display sensibly
				const maxSeconds = 99 * 3600 // 99 hours
				if secondsRemaining > maxSeconds {
					eta = "99h+"
				} else {
					eta = formatDuration(time.Duration(secondsRemaining) * time.Second)
				}
			} else {
				eta = "calculating..."
			}
		} else {
			eta = "calculating..."
		}
		lastBytes = bs
		lastTime = currentTime

		// Get file-level information
		files := t.Files()
		filesCompleted := 0
		totalFiles := len(files)
		currentFile := ""
		var firstIncompleteFile string

		for _, file := range files {
			fileBytes := file.BytesCompleted()
			fileLength := file.Length()
			if fileBytes == fileLength {
				filesCompleted++
			} else {
				// Track first incomplete file as fallback
				if firstIncompleteFile == "" {
					firstIncompleteFile = file.DisplayPath()
				}
				// Prioritize files with partial progress
				if fileBytes > 0 && currentFile == "" {
					currentFile = file.DisplayPath()
				}
			}
		}

		// If no active file found, use the first incomplete file
		if currentFile == "" {
			currentFile = firstIncompleteFile
		}

		// Display progress with file information and ETA
		if currentFile != "" {
			fmt.Printf("\rProgress: %.2f%% - %d/%d bytes | Files: %d/%d completed | ETA: %s | Current: %s",
				progress, bs, total, filesCompleted, totalFiles, eta, currentFile)
		} else {
			fmt.Printf("\rProgress: %.2f%% - %d/%d bytes | Files: %d/%d completed | ETA: %s",
				progress, bs, total, filesCompleted, totalFiles, eta)
		}

		if bs == total {
			fmt.Println("\nDownload completed!")
			done <- struct{}{} // Signal completion
			return
		}
		time.Sleep(time.Second)
	}
}

// formatDuration converts a duration into a human-readable string
func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}
