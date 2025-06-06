# Dagger Module Migration

This directory contains the migrated Dagger module for the booklit project.

## Migration Summary

The original Dagger module was written using the old "Environment" system from the primordial stages of Dagger. It has been migrated to the modern Dagger module structure.

### What Changed

#### Old Structure (Environment-based)
- Used `Environment().WithCheck()` and `Environment().Serve()` pattern
- Global `dag` variable was managed manually
- Used outdated Dagger SDK (v0.7.2 with custom replacements)
- Required manual session management and generated bindata

#### New Structure (Modern Dagger Module)
- Functions are exposed directly as Dagger functions  
- No `main()` function required
- Uses modern Dagger SDK (v0.18.9)
- Simplified dependency management
- Leverages Dagger's automatic code generation

### Functions Available

1. **Build** - Builds the booklit binary with optional version
2. **Test** - Runs all Go tests  
3. **Lint** - Runs golangci-lint against the codebase
4. **All** - Runs tests, linting, and builds (comprehensive CI)

### Usage

```bash
# Run tests
dagger call test

# Build with version  
dagger call build --version="v1.0.0"

# Run full CI pipeline
dagger call all --version="v1.0.0"

# Lint code
dagger call lint
```

### Development

To work with this module:

1. Install Dagger CLI
2. Run `dagger develop` to generate necessary code
3. Edit `main.go` to add new functions
4. Test with `dagger call <function-name>`

## Legacy Code

The old implementation used these patterns that are no longer needed:

```go
// OLD - Environment pattern
func main() {
    dag.Environment().
        WithCheck(Unit).
        WithCheck(Lint).
        WithCommand(Build).
        Serve()
}

// NEW - Direct function exposure  
func (m *Booklit) Build(ctx context.Context, version string) (*Directory, error) {
    // Implementation
}
```

The new approach is much cleaner and follows modern Dagger patterns.