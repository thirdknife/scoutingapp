package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"text/template"
)

// TestPlayer is a local struct for testing purposes
type TestPlayer struct {
	ID       int
	Name     string
	Position string
	Team     string
}

func TestGenerateTemplFile(t *testing.T) {
	// Create a temporary directory for test output
	tempDir, err := os.MkdirTemp("", "test_output")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up the temporary directory after the test

	// Set the outDir flag
	originalOutDir := *outDir
	*outDir = tempDir
	defer func() { *outDir = originalOutDir }() // Restore the original outDir after the test

	// Test case
	type TestStruct struct {
		Field1 string
		Field2 int
	}

	testTemplate := TemplTemplate{
		fileNamePrefix: "Test",
		templTemplate: `
package views

templ Test{{.StructName}}({{.LowerStructName}} *TestStruct) {
	<div>
		{{- range .Fields}}
		<p>{{.Name}}: {{.TemplSyntax}}</p>
		{{- end}}
	</div>
}
`,
	}

	// Run the function
	generateTemplFile(TestStruct{}, testTemplate)

	// Check if the file was created
	expectedFileName := filepath.Join(tempDir, "TestTestStruct.templ")
	content, err := ioutil.ReadFile(expectedFileName)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	// Check file contents
	expectedContent := `
package views

templ TestTestStruct(teststruct *TestStruct) {
	<div>
		<p>Field1: {fmt.Sprintf("%v", t.Field1)}</p>
		<p>Field2: {fmt.Sprintf("%v", t.Field2)}</p>
	</div>
}
`
	if strings.TrimSpace(string(content)) != strings.TrimSpace(expectedContent) {
		t.Errorf("Generated content does not match expected.\nExpected:\n%s\nGot:\n%s", expectedContent, string(content))
	}
}

func TestTemplData(t *testing.T) {
	player := TestPlayer{}
	typ := reflect.TypeOf(player)
	data := TemplData{
		StructName:      typ.Name(),
		LowerStructName: strings.ToLower(typ.Name()),
		FirstLetter:     strings.ToLower(typ.Name()[:1]),
		Fields:          []FieldInfo{},
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if !field.Anonymous {
			data.Fields = append(data.Fields, FieldInfo{
				Name:        field.Name,
				TemplSyntax: fmt.Sprintf("{fmt.Sprintf(\"%%v\", %s.%s)}", data.FirstLetter, field.Name),
			})
		}
	}

	if data.StructName != "TestPlayer" {
		t.Errorf("Expected StructName to be 'TestPlayer', got '%s'", data.StructName)
	}

	if data.LowerStructName != "testplayer" {
		t.Errorf("Expected LowerStructName to be 'testplayer', got '%s'", data.LowerStructName)
	}

	if data.FirstLetter != "t" {
		t.Errorf("Expected FirstLetter to be 't', got '%s'", data.FirstLetter)
	}

	// Check if all fields are present
	expectedFields := []string{"ID", "Name", "Position", "Team"}
	for _, expected := range expectedFields {
		found := false
		for _, field := range data.Fields {
			if field.Name == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected field '%s' not found in TemplData", expected)
		}
	}
}

func TestTemplateExecution(t *testing.T) {
	tmpl, err := template.New("test").Parse(listTemplTemplate.templTemplate)
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	data := TemplData{
		StructName:      "TestPlayer",
		LowerStructName: "testplayer",
		FirstLetter:     "t",
		Fields: []FieldInfo{
			{Name: "ID", TemplSyntax: "{fmt.Sprintf(\"%v\", t.ID)}"},
			{Name: "Name", TemplSyntax: "{fmt.Sprintf(\"%v\", t.Name)}"},
			{Name: "Position", TemplSyntax: "{fmt.Sprintf(\"%v\", t.Position)}"},
			{Name: "Team", TemplSyntax: "{fmt.Sprintf(\"%v\", t.Team)}"},
		},
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		t.Fatalf("Failed to execute template: %v", err)
	}

	output := buf.String()
	expectedSubstrings := []string{
		"package views",
		"templ ListTestPlayers(testplayers []*db.TestPlayer)",
		"<th>ID</th>",
		"<th>Name</th>",
		"<th>Position</th>",
		"<th>Team</th>",
		"<td>{fmt.Sprintf(\"%v\", t.ID)}</td>",
		"<td>{fmt.Sprintf(\"%v\", t.Name)}</td>",
		"<td>{fmt.Sprintf(\"%v\", t.Position)}</td>",
		"<td>{fmt.Sprintf(\"%v\", t.Team)}</td>",
	}

	for _, substr := range expectedSubstrings {
		if !strings.Contains(output, substr) {
			t.Errorf("Expected substring not found in output: %s", substr)
		}
	}
}
