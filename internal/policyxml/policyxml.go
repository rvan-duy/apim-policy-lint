package policyxml

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
)

// Node is a generic XML AST node.
type Node struct {
	Name     string            `json:"name"`
	Attrs    map[string]string `json:"attrs,omitempty"`
	Text     string            `json:"text,omitempty"`
	Children []Node            `json:"children,omitempty"`
}

// Document is the output contract for XML policy parsing.
type Document struct {
	SchemaVersion string `json:"schemaVersion"`
	SourceFile    string `json:"sourceFile"`
	Raw           Node   `json:"raw"`
}

// ParseFile parses an XML file into the raw policy document structure.
func ParseFile(filePath string) (Document, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return Document{}, fmt.Errorf("failed reading xml file: %w", err)
	}

	raw, err := ParseRawXML(data)
	if err != nil {
		return Document{}, err
	}

	return Document{
		SchemaVersion: "v1",
		SourceFile:    filePath,
		Raw:           raw,
	}, nil
}

// ParseRawXML decodes XML bytes into a generic AST while preserving node order.
func ParseRawXML(data []byte) (Node, error) {
	dec := xml.NewDecoder(bytes.NewReader(data))

	stack := make([]*Node, 0, 16)
	var root *Node

	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return Node{}, fmt.Errorf("failed parsing xml: %w", err)
		}

		switch t := tok.(type) {
		case xml.StartElement:
			n := &Node{
				Name:  t.Name.Local,
				Attrs: attrsToMap(t.Attr),
			}
			stack = append(stack, n)
		case xml.CharData:
			if len(stack) == 0 {
				continue
			}

			text := strings.TrimSpace(string(t))
			if text == "" {
				continue
			}

			current := stack[len(stack)-1]
			if current.Text == "" {
				current.Text = text
			} else {
				current.Text += " " + text
			}
		case xml.EndElement:
			if len(stack) == 0 {
				return Node{}, fmt.Errorf("failed parsing xml: unexpected closing tag %q", t.Name.Local)
			}

			current := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if len(stack) == 0 {
				if root != nil {
					return Node{}, fmt.Errorf("failed parsing xml: multiple root elements are not allowed")
				}
				root = current
				continue
			}

			parent := stack[len(stack)-1]
			parent.Children = append(parent.Children, *current)
		}
	}

	if len(stack) != 0 {
		return Node{}, fmt.Errorf("failed parsing xml: malformed xml, unclosed elements remain")
	}
	if root == nil {
		return Node{}, fmt.Errorf("failed parsing xml: empty xml document")
	}

	return *root, nil
}

func attrsToMap(attrs []xml.Attr) map[string]string {
	if len(attrs) == 0 {
		return nil
	}

	out := make(map[string]string, len(attrs))
	for _, attr := range attrs {
		out[attr.Name.Local] = attr.Value
	}
	return out
}
