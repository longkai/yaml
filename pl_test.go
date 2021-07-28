package yaml_test

import (
	"testing"

	"github.com/longkai/yaml/v2"
)

func TestWraping(t *testing.T) {
	saved1, saved2 := yaml.Tag, yaml.StructOrArrayPrefix
	defer func() {
		yaml.Tag, yaml.StructOrArrayPrefix = saved1, saved2
	}()
	yaml.Tag = "pl"
	yaml.StructOrArrayPrefix = "*.*.*.*."
	type st struct {
		Name string
	}
	v := &struct {
		A []string
		B st
		C *st
		D *int
		E int
		F *[]int
	}{F: new([]int), D: new(int), E: 100}
	b, err := yaml.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `'*.*.*.*.a': []
'*.*.*.*.b':
  name: ""
'*.*.*.*.c': null
d: 0
e: 100
'*.*.*.*.f': []
` {
		t.Fatal("unexpected Marshal result:", string(b))
	}
	v.E = 0
	v.F = nil
	if err := yaml.Unmarshal(b, &v); err != nil {
		t.Fatal(err)
	}
	if v.E != 100 {
		t.Fatal("should be 100")
	}
	if v.F == nil {
		t.Fatal("should not be nil")
	}
}

func TestYamlMarshal(t *testing.T) {
	saved1, saved2 := yaml.Tag, yaml.StructOrArrayPrefix
	defer func() {
		yaml.Tag, yaml.StructOrArrayPrefix = saved1, saved2
	}()
	yaml.Tag = "pl"
	cases := []struct {
		desc   string
		ptr    interface{}
		expect string
	}{
		{
			desc: "normal",
			ptr: struct {
				Name       string `pl:"test"`
				Age        int    `pl:"age"`
				Alias      string `pl:"alias123"`
				LastName   string `yaml:"last_name"`
				FamilyName string `json:"family_name"`
			}{
				Name:  "test",
				Age:   18,
				Alias: "alias123",
			},
			expect: `test: test
age: 18
alias123: alias123
last_name: ""
family_name: ""
`,
		},
		{
			desc: "array",
			ptr: []struct {
				Name  string `pl:"k"`
				Value string `pl:"v"`
			}{
				{
					Name:  "test",
					Value: "test",
				},
				{
					Name:  "test2",
					Value: "test2",
				},
			},
			expect: `- k: test
  v: test
- k: test2
  v: test2
`,
		},
		{
			desc: "nested",
			ptr: struct {
				Name string `pl:"test"`
				Pair struct {
					K string `pl:"k"`
					V string `pl:"v"`
				} `pl:"pair"`
			}{
				Name: "test",
				Pair: struct {
					K string `pl:"k"`
					V string `pl:"v"`
				}{K: "k", V: "v"},
			},
			expect: `test: test
pair:
  k: k
  v: v
`,
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			b, err := yaml.Marshal(c.ptr)
			if err != nil {
				t.Fatal(err)
			}
			if string(b) != c.expect {
				t.Errorf("got %v, expect %v", string(b), c.expect)
			}
		})
	}
}
