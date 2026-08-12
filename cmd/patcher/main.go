// Command patch_spec flattens anyOf/oneOf-of-strings path parameters into
// plain strings.
//
// oapi-codegen renders an inline anyOf path parameter as an anonymous struct
// with an unexported `union` field and no constructor. Such a parameter
// cannot be built outside the generated package -- the operation is
// uncallable. Every union path param in these specs is a list of string
// variants differing only in their validation pattern, so the parameter is
// just a string.
//
// Reads a spec on stdin, writes the patched spec to stdout.
//
// Delete this once the upstream specs stop declaring path params as anyOf.
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

func main() {
	var spec map[string]any
	if err := json.NewDecoder(os.Stdin).Decode(&spec); err != nil {
		fmt.Fprintf(os.Stderr, "patch_spec: failed to decode spec: %v\n", err)
		os.Exit(1)
	}

	patched := 0

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
	}

	if err := json.NewEncoder(os.Stdout).Encode(spec); err != nil {
		fmt.Fprintf(os.Stderr, "patch_spec: failed to encode spec: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "patch_spec: flattened %d union path parameters\n", patched)
}
