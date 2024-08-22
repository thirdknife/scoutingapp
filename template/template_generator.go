// Generates Templ templates for arbitrary structs.
//
// Expected usage is to modify `main()` to add or remove any template generation, and then
// run the generator to create Templ templates:
//
//	go run ./template/template_generator.go -out_dir=views/
package main

import (
	"flag"
	"fmt"
	"github.com/thirdknife/scoutingapp/database"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"
)

var (
	outDir = flag.String("out_dir", "", "Output directory to store templates in.")
)

func main() {
	flag.Parse()

	if *outDir == "" {
		fmt.Println("Error: Output directory not specified. Use -out_dir flag.")
		os.Exit(1)
	}

	// Create the output directory if it doesn't exist
	if err := os.MkdirAll(*outDir, os.ModePerm); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	// TODO: consider also providing fields which should not show up, e.g. PlayerID probably isn't worth displaying.
	generateTemplFile(database.Player{}, listTemplTemplate)
	generateTemplFile(database.PlayerAnalysis{}, listTemplTemplate)
	generateTemplFile(database.Analysis{}, listTemplTemplate)
}

type TemplTemplate struct {
	fileNamePrefix string
	templTemplate  string
}

var listTemplTemplate = TemplTemplate{
	fileNamePrefix: "List",
	templTemplate: `
package views

import (
    db "github.com/thirdknife/scoutingapp/database"
    "fmt"
)

templ List{{.StructName}}s({{.LowerStructName}}s []*db.{{.StructName}}) {
    @layout("{{.StructName}}s") {
       <table>
          {{- range .Fields}}
          <th>{{.Name}}</th>
          {{- end}}
          for _, {{.FirstLetter}} := range {{.LowerStructName}}s {
             <tr>
                {{- range .Fields}}
                <td>{{.TemplSyntax}}</td>
                {{- end}}
             </tr>
          }
       </table>
    }
}
`}

type FieldInfo struct {
	Name        string
	TemplSyntax string
}

type TemplData struct {
	StructName      string
	LowerStructName string
	FirstLetter     string
	Fields          []FieldInfo
}

func generateTemplFile(t interface{}, templTmpl TemplTemplate) {
	typ := reflect.TypeOf(t)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	structName := typ.Name()
	fields := make([]FieldInfo, 0)

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Anonymous {
			continue // Skip embedded structs
		}
		fieldInfo := FieldInfo{
			Name:        field.Name,
			TemplSyntax: fmt.Sprintf("{fmt.Sprintf(\"%%v\", %s.%s)}", strings.ToLower(structName[:1]), field.Name),
		}
		fields = append(fields, fieldInfo)
	}

	data := TemplData{
		StructName:      structName,
		LowerStructName: strings.ToLower(structName),
		FirstLetter:     strings.ToLower(structName[:1]),
		Fields:          fields,
	}

	tmpl, err := template.New("templ").Parse(templTmpl.templTemplate)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	fileName := fmt.Sprintf("%s%s.templ", templTmpl.fileNamePrefix, structName)
	filePath := filepath.Join(*outDir, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	err = tmpl.Execute(file, data)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return
	}

	fmt.Printf("%s file generated successfully in %s.\n", fileName, *outDir)
}
