# Contributing to URLMeta

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