# Changelog

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
