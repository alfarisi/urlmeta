package urlmeta

// This file contains oEmbed provider definitions
// To add a new provider, add a new OEmbedProvider entry to knownProviders

// knownProviders contains well-known oEmbed providers with their endpoints
// This list is intentionally hardcoded for:
// - Zero runtime overhead
// - No external dependencies
// - Fast pattern matching (~1μs)
// - Compile-time validation
//
// Source: https://oembed.com/providers.json (curated and verified)
// Last updated: 2025-01-XX
var knownProviders = []OEmbedProvider{
	{
		Name: "YouTube",
		URL:  "https://www.youtube.com",
		Endpoints: []OEmbedEndpoint{
			{
				Schemes: []string{
					"https://*.youtube.com/watch*",
					"https://*.youtube.com/v/*",
					"https://youtu.be/*",
					"https://*.youtube.com/shorts/*",
				},
				URL:       "https://www.youtube.com/oembed",
				Discovery: true,
			},
		},
	},
	{
		Name: "Vimeo",
		URL:  "https://vimeo.com",
		Endpoints: []OEmbedEndpoint{
			{
				Schemes: []string{
					"https://vimeo.com/*",
					"https://vimeo.com/groups/*/videos/*",
					"https://player.vimeo.com/video/*",
				},
				URL:       "https://vimeo.com/api/oembed.json",
				Discovery: true,
			},
		},
	},
	{
		Name: "Twitter",
		URL:  "https://twitter.com",
		Endpoints: []OEmbedEndpoint{
			{
				Schemes: []string{
					"https://twitter.com/*/status/*",
					"https://twitter.com/*/statuses/*",
					"https://*.twitter.com/*/status/*",
					"https://x.com/*/status/*", // New domain
				},
				URL:       "https://publish.twitter.com/oembed",
				Discovery: true,
			},
		},
	},
	{
		Name: "Instagram",
		URL:  "https://instagram.com",
		Endpoints: []OEmbedEndpoint{
			{
				Schemes: []string{
					"http://instagram.com/*/p/*",
					"http://www.instagram.com/*/p/*",
					"https://instagram.com/*/p/*",
					"https://www.instagram.com/*/p/*",
					"http://instagram.com/p/*",
					"http://www.instagram.com/p/*",
					"https://instagram.com/p/*",
					"https://www.instagram.com/p/*",
					"https://instagram.com/reel/*",
					"https://www.instagram.com/reel/*",
				},
				URL:       "https://graph.facebook.com/v16.0/instagram_oembed",
				Discovery: false,
			},
		},
	},
	{
		Name: "Flickr",
		URL:  "https://www.flickr.com",
		Endpoints: []OEmbedEndpoint{
			{
				Schemes: []string{
					"http://*.flickr.com/photos/*",
					"http://flic.kr/p/*",
					"https://*.flickr.com/photos/*",
					"https://flic.kr/p/*",
				},
				URL:       "https://www.flickr.com/services/oembed/",
				Discovery: true,
			},
		},
	},
	{
		Name: "SoundCloud",
		URL:  "https://soundcloud.com",
		Endpoints: []OEmbedEndpoint{
			{
				Schemes: []string{
					"https://soundcloud.com/*",
					"https://soundcloud.app.goo.gl/*",
					"https://on.soundcloud.com/*",
				},
				URL:       "https://soundcloud.com/oembed",
				Discovery: true,
			},
		},
	},
	{
		Name: "Spotify",
		URL:  "https://spotify.com",
		Endpoints: []OEmbedEndpoint{
			{
				Schemes: []string{
					"https://open.spotify.com/*",
					"https://play.spotify.com/*",
					"https://open.spotify.com/track/*",
					"https://open.spotify.com/album/*",
					"https://open.spotify.com/playlist/*",
					"https://open.spotify.com/artist/*",
				},
				URL:       "https://open.spotify.com/oembed",
				Discovery: true,
			},
		},
	},
	{
		Name: "TikTok",
		URL:  "https://www.tiktok.com",
		Endpoints: []OEmbedEndpoint{
			{
				Schemes: []string{
					"https://www.tiktok.com/*/video/*",
					"https://www.tiktok.com/@*/video/*",
					"https://m.tiktok.com/*/video/*",
					"https://vm.tiktok.com/*",
				},
				URL:       "https://www.tiktok.com/oembed",
				Discovery: true,
			},
		},
	},
}

// GetKnownProviders returns a copy of the known providers list
// This is useful for displaying supported providers to users
func GetKnownProviders() []OEmbedProvider {
	// Return a copy to prevent modifications
	providers := make([]OEmbedProvider, len(knownProviders))
	copy(providers, knownProviders)
	return providers
}

// AddCustomProvider allows users to add custom oEmbed providers at runtime
// This is useful for private/internal services or new providers not yet in the list
//
// Example:
//
//	provider := urlmeta.OEmbedProvider{
//	    Name: "MyService",
//	    URL: "https://myservice.com",
//	    Endpoints: []urlmeta.OEmbedEndpoint{
//	        {
//	            Schemes: []string{"https://myservice.com/videos/*"},
//	            URL: "https://myservice.com/oembed",
//	        },
//	    },
//	}
//	urlmeta.AddCustomProvider(provider)
func AddCustomProvider(provider OEmbedProvider) {
	knownProviders = append(knownProviders, provider)
}

// ProviderCount returns the number of supported oEmbed providers
func ProviderCount() int {
	return len(knownProviders)
}

// IsProviderSupported checks if a provider name is supported
func IsProviderSupported(providerName string) bool {
	for _, p := range knownProviders {
		if p.Name == providerName {
			return true
		}
	}
	return false
}

// GetProviderByName returns a provider by its name
func GetProviderByName(name string) *OEmbedProvider {
	for _, p := range knownProviders {
		if p.Name == name {
			return &p
		}
	}
	return nil
}

/*
MAINTENANCE NOTES:

To add a new provider:
1. Go to https://oembed.com/providers.json
2. Find the provider you want to add
3. Copy the schemes and endpoint URL
4. Add a new entry to knownProviders array above
5. Update "Last updated" date
6. Run tests: go test -v
7. Commit with message: "feat: add [Provider] oEmbed support"

Example commit:
  feat: add Dailymotion oEmbed support

Common providers NOT included (and why):
- Reddit: Requires OAuth (complex)
- Medium: No official oEmbed endpoint
- LinkedIn: Rate-limited, requires API key
- Facebook: Deprecated oEmbed API

If you need support for these, use HTML extraction instead.
*/
