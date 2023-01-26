package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/snapp-incubator/barat/internal/config"
)

func CheckCodeForLocalizationFunctions(tomlFiles map[string]map[string]interface{}, projectPath string) (errs []error) {

	// recursive search for all files in projectPath
	entries, err := os.ReadDir(projectPath)
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
				tomlFiles, projectPath+"/"+entry.Name())...)
		} else if len(entry.Name()) > 3 && entry.Name()[len(entry.Name())-3:] == ".go" {
			errs = append(errs, fileParser(tomlFiles, projectPath+"/"+entry.Name())...)
		}
	}

	return errs
}

func fileParser(
	tomlFiles map[string]map[string]interface{},
	filePath string) (errs []error) {
	fileSet := token.NewFileSet()

	file, err := os.ReadFile(filePath)
	if err != nil {
		return []error{err}
	}

	f, err := parser.ParseFile(fileSet, "", file, parser.ParseComments)
	if err != nil {
		return []error{err}
	}

	for _, decl := range f.Decls {
		switch decl.(type) {
		case *ast.FuncDecl: // find function declaration
			if decl.(*ast.FuncDecl).Body != nil {
				for _, spec := range decl.(*ast.FuncDecl).Body.List {
					errs = append(errs,
						parsLineOfCode(spec, tomlFiles)...)
				}
			}
		}
	}

	return errs
}

func parsLineOfCode(stmt ast.Stmt, tomlFiles map[string]map[string]interface{}) (errs []error) {
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
						errs = append(errs,
							checkKeyInTomlFiles(tomlFiles, rhs.(*ast.CallExpr).Args[index].(*ast.BasicLit).Value)...)
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
			errs = append(errs, parsLineOfCode(stmt, tomlFiles)...)
		}
	case *ast.RangeStmt:
		for _, stmt := range stmt.(*ast.RangeStmt).Body.List {
			errs = append(errs, parsLineOfCode(stmt, tomlFiles)...)
		}
	case *ast.IfStmt:
		for _, stmt := range stmt.(*ast.IfStmt).Body.List {
			errs = append(errs, parsLineOfCode(stmt, tomlFiles)...)
		}
	case *ast.SwitchStmt:
		for _, stmt := range stmt.(*ast.SwitchStmt).Body.List {
			errs = append(errs, parsLineOfCode(stmt, tomlFiles)...)
		}
	}

	return errs
}

func checkKeyInTomlFiles(tomlFiles map[string]map[string]interface{}, arg string) (errs []error) {
	arg = strings.Trim(arg, "\"")
	// TODO: check if key not in excluding list
	for lang, val := range tomlFiles {
		flag := false
		for key := range val {
			if key == arg {
				flag = true
			}
		}
		if !flag {
			errs = append(errs,
				fmt.Errorf("key \"%s\" is not valid in language \"%s\"", arg, lang),
			)
		}
	}
	return errs
}
