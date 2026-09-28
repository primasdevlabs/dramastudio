package sqlite

import (
	"regexp"
	"strings"
)

// translate rewrites the Postgres SQL dialect used by repositories into
// SQLite. It also rebinds arguments: $N placeholders become positional ? and
// each occurrence consumes args[N-1], so params referenced twice (e.g.
// `$2::text = ” OR episode_id = $2`) bind correctly.
//
// Handled dialect surface:
//   - $N placeholders (+ trailing ::type casts)      → ?
//   - schema.table                                    → schema_table
//   - col @> to_jsonb($N::text)                       → EXISTS json_each(col)
//   - now()                                           → datetime('now')
//   - RETURNING / ON CONFLICT / CTEs                  → native in SQLite 3.35+
var (
	jsonbContains = regexp.MustCompile(`(\w+)\s*@>\s*to_jsonb\(\s*\$(\d+)(?:::\w+)?\s*\)`)
	placeholder   = regexp.MustCompile(`\$(\d+)(?:::[a-zA-Z_]+)*`)
	nowFunc       = regexp.MustCompile(`(?i)\bnow\(\s*\)`)
	schemaPrefix  = regexp.MustCompile(`\b(` + strings.Join(schemaNames, "|") + `)\.`)
)

// schemaNames mirrors the CREATE SCHEMA statements in migrations/. Tables are
// flattened to schema_name.table_name on the SQLite side.
var schemaNames = []string{
	"ai", "agents", "analytics", "canon", "characters", "continuity",
	"identity", "intelligence", "media", "platform", "postproduction",
	"production", "projects", "publishing", "story", "world",
}

func translate(sqlText string, args []interface{}) (string, []interface{}) {
	// JSONB containment → EXISTS over json_each (SQLite JSON1). The $$${2}
	// escape emits a literal $N placeholder consumed by the pass below.
	q := jsonbContains.ReplaceAllString(sqlText,
		`EXISTS (SELECT 1 FROM json_each($1) WHERE value = $$${2})`)

	// Placeholders → ?, rebinding args per occurrence (a $N may repeat).
	bound := make([]interface{}, 0, len(args))
	q = placeholder.ReplaceAllStringFunc(q, func(m string) string {
		sub := placeholder.FindStringSubmatch(m)
		idx := atoi(sub[1]) - 1
		if idx < 0 || idx >= len(args) {
			return m
		}
		bound = append(bound, args[idx])
		return "?"
	})

	q = schemaPrefix.ReplaceAllString(q, `${1}_`)
	q = nowFunc.ReplaceAllString(q, `datetime('now')`)
	return q, bound
}

func atoi(s string) int {
	n := 0
	for _, c := range s {
		n = n*10 + int(c-'0')
	}
	return n
}
