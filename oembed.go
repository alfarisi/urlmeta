package urlmeta

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

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

// matchScheme checks if URL matches the scheme pattern
func matchScheme(targetURL, scheme string) bool {
	// Simple pattern matching (can be improved with regex)
	scheme = strings.Replace(scheme, "*", "", -1)
	scheme = strings.Replace(scheme, "http://", "", 1)
	scheme = strings.Replace(scheme, "https://", "", 1)

	targetURL = strings.Replace(targetURL, "http://", "", 1)
	targetURL = strings.Replace(targetURL, "https://", "", 1)

	// Basic contains check (simplified)
	parts := strings.Split(scheme, "/")
	for _, part := range parts {
		if part != "" && !strings.Contains(targetURL, part) {
			return false
		}
	}
	return true
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
