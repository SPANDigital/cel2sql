package cel2sql

import (
	"context"
	"strings"
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/spandigital/cel2sql/v3/pg"
)

func TestAnalyzeQuery_JSONPathOperations(t *testing.T) {
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "id", Type: "text"},
		{Name: "metadata", Type: "jsonb", IsJSON: true, IsJSONB: true},
	})
	provider := pg.NewTypeProvider(map[string]pg.Schema{"person": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("person", cel.ObjectType("person")),
	)
	if err != nil {
		t.Fatalf("failed to create CEL environment: %v", err)
	}

	tests := []struct {
		name           string
		expression     string
		expectedColumn string
		expectedType   string
		expectReason   string
	}{
		{
			name:           "simple JSON path access",
			expression:     `person.metadata.name == "John"`,
			expectedColumn: "person.metadata",
			expectedType:   "GIN",
			expectReason:   "JSON path operations",
		},
		{
			name:           "nested JSON path access",
			expression:     `person.metadata.profile.age > 18`,
			expectedColumn: "person.metadata",
			expectedType:   "GIN",
			expectReason:   "JSON path operations",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.expression)
			if issues != nil && issues.Err() != nil {
				t.Fatalf("failed to compile expression: %v", issues.Err())
			}

			sql, recommendations, err := AnalyzeQuery(ast, WithSchemas(provider.GetSchemas()))
			if err != nil {
				t.Fatalf("AnalyzeQuery failed: %v", err)
			}

			if sql == "" {
				t.Error("expected SQL output, got empty string")
			}

			// Check that we got a recommendation for the expected column
			found := false
			for _, rec := range recommendations {
				if rec.Column == tt.expectedColumn && rec.IndexType == tt.expectedType {
					found = true
					if !strings.Contains(rec.Reason, tt.expectReason) {
						t.Errorf("expected reason to contain %q, got %q", tt.expectReason, rec.Reason)
					}
					if !strings.Contains(rec.Expression, "CREATE INDEX") {
						t.Error("expected CREATE INDEX in expression")
					}
					if !strings.Contains(rec.Expression, "GIN") {
						t.Error("expected GIN in expression")
					}
				}
			}

			if !found {
				t.Errorf("expected recommendation for column %s with type %s, got recommendations: %+v",
					tt.expectedColumn, tt.expectedType, recommendations)
			}
		})
	}
}

func TestAnalyzeQuery_RegexOperations(t *testing.T) {
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "id", Type: "text"},
		{Name: "email", Type: "text"},
	})
	provider := pg.NewTypeProvider(map[string]pg.Schema{"person": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("person", cel.ObjectType("person")),
	)
	if err != nil {
		t.Fatalf("failed to create CEL environment: %v", err)
	}

	ast, issues := env.Compile(`person.email.matches(r"^[a-z]+@example\.com$")`)
	if issues != nil && issues.Err() != nil {
		t.Fatalf("failed to compile expression: %v", issues.Err())
	}

	sql, recommendations, err := AnalyzeQuery(ast, WithSchemas(provider.GetSchemas()))
	if err != nil {
		t.Fatalf("AnalyzeQuery failed: %v", err)
	}

	if sql == "" {
		t.Error("expected SQL output, got empty string")
	}

	// Check that we got a GIN index recommendation with pg_trgm
	found := false
	for _, rec := range recommendations {
		if rec.Column == "person.email" && rec.IndexType == IndexTypeGIN {
			found = true
			if !strings.Contains(rec.Reason, "Regex matching") {
				t.Errorf("expected reason to mention regex matching, got %q", rec.Reason)
			}
			if !strings.Contains(rec.Expression, "gin_trgm_ops") {
				t.Error("expected gin_trgm_ops in expression for regex index")
			}
		}
	}

	if !found {
		t.Errorf("expected GIN index recommendation for person.email, got: %+v", recommendations)
	}
}

func TestAnalyzeQuery_ComparisonOperations(t *testing.T) {
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "id", Type: "bigint"},
		{Name: "age", Type: "integer"},
		{Name: "score", Type: "double precision"},
	})
	provider := pg.NewTypeProvider(map[string]pg.Schema{"person": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("person", cel.ObjectType("person")),
	)
	if err != nil {
		t.Fatalf("failed to create CEL environment: %v", err)
	}

	tests := []struct {
		name           string
		expression     string
		expectedColumn string
		expectedType   string
	}{
		{
			name:           "equality comparison",
			expression:     `person.id == 123`,
			expectedColumn: "person.id",
			expectedType:   "BTREE",
		},
		{
			name:           "range comparison",
			expression:     `person.age > 18`,
			expectedColumn: "person.age",
			expectedType:   "BTREE",
		},
		{
			name:           "less than comparison",
			expression:     `person.score < 100.0`,
			expectedColumn: "person.score",
			expectedType:   "BTREE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.expression)
			if issues != nil && issues.Err() != nil {
				t.Fatalf("failed to compile expression: %v", issues.Err())
			}

			sql, recommendations, err := AnalyzeQuery(ast, WithSchemas(provider.GetSchemas()))
			if err != nil {
				t.Fatalf("AnalyzeQuery failed: %v", err)
			}

			if sql == "" {
				t.Error("expected SQL output, got empty string")
			}

			// Check that we got a B-tree index recommendation
			found := false
			for _, rec := range recommendations {
				if rec.Column == tt.expectedColumn && rec.IndexType == tt.expectedType {
					found = true
					if !strings.Contains(rec.Reason, "Comparison operations") {
						t.Errorf("expected reason to mention comparison operations, got %q", rec.Reason)
					}
				}
			}

			if !found {
				t.Errorf("expected BTREE index recommendation for %s, got: %+v", tt.expectedColumn, recommendations)
			}
		})
	}
}

func TestAnalyzeQuery_ArrayOperations(t *testing.T) {
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "id", Type: "text"},
		{Name: "tags", Type: "text", Repeated: true},
	})
	provider := pg.NewTypeProvider(map[string]pg.Schema{"article": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("article", cel.ObjectType("article")),
	)
	if err != nil {
		t.Fatalf("failed to create CEL environment: %v", err)
	}

	ast, issues := env.Compile(`"golang" in article.tags`)
	if issues != nil && issues.Err() != nil {
		t.Fatalf("failed to compile expression: %v", issues.Err())
	}

	sql, recommendations, err := AnalyzeQuery(ast, WithSchemas(provider.GetSchemas()))
	if err != nil {
		t.Fatalf("AnalyzeQuery failed: %v", err)
	}

	if sql == "" {
		t.Error("expected SQL output, got empty string")
	}

	// Check that we got a GIN index recommendation for array operations
	found := false
	for _, rec := range recommendations {
		if rec.Column == "article.tags" && rec.IndexType == IndexTypeGIN {
			found = true
			if !strings.Contains(rec.Reason, "Array membership") {
				t.Errorf("expected reason to mention array membership, got %q", rec.Reason)
			}
		}
	}

	if !found {
		t.Errorf("expected GIN index recommendation for article.tags, got: %+v", recommendations)
	}
}

func TestAnalyzeQuery_Comprehensions(t *testing.T) {
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "id", Type: "text"},
		{Name: "scores", Type: "integer", Repeated: true},
	})
	provider := pg.NewTypeProvider(map[string]pg.Schema{"student": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("student", cel.ObjectType("student")),
	)
	if err != nil {
		t.Fatalf("failed to create CEL environment: %v", err)
	}

	tests := []struct {
		name       string
		expression string
	}{
		{
			name:       "all comprehension",
			expression: `student.scores.all(s, s > 50)`,
		},
		{
			name:       "exists comprehension",
			expression: `student.scores.exists(s, s == 100)`,
		},
		{
			name:       "map comprehension",
			expression: `student.scores.map(s, s * 2)[0] > 100`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.expression)
			if issues != nil && issues.Err() != nil {
				t.Fatalf("failed to compile expression: %v", issues.Err())
			}

			sql, recommendations, err := AnalyzeQuery(ast, WithSchemas(provider.GetSchemas()))
			if err != nil {
				t.Fatalf("AnalyzeQuery failed: %v", err)
			}

			if sql == "" {
				t.Error("expected SQL output, got empty string")
			}

			// Check that we got a GIN index recommendation for array comprehensions
			found := false
			for _, rec := range recommendations {
				if rec.Column == "student.scores" && rec.IndexType == IndexTypeGIN {
					found = true
					if !strings.Contains(rec.Reason, "comprehension") {
						t.Errorf("expected reason to mention comprehension, got %q", rec.Reason)
					}
				}
			}

			if !found {
				t.Errorf("expected GIN index recommendation for student.scores, got: %+v", recommendations)
			}
		})
	}
}

func TestAnalyzeQuery_MultipleRecommendations(t *testing.T) {
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "id", Type: "bigint"},
		{Name: "email", Type: "text"},
		{Name: "age", Type: "integer"},
		{Name: "metadata", Type: "jsonb", IsJSON: true, IsJSONB: true},
	})
	provider := pg.NewTypeProvider(map[string]pg.Schema{"person": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("person", cel.ObjectType("person")),
	)
	if err != nil {
		t.Fatalf("failed to create CEL environment: %v", err)
	}

	// Complex query with multiple index-worthy patterns
	ast, issues := env.Compile(`person.age > 18 && person.email.matches(r"@example\.com$") && person.metadata.verified == true`)
	if issues != nil && issues.Err() != nil {
		t.Fatalf("failed to compile expression: %v", issues.Err())
	}

	sql, recommendations, err := AnalyzeQuery(ast, WithSchemas(provider.GetSchemas()))
	if err != nil {
		t.Fatalf("AnalyzeQuery failed: %v", err)
	}

	if sql == "" {
		t.Error("expected SQL output, got empty string")
	}

	// We should have at least 3 recommendations (age BTREE, email GIN, metadata GIN)
	if len(recommendations) < 3 {
		t.Errorf("expected at least 3 recommendations, got %d: %+v", len(recommendations), recommendations)
	}

	// Check for specific recommendations
	foundAge := false
	foundEmail := false
	foundMetadata := false

	for _, rec := range recommendations {
		switch rec.Column {
		case "person.age":
			foundAge = rec.IndexType == IndexTypeBTree
		case "person.email":
			foundEmail = rec.IndexType == IndexTypeGIN
		case "person.metadata":
			foundMetadata = rec.IndexType == IndexTypeGIN
		}
	}

	if !foundAge {
		t.Error("expected BTREE index recommendation for person.age")
	}
	if !foundEmail {
		t.Error("expected GIN index recommendation for person.email")
	}
	if !foundMetadata {
		t.Error("expected GIN index recommendation for person.metadata")
	}
}

func TestAnalyzeQuery_NoRecommendations(t *testing.T) {
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "id", Type: "text"},
		{Name: "name", Type: "text"},
	})
	provider := pg.NewTypeProvider(map[string]pg.Schema{"person": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("person", cel.ObjectType("person")),
	)
	if err != nil {
		t.Fatalf("failed to create CEL environment: %v", err)
	}

	// Simple constant expression with no field access
	ast, issues := env.Compile(`true`)
	if issues != nil && issues.Err() != nil {
		t.Fatalf("failed to compile expression: %v", issues.Err())
	}

	sql, recommendations, err := AnalyzeQuery(ast, WithSchemas(provider.GetSchemas()))
	if err != nil {
		t.Fatalf("AnalyzeQuery failed: %v", err)
	}

	if sql == "" {
		t.Error("expected SQL output, got empty string")
	}

	if len(recommendations) != 0 {
		t.Errorf("expected no recommendations for constant expression, got %d: %+v",
			len(recommendations), recommendations)
	}
}

func TestAnalyzeQuery_WithContext(t *testing.T) {
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "id", Type: "bigint"},
		{Name: "age", Type: "integer"},
	})
	provider := pg.NewTypeProvider(map[string]pg.Schema{"person": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("person", cel.ObjectType("person")),
	)
	if err != nil {
		t.Fatalf("failed to create CEL environment: %v", err)
	}

	ast, issues := env.Compile(`person.age > 18`)
	if issues != nil && issues.Err() != nil {
		t.Fatalf("failed to compile expression: %v", issues.Err())
	}

	ctx := context.Background()
	sql, recommendations, err := AnalyzeQuery(ast,
		WithSchemas(provider.GetSchemas()),
		WithContext(ctx))
	if err != nil {
		t.Fatalf("AnalyzeQuery with context failed: %v", err)
	}

	if sql == "" {
		t.Error("expected SQL output, got empty string")
	}

	if len(recommendations) == 0 {
		t.Error("expected at least one recommendation")
	}
}

func TestAnalyzeQuery_IndexRecommendationPriority(t *testing.T) {
	// Test that GIN recommendations take priority over BTREE for JSON fields
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "id", Type: "text"},
		{Name: "metadata", Type: "jsonb", IsJSON: true, IsJSONB: true},
	})
	provider := pg.NewTypeProvider(map[string]pg.Schema{"person": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("person", cel.ObjectType("person")),
	)
	if err != nil {
		t.Fatalf("failed to create CEL environment: %v", err)
	}

	// Query that could trigger both comparison (BTREE) and JSON path (GIN) recommendations
	ast, issues := env.Compile(`person.metadata.age > 18 && person.metadata.name == "John"`)
	if issues != nil && issues.Err() != nil {
		t.Fatalf("failed to compile expression: %v", issues.Err())
	}

	sql, recommendations, err := AnalyzeQuery(ast, WithSchemas(provider.GetSchemas()))
	if err != nil {
		t.Fatalf("AnalyzeQuery failed: %v", err)
	}

	if sql == "" {
		t.Error("expected SQL output, got empty string")
	}

	// We should get a GIN recommendation for metadata, not BTREE
	for _, rec := range recommendations {
		if rec.Column == "person.metadata" {
			if rec.IndexType != IndexTypeGIN {
				t.Errorf("expected GIN index for JSON field, got %s", rec.IndexType)
			}
		}
	}
}
