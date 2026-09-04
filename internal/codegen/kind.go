package codegen

import (
	"fmt"
	"strconv"
	"strings"
)

// argKind says how one API type is exposed as a flag and converted back.
//
// The flag's pflag type depends on whether the argument can be prompted for:
// internal/command/args.go only offers FlagOrDifferentPrompt helpers for
// string, bool, int and float64, so a float32 parameter that should prompt has
// to be carried on a Float64 flag and narrowed at the call site.
type argKind struct {
	Flag      string // pflag method suffix: "String", "Float64", ...
	Default   string // default literal in the flag registration
	Helper    string // command.FlagOrDifferentPrompt<X>, empty when it cannot prompt
	Wrap      string // conversion at the call site, %s is the variable
	Pointer   bool   // the target field is a pointer
	Parse     bool   // needs a parse step rather than a cast (UUID)
	Imports   []string
	Supported bool
}

// DefaultLiteral turns a recipe's plain default value into the Go literal the
// flag registration needs, so a recipe can say `default: false` or
// `default: my-space` rather than carrying Go source.
//
// An empty value means the type's zero, which is what pflag would use anyway.
func (k argKind) DefaultLiteral(value string) (string, error) {
	if value == "" {
		return k.Default, nil
	}

	switch k.Flag {
	case "String":
		return strconv.Quote(value), nil
	case "StringSlice":
		parts := strings.Split(value, ",")
		quoted := make([]string, 0, len(parts))
		for _, part := range parts {
			quoted = append(quoted, strconv.Quote(strings.TrimSpace(part)))
		}
		return "[]string{" + strings.Join(quoted, ", ") + "}", nil
	case "Bool":
		b, err := strconv.ParseBool(value)
		if err != nil {
			return "", fmt.Errorf("%q is not a boolean; use true or false", value)
		}
		return strconv.FormatBool(b), nil
	case "Int", "Int32", "Int64":
		if _, err := strconv.ParseInt(value, 10, 64); err != nil {
			return "", fmt.Errorf("%q is not a whole number", value)
		}
		return value, nil
	case "Float32", "Float64":
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			return "", fmt.Errorf("%q is not a number", value)
		}
		return value, nil
	}
	return "", fmt.Errorf("a %s flag cannot carry a default", k.Flag)
}

// Convert renders the call-site expression for a resolved variable.
func (k argKind) Convert(v string) string {
	if k.Wrap == "" {
		return v
	}
	return fmt.Sprintf(k.Wrap, v)
}

// builtinKinds are the Go types the generated clients use for scalars.
var builtinKinds = map[string]argKind{
	"string":   {Flag: "String", Default: `""`, Helper: "String", Supported: true},
	"bool":     {Flag: "Bool", Default: "false", Helper: "Bool", Supported: true},
	"int":      {Flag: "Int", Default: "0", Helper: "Int", Supported: true},
	"int32":    {Flag: "Int32", Default: "0", Supported: true},
	"int64":    {Flag: "Int64", Default: "0", Supported: true},
	"float32":  {Flag: "Float32", Default: "0", Supported: true},
	"float64":  {Flag: "Float64", Default: "0", Helper: "Float", Supported: true},
	"[]byte":   {Flag: "String", Default: `""`, Helper: "String", Wrap: "[]byte(%s)", Supported: true},
	"[]string": {Flag: "StringSlice", Default: "nil", Supported: true},
}

// promptWidths maps a kind that cannot prompt at its own width onto the wider
// flag type whose helper exists, plus the narrowing conversion.
var promptWidths = map[string]struct {
	flag, helper, wrap string
}{
	"float32": {"Float64", "Float", "float32(%s)"},
	"int32":   {"Int", "Int", "int32(%s)"},
	"int64":   {"Int", "Int", "int64(%s)"},
}

// KindFor describes how to expose an API type as a flag. clientPkg is the
// generated client's package name, needed to qualify its named types, and
// scalar is the builtin a named type resolves to (empty when it resolves to a
// struct, which no flag can carry).
func KindFor(goType, scalar string, wantPrompt bool, clientPkg string) argKind {
	pointer := strings.HasPrefix(goType, "*")
	base := strings.TrimPrefix(goType, "*")

	k, ok := resolveKind(base, scalar, clientPkg)
	if !ok {
		return argKind{}
	}
	k.Pointer = pointer

	// Widen the flag when the caller wants a prompt this width cannot serve.
	if wantPrompt && k.Helper == "" {
		if w, ok := promptWidths[base]; ok {
			k.Flag, k.Helper = w.flag, w.helper
			k.Wrap = combineWrap(w.wrap, k.Wrap)
		}
	}
	return k
}

func resolveKind(base, scalar, clientPkg string) (argKind, bool) {
	if k, ok := builtinKinds[base]; ok {
		return k, true
	}

	switch base {
	case "openapi_types.Email", "types.Email":
		return argKind{
			Flag: "String", Default: `""`, Helper: "String",
			Wrap:      "types.Email(%s)",
			Imports:   []string{"github.com/oapi-codegen/runtime/types"},
			Supported: true,
		}, true
	case "openapi_types.UUID", "types.UUID", "uuid.UUID":
		return argKind{
			Flag: "String", Default: `""`, Helper: "String",
			Parse:     true,
			Imports:   []string{"github.com/google/uuid"},
			Supported: true,
		}, true
	case "openapi_types.Date", "types.Date":
		return argKind{
			Flag: "String", Default: `""`, Helper: "String",
			Wrap:      "types.Date{Time: %s}",
			Imports:   []string{"github.com/oapi-codegen/runtime/types"},
			Supported: false,
		}, false
	}

	// A bare name declared in the client package is usable only when it
	// resolves to a builtin — a string enum such as AgentStatus does, a struct
	// such as Sorting does not and has to be reached through a raw JSON body.
	if !strings.ContainsAny(base, ".[]{}*") && scalar != "" {
		k, ok := builtinKinds[scalar]
		if !ok {
			return argKind{}, false
		}
		k.Wrap = combineWrap(clientPkg+"."+base+"(%s)", k.Wrap)
		return k, true
	}
	return argKind{}, false
}

func combineWrap(outer, inner string) string {
	if inner == "" {
		return outer
	}
	return fmt.Sprintf(outer, inner)
}
