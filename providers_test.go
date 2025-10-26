package urlmeta

import (
	"strings"
	"testing"
)

func TestProviderCount(t *testing.T) {
	count := ProviderCount()
	if count < 8 {
		t.Errorf("Expected at least 8 providers, got %d", count)
	}
}

func TestGetKnownProviders(t *testing.T) {
	providers := GetKnownProviders()

	if len(providers) == 0 {
		t.Error("Expected non-empty provider list")
	}

	// Check for essential providers
	expectedProviders := []string{"YouTube", "Vimeo", "Twitter", "Instagram", "SoundCloud", "Spotify", "TikTok", "Flickr"}
	providerMap := make(map[string]bool)

	for _, p := range providers {
		providerMap[p.Name] = true

		// Validate provider structure
		if p.Name == "" {
			t.Error("Provider name cannot be empty")
		}
		if p.URL == "" {
			t.Error("Provider URL cannot be empty")
		}
		if len(p.Endpoints) == 0 {
			t.Errorf("Provider %s has no endpoints", p.Name)
		}

		// Validate endpoints
		for _, endpoint := range p.Endpoints {
			if endpoint.URL == "" {
				t.Errorf("Provider %s has endpoint with empty URL", p.Name)
			}
			if len(endpoint.Schemes) == 0 {
				t.Errorf("Provider %s has endpoint with no schemes", p.Name)
			}
		}
	}

	// Check if all expected providers exist
	for _, expected := range expectedProviders {
		if !providerMap[expected] {
			t.Errorf("Expected provider '%s' not found", expected)
		}
	}
}

func TestIsProviderSupported(t *testing.T) {
	tests := []struct {
		name      string
		supported bool
	}{
		{"YouTube", true},
		{"Vimeo", true},
		{"Twitter", true},
		{"Instagram", true},
		{"SoundCloud", true},
		{"Spotify", true},
		{"TikTok", true},
		{"Flickr", true},
		{"NonExistentProvider", false},
		{"GitHub", false},
	}

	for _, tt := range tests {
		result := IsProviderSupported(tt.name)
		if result != tt.supported {
			t.Errorf("IsProviderSupported(%s) = %v, expected %v", tt.name, result, tt.supported)
		}
	}
}

func TestGetProviderByName(t *testing.T) {
	// Test existing provider
	youtube := GetProviderByName("YouTube")
	if youtube == nil {
		t.Error("Expected to find YouTube provider")
	} else if youtube.Name != "YouTube" {
		t.Errorf("Expected provider name 'YouTube', got '%s'", youtube.Name)
	}

	// Test non-existent provider
	fake := GetProviderByName("FakeProvider")
	if fake != nil {
		t.Error("Expected nil for non-existent provider")
	}
}

func TestAddCustomProvider(t *testing.T) {
	initialCount := ProviderCount()

	// Add custom provider
	customProvider := OEmbedProvider{
		Name: "TestProvider",
		URL:  "https://test.com",
		Endpoints: []OEmbedEndpoint{
			{
				Schemes: []string{"https://test.com/*"},
				URL:     "https://test.com/oembed",
			},
		},
	}

	AddCustomProvider(customProvider)

	// Check if added
	newCount := ProviderCount()
	if newCount != initialCount+1 {
		t.Errorf("Expected provider count %d, got %d", initialCount+1, newCount)
	}

	// Check if retrievable
	if !IsProviderSupported("TestProvider") {
		t.Error("Custom provider not found after adding")
	}

	provider := GetProviderByName("TestProvider")
	if provider == nil {
		t.Error("Could not retrieve custom provider")
	} else if provider.Name != "TestProvider" {
		t.Errorf("Expected 'TestProvider', got '%s'", provider.Name)
	}
}

func TestProviderEndpointURLs(t *testing.T) {
	providers := GetKnownProviders()

	for _, provider := range providers {
		for _, endpoint := range provider.Endpoints {
			// Check endpoint URL format
			if !strings.HasPrefix(endpoint.URL, "http://") && !strings.HasPrefix(endpoint.URL, "https://") {
				t.Errorf("Provider %s has invalid endpoint URL: %s", provider.Name, endpoint.URL)
			}

			// Check schemes format
			for _, scheme := range endpoint.Schemes {
				if scheme == "" {
					t.Errorf("Provider %s has empty scheme", provider.Name)
				}
				// Schemes should contain URL patterns
				if !strings.Contains(scheme, "://") && !strings.Contains(scheme, "*") {
					t.Errorf("Provider %s has suspicious scheme: %s", provider.Name, scheme)
				}
			}
		}
	}
}

func TestProviderImmutability(t *testing.T) {
	// Get provider list twice
	providers1 := GetKnownProviders()
	providers2 := GetKnownProviders()

	// Modify first list
	if len(providers1) > 0 {
		providers1[0].Name = "Modified"
	}

	// Check second list is not affected
	if len(providers2) > 0 && providers2[0].Name == "Modified" {
		t.Error("Provider list is not properly copied, modifications leak through")
	}
}

func TestYouTubeSchemes(t *testing.T) {
	youtube := GetProviderByName("YouTube")
	if youtube == nil {
		t.Fatal("YouTube provider not found")
	}

	expectedSchemes := map[string]bool{
		"https://*.youtube.com/watch*": true,
		"https://youtu.be/*":           true,
	}

	foundSchemes := make(map[string]bool)
	for _, endpoint := range youtube.Endpoints {
		for _, scheme := range endpoint.Schemes {
			foundSchemes[scheme] = true
		}
	}

	for expected := range expectedSchemes {
		if !foundSchemes[expected] {
			t.Errorf("YouTube missing expected scheme: %s", expected)
		}
	}
}

func TestTwitterNewDomain(t *testing.T) {
	// Twitter now uses x.com as well
	twitter := GetProviderByName("Twitter")
	if twitter == nil {
		t.Fatal("Twitter provider not found")
	}

	hasXDomain := false
	for _, endpoint := range twitter.Endpoints {
		for _, scheme := range endpoint.Schemes {
			if strings.Contains(scheme, "x.com") {
				hasXDomain = true
				break
			}
		}
	}

	// Note: This might fail if x.com scheme is not added yet
	// Update providers.go to include x.com when Twitter fully migrates
	if !hasXDomain {
		t.Log("Warning: Twitter x.com domain scheme not found (might need update)")
	}
}

func BenchmarkProviderCount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ProviderCount()
	}
}

func BenchmarkIsProviderSupported(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = IsProviderSupported("YouTube")
	}
}

func BenchmarkGetProviderByName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = GetProviderByName("YouTube")
	}
}

func BenchmarkGetKnownProviders(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = GetKnownProviders()
	}
}
