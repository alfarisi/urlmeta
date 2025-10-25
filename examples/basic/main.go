package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/alfarisi/urlmeta"
)

func main() {
	fmt.Println("URLMeta - Basic Example")
	fmt.Println("=======================\n")

	// List of URLs to extract metadata from
	urls := []string{
		"https://github.com/golang/go",
		"https://www.youtube.com/watch?v=MbyvLY8CGFM", // Has oEmbed!
		"wordpress.com", // Auto adds https://
		"https://www.nytimes.com",
		"https://vimeo.com/1234567", // Has oEmbed!
	}

	for i, url := range urls {
		fmt.Printf("\n[%d] Extracting metadata from: %s\n", i+1, url)
		fmt.Println(strings.Repeat("-", 60))

		// Extract metadata - automatically includes oEmbed if available!
		metadata, err := urlmeta.Extract(url)
		if err != nil {
			log.Printf("❌ Error: %v\n", err)
			continue
		}

		// Display results
		displayMetadata(metadata)
	}
}

func displayMetadata(m *urlmeta.Metadata) {
	fmt.Printf("✓ Title: %s\n", m.Title)
	fmt.Printf("✓ Description: %s\n", truncate(m.Description, 100))
	fmt.Printf("✓ Provider: %s\n", m.ProviderName)
	fmt.Printf("✓ Type: %s\n", m.Type)
	
	if m.Author != "" {
		fmt.Printf("✓ Author: %s\n", m.Author)
	}
	
	if len(m.Images) > 0 {
		fmt.Printf("✓ Images: %d\n", len(m.Images))
		for i, img := range m.Images {
			if i >= 2 { // Show max 2 images
				break
			}
			fmt.Printf("  - %s", img.URL)
			if img.Width > 0 && img.Height > 0 {
				fmt.Printf(" (%dx%d)", img.Width, img.Height)
			}
			fmt.Println()
		}
	}
	
	if len(m.Videos) > 0 {
		fmt.Printf("✓ Videos: %d\n", len(m.Videos))
	}
	
	if m.Favicon != "" {
		fmt.Printf("✓ Favicon: %s\n", m.Favicon)
	}
	
	// Check for oEmbed data (automatically included!)
	if m.OEmbed != nil {
		fmt.Println("\n✨ oEmbed Data Available:")
		fmt.Printf("  Type: %s\n", m.OEmbed.Type)
		if m.OEmbed.AuthorName != "" {
			fmt.Printf("  Author: %s\n", m.OEmbed.AuthorName)
		}
		if m.OEmbed.HTML != "" {
			fmt.Printf("  Embed Code: Available (%d chars)\n", len(m.OEmbed.HTML))
		}
		if m.OEmbed.ThumbnailURL != "" {
			fmt.Printf("  Thumbnail: %s\n", m.OEmbed.ThumbnailURL)
		}
		if m.OEmbed.Width > 0 && m.OEmbed.Height > 0 {
			fmt.Printf("  Dimensions: %dx%d\n", m.OEmbed.Width, m.OEmbed.Height)
		}
	}
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}