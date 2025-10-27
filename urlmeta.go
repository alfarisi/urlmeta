// Package urlmeta provides functionality to extract metadata from web URLs
// similar to Embedly service, supporting Open Graph, Twitter Cards, and standard HTML meta tags.
package urlmeta

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// Metadata represents extracted information from a web page
type Metadata struct {
	// Basic Info
	Title        string `json:"title"`
	Description  string `json:"description"`
	URL          string `json:"url"`
	CanonicalURL string `json:"canonical_url,omitempty"`

	// Provider Info
	ProviderName    string `json:"provider_name"`
	ProviderURL     string `json:"provider_url"`
	ProviderDisplay string `json:"provider_display"`

	// Media
	Images []Image `json:"images,omitempty"`
	Videos []Video `json:"videos,omitempty"`

	// OpenGraph
	Type     string `json:"type,omitempty"`
	SiteName string `json:"site_name,omitempty"`
	Locale   string `json:"locale,omitempty"`
	OGTitle  string `json:"og_title,omitempty"`

	// Additional Meta
	Author        string   `json:"author,omitempty"`
	PublishedTime string   `json:"published_time,omitempty"`
	ModifiedTime  string   `json:"modified_time,omitempty"`
	Keywords      []string `json:"keywords,omitempty"`

	// Twitter Card
	TwitterCard    string `json:"twitter_card,omitempty"`
	TwitterSite    string `json:"twitter_site,omitempty"`
	TwitterCreator string `json:"twitter_creator,omitempty"`
	TwitterTitle   string `json:"twitter_title,omitempty"`

	// Favicon
	Favicon string `json:"favicon,omitempty"`

	// oEmbed (automatically included if available)
	OEmbed *OEmbed `json:"oembed,omitempty"`
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

// ExtractionStrategy determines how metadata is extracted
type ExtractionStrategy int

const (
	// StrategyAuto automatically chooses best strategy (default)
	StrategyAuto ExtractionStrategy = iota
	// StrategyOEmbedFirst tries oEmbed first, falls back to HTML
	StrategyOEmbedFirst
	// StrategyHTMLOnly only extracts from HTML (fastest for non-embed sites)
	StrategyHTMLOnly
)

// Client handles URL metadata extraction
type Client struct {
	httpClient   *http.Client
	userAgent    string
	maxRedirects int
	autoOEmbed   bool
	strategy     ExtractionStrategy
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

// WithStrategy sets extraction strategy (default: StrategyAuto)
func WithStrategy(strategy ExtractionStrategy) Option {
	return func(c *Client) {
		c.strategy = strategy
	}
}

// NewClient creates a new metadata extraction client with options
func NewClient(opts ...Option) *Client {
	c := &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		userAgent:    "Mozilla/5.0 (compatible; URLMetaBot/1.0; +https://github.com/yourusername/urlmeta)",
		maxRedirects: 10,
		autoOEmbed:   true,
		strategy:     StrategyAuto,
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

// Extract extracts metadata from the given URL using optimal strategy
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

	// Choose extraction strategy
	strategy := c.strategy
	if strategy == StrategyAuto {
		// Auto-detect: if oEmbed supported, use oEmbed-first strategy
		if c.autoOEmbed && IsOEmbedSupported(targetURL) {
			strategy = StrategyOEmbedFirst
		} else {
			strategy = StrategyHTMLOnly
		}
	}

	// Execute strategy
	switch strategy {
	case StrategyOEmbedFirst:
		return c.extractOEmbedFirst(targetURL, parsedURL)
	case StrategyHTMLOnly:
		return c.extractHTMLOnly(targetURL, parsedURL)
	default:
		return c.extractHTMLOnly(targetURL, parsedURL)
	}
}

// extractOEmbedFirst tries oEmbed first, optionally fetches HTML for additional data
func (c *Client) extractOEmbedFirst(targetURL string, parsedURL *url.URL) (*Metadata, error) {
	// Step 1: Get oEmbed data (ONLY 1 HTTP call!)
	oembed, err := c.ExtractOEmbed(targetURL)
	if err != nil {
		// oEmbed failed, fall back to HTML
		return c.extractHTMLOnly(targetURL, parsedURL)
	}

	// Step 2: Build metadata from oEmbed (no HTML parsing needed!)
	metadata := &Metadata{
		URL:             targetURL,
		ProviderURL:     fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host),
		ProviderDisplay: parsedURL.Host,
		Images:          []Image{},
		Videos:          []Video{},
		Keywords:        []string{},
		OEmbed:          oembed,
	}

	// Fill from oEmbed data
	if oembed.Title != "" {
		metadata.Title = oembed.Title
	}
	if oembed.AuthorName != "" {
		metadata.Author = oembed.AuthorName
	}
	if oembed.ProviderName != "" {
		metadata.ProviderName = oembed.ProviderName
		metadata.SiteName = oembed.ProviderName
	} else {
		metadata.ProviderName = parsedURL.Host
	}
	if oembed.ProviderURL != "" {
		metadata.ProviderURL = oembed.ProviderURL
	}

	// Add oEmbed thumbnail as image
	if oembed.ThumbnailURL != "" {
		metadata.Images = append(metadata.Images, Image{
			URL:    oembed.ThumbnailURL,
			Width:  oembed.ThumbnailWidth,
			Height: oembed.ThumbnailHeight,
		})
	}

	// For photo type, add the photo URL
	if oembed.Type == "photo" && oembed.URL != "" {
		metadata.Images = append(metadata.Images, Image{
			URL:    oembed.URL,
			Width:  oembed.Width,
			Height: oembed.Height,
		})
	}

	// Set type based on oEmbed
	metadata.Type = oembed.Type

	// OPTIMIZATION: We already have enough data from oEmbed!
	// Skip HTML fetching unless user explicitly needs it
	// This saves 1 HTTP call and parsing time!

	return metadata, nil
}

// extractHTMLOnly extracts metadata from HTML only
func (c *Client) extractHTMLOnly(targetURL string, parsedURL *url.URL) (*Metadata, error) {
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
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			_ = closeErr
		}
	}()

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
		URL:             resp.Request.URL.String(),
		ProviderURL:     fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host),
		ProviderDisplay: parsedURL.Host,
		Images:          []Image{},
		Videos:          []Video{},
		Keywords:        []string{},
	}

	extractFromNode(doc, metadata, parsedURL)

	// Post-processing
	if metadata.OGTitle != "" {
		metadata.Title = metadata.OGTitle
	} else if metadata.TwitterTitle != "" {
		metadata.Title = metadata.TwitterTitle
	}

	metadata.Title = strings.TrimSpace(metadata.Title)
	metadata.Description = strings.TrimSpace(metadata.Description)

	if metadata.SiteName != "" {
		metadata.ProviderName = metadata.SiteName
	} else {
		metadata.ProviderName = parsedURL.Host
	}

	return metadata, nil
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
	title := ""
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

	if metadata.Title == "" {
		metadata.Title = title
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

	if property != "" {
		processOpenGraph(property, content, metadata, baseURL)
	}

	if name != "" {
		processTwitterCard(name, content, metadata, baseURL)
		processStandardMeta(name, content, metadata)
	}

	if itemProp != "" {
		processItemProp(itemProp, content, metadata)
	}
}

// processOpenGraph handles Open Graph tags
func processOpenGraph(property, content string, metadata *Metadata, baseURL *url.URL) {
	// Map of simple string assignments
	simpleAssignments := map[string]*string{
		"og:site_name":           &metadata.SiteName,
		"og:type":                &metadata.Type,
		"og:locale":              &metadata.Locale,
		"article:published_time": &metadata.PublishedTime,
		"article:modified_time":  &metadata.ModifiedTime,
	}

	// Handle simple string assignments
	if target := simpleAssignments[property]; target != nil {
		*target = content
		return
	}

	// Handle title with fallback
	if property == "og:title" {
		metadata.OGTitle = content
		if metadata.Title == "" {
			metadata.Title = content
		}
		return
	}

	// Handle description with fallback
	if property == "og:description" {
		if metadata.Description == "" {
			metadata.Description = content
		}
		return
	}

	// Handle URL/canonical
	if property == "og:url" {
		if metadata.CanonicalURL == "" {
			metadata.CanonicalURL = content
		}
		return
	}

	// Handle author with fallback
	if property == "article:author" {
		if metadata.Author == "" {
			metadata.Author = content
		}
		return
	}

	// Handle images
	if processOpenGraphImage(property, content, metadata, baseURL) {
		return
	}

	// Handle videos
	processOpenGraphVideo(property, content, metadata, baseURL)
}

// processOpenGraphImage handles image-related Open Graph properties
func processOpenGraphImage(property, content string, metadata *Metadata, baseURL *url.URL) bool {
	switch property {
	case "og:image", "og:image:url":
		metadata.Images = append(metadata.Images, Image{URL: resolveURL(content, baseURL)})
		return true
	case "og:image:width":
		processImageDimension(metadata, content, true)
		return true
	case "og:image:height":
		processImageDimension(metadata, content, false)
		return true
	}
	return false
}

// processOpenGraphVideo handles video-related Open Graph properties
func processOpenGraphVideo(property, content string, metadata *Metadata, baseURL *url.URL) bool {
	switch property {
	case "og:video", "og:video:url":
		metadata.Videos = append(metadata.Videos, Video{URL: resolveURL(content, baseURL)})
		return true
	case "og:video:type":
		if len(metadata.Videos) > 0 {
			metadata.Videos[len(metadata.Videos)-1].Type = content
		}
		return true
	}
	return false
}

// processImageDimension handles image width/height
func processImageDimension(metadata *Metadata, content string, isWidth bool) {
	if len(metadata.Images) > 0 {
		dimension := parseInt(content)
		if dimension > 0 {
			if isWidth {
				metadata.Images[len(metadata.Images)-1].Width = dimension
			} else {
				metadata.Images[len(metadata.Images)-1].Height = dimension
			}
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
		metadata.TwitterTitle = content
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
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}
