package codegen

import (
	"strings"
	"testing"
)

// scriptedWizard drives a wizard from a fixed list of answers.
func scriptedWizard(t *testing.T, answers ...string) (*Wizard, *strings.Builder) {
	t.Helper()
	ix, schema, specs, _ := loadFixtures(t)
	tree, err := LoadTree("../..")
	if err != nil {
		t.Fatal(err)
	}
	out := &strings.Builder{}
	in := strings.NewReader(strings.Join(answers, "\n") + "\n")
	return NewWizard(in, out, ix, tree, specs, schema), out
}

// TestWizardUnsetBehaviour covers the question this exists to ask: an optional
// value either stays out of the request or carries a default the command states
// itself.
func TestWizardUnsetBehaviour(t *testing.T) {
	t.Run("omit", func(t *testing.T) {
		w, out := scriptedWizard(t, "1")
		unset, def := "", ""
		if err := w.unsetBehaviour(&unset, &def, "enableDev", "*bool"); err != nil {
			t.Fatal(err)
		}
		if unset != UnsetOmit {
			t.Errorf("unset = %q, want %q", unset, UnsetOmit)
		}
		if def != "" {
			t.Errorf("omit should carry no default, got %q", def)
		}
		if !strings.Contains(out.String(), "when --enableDev is not given") {
			t.Errorf("prompt did not name the flag:\n%s", out)
		}
	})

	t.Run("always send with a default", func(t *testing.T) {
		w, _ := scriptedWizard(t, "2", "true")
		unset, def := "", ""
		if err := w.unsetBehaviour(&unset, &def, "enableDev", "*bool"); err != nil {
			t.Fatal(err)
		}
		if unset != UnsetDefault {
			t.Errorf("unset = %q, want %q", unset, UnsetDefault)
		}
		if def != "true" {
			t.Errorf("default = %q, want true", def)
		}
	})

	t.Run("the default offered is the type zero", func(t *testing.T) {
		// Answering "2" then accepting the offered default must give a usable
		// value, not an empty one.
		w, out := scriptedWizard(t, "2", "")
		unset, def := "", ""
		if err := w.unsetBehaviour(&unset, &def, "enableDev", "*bool"); err != nil {
			t.Fatal(err)
		}
		if def != "false" {
			t.Errorf("default = %q, want false", def)
		}
		if !strings.Contains(out.String(), "default value [false]") {
			t.Errorf("did not offer the type zero:\n%s", out)
		}
	})

	t.Run("a default that does not fit the type is re-asked", func(t *testing.T) {
		w, out := scriptedWizard(t, "2", "yesplease", "true")
		unset, def := "", ""
		if err := w.unsetBehaviour(&unset, &def, "enableDev", "*bool"); err != nil {
			t.Fatal(err)
		}
		if def != "true" {
			t.Errorf("default = %q, want true", def)
		}
		if !strings.Contains(out.String(), "not a valid *bool") {
			t.Errorf("did not complain about the bad value:\n%s", out)
		}
	})

	t.Run("a numeric default is accepted", func(t *testing.T) {
		w, _ := scriptedWizard(t, "2", "3600")
		unset, def := "", ""
		if err := w.unsetBehaviour(&unset, &def, "ttl", "*int"); err != nil {
			t.Fatal(err)
		}
		if unset != UnsetDefault || def != "3600" {
			t.Errorf("unset = %q, default = %q", unset, def)
		}
	})

	t.Run("an existing recipe value is offered back", func(t *testing.T) {
		// Re-running the wizard over a recipe must not silently drop what it
		// already says.
		w, out := scriptedWizard(t, "", "")
		unset, def := UnsetDefault, "true"
		if err := w.unsetBehaviour(&unset, &def, "enableDev", "*bool"); err != nil {
			t.Fatal(err)
		}
		if unset != UnsetDefault || def != "true" {
			t.Errorf("unset = %q, default = %q; existing values not preserved", unset, def)
		}
		if !strings.Contains(out.String(), "not given [2]") {
			t.Errorf("did not preselect the recorded choice:\n%s", out)
		}
	})
}

// TestWizardDefaultValue covers the flag-only case, where nothing prompts and a
// missing flag would otherwise send the type's zero silently.
func TestWizardDefaultValue(t *testing.T) {
	w, out := scriptedWizard(t, "10")
	def := ""
	if err := w.defaultValue(&def, "limit", "float32"); err != nil {
		t.Fatal(err)
	}
	if def != "10" {
		t.Errorf("default = %q, want 10", def)
	}
	if !strings.Contains(out.String(), "default value when --limit is not given") {
		t.Errorf("unexpected prompt:\n%s", out)
	}

	// A type no flag can carry has nothing to ask about.
	w2, out2 := scriptedWizard(t, "")
	skipped := ""
	if err := w2.defaultValue(&skipped, "members", "*map[string]string"); err != nil {
		t.Fatal(err)
	}
	if out2.String() != "" {
		t.Errorf("asked about a type it cannot carry:\n%s", out2)
	}
}
