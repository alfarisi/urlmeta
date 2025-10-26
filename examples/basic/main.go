package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/alfarisi/urlmeta"
)

func main() {
	fmt.Println("URLMeta - Quick Start Guide")
	fmt.Println("===========================")

	// Test URLs - mix of sites with and without oEmbed support
	urls := []string{
		"https://www.youtube.com/watch?v=dQw4w9WgXcQ", // Has oEmbed
		"https://github.com/golang/go",                // No oEmbed
		"https://vimeo.com/148751763",                 // Has oEmbed
		"https://dev.to",                              // No oEmbed
		"wordpress.com",                               // Auto adds https://
	}

	fmt.Println("🔍 Extracting metadata from multiple URLs...")
	fmt.Println("Note: One call gets BOTH metadata AND oEmbed automatically!")

	for i, url := range urls {
		fmt.Printf("[%d/%d] %s\n", i+1, len(urls), url)
		fmt.Println(strings.Repeat("-", 70))

		// Single call gets everything!
		metadata, err := urlmeta.Extract(url)
		if err != nil {
			log.Printf("❌ Error: %v\n\n", err)
			continue
		}

		displayMetadata(metadata)
		fmt.Println()
	}

	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("✨ Summary:")
	fmt.Println("- YouTube/Vimeo URLs: metadata + oEmbed embed code")
	fmt.Println("- GitHub/Dev.to URLs: metadata only (no embed)")
	fmt.Println("- All extracted in ONE call with Extract()")
}

func displayMetadata(m *urlmeta.Metadata) {
	// Basic metadata (always available)
	fmt.Printf("📄 Title: %s\n", truncate(m.Title, 60))
	fmt.Printf("📝 Description: %s\n", truncate(m.Description, 80))
	fmt.Printf("🏢 Provider: %s\n", m.ProviderName)

	if m.Type != "" {
		fmt.Printf("🏷️  Type: %s\n", m.Type)
	}

	if m.Author != "" {
		fmt.Printf("✍️  Author: %s\n", m.Author)
	}

	// Images
	if len(m.Images) > 0 {
		fmt.Printf("🖼️  Images: %d found\n", len(m.Images))
		for i, img := range m.Images {
			if i >= 2 { // Show max 2 images
				fmt.Printf("   ... and %d more\n", len(m.Images)-2)
				break
			}
			fmt.Printf("   • %s", truncate(img.URL, 50))
			if img.Width > 0 && img.Height > 0 {
				fmt.Printf(" (%dx%d)", img.Width, img.Height)
			}
			fmt.Println()
		}
	}

	if len(m.Videos) > 0 {
		fmt.Printf("🎥 Videos: %d\n", len(m.Videos))
	}

	if m.Favicon != "" {
		fmt.Printf("🎨 Favicon: %s\n", truncate(m.Favicon, 50))
	}

	// oEmbed data (automatically included when available!)
	if m.OEmbed != nil {
		fmt.Println("\n✨ EMBEDDABLE CONTENT DETECTED!")
		fmt.Printf("   Type: %s\n", m.OEmbed.Type)

		if m.OEmbed.AuthorName != "" {
			fmt.Printf("   Author: %s\n", m.OEmbed.AuthorName)
		}

		if m.OEmbed.HTML != "" {
			fmt.Printf("   Embed Code: ✅ Available (%d characters)\n", len(m.OEmbed.HTML))
			fmt.Printf("   Preview: %s\n", truncate(m.OEmbed.HTML, 60))
		}

		if m.OEmbed.ThumbnailURL != "" {
			fmt.Printf("   Thumbnail: %s\n", truncate(m.OEmbed.ThumbnailURL, 50))
		}

		if m.OEmbed.Width > 0 && m.OEmbed.Height > 0 {
			fmt.Printf("   Dimensions: %dx%d\n", m.OEmbed.Width, m.OEmbed.Height)
		}

		fmt.Println("   💡 Use metadata.OEmbed.HTML to embed this content!")
	} else {
		fmt.Println("\n📄 Standard link preview (no embed available)")
		fmt.Println("   💡 Use metadata.Title, Description, Images for preview card")
	}
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
