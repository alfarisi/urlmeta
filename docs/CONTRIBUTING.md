# Contributing to URLMeta

Thank you for your interest in contributing to URLMeta! We welcome contributions from the community.

## Table of Contents

- [How to Contribute](#how-to-contribute)
- [Development Setup](#development-setup)
- [Adding oEmbed Providers](#adding-oembed-providers)
- [Code Quality](#code-quality)
- [Testing](#testing)
- [Documentation](#documentation)
- [Pull Request Process](#pull-request-process)

## How to Contribute

### Reporting Bugs

If you find a bug, please open an issue with:
- A clear, descriptive title
- Steps to reproduce the issue
- Expected vs actual behavior
- Your Go version and OS
- Example code if applicable
- URL that's causing the issue (if not sensitive)

### Suggesting Enhancements

We welcome enhancement suggestions! Please open an issue with:
- A clear description of the enhancement
- Why this enhancement would be useful
- Example use cases
- Proposed API (if applicable)

### Adding New oEmbed Providers

**This is the most common contribution!** See [Adding oEmbed Providers](#adding-oembed-providers) below.

### Pull Requests

1. **Fork the repository** and create your branch from `main`
2. **Make your changes** following our coding standards
3. **Add tests** for any new functionality
4. **Update documentation** if needed
5. **Run tests** to ensure everything passes
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

# Install development tools
make install-tools

# Or manually:
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run tests with race detector
make test-race

# Run benchmarks
make bench
```

### Code Quality

```bash
# Format code
make fmt

# Run linter
make lint

# Run go vet
make vet

# Run all checks
make check
```

## Adding oEmbed Providers

**Most valuable contribution!** Adding a new oEmbed provider is easy.

### Step-by-Step Guide

#### 1. Check Provider Support

Visit https://oembed.com/providers.json and find your provider.

Example for Dailymotion:
```json
{
  "provider_name": "Dailymotion",
  "provider_url": "https://www.dailymotion.com",
  "endpoints": [{
    "schemes": [
      "https://www.dailymotion.com/video/*"
    ],
    "url": "https://www.dailymotion.com/services/oembed"
  }]
}
```

#### 2. Add to `providers.go`

Open `providers.go` and add your provider to the `knownProviders` array:

```go
var knownProviders = []OEmbedProvider{
    // ... existing providers ...
    
    {
        Name: "Dailymotion",
        URL:  "https://www.dailymotion.com",
        Endpoints: []OEmbedEndpoint{
            {
                Schemes: []string{
                    "https://www.dailymotion.com/video/*",
                    "https://dai.ly/*", // Short URL
                },
                URL:       "https://www.dailymotion.com/services/oembed",
                Discovery: true,
            },
        },
    },
}
```

**Tips:**
- Add all URL schemes (including short URLs)
- Set `Discovery: true` if provider supports it
- Keep alphabetical order (optional, but nice)

#### 3. Update "Last updated" Date

In `providers.go`, update the comment:
```go
// Last updated: 2025-01-26
```

#### 4. Add Tests

Add test cases in `providers_test.go` or `oembed_test.go`:

```go
func TestDailymotionSupport(t *testing.T) {
    tests := []string{
        "https://www.dailymotion.com/video/x123456",
        "https://dai.ly/x123456",
    }
    
    for _, url := range tests {
        if !IsOEmbedSupported(url) {
            t.Errorf("Expected Dailymotion URL to be supported: %s", url)
        }
    }
}
```

#### 5. Test Manually

```bash
# Create a test file
cat > test_dailymotion.go << 'EOF'
package main

import (
    "fmt"
    "log"
    "github.com/alfarisi/urlmeta"
)

func main() {
    url := "https://www.dailymotion.com/video/x123456"
    
    if !urlmeta.IsOEmbedSupported(url) {
        log.Fatal("Not supported")
    }
    
    metadata, err := urlmeta.Extract(url)
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
EOF

go run test_dailymotion.go
```

#### 6. Run Tests

```bash
# Run all tests
go test -v ./...

# Test specific provider
go test -v -run TestDailymotion
```

#### 7. Update Documentation

Add provider to README.md supported providers list:

```markdown
### oEmbed Support (⚡ **Auto-detected!**)
- **YouTube** - `youtube.com`, `youtu.be`
- **Vimeo** - `vimeo.com`
- **Dailymotion** - `dailymotion.com`, `dai.ly` ← NEW!
- ...
```

#### 8. Commit

```bash
git add providers.go providers_test.go README.md
git commit -m "feat: add Dailymotion oEmbed support"
```

### Provider Requirements

Only add providers that:
- ✅ Have official oEmbed endpoint
- ✅ Are publicly accessible (no auth required)
- ✅ Are widely used
- ✅ Don't require API keys (for basic usage)

**Do NOT add:**
- ❌ Providers requiring OAuth
- ❌ Private/internal services (use `AddCustomProvider` instead)
- ❌ Providers with rate-limiting issues
- ❌ Defunct/deprecated services

### Common Providers to Add

Popular providers not yet included:
- Dailymotion
- Twitch
- Giphy
- CodePen
- SlideShare

Check https://oembed.com/providers.json for more.

## Code Quality

### Go Style

- Follow standard Go style guidelines
- Use `gofmt` and `goimports`
- Write clear, idiomatic Go code
- Keep functions focused and small
- Maximum cyclomatic complexity: 15

### Documentation

- Add godoc comments for all exported types and functions
- Include examples in documentation where helpful
- Keep comments clear and concise
- Update README.md for user-facing changes

### Example Documentation

```go
// ExtractWithRetry extracts metadata with automatic retry on failure.
// It will retry up to maxRetries times with exponential backoff.
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

## Testing

### Test Requirements

- Write table-driven tests where appropriate
- Aim for >80% test coverage
- Include both positive and negative test cases
- Test edge cases and error conditions
- Use meaningful test names

### Example Test

```go
func TestExtractMetadata(t *testing.T) {
    tests := []struct {
        name        string
        url         string
        wantTitle   string
        wantErr     bool
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

### Benchmarks

Add benchmarks for performance-critical code:

```go
func BenchmarkExtract(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _, _ = Extract("https://example.com")
    }
}
```

## Documentation

### Update These Files

When making changes, update:
- `README.md` - User-facing changes
- `docs/API.md` - API changes
- Code comments - Implementation details

### Documentation Style

- Use Markdown
- Include code examples
- Keep it concise
- Use tables for comparisons
- Add links to related docs

## Pull Request Process

### Before Submitting

1. ✅ All tests pass (`make test`)
2. ✅ Linter passes (`make lint`)
3. ✅ Code is formatted (`make fmt`)
4. ✅ Documentation is updated
5. ✅ Commits are clear and descriptive

### PR Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update
- [ ] New oEmbed provider

## Checklist
- [ ] Tests pass locally
- [ ] Linter passes
- [ ] Documentation updated
- [ ] Added tests for new code
- [ ] Commits are clear

## Testing
How to test these changes

## Screenshots (if applicable)
```

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
```bash
feat: add Dailymotion oEmbed support

- Add Dailymotion to providers.go
- Add test cases for Dailymotion URLs
- Update README with new provider

Closes #42

---

fix: handle nil pointer in image processing

The image dimension processor could crash when images array is empty.
Added nil checks before accessing array elements.

Fixes #38

---

docs: update API documentation for strategies

- Document StrategyAuto, StrategyOEmbedFirst, StrategyHTMLOnly
- Add performance comparison table
- Include usage examples
```

### Code Review Process

All submissions require review. We use GitHub pull requests.

**Review criteria:**
- Code quality and style
- Test coverage (>80%)
- Documentation completeness
- Performance impact
- Backward compatibility
- Security considerations

### After Approval

1. Squash commits if needed
2. Merge to `main` branch
3. Delete feature branch
4. Release notes updated (by maintainers)

## Project Structure

```
urlmeta/
├── urlmeta.go           # Main metadata extraction
├── urlmeta_test.go      # Main tests
├── oembed.go            # oEmbed logic
├── oembed_test.go       # oEmbed tests
├── providers.go         # ← ADD PROVIDERS HERE!
├── providers_test.go    # Provider tests
├── README.md
├── LICENSE
├── Makefile
├── .golangci.yml
├── docs/
│   ├── API.md
│   ├── CONTRIBUTING.md  # This file
└── examples/
    ├── basic/
    ├── advanced/
    └── batch/
```

## Common Tasks

### Add oEmbed Provider
```bash
1. Edit providers.go
2. Add to knownProviders array
3. Add tests in providers_test.go
4. Update README.md
5. Commit: "feat: add [Provider] oEmbed support"
```

### Fix Bug
```bash
1. Create issue (if not exists)
2. Write failing test
3. Fix the bug
4. Verify test passes
5. Commit: "fix: [description]"
6. Reference issue in commit
```

### Improve Performance
```bash
1. Add benchmark (if not exists)
2. Make optimization
3. Run benchmarks (before/after)
4. Document improvement in PERFORMANCE.md
5. Commit: "perf: [description]"
```

### Update Documentation
```bash
1. Make changes to docs/
2. Verify markdown renders correctly
3. Check all links work
4. Commit: "docs: [description]"
```

## Development Workflow

### Typical Workflow

```bash
# 1. Create feature branch
git checkout -b feat/add-dailymotion

# 2. Make changes
vim providers.go
vim providers_test.go
vim README.md

# 3. Test locally
make check

# 4. Commit
git add .
git commit -m "feat: add Dailymotion oEmbed support"

# 5. Push
git push origin feat/add-dailymotion

# 6. Create PR on GitHub
# 7. Address review comments
# 8. Merge when approved
```

### Working with Forks

```bash
# Add upstream remote
git remote add upstream https://github.com/alfarisi/urlmeta.git

# Keep your fork updated
git fetch upstream
git checkout main
git merge upstream/main
git push origin main

# Rebase your feature branch
git checkout feat/add-dailymotion
git rebase main
```

## Coding Standards

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

// ❌ Bad - Generic error
if err != nil {
    return nil, fmt.Errorf("error occurred")
}
```

### Comments

```go
// ✅ Good - Explain why
// Skip oEmbed discovery to avoid extra HTTP request
// Users can still call ExtractOEmbed() explicitly
if c.strategy == StrategyHTMLOnly {
    return c.extractHTMLOnly(targetURL, parsedURL)
}

// ❌ Bad - State the obvious
// Check if strategy is HTML only
if c.strategy == StrategyHTMLOnly {
    return c.extractHTMLOnly(targetURL, parsedURL)
}
```

### Testing Edge Cases

Always test:
- Empty strings
- Nil pointers
- Zero values
- Very large inputs
- Malformed data
- Network errors
- Timeouts

## Performance Considerations

### Memory Allocation

```go
// ✅ Good - Preallocate
images := make([]Image, 0, 10)

// ❌ Bad - Let it grow
images := []Image{}
```

### String Concatenation

```go
// ✅ Good - Use strings.Builder for many concatenations
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

### HTTP Requests

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

### Input Validation

Always validate:
- URL format
- Protocol (http/https only)
- Response content-type
- Response size (prevent DoS)
- Redirect limits

### Sensitive Data

- Never log sensitive URLs
- Don't include auth tokens in errors
- Sanitize user input in examples

## Getting Help

- 💬 [Discussions](https://github.com/alfarisi/urlmeta/discussions) - Ask questions
- 🐛 [Issues](https://github.com/alfarisi/urlmeta/issues) - Report bugs
- 📧 Email - maintainer@example.com (for security issues)

## Recognition

Contributors will be:
- Listed in README.md
- Mentioned in release notes
- Added to AUTHORS file (if significant contributions)

## Code of Conduct

Be respectful, inclusive, and professional. We aim to foster a welcoming community.

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

## Questions?

Feel free to open an issue or discussion if you have questions about contributing!

---

**Thank you for contributing to URLMeta!** 🎉# Contributing to URLMeta

Thank you for your interest in contributing to URLMeta! We welcome contributions from the community.

## How to Contribute

### Reporting Bugs

If you find a bug, please open an issue with:
- A clear, descriptive title
- Steps to reproduce the issue
- Expected vs actual behavior
- Your Go version and OS
- Example code if applicable

### Suggesting Enhancements

We welcome enhancement suggestions! Please open an issue with:
- A clear description of the enhancement
- Why this enhancement would be useful
- Example use cases

### Pull Requests

1. **Fork the repository** and create your branch from `main`
2. **Make your changes** following our coding standards
3. **Add tests** for any new functionality
4. **Update documentation** if needed
5. **Run tests** to ensure everything passes
6. **Submit a pull request**

## Development Setup

### Prerequisites

- Go 1.21 or later
- Git

### Setup

```bash
# Clone your fork
git clone https://github.com/alfarisi/urlmeta.git
cd urlmeta

# Install dependencies
go mod download

# Install development tools
make install-tools
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run tests with race detector
make test-race

# Run benchmarks
make bench
```

### Code Quality

```bash
# Format code
make fmt

# Run linter
make lint

# Run go vet
make vet

# Run all checks
make check
```

## Coding Standards

### Go Style

- Follow standard Go style guidelines
- Use `gofmt` and `goimports`
- Write clear, idiomatic Go code
- Keep functions focused and small

### Documentation

- Add godoc comments for all exported types and functions
- Include examples in documentation where helpful
- Keep comments clear and concise

### Testing

- Write table-driven tests where appropriate
- Aim for high test coverage (>80%)
- Include both positive and negative test cases
- Test edge cases and error conditions

### Example Test

```go
func TestExtractMetadata(t *testing.T) {
    tests := []struct {
        name        string
        url         string
        want        *Metadata
        wantErr     bool
    }{
        {
            name: "valid URL",
            url:  "https://example.com",
            want: &Metadata{
                Title: "Example Domain",
            },
            wantErr: false,
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
            // Add more assertions...
        })
    }
}
```

## Commit Messages

Write clear commit messages following this format:

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

**Example:**
```
feat: add support for Schema.org microdata

- Parse itemprop attributes
- Extract structured data
- Add tests for microdata extraction

Closes #123
```

## Project Structure

```
urlmeta/
├── urlmeta.go           # Main package code
├── urlmeta_test.go      # Tests
├── examples/            # Usage examples
│   ├── basic/
│   ├── advanced/
│   └── batch/
├── docs/                # Documentation
└── .github/             # GitHub configs
```

## Adding New Features

When adding a new feature:

1. **Discuss first**: Open an issue to discuss the feature
2. **Design**: Consider the API design and user experience
3. **Implement**: Write clean, well-tested code
4. **Document**: Update README and add examples
5. **Test**: Ensure comprehensive test coverage

### Example: Adding New Meta Tag Support

```go
// 1. Update the Metadata struct
type Metadata struct {
    // ... existing fields
    NewField string `json:"new_field,omitempty"`
}

// 2. Add extraction logic
func processMeta(n *html.Node, metadata *Metadata, baseURL *url.URL) {
    // ... existing code
    case "new-meta-tag":
        metadata.NewField = content
}

// 3. Add tests
func TestExtractNewField(t *testing.T) {
    // ... test implementation
}

// 4. Update documentation
// Add to README.md and examples
```

## Code Review Process

All submissions require review. We use GitHub pull requests for this purpose.

**Review criteria:**
- Code quality and style
- Test coverage
- Documentation completeness
- Performance impact
- Backward compatibility

## Questions?

Feel free to open an issue if you have questions about contributing!

## License

By contributing, you agree that your contributions will be licensed under the MIT License.