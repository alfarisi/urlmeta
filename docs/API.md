# API Documentation

## Table of Contents

- [Quick Start](#quick-start)
- [Types](#types)
- [Functions](#functions)
- [Options](#options)
- [Extraction Strategies](#extraction-strategies)
- [Provider Management](#provider-management)
- [Error Handling](#error-handling)
- [Performance](#performance)
- [Examples](#examples)

## Quick Start

```go
import "github.com/alfarisi/urlmeta"

// Simple usage - auto-optimized!
metadata, err := urlmeta.Extract("https://youtube.com/watch?v=...")
if err != nil {
    log.Fatal(err)
}

// Standard metadata (always available)
fmt.Println(metadata.Title)
fmt.Println(metadata.Description)

// oEmbed data (auto-included if available)
if metadata.OEmbed != nil {
    fmt.Println(metadata.OEmbed.HTML) // Embed code
}
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
    
    // oEmbed (automatically included when available)
    OEmbed          *OEmbed  `json:"oembed,omitempty"`
}
```

### OEmbed

oEmbed response structure following the [oEmbed specification](https://oembed.com/).

```go
type OEmbed struct {
    Type            string `json:"type"`               // "photo", "video", "rich", "link"
    Version         string `json:"version"`            // oEmbed version (usually "1.0")
    Title           string `json:"title,omitempty"`
    AuthorName      string `json:"author_name,omitempty"`
    AuthorURL       string `json:"author_url,omitempty"`
    ProviderName    string `json:"provider_name,omitempty"`
    ProviderURL     string `json:"provider_url,omitempty"`
    CacheAge        int    `json:"cache_age,omitempty"`
    ThumbnailURL    string `json:"thumbnail_url,omitempty"`
    ThumbnailWidth  int    `json:"thumbnail_width,omitempty"`
    ThumbnailHeight int    `json:"thumbnail_height,omitempty"`
    
    // Photo type specific
    URL             string `json:"url,omitempty"`
    Width           int    `json:"width,omitempty"`
    Height          int    `json:"height,omitempty"`
    
    // Video/Rich type specific
    HTML            string `json:"html,omitempty"`     // Embed code
}
```

### Image

```go
type Image struct {
    URL    string `json:"url"`
    Width  int    `json:"width,omitempty"`
    Height int    `json:"height,omitempty"`
    Alt    string `json:"alt,omitempty"`
}
```

### Video

```go
type Video struct {
    URL    string `json:"url"`
    Type   string `json:"type,omitempty"`
    Width  int    `json:"width,omitempty"`
    Height int    `json:"height,omitempty"`
}
```

### ExtractionStrategy

```go
type ExtractionStrategy int

const (
    StrategyAuto        ExtractionStrategy = iota  // Automatically chooses best strategy
    StrategyOEmbedFirst                            // Try oEmbed first, fall back to HTML
    StrategyHTMLOnly                               // Only extract from HTML
)
```

## Functions

### Extract

```go
func Extract(targetURL string) (*Metadata, error)
```

Extract metadata from a URL using default settings and automatic strategy selection.

**Performance:**
- YouTube/Vimeo: 1 HTTP call (oEmbed only)
- GitHub/Blogs: 1 HTTP call (HTML only)
- **Optimized automatically!**

**Parameters:**
- `targetURL`: The URL to extract metadata from (with or without `https://`)

**Returns:**
- `*Metadata`: Complete metadata including oEmbed if available
- `error`: Error if extraction fails

**Example:**
```go
metadata, err := urlmeta.Extract("youtube.com/watch?v=123")
if err != nil {
    log.Fatal(err)
}
fmt.Println(metadata.Title)
if metadata.OEmbed != nil {
    fmt.Println("Embeddable:", metadata.OEmbed.Type)
}
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
    urlmeta.WithStrategy(urlmeta.StrategyOEmbedFirst),
)
```

### Client.Extract

```go
func (c *Client) Extract(targetURL string) (*Metadata, error)
```

Extract metadata using the configured client.

**Example:**
```go
client := urlmeta.NewClient(
    urlmeta.WithTimeout(5 * time.Second),
)
metadata, err := client.Extract("https://example.com")
```

### Client.ExtractOEmbed

```go
func (c *Client) ExtractOEmbed(targetURL string) (*OEmbed, error)
```

Explicitly extract only oEmbed data (bypasses automatic strategy).

**Use when:**
- You only need embed code
- Testing oEmbed endpoints
- Custom workflows

**Example:**
```go
oembed, err := client.ExtractOEmbed("https://youtube.com/watch?v=123")
if err != nil {
    log.Fatal(err)
}
fmt.Println(oembed.HTML)
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

**Best Practice:** Identify your bot with contact info:
```go
client := urlmeta.NewClient(
    urlmeta.WithUserAgent("MyBot/1.0 (+https://mywebsite.com/bot)"),
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

Use a custom HTTP client (for proxies, custom TLS, etc.).

**Example:**
```go
customClient := &http.Client{
    Transport: &http.Transport{
        Proxy: http.ProxyFromEnvironment,
    },
}

client := urlmeta.NewClient(
    urlmeta.WithHTTPClient(customClient),
)
```

### WithAutoOEmbed

```go
func WithAutoOEmbed(auto bool) Option
```

Enable/disable automatic oEmbed extraction (default: true).

**When to disable:**
- You don't need embed functionality
- Want faster extraction for non-embeddable sites
- Custom oEmbed handling

**Example:**
```go
// Disable auto oEmbed
client := urlmeta.NewClient(
    urlmeta.WithAutoOEmbed(false),
)
```

### WithStrategy

```go
func WithStrategy(strategy ExtractionStrategy) Option
```

Set extraction strategy (default: StrategyAuto).

**Strategies:**
- `StrategyAuto`: Smart selection (recommended)
- `StrategyOEmbedFirst`: Always try oEmbed first
- `StrategyHTMLOnly`: Skip oEmbed completely

**Example:**
```go
// Force HTML-only extraction (fastest for blogs)
client := urlmeta.NewClient(
    urlmeta.WithStrategy(urlmeta.StrategyHTMLOnly),
)
```

## Extraction Strategies

URLMeta uses intelligent strategies to minimize HTTP requests.

### StrategyAuto (Default, Recommended)
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

### StrategyAuto (Default, Recommended)

Automatically chooses the best strategy:
- If oEmbed supported (YouTube, Vimeo, etc.) → Use `StrategyOEmbedFirst`
- Otherwise → Use `StrategyHTMLOnly`

**Result: Always 1 HTTP call, optimal for all sites!**

### StrategyOEmbedFirst

1. Check oEmbed support (pattern match, ~0ms)
2. Fetch oEmbed data (1 HTTP call)
3. Build metadata from oEmbed response
4. **Skip HTML fetching** (saves bandwidth & time!)

**Best for:** Known embeddable content (YouTube, Vimeo, Twitter)

### StrategyHTMLOnly

1. Fetch HTML (1 HTTP call)
2. Parse metadata
3. Skip oEmbed completely

**Best for:** Blogs, news sites, documentation

## Provider Management

### IsOEmbedSupported

```go
func IsOEmbedSupported(targetURL string) bool
```

Check if a URL supports oEmbed.

**Example:**
```go
if urlmeta.IsOEmbedSupported("https://youtube.com/watch?v=123") {
    fmt.Println("oEmbed available!")
}
```

### GetSupportedProviders

```go
func GetSupportedProviders() []OEmbedProvider
```

Get list of all supported oEmbed providers.

**Example:**
```go
providers := urlmeta.GetSupportedProviders()
for _, p := range providers {
    fmt.Printf("%s - %s\n", p.Name, p.URL)
}
```

### GetKnownProviders

```go
func GetKnownProviders() []OEmbedProvider
```

Get a copy of known providers (same as `GetSupportedProviders`).

### AddCustomProvider

```go
func AddCustomProvider(provider OEmbedProvider)
```

Add custom oEmbed provider at runtime.

**Use for:**
- Private/internal video services
- New providers not yet in the list
- Testing custom endpoints

**Example:**
```go
custom := urlmeta.OEmbedProvider{
    Name: "MyVideos",
    URL:  "https://videos.mycompany.com",
    Endpoints: []urlmeta.OEmbedEndpoint{
        {
            Schemes: []string{"https://videos.mycompany.com/watch/*"},
            URL:     "https://videos.mycompany.com/oembed",
        },
    },
}

urlmeta.AddCustomProvider(custom)

// Now works!
metadata, _ := urlmeta.Extract("https://videos.mycompany.com/watch/123")
```

### ProviderCount

```go
func ProviderCount() int
```

Get the number of supported providers.

### IsProviderSupported

```go
func IsProviderSupported(providerName string) bool
```

Check if a specific provider is supported.

**Example:**
```go
if urlmeta.IsProviderSupported("YouTube") {
    fmt.Println("YouTube is supported")
}
```

### GetProviderByName

```go
func GetProviderByName(name string) *OEmbedProvider
```

Get provider details by name.

**Example:**
```go
youtube := urlmeta.GetProviderByName("YouTube")
if youtube != nil {
    fmt.Printf("Endpoints: %d\n", len(youtube.Endpoints))
}
```

## Error Handling

The library returns descriptive errors:

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

## Performance

### Best Practices

#### 1. Reuse Client
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

#### 2. Use Appropriate Timeout
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

#### 3. Concurrent Processing
```go
// Process multiple URLs concurrently
client := urlmeta.NewClient()
var wg sync.WaitGroup

for _, url := range urls {
    wg.Add(1)
    go func(u string) {
        defer wg.Done()
        metadata, _ := client.Extract(u)
        // Process metadata...
    }(url)
}

wg.Wait()
```

### Performance Metrics

| Metric | YouTube (oEmbed) | GitHub (HTML) |
|--------|-----------------|---------------|
| HTTP Calls | 1 | 1 |
| Time | ~100ms | ~150ms |
| Bandwidth | ~2KB | ~25KB |
| Memory | ~4KB | ~30KB |

## Examples

See complete examples in the [examples/](../examples/) directory.

## Supported Meta Tags

### OpenGraph Protocol
- `og:title`, `og:description`, `og:image`, `og:video`
- `og:site_name`, `og:type`, `og:url`, `og:locale`
- `article:published_time`, `article:modified_time`, `article:author`

### Twitter Cards
- `twitter:card`, `twitter:site`, `twitter:creator`
- `twitter:title`, `twitter:description`, `twitter:image`

### Standard HTML
- `<title>`, `name="description"`, `name="author"`, `name="keywords"`
- `<link rel="icon">`, `<link rel="canonical">`

### Schema.org
- `itemprop="name"`, `itemprop="description"`, `itemprop="image"`

## Limitations

- **Content Types**: Only HTML/XHTML supported
- **Protocols**: HTTP and HTTPS only
- **Body Size**: Limited to 10MB
- **JavaScript**: Not executed (static HTML only)
- **Dynamic Content**: Cannot extract AJAX-loaded content

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
Use `http://` or `https://` protocol. FTP and other protocols are not supported.

### "HTTP error: 403" error
The site may be blocking the default User-Agent. Try setting a custom one:
```go
client := urlmeta.NewClient(
    urlmeta.WithUserAgent("MyBot/1.0 (+https://mywebsite.com)"),
)
```

### "timeout" error
Increase the timeout value:
```go
client := urlmeta.NewClient(
    urlmeta.WithTimeout(30 * time.Second),
)
```

### Empty metadata
- The site may not have proper meta tags
- Check if the site requires JavaScript
- Verify the HTML structure

### No oEmbed data
- Check if provider is supported: `urlmeta.IsOEmbedSupported(url)`
- Try manual extraction: `client.ExtractOEmbed(url)`
- Provider might have rate limits

## See Also

- [README.md](../README.md) - Overview and quick start
- [CONTRIBUTING.md](./CONTRIBUTING.md) - How to contribute
- [Examples](../examples/) - Working code exampleslog"


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