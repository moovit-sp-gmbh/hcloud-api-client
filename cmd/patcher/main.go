// Command patch_spec rewrites the parts of these OpenAPI specs that
// oapi-codegen cannot render.
//
// It applies two transforms:
//
// Path parameters. oapi-codegen renders an inline anyOf path parameter as an
// anonymous struct with an unexported `union` field and no constructor. Such a
// parameter cannot be built outside the generated package -- the operation is
// uncallable. Every union path param in these specs is a list of string
// variants differing only in their validation pattern, so the parameter is
// just a string.
//
// Request bodies. A oneOf that carries a discriminator mapping but whose
// variants are inline allOf compositions rather than $refs aborts generation
// outright with "ambiguous discriminator.mapping" -- the mapping names schemas
// to switch on, but an inline variant has no $ref for it to name. The cosmo
// spec describes each searchable field this way (`assetType filters`,
// `expires filters`, ...), so every request body that reaches one is replaced
// with a free-form object. Callers of those operations send raw JSON anyway.
//
// A body is matched by what it reaches, not by its path, so an endpoint only
// loses its typed body when it genuinely contains one of these unions.
//
// Reads a spec on stdin, writes the patched spec to stdout.
//
// Delete this once the upstream specs stop declaring path params as anyOf and
// stop pairing discriminator mappings with inline oneOf variants.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// flatten returns a plain string schema if every union variant is a string,
// or nil if the schema has no such union.
func flatten(schema map[string]any) map[string]any {
	for _, key := range []string{"anyOf", "oneOf"} {
		raw, ok := schema[key]
		if !ok {
			continue
		}
		variants, ok := raw.([]any)
		if !ok || len(variants) == 0 {
			continue
		}

		allStrings := true
		for _, v := range variants {
			vm, ok := v.(map[string]any)
			if !ok || vm["type"] != "string" {
				allStrings = false
				break
			}
		}
		if !allStrings {
			continue
		}

		merged := map[string]any{"type": "string"}

		var descriptions []string
		seen := make(map[string]bool)
		for _, v := range variants {
			vm := v.(map[string]any)
			desc, ok := vm["description"].(string)
			if !ok || desc == "" || seen[desc] {
				continue
			}
			seen[desc] = true
			descriptions = append(descriptions, desc)
		}
		if len(descriptions) > 0 {
			var joined strings.Builder
			joined.WriteString(descriptions[0])
			for _, d := range descriptions[1:] {
				joined.WriteString(" | ")
				joined.WriteString(d)
			}
			merged["description"] = joined.String()
		}

		return merged
	}
	return nil
}

const schemaRefPrefix = "#/components/schemas/"

// unrenderableUnion reports whether schema is a oneOf that oapi-codegen
// refuses to generate: a discriminator mapping names schemas to switch on, but
// at least one variant is inline and so has no $ref for the mapping to name.
//
// A oneOf of inline variants with no discriminator is fine -- it becomes an
// opaque union -- so the mapping is what makes this fatal.
func unrenderableUnion(schema map[string]any) bool {
	variants, ok := schema["oneOf"].([]any)
	if !ok || len(variants) == 0 {
		return false
	}

	discriminator, ok := schema["discriminator"].(map[string]any)
	if !ok {
		return false
	}
	if mapping, ok := discriminator["mapping"].(map[string]any); !ok || len(mapping) == 0 {
		return false
	}

	for _, raw := range variants {
		variant, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		if _, hasRef := variant["$ref"]; !hasRef {
			return true
		}
	}
	return false
}

// containsUnrenderableUnion reports whether node or anything nested inside it
// is an [unrenderableUnion].
func containsUnrenderableUnion(node any) bool {
	switch n := node.(type) {
	case map[string]any:
		if unrenderableUnion(n) {
			return true
		}
		for _, v := range n {
			if containsUnrenderableUnion(v) {
				return true
			}
		}
	case []any:
		for _, v := range n {
			if containsUnrenderableUnion(v) {
				return true
			}
		}
	}
	return false
}

// collectRefs adds every component schema referenced by node to out.
func collectRefs(node any, out map[string]bool) {
	switch n := node.(type) {
	case map[string]any:
		for k, v := range n {
			if k == "$ref" {
				if ref, ok := v.(string); ok && strings.HasPrefix(ref, schemaRefPrefix) {
					out[strings.TrimPrefix(ref, schemaRefPrefix)] = true
				}
				continue
			}
			collectRefs(v, out)
		}
	case []any:
		for _, v := range n {
			collectRefs(v, out)
		}
	}
}

// reachesUnion reports whether node holds an unrenderable union directly, or
// references one through any chain of component schemas.
func reachesUnion(node any, schemas map[string]any, unions map[string]bool) bool {
	if containsUnrenderableUnion(node) {
		return true
	}

	pending := make(map[string]bool)
	collectRefs(node, pending)

	var stack []string
	for name := range pending {
		stack = append(stack, name)
	}

	// Component schemas reference each other freely and may form cycles, so
	// track what has already been expanded.
	seen := make(map[string]bool)
	for len(stack) > 0 {
		name := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if seen[name] {
			continue
		}
		seen[name] = true

		if unions[name] {
			return true
		}

		next := make(map[string]bool)
		collectRefs(schemas[name], next)
		for n := range next {
			if !seen[n] {
				stack = append(stack, n)
			}
		}
	}
	return false
}

func main() {
	var spec map[string]any
	if err := json.NewDecoder(os.Stdin).Decode(&spec); err != nil {
		fmt.Fprintf(os.Stderr, "patch_spec: failed to decode spec: %v\n", err)
		os.Exit(1)
	}

	patched := 0
	bodies := 0

	components, _ := spec["components"].(map[string]any)
	schemas, _ := components["schemas"].(map[string]any)

	// Component schemas that generation would abort on, and so that no request
	// body may reach.
	unions := make(map[string]bool)
	for name, raw := range schemas {
		if containsUnrenderableUnion(raw) {
			unions[name] = true
		}
	}

	paths, _ := spec["paths"].(map[string]any)
	for _, rawPathItem := range paths {
		pathItem, ok := rawPathItem.(map[string]any)
		if !ok {
			continue
		}

		// Path params live either on the path item or on each operation.
		var groups []any
		groups = append(groups, pathItem["parameters"])
		for _, rawOp := range pathItem {
			op, ok := rawOp.(map[string]any)
			if !ok {
				continue
			}
			groups = append(groups, op["parameters"])
		}

		for _, rawGroup := range groups {
			group, ok := rawGroup.([]any)
			if !ok {
				continue
			}
			for _, rawParam := range group {
				param, ok := rawParam.(map[string]any)
				if !ok || param["in"] != "path" {
					continue
				}
				schema, ok := param["schema"].(map[string]any)
				if !ok {
					continue
				}
				if merged := flatten(schema); merged != nil {
					param["schema"] = merged
					patched++
				}
			}
		}

		for _, rawOp := range pathItem {
			op, ok := rawOp.(map[string]any)
			if !ok {
				continue
			}
			requestBody, ok := op["requestBody"].(map[string]any)
			if !ok {
				continue
			}
			content, ok := requestBody["content"].(map[string]any)
			if !ok {
				continue
			}

			for _, rawMedia := range content {
				media, ok := rawMedia.(map[string]any)
				if !ok {
					continue
				}
				if _, ok := media["schema"]; !ok {
					continue
				}
				if !reachesUnion(media["schema"], schemas, unions) {
					continue
				}
				media["schema"] = map[string]any{
					"type":                 "object",
					"additionalProperties": true,
				}
				bodies++
			}
		}
	}

	if err := json.NewEncoder(os.Stdout).Encode(spec); err != nil {
		fmt.Fprintf(os.Stderr, "patch_spec: failed to encode spec: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "patch_spec: flattened %d union path parameters\n", patched)
	fmt.Fprintf(os.Stderr, "patch_spec: genericized %d request bodies reaching %d unrenderable unions\n", bodies, len(unions))
}
