package sprig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefault(t *testing.T) {
	tpl := `{{"" | default "foo"}}`
	if err := runt(tpl, "foo"); err != nil {
		t.Error(err)
	}
	tpl = `{{default "foo" 234}}`
	if err := runt(tpl, "234"); err != nil {
		t.Error(err)
	}
	tpl = `{{default "foo" 2.34}}`
	if err := runt(tpl, "2.34"); err != nil {
		t.Error(err)
	}

	tpl = `{{ .Nothing | default "123" }}`
	if err := runt(tpl, "123"); err != nil {
		t.Error(err)
	}
	tpl = `{{ default "123" }}`
	if err := runt(tpl, "123"); err != nil {
		t.Error(err)
	}
}

func TestSafeDefault(t *testing.T) {
	tpl := `{{"" | safeDefault "foo"}}`
	if err := runt(tpl, ""); err != nil {
		t.Error(err)
	}
	tpl = `{{safeDefault "foo" 234}}`
	if err := runt(tpl, "234"); err != nil {
		t.Error(err)
	}
	tpl = `{{safeDefault "foo" 2.34}}`
	if err := runt(tpl, "2.34"); err != nil {
		t.Error(err)
	}

	tpl = `{{ .Nothing | safeDefault "123" }}`
	if err := runt(tpl, "123"); err != nil {
		t.Error(err)
	}
	tpl = `{{ safeDefault "123" }}`
	if err := runt(tpl, "123"); err != nil {
		t.Error(err)
	}

	tpl = `{{ .Nothing | safeDefault true }}`
	if err := runt(tpl, "true"); err != nil {
		t.Error(err)
	}
	tpl = `{{ false | safeDefault true }}`
	if err := runt(tpl, "false"); err != nil {
		t.Error(err)
	}
	tpl = `{{ true | safeDefault false }}`
	if err := runt(tpl, "true"); err != nil {
		t.Error(err)
	}
	tpl = `{{ .Nothing | safeDefault 100 }}`
	if err := runt(tpl, "100"); err != nil {
		t.Error(err)
	}
	tpl = `{{ 0 | safeDefault 100 }}`
	if err := runt(tpl, "0"); err != nil {
		t.Error(err)
	}
	tpl = `{{ 55 | safeDefault 0 }}`
	if err := runt(tpl, "55"); err != nil {
		t.Error(err)
	}
}

func TestEmpty(t *testing.T) {
	tpl := `{{if empty 1}}1{{else}}0{{end}}`
	if err := runt(tpl, "0"); err != nil {
		t.Error(err)
	}

	tpl = `{{if empty 0}}1{{else}}0{{end}}`
	if err := runt(tpl, "1"); err != nil {
		t.Error(err)
	}
	tpl = `{{if empty ""}}1{{else}}0{{end}}`
	if err := runt(tpl, "1"); err != nil {
		t.Error(err)
	}
	tpl = `{{if empty 0.0}}1{{else}}0{{end}}`
	if err := runt(tpl, "1"); err != nil {
		t.Error(err)
	}
	tpl = `{{if empty false}}1{{else}}0{{end}}`
	if err := runt(tpl, "1"); err != nil {
		t.Error(err)
	}

	dict := map[string]interface{}{"top": map[string]interface{}{}}
	tpl = `{{if empty .top.NoSuchThing}}1{{else}}0{{end}}`
	if err := runtv(tpl, "1", dict); err != nil {
		t.Error(err)
	}
	tpl = `{{if empty .bottom.NoSuchThing}}1{{else}}0{{end}}`
	if err := runtv(tpl, "1", dict); err != nil {
		t.Error(err)
	}
}

func TestCoalesce(t *testing.T) {
	tests := map[string]string{
		`{{ coalesce 1 }}`:                            "1",
		`{{ coalesce "" 0 nil 2 }}`:                   "2",
		`{{ $two := 2 }}{{ coalesce "" 0 nil $two }}`: "2",
		`{{ $two := 2 }}{{ coalesce "" $two 0 0 0 }}`: "2",
		`{{ $two := 2 }}{{ coalesce "" $two 3 4 5 }}`: "2",
		`{{ coalesce }}`:                              "<no value>",
	}
	for tpl, expect := range tests {
		assert.NoError(t, runt(tpl, expect))
	}

	dict := map[string]interface{}{"top": map[string]interface{}{}}
	tpl := `{{ coalesce .top.NoSuchThing .bottom .bottom.dollar "airplane"}}`
	if err := runtv(tpl, "airplane", dict); err != nil {
		t.Error(err)
	}
}

func TestSafeCoalesce(t *testing.T) {
	tests := map[string]string{
		`{{ safeCoalesce 1 }}`:                            "1",
		`{{ safeCoalesce "" 0 nil 2 }}`:                   "",
		`{{ safeCoalesce nil 0 "" 2 }}`:				   "0",
		`{{ $two := 2 }}{{ safeCoalesce nil $two "" 0 }}`: "2",
		`{{ $two := 2 }}{{ safeCoalesce "" $two 0 0 0 }}`: "",
		`{{ $two := 2 }}{{ safeCoalesce "" $two 3 4 5 }}`: "",
		`{{ safeCoalesce }}`:                              "<no value>",
	}
	for tpl, expect := range tests {
		assert.NoError(t, runt(tpl, expect))
	}

	dict := map[string]interface{}{"top": map[string]interface{}{}}
	tpl := `{{ coalesce .top.NoSuchThing .bottom .bottom.dollar "airplane"}}`
	if err := runtv(tpl, "airplane", dict); err != nil {
		t.Error(err)
	}
}

func TestAll(t *testing.T) {
	tests := map[string]string{
		`{{ all 1 }}`:                            "true",
		`{{ all "" 0 nil 2 }}`:                   "false",
		`{{ $two := 2 }}{{ all "" 0 nil $two }}`: "false",
		`{{ $two := 2 }}{{ all "" $two 0 0 0 }}`: "false",
		`{{ $two := 2 }}{{ all "" $two 3 4 5 }}`: "false",
		`{{ all }}`:                              "true",
	}
	for tpl, expect := range tests {
		assert.NoError(t, runt(tpl, expect))
	}

	dict := map[string]interface{}{"top": map[string]interface{}{}}
	tpl := `{{ all .top.NoSuchThing .bottom .bottom.dollar "airplane"}}`
	if err := runtv(tpl, "false", dict); err != nil {
		t.Error(err)
	}
}

func TestAny(t *testing.T) {
	tests := map[string]string{
		`{{ any 1 }}`:                              "true",
		`{{ any "" 0 nil 2 }}`:                     "true",
		`{{ $two := 2 }}{{ any "" 0 nil $two }}`:   "true",
		`{{ $two := 2 }}{{ any "" $two 3 4 5 }}`:   "true",
		`{{ $zero := 0 }}{{ any "" $zero 0 0 0 }}`: "false",
		`{{ any }}`: "false",
	}
	for tpl, expect := range tests {
		assert.NoError(t, runt(tpl, expect))
	}

	dict := map[string]interface{}{"top": map[string]interface{}{}}
	tpl := `{{ any .top.NoSuchThing .bottom .bottom.dollar "airplane"}}`
	if err := runtv(tpl, "true", dict); err != nil {
		t.Error(err)
	}
}

func TestFromJson(t *testing.T) {
	dict := map[string]interface{}{"Input": `{"foo": 55}`}

	tpl := `{{.Input | fromJson}}`
	expected := `map[foo:55]`
	if err := runtv(tpl, expected, dict); err != nil {
		t.Error(err)
	}

	tpl = `{{(.Input | fromJson).foo}}`
	expected = `55`
	if err := runtv(tpl, expected, dict); err != nil {
		t.Error(err)
	}
}

func TestToJson(t *testing.T) {
	dict := map[string]interface{}{"Top": map[string]interface{}{"bool": true, "string": "test", "number": 42}}

	tpl := `{{.Top | toJson}}`
	expected := `{"bool":true,"number":42,"string":"test"}`
	if err := runtv(tpl, expected, dict); err != nil {
		t.Error(err)
	}
}

func TestToPrettyJson(t *testing.T) {
	dict := map[string]interface{}{"Top": map[string]interface{}{"bool": true, "string": "test", "number": 42}}
	tpl := `{{.Top | toPrettyJson}}`
	expected := `{
  "bool": true,
  "number": 42,
  "string": "test"
}`
	if err := runtv(tpl, expected, dict); err != nil {
		t.Error(err)
	}
}

func TestToRawJson(t *testing.T) {
	dict := map[string]interface{}{"Top": map[string]interface{}{"bool": true, "string": "test", "number": 42, "html": "<HEAD>"}}
	tpl := `{{.Top | toRawJson}}`
	expected := `{"bool":true,"html":"<HEAD>","number":42,"string":"test"}`

	if err := runtv(tpl, expected, dict); err != nil {
		t.Error(err)
	}
}

func TestTernary(t *testing.T) {
	tpl := `{{true | ternary "foo" "bar"}}`
	if err := runt(tpl, "foo"); err != nil {
		t.Error(err)
	}

	tpl = `{{ternary "foo" "bar" true}}`
	if err := runt(tpl, "foo"); err != nil {
		t.Error(err)
	}

	tpl = `{{false | ternary "foo" "bar"}}`
	if err := runt(tpl, "bar"); err != nil {
		t.Error(err)
	}

	tpl = `{{ternary "foo" "bar" false}}`
	if err := runt(tpl, "bar"); err != nil {
		t.Error(err)
	}
}
