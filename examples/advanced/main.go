package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/alfarisi/urlmeta"
)

func main() {
	fmt.Println("URLMeta - Advanced Example")
	fmt.Println("==========================\n")

	// Create a custom client with options
	client := urlmeta.NewClient(
		urlmeta.WithTimeout(15*time.Second),
		urlmeta.WithUserAgent("MyCustomBot/1.0 (+https://mywebsite.com)"),
		urlmeta.WithMaxRedirects(5),
	)

	// Extract metadata
	url := "https://www.theverge.com"
	fmt.Printf("Extracting metadata from: %s\n\n", url)

	metadata, err := client.Extract(url)
	if err != nil {
		log.Fatalf("Failed to extract metadata: %v", err)
	}

	// Display full metadata as JSON
	fmt.Println("=== Full Metadata (JSON) ===")
	displayJSON(metadata)

	fmt.Println("\n=== Structured Display ===")
	displayStructured(metadata)

	// Save to file example
	fmt.Println("\n=== Saving to JSON file ===")
	saveToFile(metadata, "metadata.json")
}

func displayJSON(m *urlmeta.Metadata) {
	jsonData, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		log.Printf("Error marshaling to JSON: %v", err)
		return
	}
	fmt.Println(string(jsonData))
}

func displayStructured(m *urlmeta.Metadata) {
	fmt.Println("📄 Basic Information:")
	fmt.Printf("   Title: %s\n", m.Title)
	fmt.Printf("   Description: %s\n", m.Description)
	fmt.Printf("   URL: %s\n", m.URL)
	if m.CanonicalURL != "" {
		fmt.Printf("   Canonical URL: %s\n", m.CanonicalURL)
	}

	fmt.Println("\n🏢 Provider Information:")
	fmt.Printf("   Name: %s\n", m.ProviderName)
	fmt.Printf("   Display: %s\n", m.ProviderDisplay)
	fmt.Printf("   URL: %s\n", m.ProviderURL)

	if m.Type != "" || m.SiteName != "" || m.Locale != "" {
		fmt.Println("\n🏷️  OpenGraph Data:")
		if m.Type != "" {
			fmt.Printf("   Type: %s\n", m.Type)
		}
		if m.SiteName != "" {
			fmt.Printf("   Site Name: %s\n", m.SiteName)
		}
		if m.Locale != "" {
			fmt.Printf("   Locale: %s\n", m.Locale)
		}
	}

	if len(m.Images) > 0 {
		fmt.Println("\n🖼️  Images:")
		for i, img := range m.Images {
			fmt.Printf("   [%d] %s", i+1, img.URL)
			if img.Width > 0 && img.Height > 0 {
				fmt.Printf(" (%dx%d)", img.Width, img.Height)
			}
			if img.Alt != "" {
				fmt.Printf(" - %s", img.Alt)
			}
			fmt.Println()
		}
	}

	if len(m.Videos) > 0 {
		fmt.Println("\n🎥 Videos:")
		for i, video := range m.Videos {
			fmt.Printf("   [%d] %s", i+1, video.URL)
			if video.Type != "" {
				fmt.Printf(" (type: %s)", video.Type)
			}
			if video.Width > 0 && video.Height > 0 {
				fmt.Printf(" (%dx%d)", video.Width, video.Height)
			}
			fmt.Println()
		}
	}

	if m.Author != "" || m.PublishedTime != "" || m.ModifiedTime != "" {
		fmt.Println("\n✍️  Author & Dates:")
		if m.Author != "" {
			fmt.Printf("   Author: %s\n", m.Author)
		}
		if m.PublishedTime != "" {
			fmt.Printf("   Published: %s\n", m.PublishedTime)
		}
		if m.ModifiedTime != "" {
			fmt.Printf("   Modified: %s\n", m.ModifiedTime)
		}
	}

	if len(m.Keywords) > 0 {
		fmt.Println("\n🔑 Keywords:")
		fmt.Printf("   %v\n", m.Keywords)
	}

	if m.TwitterCard != "" || m.TwitterSite != "" || m.TwitterCreator != "" {
		fmt.Println("\n🐦 Twitter Card:")
		if m.TwitterCard != "" {
			fmt.Printf("   Card Type: %s\n", m.TwitterCard)
		}
		if m.TwitterSite != "" {
			fmt.Printf("   Site: %s\n", m.TwitterSite)
		}
		if m.TwitterCreator != "" {
			fmt.Printf("   Creator: %s\n", m.TwitterCreator)
		}
	}

	if m.Favicon != "" {
		fmt.Println("\n🎨 Favicon:")
		fmt.Printf("   %s\n", m.Favicon)
	}
}

func saveToFile(m *urlmeta.Metadata, filename string) {
	jsonData, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		log.Printf("Error marshaling to JSON: %v", err)
		return
	}

	// In a real application, you would write to file:
	// err = os.WriteFile(filename, jsonData, 0644)
	// if err != nil {
	//     log.Printf("Error writing to file: %v", err)
	//     return
	// }

	fmt.Printf("✓ Metadata saved to %s (%d bytes)\n", filename, len(jsonData))
	fmt.Println("  (In this example, we're not actually writing to disk)")
}