package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/alfarisi/urlmeta"
)

func main() {
	fmt.Println("URLMeta - Advanced Features")
	fmt.Println("===========================")

	// Example 1: Custom Client Configuration
	demonstrateCustomConfig()

	// Example 2: Manual oEmbed Control
	demonstrateManualOEmbed()

	// Example 3: Provider Support Checking
	demonstrateProviderChecking()

	// Example 4: JSON Export
	demonstrateJSONExport()

	// Example 5: Error Handling
	demonstrateErrorHandling()
}

func demonstrateCustomConfig() {
	fmt.Println("=== 1. Custom Client Configuration ===")

	// Create client with custom options
	client := urlmeta.NewClient(
		urlmeta.WithTimeout(15*time.Second),
		urlmeta.WithUserAgent("MyCustomBot/1.0 (+https://mywebsite.com)"),
		urlmeta.WithMaxRedirects(5),
		urlmeta.WithAutoOEmbed(true), // Auto oEmbed enabled (default)
	)

	url := "https://www.theverge.com"
	fmt.Printf("Extracting: %s\n", url)

	metadata, err := client.Extract(url)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("✅ Title: %s\n", metadata.Title)
	fmt.Printf("✅ Provider: %s\n", metadata.ProviderName)
	fmt.Printf("✅ Images: %d\n", len(metadata.Images))
	fmt.Printf("✅ oEmbed: %v\n\n", metadata.OEmbed != nil)
}

func demonstrateManualOEmbed() {
	fmt.Println("=== 2. Manual oEmbed Control ===")

	youtubeURL := "https://www.youtube.com/watch?v=dQw4w9WgXcQ"

	// Option A: Auto oEmbed (default behavior)
	fmt.Println("A) Auto Mode (default):")
	metadata1, _ := urlmeta.Extract(youtubeURL)
	fmt.Printf("   Extract() called once\n")
	fmt.Printf("   metadata.OEmbed available: %v\n", metadata1.OEmbed != nil)
	if metadata1.OEmbed != nil {
		fmt.Printf("   Embed type: %s\n", metadata1.OEmbed.Type)
		fmt.Printf("   Has HTML: %v\n", metadata1.OEmbed.HTML != "")
	}

	// Option B: Disable auto oEmbed for faster extraction
	fmt.Println("\nB) Manual Mode (auto disabled):")
	client := urlmeta.NewClient(
		urlmeta.WithAutoOEmbed(false), // Disable auto oEmbed
		urlmeta.WithTimeout(5*time.Second),
	)

	metadata2, _ := client.Extract(youtubeURL)
	fmt.Printf("   Extract() called - fast mode\n")
	fmt.Printf("   metadata.OEmbed: %v (disabled)\n", metadata2.OEmbed)

	// Explicitly extract oEmbed when needed
	oembed, err := client.ExtractOEmbed(youtubeURL)
	if err != nil {
		fmt.Printf("   ExtractOEmbed() error: %v\n", err)
	} else {
		fmt.Printf("   ExtractOEmbed() called explicitly\n")
		fmt.Printf("   Embed type: %s\n", oembed.Type)
		fmt.Printf("   Author: %s\n\n", oembed.AuthorName)
	}
}

func demonstrateProviderChecking() {
	fmt.Println("=== 3. Provider Support Checking ===")

	// Check if specific URLs support oEmbed
	testURLs := []string{
		"https://www.youtube.com/watch?v=123",
		"https://vimeo.com/123456",
		"https://twitter.com/user/status/123",
		"https://github.com/golang/go",
		"https://soundcloud.com/artist/track",
		"https://example.com/random",
	}

	fmt.Println("Checking oEmbed support:")
	for _, url := range testURLs {
		supported := urlmeta.IsOEmbedSupported(url)
		status := "❌"
		if supported {
			status = "✅"
		}
		fmt.Printf("%s %s\n", status, url)
	}

	// List all supported providers
	fmt.Println("\nSupported oEmbed Providers:")
	providers := urlmeta.GetSupportedProviders()
	for i, provider := range providers {
		fmt.Printf("%d. %s - %s\n", i+1, provider.Name, provider.URL)
	}
	fmt.Println()
}

func demonstrateJSONExport() {
	fmt.Println("=== 4. JSON Export ===")

	url := "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
	metadata, err := urlmeta.Extract(url)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	// Export to pretty JSON
	jsonData, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		log.Printf("JSON marshal error: %v\n", err)
		return
	}

	fmt.Println("Full metadata as JSON:")
	fmt.Println(string(jsonData))
	fmt.Println()

	// In production, you might save to file:
	// err = os.WriteFile("metadata.json", jsonData, 0644)
}

func demonstrateErrorHandling() {
	fmt.Println("\n=== 5. Error Handling ===")

	testCases := []struct {
		url         string
		description string
	}{
		{"invalid-url", "Invalid URL format"},
		{"ftp://example.com", "Unsupported protocol"},
		{"https://this-domain-definitely-does-not-exist-12345.com", "Network error"},
		{"https://httpstat.us/404", "HTTP 404 error"},
	}

	for _, tc := range testCases {
		fmt.Printf("Testing: %s\n", tc.description)
		fmt.Printf("URL: %s\n", tc.url)

		_, err := urlmeta.Extract(tc.url)
		if err != nil {
			fmt.Printf("❌ Error (expected): %v\n", err)
		} else {
			fmt.Printf("✅ Success (unexpected)\n")
		}
		fmt.Println()
	}
}

// Production example: Conditional rendering based on metadata
func renderContent(metadata *urlmeta.Metadata) string {
	// If embeddable content is available, use it
	if metadata.OEmbed != nil && metadata.OEmbed.HTML != "" {
		return fmt.Sprintf(`
<div class="embed-container">
  %s
</div>`, metadata.OEmbed.HTML)
	}

	// Otherwise, create a link preview card
	imageHTML := ""
	if len(metadata.Images) > 0 {
		imageHTML = fmt.Sprintf(`<img src="%s" alt="%s">`,
			metadata.Images[0].URL, metadata.Title)
	}

	return fmt.Sprintf(`
<div class="link-preview">
  %s
  <h3>%s</h3>
  <p>%s</p>
  <span class="provider">%s</span>
</div>`, imageHTML, metadata.Title, metadata.Description, metadata.ProviderName)
}

// Production example: Batch processing with error handling
func processBatchURLs(urls []string) map[string]*urlmeta.Metadata {
	client := urlmeta.NewClient(
		urlmeta.WithTimeout(10 * time.Second),
	)

	results := make(map[string]*urlmeta.Metadata)

	for _, url := range urls {
		metadata, err := client.Extract(url)
		if err != nil {
			log.Printf("Failed to extract %s: %v", url, err)
			continue
		}
		results[url] = metadata
	}

	return results
}

// Production example: Smart content display
func displaySmartPreview(url string) {
	metadata, err := urlmeta.Extract(url)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Use oEmbed if available (richer experience)
	if metadata.OEmbed != nil {
		fmt.Println("Rendering embedded content...")
		fmt.Printf("Type: %s\n", metadata.OEmbed.Type)
		fmt.Printf("Embed: %s\n", metadata.OEmbed.HTML)
		return
	}

	// Fallback to standard preview
	fmt.Println("Rendering link preview...")
	fmt.Printf("Title: %s\n", metadata.Title)
	fmt.Printf("Description: %s\n", metadata.Description)
	if len(metadata.Images) > 0 {
		fmt.Printf("Image: %s\n", metadata.Images[0].URL)
	}
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
