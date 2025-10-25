# API Documentation

## Table of Contents

- [Quick Start](#quick-start)
- [Types](#types)
- [Functions](#functions)
- [Options](#options)
- [Error Handling](#error-handling)
- [Examples](#examples)

## Quick Start

```go
import "github.com/alfarisi/urlmeta"

// Simple usage
metadata, err := urlmeta.Extract("https://example.com")
if err != nil {
    log.Fatal(err)
}
fmt.Println(metadata.Title)
```

## Types

### Metadata

The main struct containing all extracted metadata from a web page.

```go
type Metadata struct {
    // Basic Information
    Title           string   `json:"title"`
    Description     string   `json:"description"`
    URL             string   `json:"url"`
    CanonicalURL    string   `json:"canonical_url,omitempty"`
    
    // Provider Information
    ProviderName    string   `json:"provider_name"`
    ProviderURL     string   `json:"provider_url"`
    ProviderDisplay string   `json:"provider_display"`
    
    // Media
    Images          []Image  `json:"images,omitempty"`
    Videos          []Video  `json:"videos,omitempty"`
    
    // OpenGraph Protocol
    Type            string   `json:"type,omitempty"`
    SiteName        string   `json:"site_name,omitempty"`
    Locale          string   `json:"locale,omitempty"`
    
    // Additional Metadata
    Author          string   `json:"author,omitempty"`
    PublishedTime   string   `json:"published_time,omitempty"`
    ModifiedTime    string   `json:"modified_time,omitempty"`
    Keywords        []string `json:"keywords,omitempty"`
    
    // Twitter Card
    TwitterCard     string   `json:"twitter_card,omitempty"`
    TwitterSite     string   `json:"twitter_site,omitempty"`
    TwitterCreator  string   `json:"twitter_creator,omitempty"`
    
    // Additional
    Favicon         string   `json:"favicon,omitempty"`
}
```

**Fields:**

- `Title`: Page title (from og:title, twitter:title, or `<title>` tag)
- `Description`: Page description (from og:description, twitter:description, or meta description)
- `URL`: Final URL after redirects
- `CanonicalURL`: Canonical URL if specified
- `ProviderName`: Website/provider name
- `ProviderURL`: Base URL of the provider
- `ProviderDisplay`: Human-readable provider name (hostname or site_name)
- `Images`: Array of images found on the page
- `Videos`: Array of videos found on the page
- `Type`: OpenGraph type (article, website, etc.)
- `SiteName`: Site name from OpenGraph
- `Locale`: Content locale (e.g., "en_US")
- `Author`: Content author
- `PublishedTime`: Publication timestamp
- `ModifiedTime`: Last modification timestamp
- `Keywords`: Array of keywords/tags
- `TwitterCard`: Twitter card type
- `TwitterSite`: Twitter site account
- `TwitterCreator`: Twitter author account
- `Favicon`: URL to the site's favicon

### Image

Represents an image extracted from the page.

```go
type Image struct {
    URL    string `json:"url"`
    Width  int    `json:"width,omitempty"`
    Height int    `json:"height,omitempty"`
    Alt    string `json:"alt,omitempty"`
}
```

### Video

Represents a video extracted from the page.

```go
type Video struct {
    URL    string `json:"url"`
    Type   string `json:"type,omitempty"`
    Width  int    `json:"width,omitempty"`
    Height int    `json:"height,omitempty"`
}
```

### Client

The client for extracting metadata with custom configuration.

```go
type Client struct {
    // unexported fields
}
```

## Functions

### Extract

```go
func Extract(targetURL string) (*Metadata, error)
```

Extract metadata from a URL using default settings.

**Parameters:**
- `targetURL`: The URL to extract metadata from

**Returns:**
- `*Metadata`: Extracted metadata
- `error`: Error if extraction fails

**Example:**
```go
metadata, err := urlmeta.Extract("https://github.com")
if err != nil {
    log.Fatal(err)
}
fmt.Println(metadata.Title)
```

### NewClient

```go
func NewClient(opts ...Option) *Client
```

Create a new client with custom options.

**Parameters:**
- `opts`: Variable number of Option functions

**Returns:**
- `*Client`: Configured client instance

**Example:**
```go
client := urlmeta.NewClient(
    urlmeta.WithTimeout(15 * time.Second),
    urlmeta.WithUserAgent("MyBot/1.0"),
)
```

### Client.Extract

```go
func (c *Client) Extract(targetURL string) (*Metadata, error)
```

Extract metadata using the configured client.

**Parameters:**
- `targetURL`: The URL to extract metadata from

**Returns:**
- `*Metadata`: Extracted metadata
- `error`: Error if extraction fails

**Example:**
```go
client := urlmeta.NewClient()
metadata, err := client.Extract("https://example.com")
```

## Options

### WithTimeout

```go
func WithTimeout(timeout time.Duration) Option
```

Set custom timeout for HTTP requests (default: 10 seconds).

**Example:**
```go
client := urlmeta.NewClient(
    urlmeta.WithTimeout(30 * time.Second),
)
```

### WithUserAgent

```go
func WithUserAgent(ua string) Option
```

Set custom User-Agent header.

**Example:**
```go
client := urlmeta.NewClient(
    urlmeta.WithUserAgent("MyBot/1.0 (+https://mywebsite.com)"),
)
```

### WithMaxRedirects

```go
func WithMaxRedirects(max int) Option
```

Set maximum number of redirects to follow (default: 10).

**Example:**
```go
client := urlmeta.NewClient(
    urlmeta.WithMaxRedirects(5),
)
```

### WithHTTPClient

```go
func WithHTTPClient(client *http.Client) Option
```

Use a custom HTTP client.

**Example:**
```go
customClient := &http.Client{
    Transport: &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    },
}

client := urlmeta.NewClient(
    urlmeta.WithHTTPClient(customClient),
)
```

## Error Handling

The library returns errors in the following cases:

### Invalid URL
```go
_, err := urlmeta.Extract("not-a-valid-url")
// err: "invalid URL: ..."
```

### Unsupported Protocol
```go
_, err := urlmeta.Extract("ftp://example.com")
// err: "unsupported protocol: ftp (only http and https are supported)"
```

### HTTP Errors
```go
_, err := urlmeta.Extract("https://example.com/404")
// err: "HTTP error: 404 Not Found"
```

### Network Errors
```go
_, err := urlmeta.Extract("https://nonexistent-domain-12345.com")
// err: "failed to fetch URL: ..."
```

### Timeout
```go
client := urlmeta.NewClient(
    urlmeta.WithTimeout(1 * time.Millisecond),
)
_, err := client.Extract("https://slow-website.com")
// err: context deadline exceeded
```

### Unsupported Content Type
```go
_, err := urlmeta.Extract("https://example.com/data.json")
// err: "unsupported content type: application/json"
```

## Examples

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/alfarisi/urlmeta"
)

func main() {
    metadata, err := urlmeta.Extract("https://github.com")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Title:", metadata.Title)
    fmt.Println("Description:", metadata.Description)
}
```

### Custom Configuration

```go
package main

import (
    "fmt"
    "log"
    "time"
    
    "github.com/alfarisi/urlmeta"
)

func main() {
    client := urlmeta.NewClient(
        urlmeta.WithTimeout(15 * time.Second),
        urlmeta.WithUserAgent("MyBot/1.0"),
        urlmeta.WithMaxRedirects(5),
    )
    
    metadata, err := client.Extract("https://example.com")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("%+v\n", metadata)
}
```

### Extracting Images

```go
metadata, err := urlmeta.Extract("https://example.com")
if err != nil {
    log.Fatal(err)
}

for _, img := range metadata.Images {
    fmt.Printf("Image: %s\n", img.URL)
    if img.Width > 0 && img.Height > 0 {
        fmt.Printf("  Dimensions: %dx%d\n", img.Width, img.Height)
    }
}
```

### Extracting Videos

```go
metadata, err := urlmeta.Extract("https://example.com")
if err != nil {
    log.Fatal(err)
}

for _, video := range metadata.Videos {
    fmt.Printf("Video: %s\n", video.URL)
    if video.Type != "" {
        fmt.Printf("  Type: %s\n", video.Type)
    }
}
```

### Error Handling with Retry

```go
func extractWithRetry(url string, maxRetries int) (*urlmeta.Metadata, error) {
    client := urlmeta.NewClient(
        urlmeta.WithTimeout(10 * time.Second),
    )
    
    var lastErr error
    for i := 0; i <= maxRetries; i++ {
        metadata, err := client.Extract(url)
        if err == nil {
            return metadata, nil
        }
        
        lastErr = err
        if i < maxRetries {
            time.Sleep(time.Second * time.Duration(i+1))
        }
    }
    
    return nil, fmt.Errorf("failed after %d retries: %w", maxRetries, lastErr)
}
```

### Batch Processing

```go
func processBatch(urls []string) map[string]*urlmeta.Metadata {
    client := urlmeta.NewClient()
    results := make(map[string]*urlmeta.Metadata)
    
    for _, url := range urls {
        metadata, err := client.Extract(url)
        if err != nil {
            log.Printf("Error extracting %s: %v", url, err)
            continue
        }
        results[url] = metadata
    }
    
    return results
}
```

### Concurrent Extraction

```go
func extractConcurrent(urls []string) []*urlmeta.Metadata {
    client := urlmeta.NewClient()
    
    type result struct {
        metadata *urlmeta.Metadata
        err      error
    }
    
    resultChan := make(chan result, len(urls))
    
    for _, url := range urls {
        go func(u string) {
            metadata, err := client.Extract(u)
            resultChan <- result{metadata, err}
        }(url)
    }
    
    var results []*urlmeta.Metadata
    for i := 0; i < len(urls); i++ {
        r := <-resultChan
        if r.err == nil {
            results = append(results, r.metadata)
        }
    }
    
    return results
}
```

### JSON Export

```go
import "encoding/json"

metadata, err := urlmeta.Extract("https://example.com")
if err != nil {
    log.Fatal(err)
}

jsonData, err := json.MarshalIndent(metadata, "", "  ")
if err != nil {
    log.Fatal(err)
}

fmt.Println(string(jsonData))
```

### Checking Specific Fields

```go
metadata, err := urlmeta.Extract("https://example.com")
if err != nil {
    log.Fatal(err)
}

// Check if it's an article
if metadata.Type == "article" {
    fmt.Println("This is an article")
    fmt.Println("Author:", metadata.Author)
    fmt.Println("Published:", metadata.PublishedTime)
}

// Check for Twitter card
if metadata.TwitterCard != "" {
    fmt.Printf("Twitter Card: %s\n", metadata.TwitterCard)
}

// Check for images
if len(metadata.Images) > 0 {
    fmt.Printf("Found %d images\n", len(metadata.Images))
}
```

### URL Normalization

The library automatically normalizes URLs:

```go
// All of these work:
urlmeta.Extract("https://example.com")
urlmeta.Extract("http://example.com")
urlmeta.Extract("example.com")  // Automatically adds https://
```

### Custom HTTP Client

```go
import (
    "crypto/tls"
    "net/http"
)

customHTTPClient := &http.Client{
    Transport: &http.Transport{
        TLSClientConfig: &tls.Config{
            InsecureSkipVerify: true, // Only for testing!
        },
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
    },
    Timeout: 30 * time.Second,
}

client := urlmeta.NewClient(
    urlmeta.WithHTTPClient(customHTTPClient),
)

metadata, err := client.Extract("https://example.com")
```

## Performance Tips

### Reuse Client

```go
// ✅ Good - Reuse client
client := urlmeta.NewClient()
for _, url := range urls {
    metadata, _ := client.Extract(url)
}

// ❌ Bad - Creates new client each time
for _, url := range urls {
    metadata, _ := urlmeta.Extract(url)
}
```

### Use Appropriate Timeout

```go
// For fast responses
client := urlmeta.NewClient(
    urlmeta.WithTimeout(5 * time.Second),
)

// For slow websites
client := urlmeta.NewClient(
    urlmeta.WithTimeout(30 * time.Second),
)
```

### Concurrent Processing

For multiple URLs, use goroutines:

```go
// Process 10 URLs with 3 workers
semaphore := make(chan struct{}, 3)
var wg sync.WaitGroup

for _, url := range urls {
    wg.Add(1)
    go func(u string) {
        defer wg.Done()
        semaphore <- struct{}{}
        defer func() { <-semaphore }()
        
        metadata, _ := client.Extract(u)
        // Process metadata...
    }(url)
}

wg.Wait()
```

## Supported Meta Tags

### OpenGraph Protocol

- `og:title`
- `og:description`
- `og:image`
- `og:image:url`
- `og:image:width`
- `og:image:height`
- `og:video`
- `og:video:url`
- `og:video:type`
- `og:type`
- `og:url`
- `og:site_name`
- `og:locale`
- `article:published_time`
- `article:modified_time`
- `article:author`

### Twitter Cards

- `twitter:card`
- `twitter:site`
- `twitter:creator`
- `twitter:title`
- `twitter:description`
- `twitter:image`
- `twitter:image:src`

### Standard HTML Meta Tags

- `<title>`
- `name="description"`
- `name="author"`
- `name="keywords"`
- `<link rel="icon">`
- `<link rel="canonical">`

### Schema.org Microdata

- `itemprop="name"`
- `itemprop="description"`
- `itemprop="image"`

## Limitations

- **Content Types**: Only HTML/XHTML content is supported
- **Protocols**: Only HTTP and HTTPS are supported
- **Body Size**: Limited to 10MB to prevent memory issues
- **JavaScript**: Does not execute JavaScript (uses static HTML only)
- **Dynamic Content**: Cannot extract content loaded via AJAX/JavaScript

## Best Practices

1. **Always handle errors** - Network requests can fail
2. **Set appropriate timeouts** - Don't wait forever
3. **Reuse clients** - Better performance
4. **Use custom User-Agent** - Identify your bot properly
5. **Respect robots.txt** - Be a good web citizen
6. **Rate limit requests** - Don't overwhelm servers
7. **Cache results** - Avoid unnecessary requests

## Troubleshooting

### "unsupported protocol" error
- Use `http://` or `https://` protocol
- FTP, file://, and other protocols are not supported

### "HTTP error: 403" error
- The site may be blocking the default User-Agent
- Try setting a custom User-Agent

### "timeout" error
- Increase the timeout value
- Check your network connection
- The target site may be slow

### Empty metadata
- The site may not have proper meta tags
- Check if the site requires JavaScript to render
- Verify the HTML structure

### Relative URLs in images/favicon
- The library automatically resolves relative URLs
- If URLs are incorrect, it may be a bug - please report it

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md) for guidelines on contributing to this project.

## License

MIT License - see [LICENSE](../LICENSE) for details.