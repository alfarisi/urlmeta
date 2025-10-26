## Quick Start

### Simple - One Call Gets Everything!

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/alfarisi/urlmeta"
)

func main() {
    // Single call - automatically extracts metadata AND oEmbed!
    metadata, err := urlmeta.Extract("https://www.youtube.com/watch?v=MbyvLY8CGFM")
    if err != nil {
        log.Fatal(err)
    }
    
    // Standard metadata
    fmt.Println("Title:", metadata.Title)
    fmt.Println("Description:", metadata.Description)
    fmt.Println("Provider:", metadata.ProviderName)
    
    // oEmbed data (automatically included if available!)
    if metadata.OEmbed != nil {
        fmt.Println("\n✨ Embed Code Available!")
        fmt.Println("Type:", metadata.OEmbed.Type)          // "video"
        fmt.Println("HTML:", metadata.OEmbed.HTML)          // <iframe>...
        fmt.Println("Thumbnail:", metadata.OEmbed.ThumbnailURL)
    }
}
```

That's it! No need to call separate functions or check provider support.# URLMeta

[![Go Reference](https://pkg.go.dev/badge/github.com/alfarisi/urlmeta.svg)](https://pkg.go.dev/github.com/alfarisi/urlmeta)
[![Go Report Card](https://goreportcard.com/badge/github.com/alfarisi/urlmeta)](https://goreportcard.com/report/github.com/alfarisi/urlmeta)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A powerful Go library for extracting metadata from web URLs, similar to Embedly and Iframely. Supports Open Graph Protocol, Twitter Cards, Schema.org microdata, standard HTML meta tags, and **oEmbed**.

## Features

### Metadata Extraction
- ✅ **Open Graph Protocol** - Full support for og: tags
- ✅ **Twitter Cards** - Extract Twitter card metadata
- ✅ **Schema.org** - Support for microdata
- ✅ **Standard Meta Tags** - Description, keywords, author, etc.
- ✅ **Images & Videos** - Extract media with dimensions
- ✅ **Favicon Detection** - Automatic favicon discovery
- ✅ **Canonical URL** - Get canonical/preferred URL

### oEmbed Support
- ✅ **oEmbed Protocol** - Extract embeddable content
- ✅ **Known Providers** - Built-in support for YouTube, Vimeo, Twitter, Instagram, SoundCloud, Spotify, TikTok, Flickr, and more
- ✅ **Discovery** - Automatic oEmbed endpoint discovery from HTML
- ✅ **Multiple Types** - Support for photo, video, rich, and link types

### Additional Features
- ✅ **Configurable** - Custom timeout, user-agent, HTTP client
- ✅ **Production Ready** - Error handling, redirect following, content-type checking
- ✅ **Well Tested** - Comprehensive test coverage

## Installation

```bash
go get github.com/alfarisi/urlmeta
```

## Quick Start

### Basic Metadata Extraction

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/alfarisi/urlmeta"
)

func main() {
    // Extract metadata
    metadata, err := urlmeta.Extract("https://github.com")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Title:", metadata.Title)
    fmt.Println("Description:", metadata.Description)
    fmt.Println("Provider:", metadata.ProviderName)
    fmt.Println("Images:", len(metadata.Images))
    
    // oEmbed is nil for sites that don't support it - that's fine!
    if metadata.OEmbed != nil {
        fmt.Println("Embed available:", metadata.OEmbed.Type)
    }
}
```

### Manual oEmbed Extraction (Optional)

If you need explicit control or want to disable auto-detection:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/alfarisi/urlmeta"
)

func main() {
    // Create client with auto oEmbed disabled
    client := urlmeta.NewClient(
        urlmeta.WithAutoOEmbed(false),
    )
    
    // Extract metadata only
    metadata, err := client.Extract("https://youtube.com/watch?v=123")
    // metadata.OEmbed will be nil
    
    // Manually extract oEmbed when needed
    oembed, err := client.ExtractOEmbed("https://youtube.com/watch?v=123")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Type:", oembed.Type)
    fmt.Println("HTML:", oembed.HTML)
}
```

## Advanced Usage

### Custom Configuration

```go
// Create client with custom options
client := urlmeta.NewClient(
    urlmeta.WithTimeout(15 * time.Second),
    urlmeta.WithUserAgent("MyBot/1.0"),
    urlmeta.WithMaxRedirects(5),
    urlmeta.WithAutoOEmbed(true), // Auto oEmbed enabled by default
)

// Extract metadata (includes oEmbed automatically)
metadata, err := client.Extract("https://example.com")
```

### Disable Auto oEmbed Detection

```go
// Disable if you don't need oEmbed or want faster extraction
client := urlmeta.NewClient(
    urlmeta.WithAutoOEmbed(false),
    urlmeta.WithTimeout(5 * time.Second),
)

metadata, err := client.Extract("https://youtube.com/watch?v=123")
// metadata.OEmbed will be nil

// Extract oEmbed manually only when needed
oembed, err := client.ExtractOEmbed("https://youtube.com/watch?v=123")
```

### Combined Extraction

```go
// One call gets everything!
url := "https://www.youtube.com/watch?v=MbyvLY8CGFM"

metadata, err := urlmeta.Extract(url)
if err != nil {
    log.Fatal(err)
}

// Standard metadata (OpenGraph, Twitter Cards, etc.)
fmt.Printf("Title: %s\n", metadata.Title)
fmt.Printf("Description: %s\n", metadata.Description)
fmt.Printf("Images: %d\n", len(metadata.Images))

// oEmbed data (automatically included!)
if metadata.OEmbed != nil {
    fmt.Printf("\n✨ Embeddable Content:\n")
    fmt.Printf("Type: %s\n", metadata.OEmbed.Type)
    fmt.Printf("Author: %s\n", metadata.OEmbed.AuthorName)
    fmt.Printf("Embed Code:\n%s\n", metadata.OEmbed.HTML)
}
```

## Metadata Response Structure

```go
type Metadata struct {
    // Basic Info
    Title           string
    Description     string
    URL             string
    CanonicalURL    string
    
    // Provider Info
    ProviderName    string
    ProviderURL     string
    ProviderDisplay string
    
    // Media
    Images          []Image
    Videos          []Video
    
    // OpenGraph
    Type            string
    SiteName        string
    Locale          string
    
    // Additional
    Author          string
    PublishedTime   string
    ModifiedTime    string
    Keywords        []string
    
    // Twitter
    TwitterCard     string
    TwitterSite     string
    TwitterCreator  string
    
    // Favicon
    Favicon         string
    
    // oEmbed (automatically included if URL supports it)
    OEmbed          *OEmbed
}
```

## oEmbed Response Structure

```go
type OEmbed struct {
    Type            string  // photo, video, link, rich
    Version         string  // oEmbed version
    Title           string
    AuthorName      string
    AuthorURL       string
    ProviderName    string
    ProviderURL     string
    CacheAge        int
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

URLMeta **automatically detects and extracts** oEmbed data for popular providers:

- **YouTube** - `youtube.com`, `youtu.be`
- **Vimeo** - `vimeo.com`
- **Twitter** - `twitter.com`
- **Instagram** - `instagram.com`
- **Flickr** - `flickr.com`, `flic.kr`
- **SoundCloud** - `soundcloud.com`
- **Spotify** - `open.spotify.com`
- **TikTok** - `tiktok.com`

Just call `Extract()` - no need to check support manually!

```go
// Works automatically for all supported providers
metadata, _ := urlmeta.Extract("https://vimeo.com/123456")

if metadata.OEmbed != nil {
    // oEmbed data is here!
    fmt.Println(metadata.OEmbed.HTML)
}

// For unsupported sites, OEmbed will be nil (that's okay!)
metadata, _ := urlmeta.Extract("https://github.com")
// metadata.OEmbed == nil, but metadata.Title, etc. still work
```

### Manual Provider Check (Optional)

```go
// Check if specific URL is supported (optional - not required)
if urlmeta.IsOEmbedSupported("https://vimeo.com/123456") {
    fmt.Println("oEmbed is supported!")
}

// Get list of all supported providers
providers := urlmeta.GetSupportedProviders()
for _, provider := range providers {
    fmt.Printf("%s - %s\n", provider.Name, provider.URL)
}
```

## Supported Protocols

- ✅ HTTP
- ✅ HTTPS
- ❌ FTP, File, etc. (only web protocols)

## Error Handling

```go
metadata, err := urlmeta.Extract("invalid-url")
if err != nil {
    switch {
    case strings.Contains(err.Error(), "unsupported protocol"):
        // Handle protocol error
    case strings.Contains(err.Error(), "HTTP error"):
        // Handle HTTP error
    case strings.Contains(err.Error(), "timeout"):
        // Handle timeout
    default:
        // Handle other errors
    }
}
```

## Examples

See [examples/](./examples/) directory for more usage examples:
- [Basic Usage](./examples/basic/main.go) - Simple metadata extraction
- [Advanced Options](./examples/advanced/main.go) - Custom configuration
- [Batch Processing](./examples/batch/main.go) - Concurrent extraction

## Performance

Tips for better performance:
- Reuse `Client` instances
- Use appropriate timeouts
- Process multiple URLs concurrently
- Cache results when possible

## Comparison with Similar Libraries

| Feature | URLMeta | otiai10/opengraph | dyatlov/go-opengraph | Iframely/Embedly |
|---------|---------|-------------------|----------------------|------------------|
| Open Graph | ✅ | ✅ | ✅ | ✅ |
| Twitter Cards | ✅ | ❌ | ❌ | ✅ |
| oEmbed | ✅ **Auto** | ❌ | ❌ | ✅ |
| Schema.org | ✅ | ❌ | ❌ | ✅ |
| Videos | ✅ | ❌ | ❌ | ✅ |
| Favicon | ✅ | ❌ | ❌ | ✅ |
| Custom HTTP Client | ✅ | ❌ | ❌ | N/A |
| **Auto Detection** | ✅ **One Call** | N/A | N/A | ✅ |
| Price | **Free** | Free | Free | $9-999/mo |

**Key Advantage:** URLMeta extracts metadata AND oEmbed in a **single call** - no need to check provider support or make separate requests!

## Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](docs/CONTRIBUTING.md) for details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

Inspired by:
- [Embedly](https://embed.ly/) - URL metadata and embed service
- [Iframely](https://iframely.com/) - oEmbed and URL meta service
- [Open Graph Protocol](https://ogp.me/) - Metadata standard
- [oEmbed](https://oembed.com/) - Embed format standard

## Similar Projects

- [otiai10/opengraph](https://github.com/otiai10/opengraph) - OpenGraph parser
- [dyatlov/go-opengraph](https://github.com/dyatlov/go-opengraph) - OpenGraph library
- [erizocosmico/pagecard](https://github.com/erizocosmico/pagecard) - Metadata extractor

## FAQ

### Does it execute JavaScript?

No, URLMeta parses static HTML only. For JavaScript-heavy sites, consider using a headless browser.

### Can I use custom HTTP client?

Yes! Use `WithHTTPClient()` option to provide your own `http.Client`.

### How do I handle rate limiting?

Implement your own rate limiting using channels or libraries like `golang.org/x/time/rate`.

### Does it support authentication?

Yes, configure your HTTP client with authentication headers as needed.

### What about CORS?

This is a server-side library, CORS doesn't apply. However, some APIs (like Instagram oEmbed) may require authentication.

## Roadmap

- [ ] Support for more oEmbed providers
- [ ] JSON-LD extraction
- [ ] Microformats support
- [ ] AMP metadata
- [ ] RSS/Atom feed detection
- [ ] Language detection
- [ ] Content readability scoring

## Support

- 📫 Open an issue for bug reports or feature requests
- 💬 Discussions for questions and ideas
- ⭐ Star the repo if you find it useful!

## License

CC0 1.0 Universal — this project includes AI-assisted code and is released to the public domain.