package policyxml

import (
	"path/filepath"
	"testing"
)

func fixturePath(name string) string {
	return filepath.Join("..", "..", "testdata", "policies", name)
}

func findChild(node Node, name string) (Node, bool) {
	for _, child := range node.Children {
		if child.Name == name {
			return child, true
		}
	}
	return Node{}, false
}

func TestParseFileBasicPassThrough(t *testing.T) {
	doc, err := ParseFile(fixturePath("basic-pass-through.xml"))
	if err != nil {
		t.Fatalf("expected parse success, got error: %v", err)
	}

	if doc.SchemaVersion != "v1" {
		t.Fatalf("expected schema version v1, got %q", doc.SchemaVersion)
	}
	if doc.Raw.Name != "policies" {
		t.Fatalf("expected root node 'policies', got %q", doc.Raw.Name)
	}
	if len(doc.Raw.Children) != 4 {
		t.Fatalf("expected policies to contain 4 sections, got %d", len(doc.Raw.Children))
	}

	expectedSections := []string{"inbound", "backend", "outbound", "on-error"}
	for _, sectionName := range expectedSections {
		section, ok := findChild(doc.Raw, sectionName)
		if !ok {
			t.Fatalf("expected section %q in raw output", sectionName)
		}
		if len(section.Children) != 1 || section.Children[0].Name != "base" {
			t.Fatalf("expected section %q to contain single <base/>, got %#v", sectionName, section.Children)
		}
	}
}

func TestParseFileAddCorrelationID(t *testing.T) {
	doc, err := ParseFile(fixturePath("add-correlation-id.xml"))
	if err != nil {
		t.Fatalf("expected parse success, got error: %v", err)
	}

	inbound, ok := findChild(doc.Raw, "inbound")
	if !ok {
		t.Fatal("expected inbound section in raw output")
	}
	if len(inbound.Children) != 2 {
		t.Fatalf("expected inbound to contain 2 children, got %d", len(inbound.Children))
	}

	setHeader := inbound.Children[1]
	if setHeader.Name != "set-header" {
		t.Fatalf("expected second inbound child to be set-header, got %q", setHeader.Name)
	}
	if setHeader.Attrs["name"] != "x-correlation-id" {
		t.Fatalf("expected name attr x-correlation-id, got %q", setHeader.Attrs["name"])
	}
	if setHeader.Attrs["exists-action"] != "override" {
		t.Fatalf("expected exists-action attr override, got %q", setHeader.Attrs["exists-action"])
	}
	if len(setHeader.Children) != 1 || setHeader.Children[0].Name != "value" {
		t.Fatalf("expected set-header to contain one <value> child, got %#v", setHeader.Children)
	}
	if setHeader.Children[0].Text != "@(context.RequestId.ToString())" {
		t.Fatalf("expected value text to match policy expression, got %q", setHeader.Children[0].Text)
	}
}

func TestParseFileRateLimitBySubscription(t *testing.T) {
	doc, err := ParseFile(fixturePath("rate-limit-by-subscription.xml"))
	if err != nil {
		t.Fatalf("expected parse success, got error: %v", err)
	}

	inbound, ok := findChild(doc.Raw, "inbound")
	if !ok {
		t.Fatal("expected inbound section in raw output")
	}
	if len(inbound.Children) != 2 {
		t.Fatalf("expected inbound to contain 2 children, got %d", len(inbound.Children))
	}

	rateLimit := inbound.Children[1]
	if rateLimit.Name != "rate-limit-by-key" {
		t.Fatalf("expected second inbound child to be rate-limit-by-key, got %q", rateLimit.Name)
	}
	if rateLimit.Attrs["calls"] != "60" {
		t.Fatalf("expected calls attr 60, got %q", rateLimit.Attrs["calls"])
	}
	if rateLimit.Attrs["renewal-period"] != "60" {
		t.Fatalf("expected renewal-period attr 60, got %q", rateLimit.Attrs["renewal-period"])
	}
	if rateLimit.Attrs["counter-key"] == "" {
		t.Fatal("expected counter-key attr to be present")
	}
}

func TestParseRawXMLMalformedReturnsError(t *testing.T) {
	_, err := ParseRawXML([]byte("<policies><inbound></policies>"))
	if err == nil {
		t.Fatal("expected malformed xml parse error")
	}
}
