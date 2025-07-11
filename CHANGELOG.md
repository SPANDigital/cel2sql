# Changelog

## [2.4.0] - 2025-01-11

### Added
- Comprehensive JSON/JSONB support for PostgreSQL columns
- JSON path operations (`->>`) for CEL field access on JSON/JSONB columns
- Support for nested JSON field access in CEL expressions (e.g., `users.preferences.theme`)
- Comprehensive test coverage for JSON/JSONB operations with real data
- Enhanced CEL-to-SQL conversion with PostgreSQL JSON path syntax

### Changed
- Updated SQL generation to automatically detect JSON/JSONB columns and use proper path syntax
- Enhanced type provider to track JSON/JSONB column types for conversion
- Improved test data with realistic JSON structures for comprehensive testing

### Technical Details
- Added `shouldUseJSONPath` function to detect JSON field access patterns
- Enhanced `visitSelect` function to handle JSON path operations
- JSON fields are converted to PostgreSQL `field->>'key'` syntax
- Maintains backward compatibility with existing CEL expressions

## [2.3.0] - 2025-01-11

### Added
- Comprehensive integration tests using testcontainers for PostgreSQL
- Support for array type detection from PostgreSQL information_schema
- Enhanced test coverage for date arithmetic and array operations
- Automated test data generation for comprehensive testing
- Integration tests validating complete CEL-to-SQL-to-results workflow

### Changed
- Updated SQL generation to use PostgreSQL-specific syntax consistently
- Updated string literals to use single quotes (PostgreSQL standard)
- Updated `contains` function to use `POSITION` instead of `CONTAINS`
- Updated array length function to use `ARRAY_LENGTH(..., 1)`
- Updated timestamp handling to use `CAST(..., AS TIMESTAMP WITH TIME ZONE)`
- Improved array type detection to handle PostgreSQL's `ARRAY` type suffix
- Enhanced boolean handling to use `IS TRUE`/`IS FALSE` for PostgreSQL

### Fixed
- Fixed array type detection in `pg/provider.go` to properly handle PostgreSQL array types
- Fixed string literal quoting in SQL generation for PostgreSQL compatibility
- Fixed timestamp function generation for PostgreSQL date/time operations
- Fixed all test expectations to match actual PostgreSQL output
- Fixed CEL boolean handling for proper PostgreSQL boolean operations

### Removed
- Removed MySQL-style backtick quoting from SQL generation
- Removed demo and debug files to clean up the codebase

### Security
- Improved SQL injection prevention through proper quoting
- Enhanced parameterized query support

## [2.2.0] - 2025-01-10

### Added
- Initial comprehensive PostgreSQL migration
- Enhanced type system integration

## [2.1.1] - 2025-01-10

### Fixed
- Bug fixes and improvements

## [2.1.0] - 2025-01-10

### Added
- Enhanced CEL expression support

## [2.0.0] - 2025-07-10

### BREAKING CHANGES
- Migrated from BigQuery to PostgreSQL as the primary database backend
- Removed all BigQuery-specific dependencies and code
- Removed the `bq/` package entirely
- Updated type system to use PostgreSQL-native types

### Added
- PostgreSQL support with modern pgx v5 driver
- New `pg/` package for PostgreSQL type provider integration
- Comprehensive PostgreSQL schema support with composite types and arrays
- Modern security scanning with govulncheck, OSV Scanner, and gosec
- Improved CI/CD pipeline with latest GitHub Actions
- Enhanced error handling and linting configuration

### Changed
- **BREAKING**: All BigQuery-specific APIs have been removed
- **BREAKING**: Type provider now uses PostgreSQL schema format
- **BREAKING**: Database connection handling now uses pgxpool.Pool
- Updated to Go 1.23+ with support for Go 1.24
- Modernized GitHub Actions workflows (checkout@v4, setup-go@v5)
- Simplified golangci-lint configuration with essential linters only
- Improved dependency management and security scanning

### Removed
- **BREAKING**: All `cloud.google.com/go/bigquery` dependencies
- **BREAKING**: BigQuery-specific type mappings and schema handling
- **BREAKING**: `bq.TypeProvider` - replaced with `pg.TypeProvider`
- Deprecated Nancy vulnerability scanner
- Outdated GitHub Actions and workflow configurations

### Fixed
- CEL API deprecation warnings properly suppressed
- Improved error handling in code generation
- Fixed security scanner installation and configuration
- Resolved golangci-lint configuration issues

### Migration Guide
To migrate from v1.x to v2.0.0:

1. **Replace BigQuery imports**:
   ```go
   // OLD
   import "github.com/SPANDigital/cel2sql/bq"
   
   // NEW
   import "github.com/SPANDigital/cel2sql/pg"
   ```

2. **Update type provider usage**:
   ```go
   // OLD
   provider := bq.NewTypeProvider(dataset)
   
   // NEW
   schema := pg.Schema{
       {Name: "field_name", Type: "text", Repeated: false},
       // ... more fields
   }
   provider := pg.NewTypeProvider(map[string]pg.Schema{"TableName": schema})
   ```

3. **Update database connections**:
   ```go
   // NEW - Use pgxpool for PostgreSQL connections
   pool, err := pgxpool.New(context.Background(), "postgresql://...")
   ```

This is a major version release that provides better performance, modern tooling, and PostgreSQL-native support while maintaining the core CEL-to-SQL conversion functionality.
