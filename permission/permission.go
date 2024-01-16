package permission

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

const (
	permTagName   = "@Permission"
	routerTagName = "@Router"
)

type doc struct {
	PermissionDesc string `json:"desc"`
	Code           string `json:"code"`
	Type           string `json:"type"`
}

// Generate 生成ULA所需的权限列表json文件
// dir: 读取的go文件
// output: json文件保存路径
func Generate(dir, output string) error {
	files := listDir(dir)

	var docs []doc
	for _, f := range files {
		l, err := parseFile(f)
		if err != nil {
			return err
		}
		docs = append(docs, l...)
	}

	content, err := json.MarshalIndent(docs, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(output, "permission.json"), content, 0644)
}

func listDir(path string) []string {
	files, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}
	var fileList []string
	for _, file := range files {
		if file.IsDir() {
			fileList = append(fileList, listDir(filepath.Join(path, file.Name()))...)
		} else {
			fileList = append(fileList, filepath.Join(path, file.Name()))
		}
	}
	return fileList
}

func parseFile(filename string) ([]doc, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var docs []doc

	ast.Inspect(f, func(n ast.Node) bool {
		if x, ok := n.(*ast.FuncDecl); ok {
			d := parseFuncComment(x)
			if d != nil {
				docs = append(docs, *d)
			}
		}
		return true
	})

	return docs, nil
}

func parseFuncComment(f *ast.FuncDecl) *doc {
	var d doc
	if f.Doc == nil {
		return nil
	}

	for _, c := range f.Doc.List {
		if strings.Contains(c.Text, permTagName) {
			d.PermissionDesc = parseTaggedComment(c.Text, permTagName)[0]
		} else if strings.Contains(c.Text, routerTagName) {
			router := parseTaggedComment(c.Text, routerTagName)
			d.Code = fmt.Sprintf("%s:%s", strings.Trim(router[1], "[]"), router[0])
		}
	}
	if d.PermissionDesc == "" || d.Code == "" {
		return nil
	}
	d.Type = "system"
	return &d
}

func parseTaggedComment(text string, tag string) []string {
	parts := strings.SplitN(text, tag, 2)
	a := strings.SplitN(strings.TrimSpace(parts[1]), " ", 2)
	return a
}
