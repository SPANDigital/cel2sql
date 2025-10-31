# Changelog

## [Unreleased]

## [3.4.0] - 2025-10-31

### Added
- **CEL String Extension Functions: split(), join(), format()** (fixes #87)
  - Implemented `split(delimiter [, limit])` → `STRING_TO_ARRAY()` with full limit support
    * Basic splitting: `"a,b,c".split(",")` → `STRING_TO_ARRAY('a,b,c', ',')`
    * With limit: `"a,b,c".split(",", 2)` → `(STRING_TO_ARRAY('a,b,c', ','))[1:2]`
    * Limit=0: Returns empty array `ARRAY[]::text[]`
    * Limit=1: Returns original string in array (no split)
    * Limit=-1: Unlimited splits (default)
  - Implemented `join([delimiter])` → `ARRAY_TO_STRING()`
    * Basic joining: `["a", "b"].join(",")` → `ARRAY_TO_STRING(ARRAY['a', 'b'], ',', '')`
    * No delimiter: `["a", "b"].join()` → Uses empty string delimiter
    * Null handling: Replaces nulls with empty strings
  - Implemented `format(args)` → `FORMAT()` with specifier support
    * Supported specifiers: `%s` (string), `%d` (decimal), `%f` (float)
    * Unsupported specifiers return clear errors (`%b`, `%x`, etc.)
    * Format string must be constant (max 1000 chars)
    * Arguments must be constant list
    * Argument count validation
  - All functions work in comprehensions (exists, all, filter, map)
  - Security validations: null byte checks, format string limits, specifier whitelisting
  - Comprehensive test suite with 100+ test cases
  - Performance benchmarks added
  - Working example in `examples/string_extensions/`
  - Note: `quote()` not implemented (not part of CEL `ext.Strings()` standard extension)

## [3.3.0] - 2025-10-31

### Fixed
- **LIKE Pattern Escaping** (fixes #40, #43, #84)
  - Fixed LIKE pattern escaping to use ESCAPE E'\\' syntax for proper backslash handling
  - Updated startsWith() and endsWith() to properly escape special LIKE characters (%, _, \)
  - Prevents SQL syntax errors from patterns containing backslashes

- **JSON Comprehensions** (fixes #48, #84)
  - Fixed comprehensions over JSON arrays to properly use jsonb_array_elements()
  - Corrects SQL generation for expressions like `data.items.all(i, i.quantity > 0)`
  - Ensures JSON array comprehensions work correctly with PostgreSQL

- **String Functions Panic** (fixes #85, #86)
  - Fixed panic: index out of range when using CEL string extension functions as methods
  - Added defensive checks in callCasting, visitCallIndex, visitCallMapIndex, visitCallListIndex, visitCallUnary
  - All string functions now properly handle both method calls (target) and function calls (args)

### Added
- **CEL String Extension Functions** (#86)
  - Implemented 10 CEL string extension functions with PostgreSQL SQL conversion:
    * `lowerAscii()` → `LOWER()`
    * `upperAscii()` → `UPPER()`
    * `trim()` → `TRIM()`
    * `charAt(index)` → `SUBSTRING(str, index+1, 1)`
    * `indexOf(search, [offset])` → `POSITION()` with -1 for not found
    * `lastIndexOf(search)` → Uses `REVERSE()` logic
    * `substring(start, [end])` → `SUBSTRING()` with proper index conversion
    * `replace(old, new, [limit])` → `REPLACE()` (limit=-1 only)
    * `reverse()` → `REVERSE()`
  - Clear error messages for unsupported functions (split, join, format, quote)
  - Comprehensive test coverage in string_functions_test.go

## [3.2.0] - 2025-10-31

### Fixed
- **Standardized Error Message Format** (fixes #38, #83)
  - Added 16 exported sentinel errors (`ErrUnsupportedExpression`, `ErrInvalidFieldName`, etc.) to enable `errors.Is()` checking
  - Improved error wrapping with consistent `fmt.Errorf()` and `%w` patterns across ~60 error sites
  - Enhanced error messages with operation context for better debugging
  - Maintained security-conscious error handling (no credential exposure in pg package)
  - Benefits: Better debugging, programmatic error handling, improved maintainability

- **Byte Array Length Validation** (fixes #36, #82)
  - Added maximum byte array length limit of 10,000 bytes to prevent resource exhaustion (CWE-400)
  - Protection applies to non-parameterized mode (hex encoding causes 4x expansion)
  - Parameterized mode bypasses limit (bytes passed directly to database driver)
  - Clear error messages guide users when limits are exceeded

### Added
- **Comprehensive Performance Benchmarks in CI/CD** (fixes #52)
  - Automated benchmark tracking with historical data storage on GitHub Pages
  - Visual performance charts at https://spandigital.github.io/cel2sql/dev/bench/
  - PR comments when performance changes exceed 150%
  - Benchmarks cover all major features: operators, comprehensions, JSON, regex, timestamps

## [3.1.0] - 2025-10-30

### Added
- **Query Analysis and Index Recommendations** (fixes #50, #81)
  - New `AnalyzeQuery()` function to analyze CEL expressions and recommend database indexes
  - Automatic detection of B-tree, GIN, and GIN with pg_trgm index opportunities
  - Support for JSON path operations, array operations, and regex pattern matching
  - Complete working example in `examples/index_analysis/`
  - Benefits: Discover missing indexes, optimize query performance, improve production monitoring

### Documentation
- **Regex Conversion Limitations** (fixes #46, #80)
  - Added detailed documentation of automatic RE2 to POSIX regex conversions
  - Listed unsupported RE2 features that will return errors (lookahead, lookbehind, named groups)
  - Added examples of supported character class conversions (`\d`, `\w`, `\s`, `\b`)
  - Clarified case-insensitive flag handling with `~*` operator

## [3.0.0] - 2025-10-30

### Breaking Changes
- **Schema API Redesign for O(1) Performance** (fixes #28)
  - Changed `pg.Schema` from type alias `[]FieldSchema` to a struct with internal indexing
  - Schema field lookups now O(1) constant time instead of O(n) linear search
  - **Migration Guide**:
    - Old: `schema := pg.Schema{{Name: "field", Type: "text"}}`
    - New: `schema := pg.NewSchema([]pg.FieldSchema{{Name: "field", Type: "text"}})`
    - Use `schema.Fields()` to iterate fields
    - Use `schema.FindField(name)` for O(1) lookups
  - Benchmarks show constant 241ns lookup time regardless of schema size (10, 100, or 1000 fields)
  - All examples, tests, and documentation updated

### Performance
- **10-100x faster schema lookups** for large schemas
  - 10 fields: ~same performance
  - 100 fields: ~10x faster
  - 1000 fields: ~100x faster
  - Critical for applications with complex database schemas

### Security
- **Nested Comprehension Depth Limits**: Added protection against resource exhaustion from deeply nested comprehensions (fixes #35)
  - Implemented maximum nesting depth of 3 for CEL comprehensions (all, exists, exists_one, filter, map)
  - Prevents DoS attacks through expensive nested UNNEST/subquery operations (CWE-400)
  - Clear error messages guide users to restructure overly complex queries
  - Example protected pattern: `list1.map(x, list2.filter(y, list3.exists(z, z > y)))` now limited to depth 3

## [2.10.0] - 2025-10-10

### Added
- **Comprehensive Fuzzing Infrastructure**: Go native fuzzing for security and reliability testing
  - Three specialized fuzz tests: `FuzzConvert`, `FuzzEscapeLikePattern`, `FuzzConvertRE2ToPOSIX`
  - Fuzzing dictionary with 100+ tokens (CEL operators, SQL injection vectors, regex patterns, Unicode)
  - CI/CD integration: 60s on PRs, 10-minute weekly runs, manual trigger support
  - Automatic crash artifact preservation and regression testing
  - Successfully discovered 3 security vulnerabilities within minutes of initial runs

### Fixed
- **Critical Security Fixes**: Null byte handling vulnerabilities discovered by fuzzing
  - Fixed null bytes in string literals causing SQL corruption (#14)
  - Fixed null bytes in LIKE patterns for `startsWith()` and `endsWith()` functions (#16)
  - Fixed null bytes in regex patterns for `matches()` function (#18)
  - All three fixes prevent null bytes (`\x00`) from reaching PostgreSQL queries
- **Panic to Error**: Replaced panic with proper error handling in timestamp operation validation

### Documentation
- Completely rewrote README.md for better first-time user experience
- Added beginner-friendly examples and clear getting started guide
- Improved API documentation with real-world use cases

### Security
This release includes critical security fixes that prevent SQL corruption from null byte injection. All users are strongly encouraged to upgrade.

## [2.9.0] - 2025-10-09

### Fixed
- **PostgreSQL 17 Compatibility**: Comprehensive fix for all PostgreSQL 17 incompatible functions
  - Fixed ternary operator (`? :`) to generate `CASE WHEN ... THEN ... ELSE ... END` instead of `IF()`
  - Fixed Unix timestamp conversion to use `EXTRACT(EPOCH FROM ...)::bigint` instead of `UNIX_SECONDS()`
  - Fixed type casting to use correct PostgreSQL type names:
    - `BOOL` → `BOOLEAN`
    - `BYTES` → `BYTEA`
    - `INT64` → `BIGINT`
    - `FLOAT64` → `DOUBLE PRECISION`
    - `STRING` → `TEXT`
  - Fixed date part extraction to use correct PostgreSQL field names:
    - `DAYOFYEAR` → `DOY`
    - `DAYOFWEEK` → `DOW`
    - `MILLISECOND` → `MILLISECONDS`
  - Fixed struct construction to use `ROW()` instead of `STRUCT()`
  - Fixed timezone conversion to use `AT TIME ZONE` operator for 2-argument `timestamp()` calls

### Added
- **PostgreSQL 17 Integration Tests**: Comprehensive test suite validating all fixes against real PostgreSQL 17 database
  - 12 integration tests covering CASE WHEN, type casting, EXTRACT operations, and complex combinations
  - All tests validate both SQL generation correctness and actual query execution results
  - Updated all test containers from PostgreSQL 15 to PostgreSQL 17

### Technical Details
- All generated SQL now validated against PostgreSQL 17 for correctness
- Maintains backward compatibility with existing CEL expressions
- All 100+ unit tests and integration tests pass with PostgreSQL 17

### Examples
```cel
// Ternary operator (conditional)
age > 30 ? "senior" : "junior"                     // → CASE WHEN age > 30 THEN 'senior' ELSE 'junior' END

// Unix timestamp conversion
int(created_at)                                    // → EXTRACT(EPOCH FROM created_at)::bigint

// Type casting
bool(age > 30)                                     // → CAST(age > 30 AS BOOLEAN)
int(score)                                         // → CAST(score AS BIGINT)
string(age)                                        // → CAST(age AS TEXT)

// Date part extraction
created_at.getDayOfYear()                          // → EXTRACT(DOY FROM created_at) - 1
created_at.getDayOfWeek()                          // → EXTRACT(DOW FROM created_at) - 1

// Timezone conversion (2-arg timestamp)
timestamp(created_at, "America/New_York")          // → created_at AT TIME ZONE 'America/New_York'
```

## [2.8.1] - 2025-10-09

### Fixed
- **String Functions**: Fixed `startsWith()` and `endsWith()` functions to generate valid PostgreSQL SQL
  - Previously generated non-existent `STARTS_WITH()` and `ENDS_WITH()` functions
  - Now correctly converts to PostgreSQL `LIKE` patterns:
    - `string.startsWith("prefix")` → `string LIKE 'prefix%'`
    - `string.endsWith("suffix")` → `string LIKE '%suffix'`
  - Added proper escaping for LIKE special characters (`%`, `_`, `\`)
  - Updated all test cases to reflect correct PostgreSQL syntax

### Examples
```cel
// String function fixes
name.startsWith("John")                             // → name LIKE 'John%'
path.endsWith("/module1")                          // → path LIKE '%/module1'
dcType == "Article" && isPartOf.endsWith("/docs") // → dcType = 'Article' AND isPartOf LIKE '%/docs'
```

## [2.8.0] - 2025-07-19

### Added
- **Regex Pattern Matching**: Comprehensive support for CEL `matches()` function with automatic RE2 to POSIX conversion
  - Support for both `string.matches(pattern)` and `matches(string, pattern)` syntax
  - Automatic conversion of RE2 regex patterns to PostgreSQL POSIX ERE format
  - Character class conversions: `\d` → `[[:digit:]]`, `\w` → `[[:alnum:]_]`, `\s` → `[[:space:]]`
  - Word boundary conversion: `\b` → `\y` (PostgreSQL extension)
  - Comprehensive test coverage with 6 unit tests and 7 PostgreSQL integration tests
- **JSON/JSONB Field Existence**: Enhanced `has()` macro support for JSON/JSONB field existence checks
  - Direct field access: `has(table.json_column.field)` using `?` operator for JSONB
  - Nested path existence: `has(table.json_column.nested.field)` using `jsonb_extract_path_text()`
  - Mixed JSON/JSONB support with appropriate operator selection
  - Comprehensive validation with real PostgreSQL database tests

### Examples
```cel
// Regex pattern matching
user.email.matches(".*@example\\.com")          // → user.email ~ '.*@example\.com'
user.code.matches("^[A-Z]{3}\\d{3}$")          // → user.code ~ '^[A-Z]{3}[[:digit:]]{3}$'
user.phone.matches("^\\d{3}-\\d{4}$")          // → user.phone ~ '^[[:digit:]]{3}-[[:digit:]]{4}$'
matches(user.description, "\\btest\\b")         // → user.description ~ '\ytest\y'

// JSON field existence checks
has(metadata.corpus)                            // → metadata ? 'corpus'
has(metadata.corpus.section)                   // → jsonb_extract_path_text(metadata, 'corpus', 'section') IS NOT NULL
has(properties.visibility)                     // → properties->'visibility' IS NOT NULL
```

### Fixed
- **Code Quality**: Resolved all golangci-lint issues for improved code maintainability
  - Simplified conditional logic (replaced else-if patterns)
  - Optimized error handling (use `errors.New` where appropriate)
  - Cleaned unused parameters and improved code flow

### Technical Details
- Uses PostgreSQL's `~` operator for regex pattern matching with POSIX ERE syntax
- Automatic RE2 to POSIX pattern conversion for common regex features
- Enhanced JSON path handling for both JSON and JSONB data types
- All existing functionality remains unchanged and fully backward compatible
- Comprehensive test coverage including real PostgreSQL database validation

## [2.7.2] - 2025-07-19

### BREAKING CHANGES
- **Module Path Update**: Updated module path to `github.com/spandigital/cel2sql/v3` for Go module versioning compliance
- **Import Changes Required**: All users must update their imports to include the `/v2` suffix

### Fixed
- **Go Module Compliance**: Fixed "module contains a go.mod file, so module path must match major version" error
- **Module Versioning**: Proper Go module versioning for v2+ releases according to Go module standards

### Migration Required
Users must update their import statements:

```go
// OLD (v2.7.1 and earlier)
import "github.com/spandigital/cel2sql"
import "github.com/spandigital/cel2sql/pg"
import "github.com/spandigital/cel2sql/sqltypes"

// NEW (v2.7.2 and later)
import "github.com/spandigital/cel2sql/v3"
import "github.com/spandigital/cel2sql/v3/pg"
import "github.com/spandigital/cel2sql/v3/sqltypes"
```

### Technical Details
- Updated `go.mod` module path to include `/v2` suffix
- Updated all internal package imports throughout codebase
- Updated documentation and examples to reflect new import paths
- All functionality remains unchanged - only import paths have changed
- Maintains full backward compatibility at API level

### Note
This change is required for proper Go module versioning compliance. Go modules v2+ must include the major version in the module path when the module contains a go.mod file.

## [2.7.1] - 2025-07-19

### Dependencies
- **Updated Dependencies**: Upgraded core dependencies for security and compatibility
  - Bumped `github.com/google/cel-go` from v0.25.0 to v0.26.0 for improved CEL expression handling
  - Bumped `github.com/testcontainers/testcontainers-go/modules/postgres` to v0.38.0 for enhanced PostgreSQL testing support
- **Improved Compatibility**: All dependency updates maintain backward compatibility
- **Security**: Updated dependencies include latest security patches and improvements

### Technical Details
- No breaking changes - fully backward compatible
- All existing APIs and functionality remain unchanged
- Enhanced performance and reliability from dependency updates
- Maintained 100% test coverage with all tests passing

## [2.7.0] - 2025-01-20

### Added
- 🚀 **Comprehensive JSON/JSONB Nested Path Support**: Full support for deep nested JSON path expressions in PostgreSQL
- **Advanced JSON Path Generation**: Automatic conversion of CEL expressions like `table.metadata.corpus.section` to PostgreSQL JSON paths `table.metadata->'corpus'->>'section'`
- **Intelligent JSON Operator Selection**: Automatic selection between `->` (JSON navigation) and `->>` (text extraction) operators
- **Automatic Type Casting**: Smart detection and casting for numeric comparisons (e.g., `(version.major)::numeric > 1`)
- **Complex Nested Structure Support**: Support for 4+ levels of JSON nesting with proper path generation
- **JSON Array Membership**: Enhanced support for array membership operations with `jsonb_array_elements_text()`
- **Mixed JSON/JSONB Support**: Seamless handling of both JSON and JSONB column types in the same query

### Enhanced
- **JSON Field Detection**: Improved detection of JSON/JSONB fields in table schemas for automatic path generation
- **Nested Path Building**: Recursive path building with proper PostgreSQL JSON operator precedence
- **Test Coverage**: Added comprehensive test suite with 15+ JSON/JSONB nested path test cases
- **Error Handling**: Better error messages for invalid JSON path expressions

### Technical Details
- Enhanced `json.go` with `shouldUseJSONPath()`, `hasJSONFieldInChain()`, and `buildJSONPath()` functions
- Updated `cel2sql.go` with automatic numeric casting detection for JSON text extractions
- Added comprehensive PostgreSQL test data with complex nested JSON/JSONB structures
- Improved type checking and validation for nested JSON field access
- Added utilities for JSON field type detection and path validation

### Testing
- Added `create_json_nested_path_test_data.sql` with comprehensive test scenarios
- Created `TestJSONNestedPathExpressions` with 15 test cases covering various nesting patterns
- Enhanced testcontainer integration for PostgreSQL JSON/JSONB testing
- Verified compatibility with PostgreSQL 15+ JSON features

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
