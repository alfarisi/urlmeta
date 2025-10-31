# Contributing to URLMeta

Thank you for your interest in contributing to URLMeta! We welcome contributions from the community.

## Table of Contents

- [How to Contribute](#how-to-contribute)
- [Development Setup](#development-setup)
- [Adding oEmbed Providers](#adding-oembed-providers)
- [Code Quality Standards](#code-quality-standards)
- [Testing Guidelines](#testing-guidelines)
- [Pull Request Process](#pull-request-process)
- [Project Structure](#project-structure)

## How to Contribute

### Reporting Bugs

If you find a bug, please open an issue with:
- Clear, descriptive title
- Steps to reproduce the issue
- Expected vs actual behavior
- Go version and OS
- Example code if applicable
- URL causing the issue (if not sensitive)

### Suggesting Enhancements

We welcome enhancement suggestions! Please open an issue with:
- Clear description of the enhancement
- Why this enhancement would be useful
- Example use cases
- Proposed API (if applicable)

### Contributing Code

1. **Fork the repository** and create your branch from `main`
2. **Make your changes** following our coding standards
3. **Add tests** for any new functionality
4. **Update documentation** if needed
5. **Run all checks** to ensure everything passes
6. **Submit a pull request**

## Development Setup

### Prerequisites

- Go 1.21 or later
- Git
- Make (optional, but recommended)

### Setup

```bash
# Clone your fork
git clone https://github.com/alfarisi/urlmeta.git
cd urlmeta

# Install dependencies
go mod download

# Install development tools (optional)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Running Tests

```bash
# Run all tests
go test -v ./...

# Run with coverage
go test -v -cover ./...

# Run with race detector
go test -v -race ./...

# Run benchmarks
go test -bench=. -benchmem ./...
```

### Code Quality Checks

```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run

# Run vet
go vet ./...
```

## Adding oEmbed Providers

**This is the most common and valuable contribution!**

### Step-by-Step Guide

#### 1. Find Provider Information

Visit https://oembed.com/providers.json and find your provider.

Example for Dailymotion:
```json
{
  "provider_name": "Dailymotion",
  "provider_url": "https://www.dailymotion.com",
  "endpoints": [{
    "schemes": ["https://www.dailymotion.com/video/*"],
    "url": "https://www.dailymotion.com/services/oembed"
  }]
}
```

#### 2. Add to `providers.go`

Add your provider to the `knownProviders` array:

```go
{
    Name: "Dailymotion",
    URL:  "https://www.dailymotion.com",
    Endpoints: []OEmbedEndpoint{
        {
            Schemes: []string{
                "https://www.dailymotion.com/video/*",
                "https://dai.ly/*", // Include short URLs
            },
            URL:       "https://www.dailymotion.com/services/oembed",
            Discovery: true,
        },
    },
},
```

**Tips:**
- Include all URL schemes (including short URLs)
- Set `Discovery: true` if provider supports it
- Keep alphabetical order for easier maintenance
- Update "Last updated" date in comment

#### 3. Add Tests

Add test cases in `providers_test.go`:

```go
func TestDailymotionSupport(t *testing.T) {
    tests := []string{
        "https://www.dailymotion.com/video/x123456",
        "https://dai.ly/x123456",
    }
    
    for _, url := range tests {
        if !IsOEmbedSupported(url) {
            t.Errorf("Expected URL to be supported: %s", url)
        }
    }
}
```

#### 4. Test Manually

```go
package main

import (
    "fmt"
    "log"
    "github.com/alfarisi/urlmeta"
)

func main() {
    metadata, err := urlmeta.Extract("https://www.dailymotion.com/video/x123456")
    if err != nil {
        log.Fatal(err)
    }
    
    if metadata.OEmbed == nil {
        log.Fatal("oEmbed not extracted")
    }
    
    fmt.Println("✅ Success!")
    fmt.Println("Title:", metadata.Title)
    fmt.Println("Type:", metadata.OEmbed.Type)
}
```

#### 5. Update README.md

Add provider to the supported providers table:

```markdown
| Provider | Domains |
|----------|---------|
| Dailymotion | `dailymotion.com`, `dai.ly` |
```

#### 6. Commit

```bash
git add providers.go providers_test.go README.md
git commit -m "feat: add Dailymotion oEmbed support"
```

### Provider Requirements

**Add providers that:**
- ✅ Have official oEmbed endpoint
- ✅ Are publicly accessible (no auth required)
- ✅ Are widely used
- ✅ Don't require API keys for basic usage

**Do NOT add:**
- ❌ Providers requiring OAuth
- ❌ Private/internal services
- ❌ Providers with severe rate-limiting
- ❌ Defunct/deprecated services

### Popular Providers to Add

Not yet included:
- Dailymotion
- Twitch
- Giphy
- CodePen
- SlideShare

Check https://oembed.com/providers.json for complete list.

## Code Quality Standards

### Go Style Guidelines

- Follow standard Go style guidelines
- Use `gofmt` and `goimports`
- Write clear, idiomatic Go code
- Keep functions focused and small
- Maximum cyclomatic complexity: 15

### Naming Conventions

```go
// ✅ Good
func ExtractMetadata(url string) (*Metadata, error)
var knownProviders = []OEmbedProvider{...}
type OEmbedProvider struct {...}

// ❌ Bad
func extract_metadata(url string) (*Metadata, error)
var KnownProviders = []OEmbedProvider{...}
type oembedProvider struct {...}
```

### Error Handling

```go
// ✅ Good - Wrap errors with context
if err != nil {
    return nil, fmt.Errorf("failed to fetch URL: %w", err)
}

// ❌ Bad - Lose error context
if err != nil {
    return nil, err
}
```

### Comments

```go
// ✅ Good - Explain why
// Skip oEmbed discovery to avoid extra HTTP request
if c.strategy == StrategyHTMLOnly {
    return c.extractHTMLOnly(targetURL, parsedURL)
}

// ❌ Bad - State the obvious
// Check if strategy is HTML only
if c.strategy == StrategyHTMLOnly {
    return c.extractHTMLOnly(targetURL, parsedURL)
}
```

### Documentation

- Add godoc comments for all exported types and functions
- Include examples in documentation where helpful
- Keep comments clear and concise
- Update README.md for user-facing changes

Example:
```go
// ExtractWithRetry extracts metadata with automatic retry on failure.
// It retries up to maxRetries times with exponential backoff.
//
// Example:
//
//	metadata, err := ExtractWithRetry("https://example.com", 3)
//	if err != nil {
//	    log.Fatal(err)
//	}
func ExtractWithRetry(url string, maxRetries int) (*Metadata, error) {
    // Implementation...
}
```

## Testing Guidelines

### Test Requirements

- Write table-driven tests where appropriate
- Aim for >80% test coverage
- Include both positive and negative test cases
- Test edge cases and error conditions
- Use meaningful test names

### Test Example

```go
func TestExtractMetadata(t *testing.T) {
    tests := []struct {
        name      string
        url       string
        wantTitle string
        wantErr   bool
    }{
        {
            name:      "valid YouTube URL",
            url:       "https://youtube.com/watch?v=123",
            wantTitle: "Test Video",
            wantErr:   false,
        },
        {
            name:    "invalid URL",
            url:     "not-a-url",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Extract(tt.url)
            if (err != nil) != tt.wantErr {
                t.Errorf("Extract() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !tt.wantErr && got.Title != tt.wantTitle {
                t.Errorf("Extract() title = %v, want %v", got.Title, tt.wantTitle)
            }
        })
    }
}
```

### Edge Cases to Test

Always test:
- Empty strings
- Nil pointers
- Zero values
- Very large inputs
- Malformed data
- Network errors
- Timeouts

### Benchmarks

Add benchmarks for performance-critical code:

```go
func BenchmarkExtract(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Extract("https://example.com")
    }
}
```

## Pull Request Process

### Before Submitting

Checklist:
- [ ] All tests pass (`go test ./...`)
- [ ] Code is formatted (`go fmt ./...`)
- [ ] Linter passes (`golangci-lint run`)
- [ ] Documentation is updated
- [ ] Commits are clear and descriptive

### Commit Messages

Follow conventional commits format:

```
<type>: <subject>

<body>

<footer>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `test`: Adding or updating tests
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `chore`: Maintenance tasks
- `style`: Code style changes

**Examples:**
```
feat: add Dailymotion oEmbed support

- Add Dailymotion to providers.go
- Add test cases for Dailymotion URLs
- Update README with new provider

Closes #42
```

```
fix: handle nil pointer in image processing

The image dimension processor could crash when images array is empty.
Added nil checks before accessing array elements.

Fixes #38
```

### PR Template

When creating a pull request, include:

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update
- [ ] New oEmbed provider

## Testing
How to test these changes

## Checklist
- [ ] Tests pass locally
- [ ] Code is formatted
- [ ] Documentation updated
- [ ] Added tests for new code
```

### Code Review

All submissions require review. Review criteria:
- Code quality and style
- Test coverage (>80%)
- Documentation completeness
- Performance impact
- Backward compatibility
- Security considerations

## Project Structure

```
urlmeta/
├── urlmeta.go           # Main metadata extraction
├── urlmeta_test.go      # Main tests
├── oembed.go            # oEmbed logic
├── oembed_test.go       # oEmbed tests
├── providers.go         # ⭐ ADD PROVIDERS HERE!
├── providers_test.go    # Provider tests
├── README.md
├── LICENSE
├── go.mod
├── go.sum
├── docs/
│   └── CONTRIBUTING.md  # This file
└── examples/
    ├── basic/
    ├── advanced/
    └── batch/
```

## Common Tasks Quick Reference

### Add oEmbed Provider
```bash
1. Edit providers.go - add to knownProviders
2. Add tests in providers_test.go
3. Update README.md supported providers table
4. Test manually
5. Commit: "feat: add [Provider] oEmbed support"
```

### Fix Bug
```bash
1. Create issue (if not exists)
2. Write failing test
3. Fix the bug
4. Verify test passes
5. Commit: "fix: [description]"
6. Reference issue in commit footer
```

### Update Documentation
```bash
1. Make changes to docs/
2. Verify markdown renders correctly
3. Check all links work
4. Commit: "docs: [description]"
```

## Performance Considerations

### Memory Allocation

```go
// ✅ Good - Preallocate
images := make([]Image, 0, 10)

// ❌ Bad - Let it grow
images := []Image{}
```

### String Operations

```go
// ✅ Good - Use strings.Builder
var b strings.Builder
for _, s := range strs {
    b.WriteString(s)
}
result := b.String()

// ❌ Bad - Multiple allocations
result := ""
for _, s := range strs {
    result += s
}
```

### HTTP Clients

```go
// ✅ Good - Reuse client
client := &http.Client{}
for _, url := range urls {
    resp, _ := client.Get(url)
}

// ❌ Bad - New client each time
for _, url := range urls {
    resp, _ := http.Get(url)
}
```

## Security Considerations

Always validate:
- URL format and protocol (http/https only)
- Response content-type
- Response size (prevent DoS)
- Redirect limits

Never:
- Log sensitive URLs
- Include auth tokens in errors
- Trust user input without validation

## Getting Help

- 🐛 [Issues](https://github.com/alfarisi/urlmeta/issues) - Report bugs

## Recognition

Contributors will be listed in README.md and release notes.

## License

By contributing, you agree that your contributions will be licensed under the Zero-Clause BSD (0BSD) License.

---

**Thank you for contributing to URLMeta!** 🎉