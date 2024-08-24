package reggi

import (
	"reflect"
	"testing"
)

func TestParser(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedNode Node
	}{
		{
			name: "simple string", input: "aBc", expectedNode: &Group{
				Nodes: []Node{
					&CharLiteral{Char: 'a'},
					&CharLiteral{Char: 'B'},
					&CharLiteral{Char: 'c'},
				},
			},
		},
		{
			name: "wildcard character", input: "ab.", expectedNode: &Group{
				Nodes: []Node{
					&CharLiteral{Char: 'a'},
					&CharLiteral{Char: 'b'},
					&Wildcard{},
				},
			},
		},
		{
			name:  "branches",
			input: "ab|cd|ef",
			expectedNode: &Branch{
				Nodes: []Node{
					&Group{Nodes: []Node{&CharLiteral{Char: 'a'}, &CharLiteral{Char: 'b'}}},
					&Group{Nodes: []Node{&CharLiteral{Char: 'c'}, &CharLiteral{Char: 'd'}}},
					&Group{Nodes: []Node{&CharLiteral{Char: 'e'}, &CharLiteral{Char: 'f'}}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := lex(tt.input)
			p := NewParser(tokens)
			result := p.Parse()

			if !reflect.DeepEqual(result, tt.expectedNode) {
				t.Fatalf("Expected:\n%+v\n\nGot:\n%+v\n", tt.expectedNode, result)
			}
		})
	}
}
