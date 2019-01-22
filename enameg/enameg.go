package enameg

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strings"

	"github.com/moznion/gowrtr/generator"
)

const annotation = "enameg"

type constantVal struct {
	Name       string
	CommentVal string
}
type constant struct {
	TypeName string
	Vals     []constantVal
}

// Generate returns packageName and generated functions by paths.
func Generate(paths []string) (string, string) {
	constMap, err := correctConstants(paths)
	if err != nil {
		log.Fatal(err)
	}

	var packageName string
	var constants []constant

	for _, path := range paths {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			log.Fatal(err)
		}

		if packageName == "" {
			packageName = f.Name.Name
		}
		if packageName != f.Name.Name {
			log.Fatal("error: multiple packages")
		}

		cmap := ast.NewCommentMap(fset, f, f.Comments)

		for node, commentGroups := range cmap {
			for _, cg := range commentGroups {
				if hasAnnotation(cg) {
					gen, ok := node.(*ast.GenDecl)
					if !ok || len(gen.Specs) <= 0 {
						continue
					}

					spec, ok := gen.Specs[0].(*ast.TypeSpec)
					if !ok {
						continue
					}

					c := newConst(spec.Name.Name, constMap)
					constants = append(constants, c)
				}
			}
		}
	}

	if len(constants) == 0 {
		return packageName, ""
	}

	generated, err := generateNameFunc(packageName, constants)
	if err != nil {
		log.Fatal(err)
	}
	return packageName, generated
}

func hasAnnotation(cg *ast.CommentGroup) bool {
	commentAnnotation := "+" + annotation
	for _, c := range cg.List {
		comment := strings.TrimSpace(strings.TrimLeft(c.Text, "//"))
		if strings.HasPrefix(comment, commentAnnotation) {
			return true
		}
	}
	return false
}

func correctConstants(paths []string) (map[string][]*ast.ValueSpec, error) {
	constMap := make(map[string][]*ast.ValueSpec)

	for _, path := range paths {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}

		constMap = correctConstantsInFile(constMap, f)
	}

	return constMap, nil
}

func correctConstantsInFile(constMap map[string][]*ast.ValueSpec, f *ast.File) map[string][]*ast.ValueSpec {
	for _, dec := range f.Decls {
		gen, ok := dec.(*ast.GenDecl)
		if !ok {
			continue
		}

		if gen.Tok != token.CONST {
			continue
		}

		for _, spec := range gen.Specs {
			val, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			valType, ok := val.Type.(*ast.Ident)
			if !ok {
				continue
			}

			typeName := valType.Name
			if constMap[typeName] == nil {
				constMap[typeName] = []*ast.ValueSpec{}
			}

			constMap[typeName] = append(constMap[typeName], val)
		}
	}

	return constMap
}

func newCommentVal(comment string) string {
	comment = strings.TrimSpace(strings.TrimLeft(comment, "//"))
	comment = strings.Fields(comment)[0]

	specialDelimiters := []string{".", "。"}
	for _, del := range specialDelimiters {
		comment = strings.Split(comment, del)[0]
	}

	return comment
}

func newConst(typeName string, constMap map[string][]*ast.ValueSpec) constant {
	nodes := constMap[typeName]
	vals := make([]constantVal, 0, len(nodes))

	for _, n := range nodes {
		vals = append(vals, constantVal{
			Name:       n.Names[0].Name,
			CommentVal: newCommentVal(n.Comment.List[0].Text),
		})
	}

	return constant{
		TypeName: typeName,
		Vals:     vals,
	}
}

func generateNameFunc(packageName string, consts []constant) (string, error) {
	g := generator.NewRoot(
		generator.NewPackage(packageName),
		generator.NewNewline(),
	)

	for _, c := range consts {
		caseStatements := make([]*generator.Case, 0, len(c.Vals))
		for _, v := range c.Vals {
			caseStatements = append(caseStatements, generator.NewCase(v.Name, generator.NewReturnStatement(fmt.Sprintf(`"%s"`, v.CommentVal))))
		}

		g = g.AddStatements(
			generator.NewComment(fmt.Sprintf(" Name returns the %s Name.", c.TypeName)),
			generator.NewFunc(
				generator.NewFuncReceiver("src", c.TypeName),
				generator.NewFuncSignature("Name").AddReturnTypes("string"),
			).AddStatements(
				generator.NewSwitch("src").
					AddCase(caseStatements...).
					Default(generator.NewDefaultCase(generator.NewReturnStatement(`fmt.Sprintf("%v", src)`))),
			),
			generator.NewNewline(),
		)
	}

	generated, err := g.Gofmt("-s").Goimports().Generate(0)
	if err != nil {
		return "", err
	}

	return generated, nil
}
