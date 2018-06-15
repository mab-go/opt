package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/mab-go/opt.v0/cmd/opt-gen/tmpl"
)

// --- Public Members ----------------------------------------------------------

// Config is a collection of ConfigItem structs.
type Config struct {
	Package           string       `json:"package"`
	ImportPath        string       `json:"import_path"`
	SeparatorComments bool         `json:"separator_comments"`
	Items             []ConfigItem `json:"items"`
}

// ConfigItem is a set of configuration values for a Generator.
type ConfigItem struct {
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	Imports    []string `json:"imports"` // TODO Separate imports for source and test?
	Implement  []string `json:"implement"`
	Tests      []string `json:"tests"`
	TestValues []string `json:"testValues"`
}

// LoadJSON loads the JSON object contained in the []byte b into a *Config.
func (c *Config) LoadJSON(b []byte) error {
	return json.Unmarshal(b, c)
}

// LoadJSONFile reads from the file pointed to by path and reads the JSON object
// contained within into a *Config.
func (c *Config) LoadJSONFile(path string) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return c.LoadJSON(b)
}

// Generate generates the source and test code for a set of optional types.
func Generate(cfg Config) (Result, error) {
	var src bytes.Buffer
	var test bytes.Buffer

	// Write file headers
	src.WriteString("// Code generated by opt-gen. DO NOT EDIT.\n\n")
	src.WriteString("package " + cfg.Package + "\n\n")
	test.WriteString("// Code generated by opt-gen. DO NOT EDIT.\n\n")
	test.WriteString("package " + cfg.Package + "_test\n\n")
	test.WriteString(fmt.Sprintf(`import (
	"testing"

	"%s"
)

`, cfg.ImportPath))

	// Generate code for each config item
	for _, item := range cfg.Items {
		if cfg.SeparatorComments {
			pad := 80 - (len("// --- ") + len(item.Name) + 1) // Add 1 for a space character
			b := []byte(fmt.Sprintf("// --- %s %s\n", item.Name, strings.Repeat("-", pad)))
			src.Write(b)
			test.Write(b)
		}

		g := generator{Package: cfg.Package}
		if err := g.generate(item); err != nil {
			return Result{}, err
		}
		src.Write(g.sourceBody.Bytes())
		test.Write(g.testBody.Bytes())
	}

	// Format code
	fmtSrc, err := format.Source(src.Bytes())
	if err != nil {
		return Result{}, fmt.Errorf("%v\n\n%s", err, src.Bytes())
	}
	src.Write(fmtSrc)

	fmtTest, err := format.Source(test.Bytes())
	if err != nil {
		return Result{}, fmt.Errorf("%v\n\n%s", err, test.String())
	}
	test.Write(fmtTest)

	return Result{Source: fmtSrc, Test: fmtTest}, nil
}

// A Result contains the source and test code results from a code generator
// operation.
type Result struct {
	Source []byte
	Test   []byte
}

// --- Private Members ---------------------------------------------------------

type generator struct {
	Package    string
	sourceBody bytes.Buffer
	testBody   bytes.Buffer
}

func (g *generator) generate(item ConfigItem) error {
	if g.Package == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		g.Package = filepath.Base(cwd)
	}

	if err := g.generateSource(item.Name, item.Type); err != nil {
		return err
	}
	return g.generateTests(item.Name, item.Type, item.TestValues)
}

func (g *generator) generateSource(n string, t string) error {
	data := struct{ Name, Type string }{Name: n, Type: t}
	return tmpl.BaseSource.Execute(&g.sourceBody, data)
}

func (g *generator) generateGetSetTests(n string, t string, vv []string) []byte {
	var body bytes.Buffer
	for i, v := range vv {
		var b bytes.Buffer
		data := struct {
			Package, Name, Val string
			ValIdx             int
		}{g.Package, n, v, i}
		if err := tmpl.BaseTestsGetSetBody.Execute(&b, data); err != nil {
			panic(err)
		}

		body.Write(b.Bytes())
		if i < len(vv)-1 {
			body.WriteString("\n")
		}
	}

	var sig bytes.Buffer
	data := struct{ Name string }{n}
	if err := tmpl.BaseTestsGetSetSignature.Execute(&sig, data); err != nil {
		panic(err)
	}
	sig.WriteString(" {")
	sig.Write(body.Bytes())
	sig.WriteString("\n}\n")

	return sig.Bytes()
}

func (g *generator) generateIsSetTests(n string, t string, vv []string) []byte {
	var body bytes.Buffer
	for i, v := range vv {
		var b bytes.Buffer
		data := struct {
			Package, Name, Val string
			ValIdx             int
		}{g.Package, n, v, i}
		if err := tmpl.BaseTestsIsSetBody.Execute(&b, data); err != nil {
			panic(err)
		}

		body.Write(b.Bytes())
		if i < len(vv)-1 {
			body.WriteString("\n")
		}
	}

	var sig bytes.Buffer
	data := struct{ Name string }{n}
	if err := tmpl.BaseTestsIsSetSignature.Execute(&sig, data); err != nil {
		panic(err)
	}
	sig.WriteString(" {")
	sig.Write(body.Bytes())
	sig.WriteString("\n}\n")

	return sig.Bytes()
}

func (g *generator) generateMakeTests(n string, t string, vv []string) []byte {
	var body bytes.Buffer
	for _, v := range vv {
		var b bytes.Buffer
		data := struct{ Package, Name, Val string }{g.Package, n, v}
		if err := tmpl.BaseTestsMakeBody.Execute(&b, data); err != nil {
			panic(err)
		}

		body.Write(b.Bytes())
	}

	var sig bytes.Buffer
	data := struct{ Name string }{n}
	if err := tmpl.BaseTestsMakeSignature.Execute(&sig, data); err != nil {
		panic(err)
	}
	sig.WriteString(" {")
	sig.Write(body.Bytes())
	sig.WriteString("\n}\n")

	return sig.Bytes()
}

func (g *generator) generateTests(n string, t string, vv []string) error {
	// FIXME Catch panics from generateTests* functions

	if _, err := g.testBody.Write(g.generateMakeTests(n, t, vv)); err != nil {
		return err
	}

	if _, err := g.testBody.Write(g.generateGetSetTests(n, t, vv)); err != nil {
		return err
	}

	if _, err := g.testBody.Write(g.generateIsSetTests(n, t, vv)); err != nil {
		return err
	}

	return nil
}
