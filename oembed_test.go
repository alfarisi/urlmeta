package urlmeta

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	mockOEmbedResponse = `{
		"type": "video",
		"version": "1.0",
		"title": "Test Video",
		"author_name": "Test Author",
		"author_url": "https://example.com/author",
		"provider_name": "Test Provider",
		"provider_url": "https://example.com",
		"thumbnail_url": "https://example.com/thumb.jpg",
		"thumbnail_width": 480,
		"thumbnail_height": 360,
		"html": "<iframe src=\"https://example.com/embed/123\"></iframe>",
		"width": 640,
		"height": 480
	}`

	mockHTMLWithOEmbed = `
<!DOCTYPE html>
<html>
<head>
	<title>Test Page</title>
	<link rel="alternate" type="application/json+oembed" 
		  href="https://example.com/oembed?url=https://example.com/video/123" 
		  title="Test Video oEmbed">
</head>
<body>
	<h1>Test Content</h1>
</body>
</html>
`

	mockHTMLWithoutOEmbed = `
<!DOCTYPE html>
<html>
<head>
	<title>Test Page</title>
</head>
<body>
	<h1>No oEmbed here</h1>
</body>
</html>
`
)

func TestExtractOEmbed(t *testing.T) {
	// Mock oEmbed endpoint
	oembedServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.String(), "/oembed") {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(mockOEmbedResponse))
	}))
	defer oembedServer.Close()

	// Mock content page with oEmbed discovery
	contentServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		html := strings.Replace(mockHTMLWithOEmbed, "https://example.com/oembed", oembedServer.URL+"/oembed", 1)
		w.Write([]byte(html))
	}))
	defer contentServer.Close()

	client := NewClient()

	oembed, err := client.ExtractOEmbed(contentServer.URL)
	if err != nil {
		t.Fatalf("ExtractOEmbed failed: %v", err)
	}

	if oembed.Type != "video" {
		t.Errorf("Expected type 'video', got '%s'", oembed.Type)
	}

	if oembed.Title != "Test Video" {
		t.Errorf("Expected title 'Test Video', got '%s'", oembed.Title)
	}

	if oembed.AuthorName != "Test Author" {
		t.Errorf("Expected author 'Test Author', got '%s'", oembed.AuthorName)
	}

	if oembed.ProviderName != "Test Provider" {
		t.Errorf("Expected provider 'Test Provider', got '%s'", oembed.ProviderName)
	}

	if oembed.Width != 640 {
		t.Errorf("Expected width 640, got %d", oembed.Width)
	}

	if oembed.Height != 480 {
		t.Errorf("Expected height 480, got %d", oembed.Height)
	}

	if oembed.ThumbnailURL != "https://example.com/thumb.jpg" {
		t.Errorf("Expected thumbnail URL, got '%s'", oembed.ThumbnailURL)
	}
}

func TestIsOEmbedSupported(t *testing.T) {
	tests := []struct {
		url       string
		supported bool
	}{
		{"https://www.youtube.com/watch?v=MbyvLY8CGFM", true},
		{"https://youtu.be/MbyvLY8CGFM", true},
		{"https://vimeo.com/123456", true},
		{"https://twitter.com/user/status/123456", true},
		{"https://soundcloud.com/artist/track", true},
		{"https://open.spotify.com/track/123", true},
		{"https://example.com/random", false},
		{"https://github.com/user/repo", false},
	}

	for _, tt := range tests {
		result := IsOEmbedSupported(tt.url)
		if result != tt.supported {
			t.Errorf("IsOEmbedSupported(%s) = %v, expected %v", tt.url, result, tt.supported)
		}
	}
}

func TestFindOEmbedEndpoint(t *testing.T) {
	tests := []struct {
		url      string
		expected string
	}{
		{
			"https://www.youtube.com/watch?v=123",
			"https://www.youtube.com/oembed",
		},
		{
			"https://vimeo.com/123456",
			"https://vimeo.com/api/oembed.json",
		},
		{
			"https://twitter.com/user/status/123",
			"https://publish.twitter.com/oembed",
		},
		{
			"https://example.com/random",
			"",
		},
	}

	for _, tt := range tests {
		result := findOEmbedEndpoint(tt.url)
		if result != tt.expected {
			t.Errorf("findOEmbedEndpoint(%s) = %s, expected %s", tt.url, result, tt.expected)
		}
	}
}

func TestDiscoverOEmbedEndpoint(t *testing.T) {
	// Test with oEmbed link
	serverWithOEmbed := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(mockHTMLWithOEmbed))
	}))
	defer serverWithOEmbed.Close()

	client := NewClient()
	endpoint, err := client.discoverOEmbedEndpoint(serverWithOEmbed.URL)
	if err != nil {
		t.Fatalf("discoverOEmbedEndpoint failed: %v", err)
	}

	if !strings.Contains(endpoint, "/oembed") {
		t.Errorf("Expected endpoint to contain '/oembed', got '%s'", endpoint)
	}

	// Test without oEmbed link
	serverWithoutOEmbed := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(mockHTMLWithoutOEmbed))
	}))
	defer serverWithoutOEmbed.Close()

	endpoint, err = client.discoverOEmbedEndpoint(serverWithoutOEmbed.URL)
	if err != nil {
		t.Fatalf("discoverOEmbedEndpoint failed: %v", err)
	}

	if endpoint != "" {
		t.Errorf("Expected empty endpoint, got '%s'", endpoint)
	}
}

func TestOEmbedJSONMarshaling(t *testing.T) {
	oembed := &OEmbed{
		Type:         "video",
		Version:      "1.0",
		Title:        "Test",
		ProviderName: "Provider",
		Width:        640,
		Height:       480,
	}

	jsonData, err := json.Marshal(oembed)
	if err != nil {
		t.Fatalf("Failed to marshal oEmbed: %v", err)
	}

	var decoded OEmbed
	err = json.Unmarshal(jsonData, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal oEmbed: %v", err)
	}

	if decoded.Type != oembed.Type {
		t.Errorf("Type mismatch after marshal/unmarshal")
	}

	if decoded.Width != oembed.Width {
		t.Errorf("Width mismatch after marshal/unmarshal")
	}
}

func TestGetSupportedProviders(t *testing.T) {
	providers := GetSupportedProviders()

	if len(providers) == 0 {
		t.Error("Expected at least one provider")
	}

	// Check for some well-known providers
	providerNames := make(map[string]bool)
	for _, p := range providers {
		providerNames[p.Name] = true
	}

	expectedProviders := []string{"YouTube", "Vimeo", "Twitter", "Instagram", "SoundCloud", "Spotify"}
	for _, expected := range expectedProviders {
		if !providerNames[expected] {
			t.Errorf("Expected provider '%s' not found", expected)
		}
	}
}

func TestMatchScheme(t *testing.T) {
	// Clear cache before testing
	clearRegexCache()

	tests := []struct {
		name   string
		url    string
		scheme string
		match  bool
	}{
		{
			name:   "YouTube www with wildcard",
			url:    "https://www.youtube.com/watch?v=123",
			scheme: "https://*.youtube.com/watch*",
			match:  true,
		},
		{
			name:   "YouTube mobile with wildcard",
			url:    "https://m.youtube.com/watch?v=abc",
			scheme: "https://*.youtube.com/watch*",
			match:  true,
		},
		{
			name:   "YouTube short URL",
			url:    "https://youtu.be/123",
			scheme: "https://youtu.be/*",
			match:  true,
		},
		{
			name:   "YouTube shorts",
			url:    "https://www.youtube.com/shorts/abc123",
			scheme: "https://*.youtube.com/shorts/*",
			match:  true,
		},
		{
			name:   "Vimeo video",
			url:    "https://vimeo.com/123456",
			scheme: "https://vimeo.com/*",
			match:  true,
		},
		{
			name:   "Vimeo groups",
			url:    "https://vimeo.com/groups/test/videos/123",
			scheme: "https://vimeo.com/groups/*/videos/*",
			match:  true,
		},
		{
			name:   "Twitter status",
			url:    "https://twitter.com/user/status/123456",
			scheme: "https://twitter.com/*/status/*",
			match:  true,
		},
		{
			name:   "X.com (new Twitter)",
			url:    "https://x.com/user/status/123456",
			scheme: "https://x.com/*/status/*",
			match:  true,
		},
		{
			name:   "Wrong domain",
			url:    "https://example.com/test",
			scheme: "https://youtube.com/*",
			match:  false,
		},
		{
			name:   "Wrong path",
			url:    "https://youtube.com/about",
			scheme: "https://youtube.com/watch*",
			match:  false,
		},
		{
			name:   "HTTP vs HTTPS",
			url:    "http://youtube.com/watch",
			scheme: "https://youtube.com/watch*",
			match:  false,
		},
		{
			name:   "Subdomain mismatch",
			url:    "https://api.youtube.com/watch",
			scheme: "https://www.youtube.com/watch*",
			match:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchScheme(tt.url, tt.scheme)
			if result != tt.match {
				t.Errorf("matchScheme(%s, %s) = %v, expected %v",
					tt.url, tt.scheme, result, tt.match)
			}
		})
	}
}

func TestMatchSchemeEdgeCases(t *testing.T) {
	clearRegexCache()

	tests := []struct {
		name   string
		url    string
		scheme string
		match  bool
	}{
		{
			name:   "Empty URL",
			url:    "",
			scheme: "https://youtube.com/*",
			match:  false,
		},
		{
			name:   "Empty scheme",
			url:    "https://youtube.com/watch",
			scheme: "",
			match:  false,
		},
		{
			name:   "Invalid URL format",
			url:    "not-a-valid-url",
			scheme: "https://youtube.com/*",
			match:  false,
		},
		{
			name:   "Scheme without wildcard",
			url:    "https://youtube.com/watch",
			scheme: "https://youtube.com/watch",
			match:  true,
		},
		{
			name:   "Multiple wildcards",
			url:    "https://www.youtube.com/watch?v=123&t=10",
			scheme: "https://*.youtube.com/watch*",
			match:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchScheme(tt.url, tt.scheme)
			if result != tt.match {
				t.Errorf("matchScheme(%s, %s) = %v, expected %v",
					tt.url, tt.scheme, result, tt.match)
			}
		})
	}
}

func TestRegexCaching(t *testing.T) {
	clearRegexCache()

	scheme := "https://*.youtube.com/watch*"
	url := "https://www.youtube.com/watch?v=123"

	// First call - should compile regex
	result1 := matchScheme(url, scheme)
	if !result1 {
		t.Error("First match should succeed")
	}

	// Second call - should use cached regex
	result2 := matchScheme(url, scheme)
	if !result2 {
		t.Error("Cached match should succeed")
	}

	// Verify cache contains the scheme
	regexCacheMutex.RLock()
	_, exists := regexCache[scheme]
	regexCacheMutex.RUnlock()

	if !exists {
		t.Error("Regex should be cached after first use")
	}

	// Clear cache and verify
	clearRegexCache()

	regexCacheMutex.RLock()
	cacheSize := len(regexCache)
	regexCacheMutex.RUnlock()

	if cacheSize != 0 {
		t.Errorf("Cache should be empty after clear, got %d items", cacheSize)
	}
}

func TestSchemeToRegex(t *testing.T) {
	tests := []struct {
		name     string
		scheme   string
		testURL  string
		expected bool
	}{
		{
			name:     "Domain wildcard",
			scheme:   "https://*.youtube.com/watch",
			testURL:  "https://www.youtube.com/watch",
			expected: true,
		},
		{
			name:     "Path wildcard",
			scheme:   "https://youtu.be/*",
			testURL:  "https://youtu.be/abc123",
			expected: true,
		},
		{
			name:     "Both wildcards",
			scheme:   "https://*.youtube.com/watch*",
			testURL:  "https://m.youtube.com/watch?v=123",
			expected: true,
		},
		{
			name:     "No wildcard",
			scheme:   "https://youtube.com/watch",
			testURL:  "https://youtube.com/watch",
			expected: true,
		},
		{
			name:     "Special chars in path",
			scheme:   "https://example.com/path-with.dots/*",
			testURL:  "https://example.com/path-with.dots/file",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearRegexCache()
			result := matchScheme(tt.testURL, tt.scheme)
			if result != tt.expected {
				pattern := schemeToRegex(tt.scheme)
				t.Errorf("matchScheme failed:\nScheme: %s\nRegex: %s\nURL: %s\nGot: %v, Expected: %v",
					tt.scheme, pattern, tt.testURL, result, tt.expected)
			}
		})
	}
}

func BenchmarkExtractOEmbed(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.String(), "/oembed") {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(mockOEmbedResponse))
		} else {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(mockHTMLWithOEmbed))
		}
	}))
	defer server.Close()

	client := NewClient()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.ExtractOEmbed(server.URL)
	}
}

func BenchmarkMatchScheme(b *testing.B) {
	clearRegexCache()

	url := "https://www.youtube.com/watch?v=123"
	scheme := "https://*.youtube.com/watch*"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		matchScheme(url, scheme)
	}
}

func BenchmarkMatchSchemeCached(b *testing.B) {
	clearRegexCache()

	url := "https://www.youtube.com/watch?v=123"
	scheme := "https://*.youtube.com/watch*"

	// Warm up cache
	matchScheme(url, scheme)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		matchScheme(url, scheme)
	}
}

func BenchmarkMatchSchemeMultiplePatterns(b *testing.B) {
	clearRegexCache()

	testCases := []struct {
		url    string
		scheme string
	}{
		{"https://www.youtube.com/watch?v=123", "https://*.youtube.com/watch*"},
		{"https://vimeo.com/123456", "https://vimeo.com/*"},
		{"https://twitter.com/user/status/123", "https://twitter.com/*/status/*"},
		{"https://youtu.be/abc", "https://youtu.be/*"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tc := testCases[i%len(testCases)]
		matchScheme(tc.url, tc.scheme)
	}
}
