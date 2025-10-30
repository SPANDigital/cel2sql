# Index Analysis Example

This example demonstrates how to use `cel2sql.AnalyzeQuery()` to get index recommendations for your CEL queries.

## What It Does

The `AnalyzeQuery()` function analyzes CEL expressions and recommends PostgreSQL indexes that would optimize query performance. It detects patterns like:

- **B-tree indexes** for comparison operations (`>, <, =, etc.`)
- **GIN indexes** for JSON/JSONB path operations and array operations
- **GIN indexes with pg_trgm** for regex pattern matching

## Running the Example

```bash
cd examples/index_analysis
go run main.go
```

## Example Output

```
cel2sql - Index Analysis Example
===================================

1. Simple Comparison
   Description: Basic comparison operations that benefit from B-tree indexes
   CEL Expression: users.age > 18 && users.active == true

   Generated SQL:
   users.age > 18 AND users.active IS TRUE

   Index Recommendations (1):
   [1] Column: users.age
       Type: BTREE
       Reason: Comparison operations on 'users.age' benefit from B-tree index...
       SQL: CREATE INDEX idx_users_age_btree ON table_name (users.age);

2. Regex Matching
   Description: Pattern matching that benefits from GIN index with pg_trgm
   CEL Expression: users.email.matches(r"@example\.com$")

   Generated SQL:
   users.email ~ '@example\.com$'

   Index Recommendations (1):
   [1] Column: users.email
       Type: GIN
       Reason: Regex matching on 'users.email' benefits from GIN index with pg_trgm...
       SQL: CREATE INDEX idx_users_email_gin_trgm ON table_name USING GIN (users.email gin_trgm_ops);
```

## Key Concepts

### Index Types

- **B-tree** (default): Best for equality and range queries on ordered data
- **GIN** (Generalized Inverted Index): Best for complex data types (JSON, arrays, text search)
- **GIN with pg_trgm**: Optimized for pattern matching and fuzzy text search

### When to Use

- **Development**: Discover what indexes your application needs
- **Performance tuning**: Identify missing indexes causing slow queries
- **Production monitoring**: Analyze user-generated filter expressions to understand usage patterns

### Applying Recommendations

1. Review each recommendation's reason to understand the performance benefit
2. Replace `table_name` with your actual table name in the CREATE INDEX statement
3. Execute the CREATE INDEX statements on your database
4. Use `EXPLAIN ANALYZE` to verify the indexes are being used

```sql
-- Example: Applying a recommendation
CREATE INDEX idx_users_age_btree ON users (users.age);

-- Verify index usage
EXPLAIN ANALYZE SELECT * FROM users WHERE users.age > 18;
```

## Notes

- Index recommendations are based on query patterns, not data distribution
- Consider your specific use case, data size, and query frequency when applying indexes
- Too many indexes can slow down INSERT/UPDATE operations
- Use PostgreSQL's query planner (`EXPLAIN`) to validate index effectiveness
