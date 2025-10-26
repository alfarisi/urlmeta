package urlmeta

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// Mock HTML responses for testing
const (
	mockHTMLBasic = `
<!DOCTYPE html>
<html>
<head>
	<title>Test Page Title</title>
	<meta name="description" content="This is a test description">
	<meta name="keywords" content="test, golang, metadata">
	<meta name="author" content="Test Author">
</head>
<body>
	<h1>Test Content</h1>
</body>
</html>
`

	mockHTMLOpenGraph = `
<!DOCTYPE html>
<html>
<head>
	<title>Fallback Title</title>
	<meta property="og:title" content="OG Test Title">
	<meta property="og:description" content="OG Test Description">
	<meta property="og:image" content="https://example.com/image.jpg">
	<meta property="og:image:width" content="1200">
	<meta property="og:image:height" content="630">
	<meta property="og:type" content="article">
	<meta property="og:site_name" content="Test Site">
	<meta property="og:url" content="https://example.com/article">
	<meta property="og:locale" content="en_US">
	<meta property="article:published_time" content="2025-01-01T00:00:00Z">
	<meta property="article:author" content="OG Author">
</head>
<body></body>
</html>
`

	mockHTMLTwitterCard = `
<!DOCTYPE html>
<html>
<head>
	<title>Twitter Test</title>
	<meta name="twitter:card" content="summary_large_image">
	<meta name="twitter:site" content="@testsite">
	<meta name="twitter:creator" content="@testcreator">
	<meta name="twitter:title" content="Twitter Title">
	<meta name="twitter:description" content="Twitter Description">
	<meta name="twitter:image" content="https://example.com/twitter-image.jpg">
</head>
<body></body>
</html>
`

	mockHTMLComplete = `
<!DOCTYPE html>
<html>
<head>
	<title>Complete Test</title>
	<meta property="og:title" content="Complete OG Title">
	<meta property="og:description" content="Complete Description">
	<meta property="og:image" content="https://example.com/og-image.jpg">
	<meta property="og:video" content="https://example.com/video.mp4">
	<meta property="og:video:type" content="video/mp4">
	<meta name="twitter:card" content="summary">
	<meta name="keywords" content="test, metadata, complete">
	<link rel="icon" href="/favicon.ico">
	<link rel="canonical" href="https://example.com/canonical">
</head>
<body></body>
</html>
`

	mockHTMLRelativeURLs = `
<!DOCTYPE html>
<html>
<head>
	<title>Relative URLs Test</title>
	<meta property="og:image" content="/images/test.jpg">
	<link rel="icon" href="/favicon.ico">
	<link rel="canonical" href="/canonical-path">
</head>
<body></body>
</html>
`
)

func TestExtractBasicMetadata(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(mockHTMLBasic))
	}))
	defer server.Close()

	metadata, err := Extract(server.URL)
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	if metadata.Title != "Test Page Title" {
		t.Errorf("Expected title 'Test Page Title', got '%s'", metadata.Title)
	}

	if metadata.Description != "This is a test description" {
		t.Errorf("Expected description 'This is a test description', got '%s'", metadata.Description)
	}

	if metadata.Author != "Test Author" {
		t.Errorf("Expected author 'Test Author', got '%s'", metadata.Author)
	}

	if len(metadata.Keywords) != 3 {
		t.Errorf("Expected 3 keywords, got %d", len(metadata.Keywords))
	}

	expectedKeywords := []string{"test", "golang", "metadata"}
	for i, kw := range expectedKeywords {
		if metadata.Keywords[i] != kw {
			t.Errorf("Expected keyword '%s', got '%s'", kw, metadata.Keywords[i])
		}
	}
}

func TestExtractOpenGraphMetadata(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(mockHTMLOpenGraph))
	}))
	defer server.Close()

	metadata, err := Extract(server.URL)
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	// OG tags should override standard title
	if metadata.Title != "OG Test Title" {
		t.Errorf("Expected title 'OG Test Title', got '%s'", metadata.Title)
	}

	if metadata.Description != "OG Test Description" {
		t.Errorf("Expected OG description, got '%s'", metadata.Description)
	}

	if metadata.Type != "article" {
		t.Errorf("Expected type 'article', got '%s'", metadata.Type)
	}

	if metadata.SiteName != "Test Site" {
		t.Errorf("Expected site name 'Test Site', got '%s'", metadata.SiteName)
	}

	if metadata.Locale != "en_US" {
		t.Errorf("Expected locale 'en_US', got '%s'", metadata.Locale)
	}

	if metadata.PublishedTime != "2025-01-01T00:00:00Z" {
		t.Errorf("Expected published time '2025-01-01T00:00:00Z', got '%s'", metadata.PublishedTime)
	}

	if len(metadata.Images) != 1 {
		t.Fatalf("Expected 1 image, got %d", len(metadata.Images))
	}

	img := metadata.Images[0]
	if img.URL != "https://example.com/image.jpg" {
		t.Errorf("Expected image URL 'https://example.com/image.jpg', got '%s'", img.URL)
	}

	if img.Width != 1200 {
		t.Errorf("Expected image width 1200, got %d", img.Width)
	}

	if img.Height != 630 {
		t.Errorf("Expected image height 630, got %d", img.Height)
	}
}

func TestExtractTwitterCardMetadata(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(mockHTMLTwitterCard))
	}))
	defer server.Close()

	metadata, err := Extract(server.URL)
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	if metadata.TwitterCard != "summary_large_image" {
		t.Errorf("Expected twitter card 'summary_large_image', got '%s'", metadata.TwitterCard)
	}

	if metadata.TwitterSite != "@testsite" {
		t.Errorf("Expected twitter site '@testsite', got '%s'", metadata.TwitterSite)
	}

	if metadata.TwitterCreator != "@testcreator" {
		t.Errorf("Expected twitter creator '@testcreator', got '%s'", metadata.TwitterCreator)
	}

	if metadata.Title != "Twitter Title" {
		t.Errorf("Expected title 'Twitter Title', got '%s'", metadata.Title)
	}
}

func TestExtractCompleteMetadata(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(mockHTMLComplete))
	}))
	defer server.Close()

	metadata, err := Extract(server.URL)
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	if len(metadata.Images) != 1 {
		t.Errorf("Expected 1 image, got %d", len(metadata.Images))
	}

	if len(metadata.Videos) != 1 {
		t.Fatalf("Expected 1 video, got %d", len(metadata.Videos))
	}

	video := metadata.Videos[0]
	if video.URL != "https://example.com/video.mp4" {
		t.Errorf("Expected video URL 'https://example.com/video.mp4', got '%s'", video.URL)
	}

	if video.Type != "video/mp4" {
		t.Errorf("Expected video type 'video/mp4', got '%s'", video.Type)
	}

	if !strings.HasSuffix(metadata.Favicon, "/favicon.ico") {
		t.Errorf("Expected favicon to end with '/favicon.ico', got '%s'", metadata.Favicon)
	}

	if metadata.CanonicalURL != "https://example.com/canonical" {
		t.Errorf("Expected canonical URL 'https://example.com/canonical', got '%s'", metadata.CanonicalURL)
	}
}

func TestExtractRelativeURLs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(mockHTMLRelativeURLs))
	}))
	defer server.Close()

	metadata, err := Extract(server.URL)
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	// Check that relative URLs are resolved to absolute
	if len(metadata.Images) != 1 {
		t.Fatalf("Expected 1 image, got %d", len(metadata.Images))
	}

	expectedImageURL := server.URL + "/images/test.jpg"
	if metadata.Images[0].URL != expectedImageURL {
		t.Errorf("Expected image URL '%s', got '%s'", expectedImageURL, metadata.Images[0].URL)
	}

	expectedFavicon := server.URL + "/favicon.ico"
	if metadata.Favicon != expectedFavicon {
		t.Errorf("Expected favicon '%s', got '%s'", expectedFavicon, metadata.Favicon)
	}

	expectedCanonical := server.URL + "/canonical-path"
	if metadata.CanonicalURL != expectedCanonical {
		t.Errorf("Expected canonical URL '%s', got '%s'", expectedCanonical, metadata.CanonicalURL)
	}
}

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"example.com", "https://example.com"},
		{"http://example.com", "http://example.com"},
		{"https://example.com", "https://example.com"},
		{"example.com/path", "https://example.com/path"},
	}

	for _, tt := range tests {
		result := normalizeURL(tt.input)
		if result != tt.expected {
			t.Errorf("normalizeURL(%s) = %s, expected %s", tt.input, result, tt.expected)
		}
	}
}

func TestUnsupportedProtocol(t *testing.T) {
	_, err := Extract("ftp://example.com")
	if err == nil {
		t.Error("Expected error for FTP protocol, got nil")
	}

	if !strings.Contains(err.Error(), "unsupported protocol") {
		t.Errorf("Expected 'unsupported protocol' error, got: %v", err)
	}
}

func TestInvalidURL(t *testing.T) {
	_, err := Extract("ht!tp://invalid url")
	if err == nil {
		t.Error("Expected error for invalid URL, got nil")
	}
}

func TestHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	_, err := Extract(server.URL)
	if err == nil {
		t.Error("Expected error for 404 response, got nil")
	}

	if !strings.Contains(err.Error(), "HTTP error: 404") {
		t.Errorf("Expected '404' error, got: %v", err)
	}
}

func TestUnsupportedContentType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error": "not html"}`))
	}))
	defer server.Close()

	_, err := Extract(server.URL)
	if err == nil {
		t.Error("Expected error for non-HTML content, got nil")
	}

	if !strings.Contains(err.Error(), "unsupported content type") {
		t.Errorf("Expected 'unsupported content type' error, got: %v", err)
	}
}

func TestClientWithTimeout(t *testing.T) {
	// Server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.Write([]byte(mockHTMLBasic))
	}))
	defer server.Close()

	client := NewClient(WithTimeout(500 * time.Millisecond))
	_, err := client.Extract(server.URL)

	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
}

func TestClientWithCustomUserAgent(t *testing.T) {
	customUA := "CustomBot/1.0"
	var receivedUA string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedUA = r.Header.Get("User-Agent")
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(mockHTMLBasic))
	}))
	defer server.Close()

	client := NewClient(WithUserAgent(customUA))
	_, err := client.Extract(server.URL)

	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	if receivedUA != customUA {
		t.Errorf("Expected User-Agent '%s', got '%s'", customUA, receivedUA)
	}
}

func TestAutoOEmbedDisabled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(mockHTMLBasic))
	}))
	defer server.Close()

	// Client with auto oEmbed disabled
	client := NewClient(WithAutoOEmbed(false))
	metadata, err := client.Extract(server.URL)

	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	// oEmbed should be nil when disabled
	if metadata.OEmbed != nil {
		t.Error("Expected OEmbed to be nil when auto detection is disabled")
	}
}

func TestAutoOEmbedEnabled(t *testing.T) {
	// This test would need mock oEmbed endpoint
	// For now, we just test that the field exists
	client := NewClient(WithAutoOEmbed(true))

	if !client.autoOEmbed {
		t.Error("Expected autoOEmbed to be true")
	}
}

func TestClientWithMaxRedirects(t *testing.T) {
	redirectCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if redirectCount < 3 {
			redirectCount++
			http.Redirect(w, r, "/redirect", http.StatusFound)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(mockHTMLBasic))
	}))
	defer server.Close()

	// Should fail with 2 max redirects
	client := NewClient(WithMaxRedirects(2))
	_, err := client.Extract(server.URL)

	if err == nil {
		t.Error("Expected error for too many redirects, got nil")
	}
}

func TestEmptyMetadata(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!DOCTYPE html><html><head></head><body></body></html>`))
	}))
	defer server.Close()

	metadata, err := Extract(server.URL)
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	// Should not crash, just return empty metadata
	if metadata == nil {
		t.Error("Expected metadata object, got nil")
	} else if metadata.Title != "" {
		t.Errorf("Expected empty title, got '%s'", metadata.Title)
	}
}

func TestParseIntHelper(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"123", 123},
		{"0", 0},
		{"-1", -1},
		{"invalid", 0},
		{"", 0},
	}

	for _, tt := range tests {
		result := parseInt(tt.input)
		if result != tt.expected {
			t.Errorf("parseInt(%s) = %d, expected %d", tt.input, result, tt.expected)
		}
	}
}

func BenchmarkExtract(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(mockHTMLComplete))
	}))
	defer server.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Extract(server.URL)
		if err != nil {
			b.Fatalf("Extract failed: %v", err)
		}
	}
}

func BenchmarkExtractWithClient(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(mockHTMLComplete))
	}))
	defer server.Close()

	client := NewClient()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.Extract(server.URL)
		if err != nil {
			b.Fatalf("Extract failed: %v", err)
		}
	}
}
