package main

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/davyxu/gosproto/meta"
)

const codeTemplate = `// Generated by github.com/davyxu/gosproto/sprotogen
// DO NOT EDIT!{{range .Protos}}
// Source: {{.FileName}}{{end}}

package {{.PackageName}}

import (
	"reflect"
)

{{range .Protos}}
// {{.FileName}}
{{range .Structs}}
type {{.Name}} struct{
	{{range .GoFields}}
		{{.FieldName}} {{.GoTypeName}} {{.GoTags}} 
	{{end}}
}
{{end}}
{{end}}

var SProtoStructs = []reflect.Type{
{{range .Protos}}// {{.FileName}}{{range .Structs}}
	reflect.TypeOf((*{{.Name}})(nil)).Elem(),{{end}}
{{end}}
}

`

// 字段首字母大写
func publicFieldName(name string) string {
	return strings.ToUpper(string(name[0])) + name[1:]
}

type fieldModel struct {
	*meta.FieldDescriptor
}

func (self *fieldModel) FieldName() string {
	pname := publicFieldName(self.Name)

	// 碰到关键字在尾部加_
	if token.Lookup(pname).IsKeyword() {
		return pname + "_"
	}

	return pname
}

func (self *fieldModel) GoTypeName() string {

	var b bytes.Buffer
	if self.Repeatd {
		b.WriteString("[]")
	}

	if self.Complex != nil {
		b.WriteString("*")
	}

	// 字段类型映射go的类型
	switch self.Type {
	case meta.FieldType_Integer:
		b.WriteString("int")
	case meta.FieldType_Struct:
		b.WriteString(self.Complex.Name)
	default:
		b.WriteString(self.Type.String())
	}

	return b.String()
}

func (self *fieldModel) GoTags() string {

	var b bytes.Buffer

	b.WriteString("`sproto:\"")

	// 整形类型对解码层都视为整形
	switch self.Type {
	case meta.FieldType_Int32,
		meta.FieldType_Int64,
		meta.FieldType_UInt32,
		meta.FieldType_UInt64:
		b.WriteString("integer")
	default:
		b.WriteString(self.Kind())
	}

	b.WriteString(",")

	b.WriteString(fmt.Sprintf("%d", self.Tag))
	b.WriteString(",")

	if self.Repeatd {
		b.WriteString("array,")
	}

	b.WriteString(fmt.Sprintf("name=%s", self.FieldName()))

	b.WriteString("\"`")

	return b.String()
}

type structModel struct {
	*meta.Descriptor

	GoFields []fieldModel
}

type protoModel struct {
	*meta.FileDescriptor

	Structs []*structModel
}

type fileModel struct {
	Protos      []*protoModel
	PackageName string
}

func gen_go(pool []*meta.FileDescriptor, packageName, filename string) {

	tpl, err := template.New("sproto_go").Parse(codeTemplate)
	if err != nil {
		fmt.Println("template error ", err.Error())
		os.Exit(1)
	}

	var fm fileModel
	fm.PackageName = packageName

	for _, fileD := range pool {

		pm := &protoModel{
			FileDescriptor: fileD,
		}

		for _, st := range fileD.Structs {

			stModel := &structModel{
				Descriptor: st,
			}

			for _, fd := range st.Fields {

				fdModel := fieldModel{
					FieldDescriptor: fd,
				}

				stModel.GoFields = append(stModel.GoFields, fdModel)

			}

			pm.Structs = append(pm.Structs, stModel)

		}

		fm.Protos = append(fm.Protos, pm)

	}

	var bf bytes.Buffer

	err = tpl.Execute(&bf, &fm)
	if err != nil {
		fmt.Println("template error ", err.Error())
		os.Exit(1)
	}

	err = formatCode(&bf)

	if err != nil {
		fmt.Println("format error ", err.Error())
	}

	if fileErr := ioutil.WriteFile(filename, bf.Bytes(), 666); fileErr != nil {
		fmt.Println("write file error ", fileErr.Error())
		os.Exit(1)
	}
}

// Reformat generated code.
func formatCode(bf *bytes.Buffer) error {

	fset := token.NewFileSet()

	ast, err := parser.ParseFile(fset, "", bf, parser.ParseComments)
	if err != nil {
		return err
	}

	bf.Reset()

	err = (&printer.Config{Mode: printer.TabIndent | printer.UseSpaces, Tabwidth: 8}).Fprint(bf, fset, ast)
	if err != nil {
		return err
	}

	return nil
}