package checkpoint_test

import (
	"testing"

	"github.com/aidenfine/checkpoint"
)

func TestMatchPathPattern(t *testing.T) {
	tests := []struct {
		name        string
		pattern     string
		path        string
		shouldMatch bool
	}{
		{
			name:        "exact match - simple path",
			pattern:     "/api/users",
			path:        "/api/users",
			shouldMatch: true,
		},
		{
			name:        "exact match - no match",
			pattern:     "/api/users",
			path:        "/api/posts",
			shouldMatch: false,
		},
		{
			name:        "exact match - case sensitive",
			pattern:     "/api/users",
			path:        "/api/Users",
			shouldMatch: false,
		},
		{
			name:        "trailing wildcard - matches subdirectory",
			pattern:     "logs/*",
			path:        "logs/error.log",
			shouldMatch: true,
		},
		{
			name:        "trailing wildcard - matches nested path",
			pattern:     "logs/*",
			path:        "logs/2024/01/error.log",
			shouldMatch: true,
		},
		{
			name:        "trailing wildcard - matches empty after slash",
			pattern:     "logs/*",
			path:        "logs/",
			shouldMatch: true,
		},
		{
			name:        "trailing wildcard - no match wrong prefix",
			pattern:     "logs/*",
			path:        "logs2/error.log",
			shouldMatch: false,
		},
		{
			name:        "trailing wildcard - no match partial prefix",
			pattern:     "logs/*",
			path:        "log/error.log",
			shouldMatch: false,
		},
		{
			name:        "trailing wildcard - no match just prefix",
			pattern:     "logs/*",
			path:        "logs",
			shouldMatch: false,
		},
		{
			name:        "trailing wildcard - with leading slash",
			pattern:     "/tmp/tmpA/*",
			path:        "/tmp/tmpA/session123",
			shouldMatch: true,
		},
		{
			name:        "trailing wildcard - deep nesting",
			pattern:     "/tmp/tmpA/*",
			path:        "/tmp/tmpA/deep/nested/path/file.txt",
			shouldMatch: true,
		},
		{
			name:        "leading wildcard - matches with suffix",
			pattern:     "*/admin",
			path:        "/api/v1/admin",
			shouldMatch: true,
		},
		{
			name:        "leading wildcard - matches simple",
			pattern:     "*/admin",
			path:        "users/admin",
			shouldMatch: true,
		},
		{
			name:        "leading wildcard - no match wrong suffix",
			pattern:     "*/admin",
			path:        "/api/v1/user",
			shouldMatch: false,
		},
		{
			name:        "leading wildcard - no match partial suffix",
			pattern:     "*/admin",
			path:        "/api/v1/administrator",
			shouldMatch: false,
		},
		{
			name:        "middle wildcard - single wildcard",
			pattern:     "/api/*/users",
			path:        "/api/v1/users",
			shouldMatch: true,
		},
		{
			name:        "middle wildcard - matches multiple segments",
			pattern:     "/api/*/users",
			path:        "/api/v1/public/users",
			shouldMatch: true,
		},
		{
			name:        "middle wildcard - no match wrong prefix",
			pattern:     "/api/*/users",
			path:        "/app/v1/users",
			shouldMatch: false,
		},
		{
			name:        "middle wildcard - no match wrong suffix",
			pattern:     "/api/*/users",
			path:        "/api/v1/posts",
			shouldMatch: false,
		},
		{
			name:        "middle wildcard - matches empty middle",
			pattern:     "/api/*/users",
			path:        "/api//users",
			shouldMatch: true,
		},
		{
			name:        "multiple wildcards - two wildcards",
			pattern:     "/api/*/users/*",
			path:        "/api/v1/users/123",
			shouldMatch: true,
		},
		{
			name:        "multiple wildcards - complex path",
			pattern:     "/api/*/users/*",
			path:        "/api/v1/public/users/123/profile",
			shouldMatch: true,
		},
		{
			name:        "multiple wildcards - three wildcards",
			pattern:     "*/logs/*/error/*",
			path:        "app/logs/2024/error/critical.log",
			shouldMatch: true,
		},
		{
			name:        "empty pattern and path",
			pattern:     "",
			path:        "",
			shouldMatch: true,
		},
		{
			name:        "empty pattern non-empty path",
			pattern:     "",
			path:        "/api/users",
			shouldMatch: false,
		},
		{
			name:        "wildcard only pattern",
			pattern:     "*",
			path:        "/anything/goes/here",
			shouldMatch: true,
		},
		{
			name:        "wildcard only pattern - empty path",
			pattern:     "*",
			path:        "",
			shouldMatch: true,
		},
		{
			name:        "path with special characters",
			pattern:     "/api/*/data",
			path:        "/api/v1.2-beta/data",
			shouldMatch: true,
		},
		{
			name:        "path with query-like string",
			pattern:     "/api/users/*",
			path:        "/api/users/search?name=john",
			shouldMatch: true,
		},
		{
			name:        "consecutive wildcards",
			pattern:     "logs/*/*",
			path:        "logs/2024/01/error.log",
			shouldMatch: true,
		},
		{
			name:        "pattern longer than path",
			pattern:     "/api/v1/users/profile",
			path:        "/api/v1",
			shouldMatch: false,
		},
		{
			name:        "path longer than pattern",
			pattern:     "/api/v1",
			path:        "/api/v1/users/profile",
			shouldMatch: false,
		},
		{
			name:        "similar but different paths",
			pattern:     "/api/user/*",
			path:        "/api/users/123",
			shouldMatch: false,
		},
		{
			name:        "wildcard at end without slash",
			pattern:     "/api/logs*",
			path:        "/api/logs-archive",
			shouldMatch: true,
		},
		{
			name:        "wildcard at start without slash",
			pattern:     "*logs",
			path:        "error-logs",
			shouldMatch: true,
		},
		{
			name:        "health check ignore",
			pattern:     "/health/*",
			path:        "/health/status",
			shouldMatch: true,
		},
		{
			name:        "versioned API",
			pattern:     "/api/v*/internal/*",
			path:        "/api/v2/internal/metrics",
			shouldMatch: true,
		},
		{
			name:        "static assets",
			pattern:     "/static/*",
			path:        "/static/css/main.css",
			shouldMatch: true,
		},
		{
			name:        "user-specific paths",
			pattern:     "/users/*/private/*",
			path:        "/users/john123/private/settings",
			shouldMatch: true,
		},
		{
			name:        "admin panel",
			pattern:     "*/admin/*",
			path:        "/dashboard/admin/users",
			shouldMatch: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkpoint.MatchPathPattern(tt.pattern, tt.path)
			if result != tt.shouldMatch {
				t.Errorf("checkpoint.MatchPathPattern(%q, %q) = %v; want %v",
					tt.pattern, tt.path, result, tt.shouldMatch)
			}
		})
	}
}

func BenchMarkMatchPathPattern(b *testing.B) {
	patterns := []string{
		"logs/*",
		"/api/*/users",
		"/tmp/tmpA/*",
		"*/admin/*",
	}
	paths := []string{
		"logs/error.log",
		"/api/v1/users",
		"/tmp/tmpA/session123",
		"/app/admin/settings",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, pattern := range patterns {
			for _, path := range paths {
				checkpoint.MatchPathPattern(pattern, path)
			}
		}
	}
}

// func TestTokenBucket_MatchPath(t *testing.T) {
// 	tests := []struct {
// 		name        string
// 		ignorePaths []string
// 		testPath    string
// 		shouldMatch bool
// 	}{
// 		{
// 			name:        "matches first pattern",
// 			ignorePaths: []string{"logs/*", "/tmp/*"},
// 			testPath:    "logs/error.log",
// 			shouldMatch: true,
// 		},
// 		{
// 			name:        "matches second pattern",
// 			ignorePaths: []string{"logs/*", "/tmp/*"},
// 			testPath:    "/tmp/session",
// 			shouldMatch: true,
// 		},
// 		{
// 			name:        "no match",
// 			ignorePaths: []string{"logs/*", "/tmp/*"},
// 			testPath:    "/api/users",
// 			shouldMatch: false,
// 		},
// 		{
// 			name:        "empty ignore list",
// 			ignorePaths: []string{},
// 			testPath:    "/any/path",
// 			shouldMatch: false,
// 		},
// 		{
// 			name:        "multiple patterns with overlap",
// 			ignorePaths: []string{"/api/*", "/api/v1/*", "/api/v1/admin/*"},
// 			testPath:    "/api/v1/admin/users",
// 			shouldMatch: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tb := &TokenBucket{
// 				ignorePaths: tt.ignorePaths,
// 			}
// 			result := tb.matchPath(tt.testPath)
// 			if result != tt.shouldMatch {
// 				t.Errorf("matchPath(%q) = %v; want %v",
// 					tt.testPath, result, tt.shouldMatch)
// 			}
// 		})
// 	}
// }
