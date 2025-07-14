# Changelog

## [2.6.1] - 2025-07-14

### Improved
- **Code Architecture**: Refactored large `cel2sql.go` file into logical, maintainable modules for better code organization
- **Modular Design**: Created dedicated modules for specific concerns:
  - `json.go` - JSON/JSONB handling functions and constants (268 lines)
  - `operators.go` - Operator mappings (26 lines)  
  - `timestamps.go` - Timestamp and duration handling (207 lines)
  - `utils.go` - Utility and type-checking functions (107 lines)
- **Maintainability**: Reduced main `cel2sql.go` from ~1,700 lines to 1,094 lines with focused responsibilities
- **Developer Experience**: Improved code navigation and readability for better maintenance
- **Testing**: Maintained 100% backward compatibility with all existing tests passing
- **Code Quality**: Clean linting with 0 issues and 60.5% code coverage maintained

### Technical Details
- Extracted JSON/JSONB-related functionality to dedicated module
- Separated operator mappings for cleaner organization
- Isolated timestamp/duration logic for better maintainability
- Centralized utility functions for code reuse
- Preserved all existing functionality with zero breaking changes
- Enhanced code modularity while maintaining API compatibility

## [2.6.0] - 2025-07-14

### Added
- 🔥 **JSON/JSONB Comprehensions Support**: Full support for CEL comprehensions on JSON/JSONB arrays
- **Advanced JSON Array Operations**: Support for `exists()`, `all()`, `exists_one()` on JSON/JSONB arrays
- **Numeric JSON Field Casting**: Automatic casting of numeric JSON fields (e.g., `(score)::numeric`)
- **Nested JSON Array Access**: Support for comprehensions on nested JSON arrays (e.g., `settings.permissions`)
- **JSON Type Safety**: Null and type checks for JSON/JSONB comprehensions using `jsonb_typeof()`
- **Mixed JSON/JSONB Support**: Proper handling of both JSON and JSONB column types
- **Complex JSON Queries**: Support for complex expressions combining multiple comprehensions

### Enhanced
- **JSON Array Function Selection**: Intelligent selection between `jsonb_array_elements_text` and `json_array_elements_text`
- **JSON Path Operations**: Enhanced nested JSON access with proper `->` and `->>` operators
- **Comprehension Type Detection**: Improved detection of JSON vs regular array comprehensions
- **SQL Generation**: Optimized SQL generation for JSON/JSONB array operations
- **Error Handling**: Better error messages for JSON/JSONB comprehension issues

### Technical Details
- Added `isJSONArrayField()` function to detect JSON/JSONB array fields
- Added `getJSONArrayFunction()` to select appropriate PostgreSQL JSON array functions
- Added `isNestedJSONAccess()` for handling nested JSON field access
- Added `needsNumericCasting()` for automatic numeric casting in JSON comprehensions
- Enhanced `visitAllComprehension()`, `visitExistsComprehension()`, `visitExistsOneComprehension()` with JSON support
- Added comprehensive test suite with real PostgreSQL JSON/JSONB data
- Fixed TODO comment: "Comprehensions are now supported (all, exists, exists_one, filter, map)"

### Examples
```sql
-- CEL: json_users.tags.exists(tag, tag == "developer")
-- SQL: EXISTS (SELECT 1 FROM jsonb_array_elements_text(json_users.tags) AS tag WHERE json_users.tags IS NOT NULL AND jsonb_typeof(json_users.tags) = 'array' AND tag = 'developer')

-- CEL: json_users.scores.all(score, score > 70)
-- SQL: NOT EXISTS (SELECT 1 FROM jsonb_array_elements_text(json_users.scores) AS score WHERE json_users.scores IS NOT NULL AND jsonb_typeof(json_users.scores) = 'array' AND NOT ((score)::numeric > 70))

-- CEL: json_users.attributes.exists_one(attr, attr.skill == "JavaScript" && attr.level >= 9)
-- SQL: (SELECT COUNT(*) FROM json_array_elements(json_users.attributes) AS attr WHERE json_users.attributes IS NOT NULL AND json_typeof(json_users.attributes) = 'array' AND attr->>'skill' = 'JavaScript' AND (attr->>'level')::numeric >= 9) = 1
```

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
