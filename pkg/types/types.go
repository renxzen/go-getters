package types

// ParseResult contains the parsing results including package name and structs.
type ParseResult struct {
	PackageName string
	Structs     map[string]*StructInfo
	Imports     map[string]*ImportInfo
}

// StructInfo contains information about a struct.
type StructInfo struct {
	Name   string
	Fields []FieldInfo
}

type ImportInfo struct {
	Alias     string
	IsAliased bool
	Path      string
}

type FieldInfo struct {
	Name           string // Field name
	Type           string // Field type as string
	UnderlyingType string // Underlying type for pointers
	IsPointer      bool   // Whether the field is a pointer
	IsExported     bool   // Whether the field is exported
	IsSlice        bool   // Whether the field is a slice
	RequiredImport string // Import alias for package-qualified types
}

func (f FieldInfo) IsPrimitive() bool {
	switch f.UnderlyingType {
	case "string", "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64", "bool", "byte", "rune":
		return true
	default:
		return false
	}
}

func (f FieldInfo) GetZerovalue() string {
	fieldType := f.Type
	if f.IsPointer {
		// Only dereference for primitives and specific types that should be dereferenced
		if f.IsPrimitive() {
			fieldType = f.UnderlyingType
		} else {
			// For other pointer types (like *Example), return nil
			return "nil"
		}
	}
	// Handle slice types
	if len(fieldType) > 2 && fieldType[:2] == "[]" {
		return "nil"
	}

	switch fieldType {
	case "string":
		return `""`
	case "int", "int8", "int16", "int32", "int64":
		return "0"
	case "uint", "uint8", "uint16", "uint32", "uint64":
		return "0"
	case "float32", "float64":
		return "0.0"
	case "bool":
		return "false"
	case "byte":
		return "0"
	case "rune":
		return "0"
	default:
		// For custom types and package-qualified types, use the zero value syntax
		// This works for structs, interfaces, and other custom types
		return fieldType + "{}"
	}
}

// ShouldDereference returns true if this pointer field should be dereferenced in getters
// Currently applies to time.Time and other value types from standard library
func (f FieldInfo) ShouldDereference() bool {
	if !f.IsPointer {
		return false
	}

	// TODO: make this configurable
	// Check for specific types that should be dereferenced
	switch f.UnderlyingType {
	case "time.Time":
		return true
	default:
		return false
	}
}
