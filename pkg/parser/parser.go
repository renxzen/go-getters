package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/renxzen/go-getters/pkg/strutils"
	"github.com/renxzen/go-getters/pkg/types"
)

// Parser handles parsing of Go source files.
type Parser struct {
	fset *token.FileSet
}

// New creates a new Parser instance.
func New() *Parser {
	return &Parser{
		fset: token.NewFileSet(),
	}
}

// ParseDirectory parses all Go files in the specified directory and returns struct information.
func (p *Parser) ParseDirectory(path string) (*types.ParseResult, error) {
	pkgs, err := parser.ParseDir(p.fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse directory %s: %w", path, err)
	}

	imports := make(map[string]*types.ImportInfo)
	structs := make(map[string]*types.StructInfo)

	var packageName string
	for _, pkg := range pkgs {
		if packageName == "" {
			packageName = pkg.Name
		}

		for _, file := range pkg.Files {
			p.parseImports(file, imports)

			ast.Inspect(file, func(n ast.Node) bool {
				typeSpec, ok := n.(*ast.TypeSpec)
				if !ok {
					return true
				}

				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					return true
				}

				structInfo := p.parseStruct(typeSpec.Name.Name, structType)
				structs[structInfo.Name] = structInfo

				return true
			})
		}
	}

	return &types.ParseResult{
		PackageName: packageName,
		Structs:     structs,
		Imports:     imports,
	}, nil
}

// parseStruct parses a single struct and returns its information.
func (p *Parser) parseStruct(structName string, structType *ast.StructType) *types.StructInfo {
	structInfo := &types.StructInfo{
		Name:   structName,
		Fields: make([]types.FieldInfo, 0, len(structType.Fields.List)),
	}

	for _, field := range structType.Fields.List {
		// Skip fields without names (embedded fields)
		if len(field.Names) == 0 {
			continue
		}

		fieldName := field.Names[0].Name

		// Parse field type information
		fieldInfo := p.parseFieldType(fieldName, field.Type)
		structInfo.Fields = append(structInfo.Fields, fieldInfo)
	}

	return structInfo
}

// parseFieldType parses field type information.
func (p *Parser) parseFieldType(fieldName string, fieldType ast.Expr) types.FieldInfo {
	fieldInfo := types.FieldInfo{
		Name:       fieldName,
		IsExported: strutils.IsCapitalized(fieldName),
	}

	switch t := fieldType.(type) {
	case *ast.StarExpr:
		// Recursively parse the underlying type
		underlyingField := p.parseFieldType("", t.X)
		fieldInfo.Type = "*" + underlyingField.Type
		fieldInfo.UnderlyingType = underlyingField.UnderlyingType
		fieldInfo.IsPointer = true

		// Preserve other properties from the underlying type
		fieldInfo.IsSlice = underlyingField.IsSlice
		fieldInfo.IsMap = underlyingField.IsMap
		for _, requiredImport := range underlyingField.RequiredImports {
			fieldInfo.AddRequiredImport(requiredImport)
		}
	case *ast.Ident:
		fieldInfo.Type = t.Name
		fieldInfo.UnderlyingType = t.Name
	case *ast.SelectorExpr:
		// Recursively handle the X part of the selector
		xType := p.parseExprType(t.X)
		fieldInfo.Type = xType + "." + t.Sel.Name
		fieldInfo.UnderlyingType = xType + "." + t.Sel.Name
		fieldInfo.AddRequiredImport(xType)
	case *ast.ArrayType:
		// Handle slices and arrays
		elementType := p.parseExprType(t.Elt)
		fieldInfo.Type = "[]" + elementType
		fieldInfo.UnderlyingType = "[]" + elementType
		fieldInfo.IsSlice = true

		// Recursively handle the element type
		elementField := p.parseFieldType("", t.Elt)
		fieldInfo.RequiredImports = elementField.RequiredImports

		// if Slices and Arrays need to be handled differently
		// we need to check t.Len. if it is nil, it's a slice,
		// otherwise it's an array
	case *ast.MapType:
		mapType := p.parseExprType(t)
		fieldInfo.Type = mapType
		fieldInfo.UnderlyingType = mapType
		fieldInfo.IsMap = true

		// Recursively check for imports in map key and value types
		keyField := p.parseFieldType("", t.Key)
		for _, requiredImport := range keyField.RequiredImports {
			fieldInfo.AddRequiredImport(requiredImport)
		}

		valueField := p.parseFieldType("", t.Value)
		for _, requiredImport := range valueField.RequiredImports {
			fieldInfo.AddRequiredImport(requiredImport)
		}
	default:
		// Handle other complex types
		fieldInfo.Type = "any"
		fieldInfo.UnderlyingType = "any"
	}

	return fieldInfo
}

// parseExprType parses an expression and returns its string representation.
// This method unifies the functionality of parseElementType and parseExprType.
func (p *Parser) parseExprType(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + p.parseExprType(t.X)
	case *ast.SelectorExpr:
		xType := p.parseExprType(t.X)
		return xType + "." + t.Sel.Name
	case *ast.ArrayType:
		elementType := p.parseExprType(t.Elt)
		return "[]" + elementType
	case *ast.MapType:
		keyType := p.parseExprType(t.Key)
		valueType := p.parseExprType(t.Value)
		return fmt.Sprintf("map[%s]%s", keyType, valueType)
	default:
		return "any"
	}
}

// parseImports parses import declarations from an AST file and adds them to the imports map.
func (p *Parser) parseImports(file *ast.File, imports map[string]*types.ImportInfo) {
	for _, imp := range file.Imports {
		// Remove quotes
		importPath := imp.Path.Value[1 : len(imp.Path.Value)-1]

		var (
			alias     string
			isAliased bool
		)
		if imp.Name != nil {
			alias = imp.Name.Name
			isAliased = true
		} else {
			// Extract package name from import path
			parts := strings.Split(importPath, "/")
			alias = parts[len(parts)-1]
		}

		imports[alias] = &types.ImportInfo{
			Alias:     alias,
			IsAliased: isAliased,
			Path:      importPath,
		}
	}
}
