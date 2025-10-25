// Package urlmeta provides functionality to extract metadata from web URLs
// similar to Embedly service, supporting Open Graph, Twitter Cards, and standard HTML meta tags.
package urlmeta

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// Metadata represents extracted information from a web page
type Metadata struct {
	// Basic Info
	Title           string `json:"title"`
	Description     string `json:"description"`
	URL             string `json:"url"`
	CanonicalURL    string `json:"canonical_url,omitempty"`
	
	// Provider Info
	ProviderName    string `json:"provider_name"`
	ProviderURL     string `json:"provider_url"`
	ProviderDisplay string `json:"provider_display"`
	
	// Media
	Images          []Image `json:"images,omitempty"`
	Videos          []Video `json:"videos,omitempty"`
	
	// OpenGraph
	Type            string `json:"type,omitempty"`
	SiteName        string `json:"site_name,omitempty"`
	Locale          string `json:"locale,omitempty"`
	
	// Additional Meta
	Author          string   `json:"author,omitempty"`
	PublishedTime   string   `json:"published_time,omitempty"`
	ModifiedTime    string   `json:"modified_time,omitempty"`
	Keywords        []string `json:"keywords,omitempty"`
	
	// Twitter Card
	TwitterCard     string `json:"twitter_card,omitempty"`
	TwitterSite     string `json:"twitter_site,omitempty"`
	TwitterCreator  string `json:"twitter_creator,omitempty"`
	
	// Favicon
	Favicon         string `json:"favicon,omitempty"`
	
	// oEmbed data (automatically included if available)
	OEmbed          *OEmbed `json:"oembed,omitempty"`
}

// Image represents an image from the page
type Image struct {
	URL    string `json:"url"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
	Alt    string `json:"alt,omitempty"`
}

// Video represents a video from the page
type Video struct {
	URL    string `json:"url"`
	Type   string `json:"type,omitempty"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}

// Client handles URL metadata extraction
type Client struct {
	httpClient    *http.Client
	userAgent     string
	maxRedirects  int
	followMetaRef bool
	autoOEmbed    bool // Automatically fetch oEmbed if available
}

// Option is a function that configures a Client
type Option func(*Client)

// WithTimeout sets custom timeout for HTTP requests (default: 10s)
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// WithUserAgent sets custom User-Agent header
func WithUserAgent(ua string) Option {
	return func(c *Client) {
		c.userAgent = ua
	}
}

// WithMaxRedirects sets maximum number of redirects to follow (default: 10)
func WithMaxRedirects(max int) Option {
	return func(c *Client) {
		c.maxRedirects = max
	}
}

// WithHTTPClient sets custom HTTP client
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		c.httpClient = client
	}
}

// WithAutoOEmbed enables/disables automatic oEmbed extraction (default: true)
func WithAutoOEmbed(auto bool) Option {
	return func(c *Client) {
		c.autoOEmbed = auto
	}
}

// NewClient creates a new metadata extraction client with options
func NewClient(opts ...Option) *Client {
	c := &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		userAgent:     "Mozilla/5.0 (compatible; URLMetaBot/1.0; +https://github.com/yourusername/urlmeta)",
		maxRedirects:  10,
		followMetaRef: true,
		autoOEmbed:    true, // Enable by default
	}

	for _, opt := range opts {
		opt(c)
	}

	// Configure redirect policy
	c.httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) >= c.maxRedirects {
			return fmt.Errorf("stopped after %d redirects", c.maxRedirects)
		}
		return nil
	}

	return c
}

// Extract extracts metadata from the given URL
func (c *Client) Extract(targetURL string) (*Metadata, error) {
	// Normalize URL
	targetURL = normalizeURL(targetURL)

	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return nil, fmt.Errorf("unsupported protocol: %s (only http and https are supported)", parsedURL.Scheme)
	}

	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	// Check content type
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") && !strings.Contains(contentType, "application/xhtml") {
		return nil, fmt.Errorf("unsupported content type: %s", contentType)
	}

	// Limit response body size to prevent memory issues
	limitedBody := io.LimitReader(resp.Body, 10*1024*1024) // 10MB limit

	doc, err := html.Parse(limitedBody)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	metadata := &Metadata{
		URL:             resp.Request.URL.String(), // Use final URL after redirects
		ProviderURL:     fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host),
		ProviderDisplay: parsedURL.Host,
		Images:          []Image{},
		Videos:          []Video{},
		Keywords:        []string{},
	}

	extractFromNode(doc, metadata, parsedURL)

	// Post-processing
	metadata.Title = strings.TrimSpace(metadata.Title)
	metadata.Description = strings.TrimSpace(metadata.Description)
	
	// Set provider name from site name if available
	if metadata.SiteName != "" {
		metadata.ProviderName = metadata.SiteName
	} else {
		metadata.ProviderName = parsedURL.Host
	}

	// Auto-fetch oEmbed if enabled and supported
	if c.autoOEmbed {
		oembed, err := c.extractOEmbedQuietly(resp.Request.URL.String())
		if err == nil && oembed != nil {
			metadata.OEmbed = oembed
		}
	}

	return metadata, nil
}

// extractOEmbedQuietly attempts to extract oEmbed without returning errors
func (c *Client) extractOEmbedQuietly(targetURL string) (*OEmbed, error) {
	// Try to find oEmbed endpoint from known providers
	endpoint := findOEmbedEndpoint(targetURL)
	if endpoint != "" {
		oembed, err := c.fetchOEmbed(endpoint, targetURL)
		if err == nil {
			return oembed, nil
		}
	}

	// Try oEmbed discovery from HTML (already parsed, but we'll skip for performance)
	// Discovery would require another HTTP request, so we skip it in auto mode
	// Users can still call ExtractOEmbed() explicitly if needed

	return nil, fmt.Errorf("oEmbed not available")
}

// Extract is a convenience function using default client
func Extract(targetURL string) (*Metadata, error) {
	client := NewClient()
	return client.Extract(targetURL)
}

// normalizeURL adds https:// if no scheme is provided
func normalizeURL(targetURL string) string {
	if !strings.Contains(targetURL, "://") {
		return "https://" + targetURL
	}
	return targetURL
}

// extractFromNode traverses HTML nodes to find meta tags
func extractFromNode(n *html.Node, metadata *Metadata, baseURL *url.URL) {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "title":
			if metadata.Title == "" && n.FirstChild != nil {
				metadata.Title = n.FirstChild.Data
			}
		case "meta":
			processMeta(n, metadata, baseURL)
		case "link":
			processLink(n, metadata, baseURL)
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractFromNode(c, metadata, baseURL)
	}
}

// processMeta processes meta tags
func processMeta(n *html.Node, metadata *Metadata, baseURL *url.URL) {
	var property, name, content, itemProp string

	for _, attr := range n.Attr {
		switch attr.Key {
		case "property":
			property = attr.Val
		case "name":
			name = attr.Val
		case "content":
			content = attr.Val
		case "itemprop":
			itemProp = attr.Val
		}
	}

	content = strings.TrimSpace(content)
	if content == "" {
		return
	}

	// Open Graph Protocol
	if property != "" {
		processOpenGraph(property, content, metadata, baseURL)
	}

	// Twitter Cards
	if name != "" {
		processTwitterCard(name, content, metadata, baseURL)
	}

	// Standard meta tags
	if name != "" {
		processStandardMeta(name, content, metadata)
	}

	// Schema.org microdata
	if itemProp != "" {
		processItemProp(itemProp, content, metadata)
	}
}

// processOpenGraph handles Open Graph tags
func processOpenGraph(property, content string, metadata *Metadata, baseURL *url.URL) {
	switch property {
	case "og:title":
		if metadata.Title == "" {
			metadata.Title = content
		}
	case "og:description":
		if metadata.Description == "" {
			metadata.Description = content
		}
	case "og:image", "og:image:url":
		metadata.Images = append(metadata.Images, Image{URL: resolveURL(content, baseURL)})
	case "og:image:width":
		if len(metadata.Images) > 0 {
			if w := parseInt(content); w > 0 {
				metadata.Images[len(metadata.Images)-1].Width = w
			}
		}
	case "og:image:height":
		if len(metadata.Images) > 0 {
			if h := parseInt(content); h > 0 {
				metadata.Images[len(metadata.Images)-1].Height = h
			}
		}
	case "og:video", "og:video:url":
		metadata.Videos = append(metadata.Videos, Video{URL: resolveURL(content, baseURL)})
	case "og:video:type":
		if len(metadata.Videos) > 0 {
			metadata.Videos[len(metadata.Videos)-1].Type = content
		}
	case "og:site_name":
		metadata.SiteName = content
	case "og:type":
		metadata.Type = content
	case "og:url":
		if metadata.CanonicalURL == "" {
			metadata.CanonicalURL = content
		}
	case "og:locale":
		metadata.Locale = content
	case "article:published_time":
		metadata.PublishedTime = content
	case "article:modified_time":
		metadata.ModifiedTime = content
	case "article:author":
		if metadata.Author == "" {
			metadata.Author = content
		}
	}
}

// processTwitterCard handles Twitter Card tags
func processTwitterCard(name, content string, metadata *Metadata, baseURL *url.URL) {
	switch name {
	case "twitter:card":
		metadata.TwitterCard = content
	case "twitter:site":
		metadata.TwitterSite = content
	case "twitter:creator":
		metadata.TwitterCreator = content
	case "twitter:title":
		if metadata.Title == "" {
			metadata.Title = content
		}
	case "twitter:description":
		if metadata.Description == "" {
			metadata.Description = content
		}
	case "twitter:image", "twitter:image:src":
		metadata.Images = append(metadata.Images, Image{URL: resolveURL(content, baseURL)})
	}
}

// processStandardMeta handles standard HTML meta tags
func processStandardMeta(name, content string, metadata *Metadata) {
	switch strings.ToLower(name) {
	case "description":
		if metadata.Description == "" {
			metadata.Description = content
		}
	case "author":
		if metadata.Author == "" {
			metadata.Author = content
		}
	case "keywords":
		keywords := strings.Split(content, ",")
		for _, kw := range keywords {
			kw = strings.TrimSpace(kw)
			if kw != "" {
				metadata.Keywords = append(metadata.Keywords, kw)
			}
		}
	}
}

// processItemProp handles Schema.org microdata
func processItemProp(itemProp, content string, metadata *Metadata) {
	switch itemProp {
	case "name":
		if metadata.Title == "" {
			metadata.Title = content
		}
	case "description":
		if metadata.Description == "" {
			metadata.Description = content
		}
	case "image":
		metadata.Images = append(metadata.Images, Image{URL: content})
	}
}

// processLink handles link tags (favicon, canonical)
func processLink(n *html.Node, metadata *Metadata, baseURL *url.URL) {
	var rel, href string

	for _, attr := range n.Attr {
		switch attr.Key {
		case "rel":
			rel = attr.Val
		case "href":
			href = attr.Val
		}
	}

	href = strings.TrimSpace(href)
	if href == "" {
		return
	}

	switch strings.ToLower(rel) {
	case "icon", "shortcut icon":
		if metadata.Favicon == "" {
			metadata.Favicon = resolveURL(href, baseURL)
		}
	case "canonical":
		if metadata.CanonicalURL == "" {
			metadata.CanonicalURL = resolveURL(href, baseURL)
		}
	}
}

// resolveURL resolves relative URLs to absolute
func resolveURL(href string, baseURL *url.URL) string {
	if href == "" {
		return ""
	}

	parsedURL, err := url.Parse(href)
	if err != nil {
		return href
	}

	if parsedURL.IsAbs() {
		return href
	}

	return baseURL.ResolveReference(parsedURL).String()
}

// parseInt safely converts string to int
func parseInt(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}