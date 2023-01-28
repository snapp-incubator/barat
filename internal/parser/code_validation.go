package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/snapp-incubator/barat/internal/config"
)

// CheckCodeForLocalizationFunctions is function for finding go files and recursive search in directories.
func CheckCodeForLocalizationFunctions(mapLangToToml map[Language]TomlFile, path string) (errs []error) {
	// recursive search for all files in projectPath
	entries, err := os.ReadDir(path)
	if err != nil {
		return []error{err}
	}

ENTRY:
	for _, entry := range entries {
		if entry.IsDir() {
			// check list of exclude folders from search
			for _, excludeFolder := range config.C.Exclude.Folders {
				if entry.Name() == excludeFolder {
					continue ENTRY
				}
			}
			errs = append(errs, CheckCodeForLocalizationFunctions(
				mapLangToToml, filepath.Join(path, entry.Name()))...)
		} else if filepath.Ext(entry.Name()) == ".go" {
			errs = append(errs, fileParser(mapLangToToml, filepath.Join(path, entry.Name()))...)
		}
	}

	return errs
}

// fileParser is function for parsing go files and finding functions.
func fileParser(
	mapLangToToml map[Language]TomlFile,
	filePath string) (errs []error) {
	fileSet := token.NewFileSet()

	file, err := os.ReadFile(filePath)
	if err != nil {
		err = fmt.Errorf("error in reading file: %s", err)
		return []error{err}
	}

	f, err := parser.ParseFile(fileSet, "", file, parser.ParseComments)
	if err != nil {
		err = fmt.Errorf("error in parsing file: %s", err)
		return []error{err}
	}

	for _, decl := range f.Decls {
		switch decl.(type) {
		case *ast.FuncDecl: // find function declaration
			if decl.(*ast.FuncDecl).Body != nil {
				for _, spec := range decl.(*ast.FuncDecl).Body.List {
					errs = append(errs,
						parseCode(spec, mapLangToToml)...)
				}
			}
		}
	}

	return errs
}

// parseCode is function for parsing each token of code and finding localization functions.
func parseCode(stmt ast.Stmt, mapLangToToml map[Language]TomlFile) (errs []error) {
	switch stmt.(type) {
	case *ast.AssignStmt: // find assignment statement (e.g. var a = "hello")
		for _, rhs := range stmt.(*ast.AssignStmt).Rhs { // iterate right side of assignment statement
			switch rhs.(type) {
			case *ast.CallExpr: // find call expression (e.g. p.getMessage("hello"))
				isSelectedFunction := false
				index := 0 // index of MessageID in args of function that found in MessageFunc
				fn := ""   // name of function that found in MapFunctionNamesToArgNo
				switch rhs.(*ast.CallExpr).Fun.(type) {
				case *ast.SelectorExpr: // find selector expression (e.g. package.getMessage)
					for _, m := range config.C.MessageFuncs {
						if rhs.(*ast.CallExpr).Fun.(*ast.SelectorExpr).Sel.Name == m.Name {
							isSelectedFunction = true
							index = m.MessageIDNo
							fn = m.Name
							break
						}
					}
				case *ast.Ident: // find identifier (e.g. getMessage)
					for _, m := range config.C.MessageFuncs {
						if rhs.(*ast.CallExpr).Fun.(*ast.Ident).Name == m.Name {
							isSelectedFunction = true
							index = m.MessageIDNo
							fn = m.Name
							break
						}
					}
				}

				if isSelectedFunction {
					// validate number of arguments
					if len(rhs.(*ast.CallExpr).Args) < index+1 {
						return []error{
							fmt.Errorf(
								"not enough arguments: %s function needs at least %d arguments but %d arguments found",
								fn, index+1, len(rhs.(*ast.CallExpr).Args))}
					}
					switch rhs.(*ast.CallExpr).Args[index].(type) {
					case *ast.Ident: // find identifier (e.g. hello in getMessage(hello))
						// TODO: check if value is valid in toml files
					case *ast.BasicLit: // find basic literal (e.g. "world" in getMessage("world"))
						messageID := rhs.(*ast.CallExpr).Args[index].(*ast.BasicLit).Value
						errs = append(errs,
							checkKeyInTomlFiles(mapLangToToml, messageID)...)
					case *ast.SelectorExpr: // find selector expression (e.g. p.hello in getMessage(p.hello))
						// TODO: check if value is valid in toml files
					}

				}
			}
		}
	}

	switch stmt.(type) {
	case *ast.ForStmt:
		for _, stmt := range stmt.(*ast.ForStmt).Body.List {
			errs = append(errs, parseCode(stmt, mapLangToToml)...)
		}
	case *ast.RangeStmt:
		for _, stmt := range stmt.(*ast.RangeStmt).Body.List {
			errs = append(errs, parseCode(stmt, mapLangToToml)...)
		}
	case *ast.IfStmt:
		for _, stmt := range stmt.(*ast.IfStmt).Body.List {
			errs = append(errs, parseCode(stmt, mapLangToToml)...)
		}
	case *ast.SwitchStmt:
		for _, stmt := range stmt.(*ast.SwitchStmt).Body.List {
			errs = append(errs, parseCode(stmt, mapLangToToml)...)
		}
	}

	return errs
}

// checkKeyInTomlFiles is function for checking if key exists in toml files.
func checkKeyInTomlFiles(mapLangToToml map[Language]TomlFile, key string) (errs []error) {
	key = strings.Trim(key, "\"")
	messageID := MessageID(key)

	// TODO: check if key not in excluding list
	for language, tomlFiles := range mapLangToToml {
		flag := false
		for mID := range tomlFiles {
			if mID == messageID {
				flag = true
			}
		}
		if !flag {
			errs = append(errs,
				fmt.Errorf("MessageID \"%s\" is not valid in language \"%s\"", messageID, language),
			)
		}
	}
	return errs
}
