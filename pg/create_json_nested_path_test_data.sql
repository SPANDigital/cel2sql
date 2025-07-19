-- Create test tables for JSON/JSONB nested path expressions
-- This includes complex nested structures to test CEL expressions like:
-- informationAsset.metadata.corpus.section == 'Getting Started'

-- Create information_assets table with complex nested JSON/JSONB structures
CREATE TABLE information_assets (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    metadata JSONB DEFAULT '{}',
    properties JSON DEFAULT '{}',
    classification JSONB DEFAULT '{}',
    content_structure JSON DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create documents table with nested metadata
CREATE TABLE documents (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    content JSONB DEFAULT '{}',
    structure JSON DEFAULT '{}',
    taxonomy JSONB DEFAULT '{}',
    analytics JSON DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Insert test data for information_assets with complex nested structures
INSERT INTO information_assets (name, metadata, properties, classification, content_structure) VALUES
    ('User Guide Documentation',
     '{
        "corpus": {
            "section": "Getting Started",
            "chapter": 1,
            "subsection": "Installation",
            "path": ["docs", "user-guide", "getting-started"],
            "tags": ["documentation", "tutorial", "beginner"]
        },
        "version": {
            "major": 2,
            "minor": 1,
            "patch": 0,
            "stage": "stable"
        },
        "author": {
            "name": "Documentation Team",
            "department": "Product",
            "email": "docs@company.com"
        }
     }',
     '{
        "visibility": "public",
        "language": "en",
        "format": "markdown",
        "size": {
            "bytes": 45678,
            "lines": 892,
            "pages": 12
        },
        "links": {
            "internal": ["api-reference", "tutorials"],
            "external": ["github.com/example", "docs.external.com"]
        }
     }',
     '{
        "security": {
            "level": "public",
            "clearance": "none",
            "restricted": false
        },
        "category": {
            "primary": "documentation",
            "secondary": ["tutorial", "guide"],
            "tags": ["user-facing", "beginner-friendly"]
        }
     }',
     '{
        "hierarchy": {
            "level": 1,
            "parent": null,
            "children": ["installation", "configuration", "first-steps"]
        },
        "sections": [
            {"name": "Overview", "order": 1, "type": "introduction"},
            {"name": "Prerequisites", "order": 2, "type": "requirements"},
            {"name": "Installation Steps", "order": 3, "type": "procedure"}
        ]
     }'),

    ('API Reference Manual',
     '{
        "corpus": {
            "section": "Reference",
            "chapter": 3,
            "subsection": "REST API",
            "path": ["docs", "api", "reference"],
            "tags": ["api", "reference", "technical"]
        },
        "version": {
            "major": 2,
            "minor": 1,
            "patch": 0,
            "stage": "stable"
        },
        "author": {
            "name": "Engineering Team",
            "department": "Engineering",
            "email": "engineering@company.com"
        }
     }',
     '{
        "visibility": "public",
        "language": "en",
        "format": "openapi",
        "size": {
            "bytes": 123456,
            "lines": 2890,
            "pages": 45
        },
        "links": {
            "internal": ["getting-started", "tutorials"],
            "external": ["swagger.io", "openapi.org"]
        }
     }',
     '{
        "security": {
            "level": "public",
            "clearance": "none",
            "restricted": false
        },
        "category": {
            "primary": "reference",
            "secondary": ["api", "technical"],
            "tags": ["developer-facing", "comprehensive"]
        }
     }',
     '{
        "hierarchy": {
            "level": 2,
            "parent": "documentation",
            "children": ["endpoints", "schemas", "examples"]
        },
        "sections": [
            {"name": "Authentication", "order": 1, "type": "security"},
            {"name": "Endpoints", "order": 2, "type": "reference"},
            {"name": "Data Models", "order": 3, "type": "schemas"}
        ]
     }'),

    ('Advanced Tutorial Series',
     '{
        "corpus": {
            "section": "Tutorials",
            "chapter": 2,
            "subsection": "Advanced Features",
            "path": ["docs", "tutorials", "advanced"],
            "tags": ["tutorial", "advanced", "features"]
        },
        "version": {
            "major": 1,
            "minor": 8,
            "patch": 2,
            "stage": "beta"
        },
        "author": {
            "name": "Technical Writers",
            "department": "Documentation",
            "email": "writers@company.com"
        }
     }',
     '{
        "visibility": "members",
        "language": "en",
        "format": "interactive",
        "size": {
            "bytes": 78901,
            "lines": 1456,
            "pages": 23
        },
        "links": {
            "internal": ["getting-started", "api-reference"],
            "external": ["community.example.com", "forum.example.com"]
        }
     }',
     '{
        "security": {
            "level": "member",
            "clearance": "basic",
            "restricted": true
        },
        "category": {
            "primary": "tutorial",
            "secondary": ["advanced", "interactive"],
            "tags": ["member-only", "hands-on"]
        }
     }',
     '{
        "hierarchy": {
            "level": 1,
            "parent": "tutorials",
            "children": ["webhooks", "integrations", "optimization"]
        },
        "sections": [
            {"name": "Webhook Configuration", "order": 1, "type": "tutorial"},
            {"name": "Third-party Integrations", "order": 2, "type": "guide"},
            {"name": "Performance Optimization", "order": 3, "type": "best-practices"}
        ]
     }'),

    ('Migration Guide',
     '{
        "corpus": {
            "section": "Getting Started",
            "chapter": 4,
            "subsection": "Migration",
            "path": ["docs", "migration", "v2"],
            "tags": ["migration", "upgrade", "breaking-changes"]
        },
        "version": {
            "major": 2,
            "minor": 0,
            "patch": 0,
            "stage": "stable"
        },
        "author": {
            "name": "Platform Team",
            "department": "Engineering",
            "email": "platform@company.com"
        }
     }',
     '{
        "visibility": "public",
        "language": "en",
        "format": "step-by-step",
        "size": {
            "bytes": 34567,
            "lines": 678,
            "pages": 8
        },
        "links": {
            "internal": ["api-reference", "changelog"],
            "external": ["migration-tools.example.com"]
        }
     }',
     '{
        "security": {
            "level": "public",
            "clearance": "none",
            "restricted": false
        },
        "category": {
            "primary": "guide",
            "secondary": ["migration", "upgrade"],
            "tags": ["version-specific", "critical"]
        }
     }',
     '{
        "hierarchy": {
            "level": 1,
            "parent": "getting-started",
            "children": ["breaking-changes", "compatibility", "rollback"]
        },
        "sections": [
            {"name": "Pre-migration Checklist", "order": 1, "type": "checklist"},
            {"name": "Step-by-step Migration", "order": 2, "type": "procedure"},
            {"name": "Post-migration Verification", "order": 3, "type": "validation"}
        ]
     }');

-- Insert test data for documents with nested metadata
INSERT INTO documents (title, content, structure, taxonomy, analytics) VALUES
    ('Introduction to APIs',
     '{
        "body": {
            "sections": [
                {
                    "heading": "What is an API?",
                    "content": "An API (Application Programming Interface)...",
                    "metadata": {"wordCount": 245, "readingTime": 2}
                },
                {
                    "heading": "Types of APIs",
                    "content": "There are several types of APIs...",
                    "metadata": {"wordCount": 189, "readingTime": 1}
                }
            ]
        },
        "metadata": {
            "corpus": {
                "section": "Getting Started",
                "topic": "fundamentals",
                "difficulty": "beginner"
            },
            "stats": {
                "totalWords": 434,
                "estimatedReadTime": 3,
                "lastUpdated": "2024-07-15"
            }
        }
     }',
     '{
        "outline": {
            "chapters": [
                {"id": 1, "title": "Introduction", "pages": [1, 2, 3]},
                {"id": 2, "title": "Core Concepts", "pages": [4, 5, 6, 7]}
            ]
        },
        "navigation": {
            "previous": null,
            "next": "api-authentication",
            "breadcrumbs": ["docs", "getting-started", "introduction"]
        }
     }',
     '{
        "categories": {
            "primary": "documentation",
            "secondary": ["introduction", "concepts"],
            "audience": ["beginners", "new-developers"]
        },
        "keywords": {
            "primary": ["api", "introduction", "fundamentals"],
            "secondary": ["rest", "http", "web-services"],
            "tags": ["beginner-friendly", "foundational"]
        }
     }',
     '{
        "engagement": {
            "views": {"total": 1250, "unique": 980, "returning": 270},
            "feedback": {"positive": 45, "negative": 3, "neutral": 12},
            "timeOnPage": {"average": 180, "median": 165, "mode": 120}
        },
        "performance": {
            "loadTime": {"average": 1.2, "p95": 2.1, "p99": 3.4},
            "bounceRate": 0.15,
            "conversionRate": 0.78
        }
     }'),

    ('Authentication Best Practices',
     '{
        "body": {
            "sections": [
                {
                    "heading": "Security Overview",
                    "content": "Authentication is critical for API security...",
                    "metadata": {"wordCount": 312, "readingTime": 2}
                },
                {
                    "heading": "Implementation Guide",
                    "content": "Follow these steps to implement secure authentication...",
                    "metadata": {"wordCount": 456, "readingTime": 3}
                }
            ]
        },
        "metadata": {
            "corpus": {
                "section": "Security",
                "topic": "authentication",
                "difficulty": "intermediate"
            },
            "stats": {
                "totalWords": 768,
                "estimatedReadTime": 5,
                "lastUpdated": "2024-07-18"
            }
        }
     }',
     '{
        "outline": {
            "chapters": [
                {"id": 1, "title": "Security Fundamentals", "pages": [1, 2]},
                {"id": 2, "title": "OAuth Implementation", "pages": [3, 4, 5]},
                {"id": 3, "title": "Best Practices", "pages": [6, 7]}
            ]
        },
        "navigation": {
            "previous": "api-introduction",
            "next": "rate-limiting",
            "breadcrumbs": ["docs", "security", "authentication"]
        }
     }',
     '{
        "categories": {
            "primary": "security",
            "secondary": ["authentication", "oauth", "best-practices"],
            "audience": ["developers", "security-engineers"]
        },
        "keywords": {
            "primary": ["authentication", "oauth", "security"],
            "secondary": ["tokens", "jwt", "authorization"],
            "tags": ["security-critical", "implementation-guide"]
        }
     }',
     '{
        "engagement": {
            "views": {"total": 890, "unique": 720, "returning": 170},
            "feedback": {"positive": 62, "negative": 2, "neutral": 8},
            "timeOnPage": {"average": 240, "median": 220, "mode": 200}
        },
        "performance": {
            "loadTime": {"average": 1.1, "p95": 1.9, "p99": 2.8},
            "bounceRate": 0.12,
            "conversionRate": 0.85
        }
     }'),

    ('Troubleshooting Common Issues',
     '{
        "body": {
            "sections": [
                {
                    "heading": "Error Diagnosis",
                    "content": "When encountering API errors...",
                    "metadata": {"wordCount": 198, "readingTime": 1}
                },
                {
                    "heading": "Common Solutions",
                    "content": "Here are solutions to frequent problems...",
                    "metadata": {"wordCount": 523, "readingTime": 4}
                }
            ]
        },
        "metadata": {
            "corpus": {
                "section": "Troubleshooting",
                "topic": "common-issues",
                "difficulty": "intermediate"
            },
            "stats": {
                "totalWords": 721,
                "estimatedReadTime": 5,
                "lastUpdated": "2024-07-20"
            }
        }
     }',
     '{
        "outline": {
            "chapters": [
                {"id": 1, "title": "Error Types", "pages": [1, 2]},
                {"id": 2, "title": "Debugging Techniques", "pages": [3, 4]},
                {"id": 3, "title": "Resolution Steps", "pages": [5, 6, 7]}
            ]
        },
        "navigation": {
            "previous": "authentication",
            "next": "performance-optimization",
            "breadcrumbs": ["docs", "troubleshooting", "common-issues"]
        }
     }',
     '{
        "categories": {
            "primary": "support",
            "secondary": ["troubleshooting", "debugging", "solutions"],
            "audience": ["developers", "support-engineers"]
        },
        "keywords": {
            "primary": ["troubleshooting", "errors", "debugging"],
            "secondary": ["solutions", "fixes", "diagnosis"],
            "tags": ["problem-solving", "reference"]
        }
     }',
     '{
        "engagement": {
            "views": {"total": 1456, "unique": 1123, "returning": 333},
            "feedback": {"positive": 89, "negative": 5, "neutral": 15},
            "timeOnPage": {"average": 300, "median": 275, "mode": 250}
        },
        "performance": {
            "loadTime": {"average": 1.3, "p95": 2.2, "p99": 3.1},
            "bounceRate": 0.08,
            "conversionRate": 0.92
        }
     }');

-- Add some test data with different section values for comprehensive testing
INSERT INTO information_assets (name, metadata, properties, classification, content_structure) VALUES
    ('Quick Start Guide',
     '{
        "corpus": {
            "section": "Quick Start",
            "chapter": 1,
            "subsection": "Overview",
            "path": ["docs", "quick-start"],
            "tags": ["quick", "overview", "summary"]
        },
        "version": {"major": 1, "minor": 0, "patch": 0, "stage": "stable"},
        "author": {"name": "Product Team", "department": "Product", "email": "product@company.com"}
     }',
     '{"visibility": "public", "language": "en", "format": "markdown"}',
     '{"security": {"level": "public", "clearance": "none", "restricted": false}}',
     '{"hierarchy": {"level": 1, "parent": null, "children": ["overview", "setup"]}}'),

    ('Developer Resources',
     '{
        "corpus": {
            "section": "Reference",
            "chapter": 5,
            "subsection": "Resources",
            "path": ["docs", "developer", "resources"],
            "tags": ["developer", "resources", "tools"]
        },
        "version": {"major": 1, "minor": 5, "patch": 1, "stage": "stable"},
        "author": {"name": "Dev Relations", "department": "Developer Relations", "email": "devrel@company.com"}
     }',
     '{"visibility": "public", "language": "en", "format": "collection"}',
     '{"security": {"level": "public", "clearance": "none", "restricted": false}}',
     '{"hierarchy": {"level": 2, "parent": "reference", "children": ["tools", "libraries", "sdks"]}}');
