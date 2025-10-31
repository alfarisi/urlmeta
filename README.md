# URLMeta

[![Go Reference](https://pkg.go.dev/badge/github.com/alfarisi/urlmeta.svg)](https://pkg.go.dev/github.com/alfarisi/urlmeta)
[![Go Report Card](https://goreportcard.com/badge/github.com/alfarisi/urlmeta)](https://goreportcard.com/report/github.com/alfarisi/urlmeta)
[![License: 0BSD](https://img.shields.io/badge/License-0BSD-yellow.svg)](https://opensource.org/license/0bsd)

A powerful Go library for extracting metadata from web URLs, similar to Embedly and Iframely. Supports Open Graph Protocol, Twitter Cards, standard HTML meta tags, and **oEmbed**.

**Key Feature:** Extract metadata AND oEmbed in a single call—no need to check provider support or make separate requests!

## Features

- ✅ **Open Graph Protocol** - Full support for og: tags
- ✅ **Twitter Cards** - Extract Twitter card metadata
- ✅ **Standard Meta Tags** - Description, keywords, author, etc.
- ✅ **oEmbed Protocol** - Automatic extraction for YouTube, Vimeo, Twitter, Instagram, SoundCloud, Spotify, TikTok, Flickr
- ✅ **Images & Videos** - Extract media with dimensions
- ✅ **Favicon & Canonical URL** - Automatic discovery
- ✅ **Configurable** - Custom timeout, user-agent, HTTP client
- ✅ **Production Ready** - Error handling, redirect following, comprehensive tests

## Installation

```bash
go get github.com/alfarisi/urlmeta
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/alfarisi/urlmeta"
)

func main() {
    // Single call extracts metadata AND oEmbed automatically!
    metadata, err := urlmeta.Extract("https://www.youtube.com/watch?v=123")
    if err != nil {
        log.Fatal(err)
    }
    
    // Standard metadata (works for ALL websites)
    fmt.Println("Title:", metadata.Title)
    fmt.Println("Description:", metadata.Description)
    fmt.Println("Provider:", metadata.ProviderName)
    fmt.Println("Images:", len(metadata.Images))
    
    // oEmbed data (automatically included if available)
    if metadata.OEmbed != nil {
        fmt.Println("\n✨ Embeddable Content:")
        fmt.Println("Type:", metadata.OEmbed.Type)
        fmt.Println("HTML:", metadata.OEmbed.HTML)
        fmt.Println("Thumbnail:", metadata.OEmbed.ThumbnailURL)
    }
}
```

## Usage Examples

### Custom Configuration

```go
client := urlmeta.NewClient(
    urlmeta.WithTimeout(15 * time.Second),
    urlmeta.WithUserAgent("MyBot/1.0"),
    urlmeta.WithMaxRedirects(5),
)

metadata, err := client.Extract("https://example.com")
```

### Disable Auto oEmbed (Faster for Non-Embed Sites)

```go
// Skip oEmbed detection for better performance
client := urlmeta.NewClient(
    urlmeta.WithAutoOEmbed(false),
)

metadata, err := client.Extract("https://github.com")
// metadata.OEmbed will be nil

// Extract oEmbed manually when needed
oembed, err := client.ExtractOEmbed("https://youtube.com/watch?v=123")
```

### Batch Processing

```go
urls := []string{
    "https://youtube.com/watch?v=123",
    "https://vimeo.com/456",
    "https://github.com",
}

client := urlmeta.NewClient()

for _, url := range urls {
    metadata, err := client.Extract(url)
    if err != nil {
        log.Printf("Error extracting %s: %v", url, err)
        continue
    }
    
    fmt.Printf("%s: %s\n", url, metadata.Title)
}
```

## Response Structure

### Metadata

```go
type Metadata struct {
    // Basic
    Title           string
    Description     string
    URL             string
    CanonicalURL    string
    
    // Provider
    ProviderName    string
    ProviderURL     string
    
    // Media
    Images          []Image
    Videos          []Video
    Favicon         string
    
    // OpenGraph
    Type            string
    SiteName        string
    Locale          string
    
    // Meta
    Author          string
    PublishedTime   string
    ModifiedTime    string
    Keywords        []string
    
    // Twitter
    TwitterCard     string
    TwitterSite     string
    TwitterCreator  string
    
    // oEmbed (auto-included if supported)
    OEmbed          *OEmbed
}
```

### OEmbed

```go
type OEmbed struct {
    Type            string  // photo, video, link, rich
    Version         string
    Title           string
    AuthorName      string
    AuthorURL       string
    ProviderName    string
    ProviderURL     string
    ThumbnailURL    string
    ThumbnailWidth  int
    ThumbnailHeight int
    
    // Photo type
    URL             string
    Width           int
    Height          int
    
    // Video/Rich type
    HTML            string  // Embed code
}
```

## Supported oEmbed Providers

URLMeta automatically extracts oEmbed for:

| Provider | Domains |
|----------|---------|
| YouTube | `youtube.com`, `youtu.be` |
| Vimeo | `vimeo.com` |
| Twitter/X | `twitter.com`, `x.com` |
| Instagram | `instagram.com` |
| Flickr | `flickr.com`, `flic.kr` |
| SoundCloud | `soundcloud.com` |
| Spotify | `open.spotify.com` |
| TikTok | `tiktok.com` |

**Note:** For unsupported sites, `metadata.OEmbed` will be `nil`, but standard metadata extraction still works!

### Check Provider Support

```go
// Check if URL supports oEmbed (optional)
if urlmeta.IsOEmbedSupported("https://vimeo.com/123") {
    fmt.Println("oEmbed supported!")
}

// List all supported providers
providers := urlmeta.GetSupportedProviders()
```

## Error Handling

```go
metadata, err := urlmeta.Extract("https://example.com")
if err != nil {
    switch {
    case strings.Contains(err.Error(), "unsupported protocol"):
        // Only HTTP/HTTPS supported
    case strings.Contains(err.Error(), "HTTP error"):
        // Server returned error (404, 500, etc)
    case strings.Contains(err.Error(), "timeout"):
        // Request timed out
    case strings.Contains(err.Error(), "unsupported content type"):
        // Not an HTML page
    default:
        log.Printf("Extraction failed: %v", err)
    }
}
```

## Performance Tips

1. **Reuse Client** - Create once, use many times
2. **Set Appropriate Timeouts** - Balance speed vs reliability
3. **Concurrent Processing** - Use goroutines for batch extraction
4. **Cache Results** - Store metadata to avoid repeated requests
5. **Disable Auto-oEmbed** - Skip oEmbed detection for non-embed sites

```go
// Good: Reuse client
client := urlmeta.NewClient(urlmeta.WithTimeout(10 * time.Second))
for _, url := range urls {
    metadata, _ := client.Extract(url)
}

// Bad: Create new client each time
for _, url := range urls {
    metadata, _ := urlmeta.Extract(url) // Creates new client internally
}
```

## Cache Behavior

URLMeta uses an internal regex cache for performance optimization. This is **safe and recommended** for most use cases.

## Examples

Complete examples available in [examples/](./examples/):
- [basic/](./examples/basic/main.go) - Simple metadata extraction
- [advanced/](./examples/advanced/main.go) - Custom configuration
- [batch/](./examples/batch/main.go) - Concurrent processing

## FAQ

**Q: Does it execute JavaScript?**  
A: No, URLMeta parses static HTML only. For JavaScript-heavy sites, use a headless browser.

**Q: Can I use a custom HTTP client?**  
A: Yes! Use `WithHTTPClient()` option.

**Q: How do I handle rate limiting?**  
A: Implement rate limiting using `golang.org/x/time/rate` or similar libraries.

**Q: Does it support authentication?**  
A: Yes, configure your HTTP client with auth headers.

**Q: Why is oEmbed nil for some URLs?**  
A: oEmbed is only available for supported providers (YouTube, Vimeo, etc). Standard metadata still works for all sites.

**Q: How do I add custom oEmbed providers?**  
A: Use `urlmeta.AddCustomProvider()` to register your own provider at runtime.

## Supported Protocols

- ✅ HTTP, HTTPS
- ❌ FTP, File, etc.

## Contributing

Contributions welcome! Please see [CONTRIBUTING.md](docs/CONTRIBUTING.md).

## Acknowledgments

Inspired by [Embedly](https://embed.ly/), [Iframely](https://iframely.com/), [Open Graph Protocol](https://ogp.me/), and [oEmbed](https://oembed.com/).

## Similar Projects

- [otiai10/opengraph](https://github.com/otiai10/opengraph) - OpenGraph parser
- [dyatlov/go-opengraph](https://github.com/dyatlov/go-opengraph) - OpenGraph library
- [erizocosmico/pagecard](https://github.com/erizocosmico/pagecard) - Metadata extractor

## License

Zero-Clause BSD (0BSD) - Functionally equivalent to public domain. See [LICENSE](LICENSE) file.

---

⭐ **Star this repo if you find it useful!**  
📫 **Open an issue** for bugs or feature requests  
💬 **Start a discussion** for questions and ideas