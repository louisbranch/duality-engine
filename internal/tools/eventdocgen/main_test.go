package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFindModuleRoot(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.com/test"), 0o644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}
	nested := filepath.Join(root, "a", "b")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatalf("mkdir nested: %v", err)
	}

	got, err := findModuleRoot(nested)
	if err != nil {
		t.Fatalf("findModuleRoot returned error: %v", err)
	}
	if got != root {
		t.Fatalf("expected root %s, got %s", root, got)
	}
}

func TestFindModuleRootMissing(t *testing.T) {
	root := t.TempDir()
	_, err := findModuleRoot(root)
	if err == nil {
		t.Fatal("expected error when go.mod is missing")
	}
}

func TestParsePackage(t *testing.T) {
	root := t.TempDir()
	pkgDir := filepath.Join(root, "pkg")
	if err := os.MkdirAll(pkgDir, 0o755); err != nil {
		t.Fatalf("mkdir pkg: %v", err)
	}
	src := strings.Join([]string{
		"package sample",
		"",
		"import \"example.com/event\"",
		"",
		"type Type string",
		"",
		"const (",
		"\tTypeFoo Type = \"foo\"",
		"\tTypeBar event.Type = \"bar\"",
		"\tTypeIgnored string = \"ignored\"",
		")",
		"",
		"type FooPayload struct {",
		"\tID string `json:\"id\"`",
		"\tName string",
		"}",
		"",
		"type Ignored struct {",
		"\tValue string",
		"}",
	}, "\n")
	if err := os.WriteFile(filepath.Join(pkgDir, "sample.go"), []byte(src), 0o644); err != nil {
		t.Fatalf("write sample.go: %v", err)
	}

	defs, err := parsePackage(pkgDir, root, "Core")
	if err != nil {
		t.Fatalf("parsePackage returned error: %v", err)
	}
	if len(defs.Events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(defs.Events))
	}
	payload, ok := defs.Payloads["FooPayload"]
	if !ok {
		t.Fatal("expected FooPayload in payloads")
	}
	if payload.Owner != "Core" {
		t.Fatalf("expected payload owner Core, got %s", payload.Owner)
	}
	if len(payload.Fields) != 2 {
		t.Fatalf("expected 2 payload fields, got %d", len(payload.Fields))
	}
	if payload.Fields[0].JSONTag != "json:\"id\"" {
		t.Fatalf("expected json tag on first field, got %s", payload.Fields[0].JSONTag)
	}
}

func TestPayloadNameForEvent(t *testing.T) {
	if got := payloadNameForEvent("TypeCampaignCreated", "Core"); got != "CampaignCreatedPayload" {
		t.Fatalf("unexpected payload name: %s", got)
	}
	if got := payloadNameForEvent("EventTypeSessionStarted", "Daggerheart"); got != "SessionStartedPayload" {
		t.Fatalf("unexpected payload name: %s", got)
	}
	if got := payloadNameForEvent("EventTypeSessionStarted", "Unknown"); got != "" {
		t.Fatalf("expected empty payload name, got %s", got)
	}
}

func TestRenderCatalog(t *testing.T) {
	defs := packageDefs{
		Events: []eventDef{{
			Owner:     "Core",
			Name:      "TypeFoo",
			Value:     "foo",
			DefinedAt: "internal/foo.go:10",
		}},
		Payloads: map[string]payloadDef{
			"FooPayload": {
				Owner:     "Core",
				Name:      "FooPayload",
				DefinedAt: "internal/foo.go:20",
				Fields: []payloadField{{
					Name:    "ID",
					Type:    "string",
					JSONTag: "json:\"id\"",
				}},
			},
			"UnusedPayload": {
				Owner:     "Core",
				Name:      "UnusedPayload",
				DefinedAt: "internal/foo.go:30",
			},
		},
	}
	emitters := map[string][]string{
		"TypeFoo": {"internal/emit.go:12"},
	}
	output, err := renderCatalog([]packageDefs{defs}, emitters)
	if err != nil {
		t.Fatalf("renderCatalog returned error: %v", err)
	}
	checks := []string{
		"## Core Events",
		"### `foo` (`TypeFoo`)",
		"Payload: `FooPayload`",
		"`ID (json:\"id\")`: `string`",
		"### Unmapped Payloads",
		"`UnusedPayload`",
		"Emitters:",
		"`internal/emit.go:12`",
	}
	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Fatalf("expected output to contain %q", check)
		}
	}
}

func TestScanEmitters(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "emitters")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir emitters: %v", err)
	}
	src := strings.Join([]string{
		"package sample",
		"",
		"import \"example.com/event\"",
		"",
		"func emit() {",
		"\t_ = event.Event{Type: event.TypeFoo}",
		"\t_ = event.Event{Type: TypeBar}",
		"}",
	}, "\n")
	path := filepath.Join(dir, "emit.go")
	if err := os.WriteFile(path, []byte(src), 0o644); err != nil {
		t.Fatalf("write emit.go: %v", err)
	}

	emitters, err := scanEmitters(root, root)
	if err != nil {
		t.Fatalf("scanEmitters returned error: %v", err)
	}
	if len(emitters["TypeFoo"]) != 1 {
		t.Fatalf("expected one emitter for TypeFoo, got %d", len(emitters["TypeFoo"]))
	}
	if len(emitters["TypeBar"]) != 1 {
		t.Fatalf("expected one emitter for TypeBar, got %d", len(emitters["TypeBar"]))
	}
	if !strings.HasPrefix(emitters["TypeFoo"][0], "emitters/emit.go:") {
		t.Fatalf("unexpected emitter path: %s", emitters["TypeFoo"][0])
	}
}

func TestFormatPosition(t *testing.T) {
	pos := formatPosition(tokenPosition("/root/pkg/file.go", 12), "/root")
	if pos != "pkg/file.go:12" {
		t.Fatalf("expected formatted position, got %s", pos)
	}
}

func tokenPosition(file string, line int) token.Position {
	return token.Position{Filename: file, Line: line}
}

func TestResolveRoot(t *testing.T) {
	t.Run("explicit flag", func(t *testing.T) {
		got, err := resolveRoot("/some/path")
		if err != nil {
			t.Fatal(err)
		}
		if got != "/some/path" {
			t.Errorf("got %q, want /some/path", got)
		}
	})

	t.Run("empty flag uses cwd", func(t *testing.T) {
		// From the project root, findModuleRoot should succeed.
		got, err := resolveRoot("")
		if err != nil {
			t.Fatal(err)
		}
		if got == "" {
			t.Error("expected non-empty root")
		}
	})
}

func TestSelectValueExpr(t *testing.T) {
	a := &ast.BasicLit{Kind: token.STRING, Value: `"a"`}
	b := &ast.BasicLit{Kind: token.STRING, Value: `"b"`}

	t.Run("empty list", func(t *testing.T) {
		if got := selectValueExpr(nil, 0); got != nil {
			t.Error("expected nil for empty list")
		}
	})

	t.Run("single value any index", func(t *testing.T) {
		if got := selectValueExpr([]ast.Expr{a}, 5); got != a {
			t.Error("expected single value to always be returned")
		}
	})

	t.Run("multi value in range", func(t *testing.T) {
		if got := selectValueExpr([]ast.Expr{a, b}, 1); got != b {
			t.Error("expected second element")
		}
	})

	t.Run("multi value out of range", func(t *testing.T) {
		if got := selectValueExpr([]ast.Expr{a, b}, 5); got != nil {
			t.Error("expected nil for out-of-range index")
		}
	})
}

func TestEventNameFromExpr(t *testing.T) {
	t.Run("selector expr", func(t *testing.T) {
		e := &ast.SelectorExpr{
			X:   &ast.Ident{Name: "event"},
			Sel: &ast.Ident{Name: "TypeFoo"},
		}
		if got := eventNameFromExpr(e); got != "TypeFoo" {
			t.Errorf("got %q, want TypeFoo", got)
		}
	})

	t.Run("ident expr", func(t *testing.T) {
		e := &ast.Ident{Name: "TypeBar"}
		if got := eventNameFromExpr(e); got != "TypeBar" {
			t.Errorf("got %q, want TypeBar", got)
		}
	})

	t.Run("other expr", func(t *testing.T) {
		e := &ast.BasicLit{Kind: token.STRING, Value: `"literal"`}
		if got := eventNameFromExpr(e); got != "" {
			t.Errorf("got %q, want empty string", got)
		}
	})
}

func TestUnmappedPayloads(t *testing.T) {
	t.Run("nil payloads", func(t *testing.T) {
		result := unmappedPayloads(nil, nil)
		if result != nil {
			t.Error("expected nil for nil payloads")
		}
	})

	t.Run("all used", func(t *testing.T) {
		payloads := map[string]payloadDef{"A": {Name: "A"}}
		used := map[string]struct{}{"A": {}}
		result := unmappedPayloads(payloads, used)
		if len(result) != 0 {
			t.Errorf("expected 0 unmapped, got %d", len(result))
		}
	})

	t.Run("some unmapped", func(t *testing.T) {
		payloads := map[string]payloadDef{
			"A": {Name: "A"},
			"B": {Name: "B"},
			"C": {Name: "C"},
		}
		used := map[string]struct{}{"A": {}}
		result := unmappedPayloads(payloads, used)
		if len(result) != 2 {
			t.Fatalf("expected 2 unmapped, got %d", len(result))
		}
		// Should be sorted by name.
		if result[0].Name != "B" || result[1].Name != "C" {
			t.Errorf("expected [B, C], got [%s, %s]", result[0].Name, result[1].Name)
		}
	})
}

func TestExprString(t *testing.T) {
	fset := token.NewFileSet()
	e := &ast.Ident{Name: "MyType"}
	got := exprString(fset, e)
	if got != "MyType" {
		t.Errorf("got %q, want MyType", got)
	}
}

func TestParsePayloadFields(t *testing.T) {
	t.Run("nil fields", func(t *testing.T) {
		result := parsePayloadFields(nil, token.NewFileSet())
		if result != nil {
			t.Error("expected nil for nil fields")
		}
	})

	t.Run("embedded field skipped", func(t *testing.T) {
		// Parse a struct with an embedded field (no names).
		src := `package x; type S struct { Embedded; Name string }` //nolint
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, "test.go", src, 0)
		if err != nil {
			t.Fatal(err)
		}
		var structType *ast.StructType
		ast.Inspect(file, func(n ast.Node) bool {
			if st, ok := n.(*ast.StructType); ok {
				structType = st
				return false
			}
			return true
		})
		if structType == nil {
			t.Fatal("struct not found")
		}
		fields := parsePayloadFields(structType.Fields, fset)
		// Only "Name" should be returned (embedded "Embedded" is skipped).
		if len(fields) != 1 {
			t.Fatalf("expected 1 field, got %d", len(fields))
		}
		if fields[0].Name != "Name" {
			t.Errorf("got field %q, want Name", fields[0].Name)
		}
	})
}

func TestRenderCatalog_EmptyPackage(t *testing.T) {
	// A package with no events should produce no section.
	defs := packageDefs{Events: nil, Payloads: map[string]payloadDef{}}
	output, err := renderCatalog([]packageDefs{defs}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(output, "## ") {
		t.Error("expected no section header for empty package")
	}
}

func TestRenderCatalog_NoPayload(t *testing.T) {
	defs := packageDefs{
		Events: []eventDef{{
			Owner:     "Core",
			Name:      "TypeOrphan",
			Value:     "orphan",
			DefinedAt: "foo.go:1",
		}},
		Payloads: map[string]payloadDef{},
	}
	output, err := renderCatalog([]packageDefs{defs}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output, "Payload: not found") {
		t.Error("expected 'Payload: not found' for event without matching payload")
	}
}
