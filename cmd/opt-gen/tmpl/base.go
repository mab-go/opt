package tmpl

import (
	"text/template"
)

// BaseSource is the template for an optional type's base source code.
var BaseSource *template.Template

// BaseTestsMakeSignature is the template for an optional type's "Make*" test
// function signature.
var BaseTestsMakeSignature *template.Template

// BaseTestsMakeBody is the template for an optional type's "Make*" test
// // function body.
var BaseTestsMakeBody *template.Template

// BaseTestsGetSetSignature is the template for an optional type's "Get*" and
// "Set*" test function signatures.
var BaseTestsGetSetSignature *template.Template

// BaseTestsGetSetBody is the template for an optional type's "Get*" and "Set*"
// test function bodies.
var BaseTestsGetSetBody *template.Template

// BaseTestsIsSetSignature is the template for an optional type's "IsSet" test
// function signature.
var BaseTestsIsSetSignature *template.Template

// BaseTestsIsSetBody is the template for an optional type's "IsSet" test
// function body.
var BaseTestsIsSetBody *template.Template

var baseSourceText = `
// {{.Name}} is an optional type that wraps a {{.Type}}.
type {{.Name}} struct {
	isSet bool
	val   {{.Type}}
}

// Make{{.Name}} creates a new {{.Name}} with the specified value. 
func Make{{.Name}}(v {{.Type}}) {{.Name}} {
	return {{.Name}}{isSet: true, val: v}
}

// IsSet returns a value indicating whether the optional type's value is set.
func (p {{.Name}}) IsSet() bool {
	return p.isSet
}

// Set sets the optional type's value.
func (p *{{.Name}}) Set(v {{.Type}}) {
	p.val = v
	p.isSet = true
}

// Get returns the underlying value wrapped by the optional type. If IsSet returns
// false, then Get's return value will be the zero value for the underlying type.
func (p {{.Name}}) Get() {{.Type}} {
	return p.val
}
`

var baseTestsMakeSignatureText = `
func TestMake{{.Name}}(t *testing.T)`

var baseTestsMakeBodyText = "\n_ = {{.Package}}.Make{{.Name}}({{.Val}})"

var baseTestsGetSetSignatureText = `
func Test{{.Name}}_GetSet(t *testing.T)`

var baseTestsGetSetBodyText = `
val{{.ValIdx}} := {{.Package}}.{{.Name}}{}
val{{.ValIdx}}.Set({{.Val}})
if val{{.ValIdx}}.Get() != {{.Val}} {
	t.Fatalf("val{{.ValIdx}}.Get; got: %v, want: %v", val{{.ValIdx}}.Get(), {{.Val}})
}`

var baseTestsIsSetSignatureText = `
func Test{{.Name}}_IsSet(t *testing.T)`

var baseTestsIsSetBodyText = `
val{{.ValIdx}}A := {{.Package}}.{{.Name}}{}
if val{{.ValIdx}}A.IsSet() {
	t.Fatalf("val{{.ValIdx}}A.IsSet; got: %v, want: %v", val{{.ValIdx}}A.IsSet(), false)
}
val{{.ValIdx}}A.Set({{.Val}})
if !val{{.ValIdx}}A.IsSet() {
	t.Fatalf("val{{.ValIdx}}A.IsSet; got: %v, want: %v", val{{.ValIdx}}A.IsSet(), true)
}

val{{.ValIdx}}B := {{.Package}}.Make{{.Name}}({{.Val}})
if !val{{.ValIdx}}B.IsSet() {
	t.Fatalf("val{{.ValIdx}}B.IsSet; got: %v, want: %v", val{{.ValIdx}}B.IsSet(), true)
}`

func init() {
	BaseSource = template.Must(template.New("BaseSource").Parse(baseSourceText))
	BaseTestsMakeSignature = template.Must(template.New("BaseTestsMakeSignature").Parse(baseTestsMakeSignatureText))
	BaseTestsMakeBody = template.Must(template.New("BaseTestsMakeBody").Parse(baseTestsMakeBodyText))
	BaseTestsGetSetSignature = template.Must(template.New("BaseTestsGetSetSignature").Parse(baseTestsGetSetSignatureText))
	BaseTestsGetSetBody = template.Must(template.New("BaseTestsGetSetBody").Parse(baseTestsGetSetBodyText))
	BaseTestsIsSetSignature = template.Must(template.New("BaseTestsIsSetSignature").Parse(baseTestsIsSetSignatureText))
	BaseTestsIsSetBody = template.Must(template.New("BaseTestsIsSetBody").Parse(baseTestsIsSetBodyText))
}
