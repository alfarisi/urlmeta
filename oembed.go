package urlmeta

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

// OEmbed represents oEmbed response data
// Specification: https://oembed.com/
type OEmbed struct {
	Type            string `json:"type"`                       // photo, video, link, rich
	Version         string `json:"version"`                    // oEmbed version (usually "1.0")
	Title           string `json:"title,omitempty"`            // Resource title
	AuthorName      string `json:"author_name,omitempty"`      // Author/owner name
	AuthorURL       string `json:"author_url,omitempty"`       // Author/owner URL
	ProviderName    string `json:"provider_name,omitempty"`    // Provider name
	ProviderURL     string `json:"provider_url,omitempty"`     // Provider URL
	CacheAge        int    `json:"cache_age,omitempty"`        // Suggested cache lifetime in seconds
	ThumbnailURL    string `json:"thumbnail_url,omitempty"`    // Thumbnail URL
	ThumbnailWidth  int    `json:"thumbnail_width,omitempty"`  // Thumbnail width
	ThumbnailHeight int    `json:"thumbnail_height,omitempty"` // Thumbnail height

	// Photo type specific
	URL    string `json:"url,omitempty"`    // Photo URL
	Width  int    `json:"width,omitempty"`  // Photo width
	Height int    `json:"height,omitempty"` // Photo height

	// Video/Rich type specific
	HTML string `json:"html,omitempty"` // HTML embed code
}

// OEmbedProvider represents an oEmbed provider configuration
type OEmbedProvider struct {
	Name      string
	URL       string
	Endpoints []OEmbedEndpoint
}

// OEmbedEndpoint represents an oEmbed endpoint
type OEmbedEndpoint struct {
	Schemes   []string
	URL       string
	Discovery bool
}

// Note: Provider list is defined in providers.go for better organization

// ExtractOEmbed attempts to extract oEmbed data from a URL
func (c *Client) ExtractOEmbed(targetURL string) (*OEmbed, error) {
	// Normalize URL
	targetURL = normalizeURL(targetURL)

	// 1. Try to find oEmbed endpoint from known providers
	endpoint := findOEmbedEndpoint(targetURL)
	if endpoint != "" {
		oembed, err := c.fetchOEmbed(endpoint, targetURL)
		if err == nil {
			return oembed, nil
		}
	}

	// 2. Try oEmbed discovery from HTML
	discoveredEndpoint, err := c.discoverOEmbedEndpoint(targetURL)
	if err == nil && discoveredEndpoint != "" {
		oembed, err := c.fetchOEmbed(discoveredEndpoint, targetURL)
		if err == nil {
			return oembed, nil
		}
	}

	return nil, fmt.Errorf("oEmbed endpoint not found for URL: %s", targetURL)
}

// ExtractOEmbed is a convenience function using default client
func ExtractOEmbed(targetURL string) (*OEmbed, error) {
	client := NewClient()
	return client.ExtractOEmbed(targetURL)
}

// findOEmbedEndpoint finds oEmbed endpoint from known providers
func findOEmbedEndpoint(targetURL string) string {
	for _, provider := range knownProviders {
		for _, endpoint := range provider.Endpoints {
			for _, scheme := range endpoint.Schemes {
				if matchScheme(targetURL, scheme) {
					return endpoint.URL
				}
			}
		}
	}
	return ""
}

// Cache compiled regexes for performance
var (
	regexCache      = make(map[string]*regexp.Regexp)
	regexCacheMutex sync.RWMutex
)

// matchScheme checks if URL matches the scheme pattern using regex
// Supports wildcards: *, *.domain.com, /path/*
// Examples:
//   - "https://*.youtube.com/watch*" matches "https://www.youtube.com/watch?v=123"
//   - "https://youtu.be/*" matches "https://youtu.be/abc123"
func matchScheme(targetURL, scheme string) bool {
	// Get or compile regex for this scheme
	re := getCompiledRegex(scheme)
	if re == nil {
		return false
	}

	return re.MatchString(targetURL)
}

// getCompiledRegex gets cached regex or compiles new one
func getCompiledRegex(scheme string) *regexp.Regexp {
	// Try to get from cache first (read lock)
	regexCacheMutex.RLock()
	if re, exists := regexCache[scheme]; exists {
		regexCacheMutex.RUnlock()
		return re
	}
	regexCacheMutex.RUnlock()

	// Compile new regex (write lock)
	regexCacheMutex.Lock()
	defer regexCacheMutex.Unlock()

	// Double-check after acquiring write lock
	if re, exists := regexCache[scheme]; exists {
		return re
	}

	// Convert scheme pattern to regex
	pattern := schemeToRegex(scheme)
	re, err := regexp.Compile(pattern)
	if err != nil {
		// Invalid pattern, return nil
		return nil
	}

	// Cache for future use
	regexCache[scheme] = re
	return re
}

// schemeToRegex converts oEmbed scheme pattern to regex pattern
// Scheme format: "https://*.youtube.com/watch*"
// Regex output: "^https://[^/]*\.youtube\.com/watch.*$"
func schemeToRegex(scheme string) string {
	// Escape special regex characters except *
	pattern := regexp.QuoteMeta(scheme)

	// Replace escaped \* with regex equivalents
	// *.domain.com -> [^/]* (any chars except /)
	// /path/* -> .* (any chars)

	// Replace \* at domain level (before first /)
	parts := strings.SplitN(pattern, "/", 4) // Split: scheme, "", domain, path
	if len(parts) >= 3 {
		// Handle domain wildcards: *.youtube.com
		parts[2] = strings.Replace(parts[2], "\\*", "[^/]*", -1)

		// Handle path wildcards: /watch*
		if len(parts) >= 4 {
			parts[3] = strings.Replace(parts[3], "\\*", ".*", -1)
		}

		pattern = strings.Join(parts, "/")
	} else {
		// Fallback: just replace all \*
		pattern = strings.Replace(pattern, "\\*", ".*", -1)
	}

	// Anchor to match full URL
	return "^" + pattern + "$"
}

// clearRegexCache clears the regex cache (useful for testing)
func clearRegexCache() {
	regexCacheMutex.Lock()
	defer regexCacheMutex.Unlock()
	regexCache = make(map[string]*regexp.Regexp)
}

// discoverOEmbedEndpoint discovers oEmbed endpoint from HTML
func (c *Client) discoverOEmbedEndpoint(targetURL string) (string, error) {
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			// Ignore close error
			_ = closeErr
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}

	endpoint := findOEmbedLink(doc)
	if endpoint != "" {
		// Resolve relative URLs
		baseURL, parseErr := url.Parse(targetURL)
		if parseErr != nil {
			return endpoint, nil
		}
		endpointURL, parseErr := url.Parse(endpoint)
		if parseErr == nil && !endpointURL.IsAbs() {
			endpoint = baseURL.ResolveReference(endpointURL).String()
		}
	}

	return endpoint, nil
}

// findOEmbedLink searches for oEmbed link in HTML
func findOEmbedLink(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "link" {
		var rel, href, typeAttr string
		for _, attr := range n.Attr {
			switch attr.Key {
			case "rel":
				rel = attr.Val
			case "href":
				href = attr.Val
			case "type":
				typeAttr = attr.Val
			}
		}

		// Look for oEmbed link
		if rel == "alternate" && (typeAttr == "application/json+oembed" || typeAttr == "text/json+oembed") {
			return href
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := findOEmbedLink(c); result != "" {
			return result
		}
	}

	return ""
}

// fetchOEmbed fetches oEmbed data from endpoint
func (c *Client) fetchOEmbed(endpoint, targetURL string) (*OEmbed, error) {
	// Build oEmbed request URL
	oembedURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	query := oembedURL.Query()
	query.Set("url", targetURL)
	query.Set("format", "json")
	oembedURL.RawQuery = query.Encode()

	req, err := http.NewRequest("GET", oembedURL.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			// Ignore close error
			_ = closeErr
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("oEmbed endpoint returned HTTP %d", resp.StatusCode)
	}

	var oembed OEmbed
	if err := json.NewDecoder(resp.Body).Decode(&oembed); err != nil {
		return nil, fmt.Errorf("failed to decode oEmbed response: %w", err)
	}

	return &oembed, nil
}

// IsOEmbedSupported checks if a URL is likely to support oEmbed
func IsOEmbedSupported(targetURL string) bool {
	return findOEmbedEndpoint(targetURL) != ""
}

// GetSupportedProviders returns list of known oEmbed providers
// Provider list is defined in providers.go
func GetSupportedProviders() []OEmbedProvider {
	return GetKnownProviders()
}
