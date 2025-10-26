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

	// Test with discovery URL pointing to our mock server
	/*client.httpClient = &http.Client{
		Transport: &mockTransport{
			contentURL: contentServer.URL,
			oembedURL:  oembedServer.URL + "/oembed",
		},
	}*/

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
	tests := []struct {
		url    string
		scheme string
		match  bool
	}{
		{
			"https://www.youtube.com/watch?v=123",
			"https://*.youtube.com/watch*",
			true,
		},
		{
			"https://youtu.be/123",
			"https://youtu.be/*",
			true,
		},
		{
			"https://vimeo.com/123",
			"https://vimeo.com/*",
			true,
		},
		{
			"https://example.com/test",
			"https://youtube.com/*",
			false,
		},
	}

	for _, tt := range tests {
		result := matchScheme(tt.url, tt.scheme)
		if result != tt.match {
			t.Errorf("matchScheme(%s, %s) = %v, expected %v", tt.url, tt.scheme, result, tt.match)
		}
	}
}

// mockTransport is a custom RoundTripper for testing
/*type mockTransport struct {
	contentURL string
	oembedURL  string
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if strings.Contains(req.URL.String(), "/oembed") {
		return &http.Response{
			StatusCode: 200,
			Body:       http.NoBody,
			Header:     make(http.Header),
		}, nil
	}
	return http.DefaultTransport.RoundTrip(req)
}*/

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
