// Package catalog loads production-intelligence definitions from YAML files.
// Agents, skills, policies, rules, and evaluators are editable data — the
// catalog is the registry source, not compiled code.
package catalog

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/yaml.v3"

	"dramastudio/internal/intelligence/domain"
)

// Catalog is the loaded definition set.
type Catalog struct {
	Agents     map[string]*domain.Agent
	Skills     map[string]*domain.Skill
	Policies   map[string][]*domain.Policy // policy ID → layers/versions
	Rules      map[string]*domain.Rule
	Evaluators map[string]*domain.Evaluator
}

// LoadDir reads <root>/{agents,skills,policies,rules,evaluators}/**/*.yaml.
// Missing subdirectories are not errors — the corresponding map stays empty.
func LoadDir(root string) (*Catalog, error) {
	c := &Catalog{
		Agents: map[string]*domain.Agent{}, Skills: map[string]*domain.Skill{},
		Policies: map[string][]*domain.Policy{}, Rules: map[string]*domain.Rule{},
		Evaluators: map[string]*domain.Evaluator{},
	}
	if err := loadKind(root, "agents", func(d *domain.Agent) {
		c.Agents[d.ID] = d
	}); err != nil {
		return nil, err
	}
	if err := loadKind(root, "skills", func(d *domain.Skill) {
		c.Skills[d.ID] = d
	}); err != nil {
		return nil, err
	}
	if err := loadKind(root, "policies", func(d *domain.Policy) {
		c.Policies[d.ID] = append(c.Policies[d.ID], d)
	}); err != nil {
		return nil, err
	}
	if err := loadKind(root, "rules", func(d *domain.Rule) {
		c.Rules[d.ID] = d
	}); err != nil {
		return nil, err
	}
	if err := loadKind(root, "evaluators", func(d *domain.Evaluator) {
		c.Evaluators[d.ID] = d
	}); err != nil {
		return nil, err
	}
	for _, layers := range c.Policies {
		sort.Slice(layers, func(i, j int) bool { return layers[i].Version < layers[j].Version })
	}
	return c, nil
}

// loadKind decodes every YAML file under <root>/<kind>/ and hands each to
// store. Missing directories are skipped.
func loadKind[D interface{ GetID() string }](root, kind string, store func(d D)) error {
	dir := filepath.Join(root, kind)
	for _, f := range yamlFiles(dir) {
		data, err := os.ReadFile(f)
		if err != nil {
			return fmt.Errorf("catalog %s: read %s: %w", kind, f, err)
		}
		var d D
		if err := yaml.Unmarshal(data, &d); err != nil {
			return fmt.Errorf("catalog %s: %s: %w", kind, f, err)
		}
		if d.GetID() == "" {
			return fmt.Errorf("catalog %s: %s: missing id", kind, f)
		}
		store(d)
	}
	return nil
}

// yamlFiles lists *.yaml|*.yml under dir recursively.
func yamlFiles(dir string) []string {
	var files []string
	_ = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err == nil && !d.IsDir() && (filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml") {
			files = append(files, path)
		}
		return nil
	})
	sort.Strings(files)
	return files
}
