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
	for {
		bs := t.BytesCompleted()
		total := t.Length()
		progress := float64(bs) / float64(total) * 100

		fmt.Printf("\rProgress: %.2f%% - %d/%d bytes", progress, bs, total)
		if bs == total {
			fmt.Println("\nDownload completed!")
			done <- struct{}{} // Signal completion
			return
		}
		time.Sleep(time.Second)
	}
}