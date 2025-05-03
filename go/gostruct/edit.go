package gostruct

import (
	"go/ast"
	"go/token"
)

// UpdateStructFields parses the AST struct and merges it with the desired fields.
// It returns the merged struct definition without modifying the AST.
//
// Parameters:
// - structType: AST representation of the struct to analyze
// - name: name of the struct being updated
// - fields: the complete set of fields the struct should have after update
// - reserveFields: map of field names that should be preserved even if not in fields list
//
// Returns the merged struct definition.
func UpdateStructFields(fset *token.FileSet, structType *ast.StructType, name string, fields []FieldDef, reserveFields map[string]bool) StructDef {
	// Parse AST into StructDef
	current := ParseStruct(fset, structType, name)

	// Create desired StructDef
	desired := StructDef{
		Name:   name,
		Fields: fields,
	}

	// Merge the structs
	result := MergeStructs(current, desired, reserveFields)

	return result
}
