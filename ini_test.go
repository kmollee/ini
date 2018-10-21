package ini

import (
	"bufio"
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func TestParseError(t *testing.T) {
	tt := []struct {
		name string
		data string
		err  error
	}{
		{"invalid ? mark", "?", ErrFormat},
		{"invalid @ mark", "@", ErrFormat},
		{"invalid ! mark", "!", ErrFormat},
		{"invalid = mark", "=", ErrFormat},
		{"invalid $ mark", "$", ErrFormat},
		{"invalid % mark", "%", ErrFormat},
		{"invalid & mark", "&", ErrFormat},
		{"invalid ; mark", ";", ErrFormat},
		{"invalid [] mark", "[]", ErrFormat},
		{"valid assgin", "a = b", nil},
		{"valid comment", "#a = b", nil},
		{"valid comment with space", "     #a = b", nil},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r := strings.NewReader(tc.data)
			conf := New()
			err := conf.parse(r)
			if err != tc.err {
				t.Errorf("parse data: %v expect error '%v' got '%v'", tc.data, tc.err, err)
			}
		})

	}
}

func TestParseSection(t *testing.T) {
	tt := []struct {
		name   string
		data   string
		expect INI
	}{
		{
			name: "case1",
			data: `age = 18
			gender = male
			
			[first]
			firstk = firstv
			`,
			expect: INI{
				defaultSection: {"age": "18", "gender": "male"},
				"first":        {"firstk": "firstv"},
			},
		},
		{
			name: "case2",
			data: `age = 18
			gender = male
			
			[first]
			firstk = firstv
			
			[second]
			secondk = secondv

			`,
			expect: INI{
				defaultSection: {"age": "18", "gender": "male"},
				"first":        {"firstk": "firstv"},
				"second":       {"secondk": "secondv"},
			},
		},
		{
			name: "case3 ignore comment",
			data: `age = 18
			gender = male
			
			[first]
			firstk = firstv
			# first comment
			####
			[second]
			secondk = secondv

			## this is comment

			`,
			expect: INI{
				defaultSection: {"age": "18", "gender": "male"},
				"first":        {"firstk": "firstv"},
				"second":       {"secondk": "secondv"},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r := strings.NewReader(tc.data)
			conf := New()
			err := conf.parse(r)
			if err != nil {
				t.Errorf("parse content %v error: %v", tc.data, err)
			}

			if !reflect.DeepEqual(conf, tc.expect) {
				t.Errorf("parse ini expect '%#+v' got '%#+v'", tc.expect, conf)
			}
		})
	}
}

func TestSectionOperate(t *testing.T) {
	t.Run("section key/val set", func(t *testing.T) {
		expect := INI{
			defaultSection: {"hello": "world"},
			"section1":     {"key": "val"},
		}

		conf := New()
		conf.DefaultSectionSetKey("hello", "world")
		conf.SectionSetKey("section1", "key", "val")

		if !reflect.DeepEqual(conf, expect) {
			t.Errorf("parse ini expect '%#+v' got '%#+v'", expect, conf)
		}
	})

	t.Run("section get", func(t *testing.T) {
		data := `
		century = 21

		[section1]
		key = 1

		[section2]
		key = 2

		[section3]
		key = 3
		`

		tt := []struct {
			section string
			key     string
			expect  string
		}{
			{defaultSection, "century", "21"},
			{"section1", "key", "1"},
			{"section2", "key", "2"},
			{"section3", "key", "3"},
		}

		conf, err := ParseString(data)
		if err != nil {
			t.Errorf("could not parse config: %v", err)
		}

		for _, tc := range tt {

			val, err := conf.SectionGetKey(tc.section, tc.key)
			if err != nil {
				t.Errorf("could not get section:'%v' key'%v' value: %v", tc.section, tc.key, err)
			}
			if val != tc.expect {
				t.Errorf("get section: %s key: %s expect value:'%s' got '%s'", tc.section, tc.key, tc.expect, val)
			}
		}

	})
	t.Run("section key remove", func(t *testing.T) {
		expect := INI{
			defaultSection: {"century": "21"},
			"section1":     {},
		}
		data := `
		hello = world
		century = 21

		[section1]
		key = val
		`

		conf, err := ParseString(data)
		if err != nil {
			t.Errorf("parse ini content fail: %v", err)
		}

		err = conf.SectionDelKey("section1", "key")
		if err != nil {
			t.Errorf("remove key:%v fail with error: %v", "section1:key", err)
		}

		err = conf.DefaultSectionDelKey("hello")
		if err != nil {
			t.Errorf("remove key:%v fail with error: %v", "hello", err)
		}

		if !reflect.DeepEqual(conf, expect) {
			t.Errorf("parse ini expect '%#+v' got '%#+v'", expect, conf)
		}
	})

	t.Run("remove section", func(t *testing.T) {
		expect := INI{
			defaultSection: {"century": "21"},
		}
		data := `
		century = 21

		[section1]
		key = val

		[section2]
		key = val

		[section3]
		key = val
		`

		conf, err := ParseString(data)
		if err != nil {
			t.Errorf("parse ini content fail: %v", err)
		}

		// remove all extra section
		for _, name := range []string{"section1", "section2", "section3"} {
			if err = conf.SectionDel(name); err != nil {
				t.Errorf("remove section:%v fail with error: %v", name, err)
			}
		}

		if !reflect.DeepEqual(conf, expect) {
			t.Errorf("parse ini expect '%#+v' got '%#+v'", expect, conf)
		}
	})

	t.Run("update section not exist", func(t *testing.T) {
		expect := INI{
			defaultSection: {"century": "21"},
			"LOL":          {"year": "2018", "G2": "EU", "C9": "NA", "IG": "CN", "FNC": "EU"},
		}
		data := `
		century = 21
		`

		quarterfinals := map[string]string{
			"year": "2018",
			"G2":   "EU",
			"C9":   "NA",
			"IG":   "CN",
			"FNC":  "EU",
		}

		conf, err := ParseString(data)
		if err != nil {
			t.Errorf("parse ini content fail: %v", err)
		}

		conf.SectionUpdate("LOL", quarterfinals)
		if err != nil {
			t.Errorf("could not update exist section:%v", err)
		}

		if !reflect.DeepEqual(conf, expect) {
			t.Errorf("parse ini expect '%#+v' got '%#+v'", expect, conf)
		}

	})

	t.Run("update exist section", func(t *testing.T) {
		expect := INI{
			defaultSection: {"century": "21"},
			"LOL":          {"year": "2018", "G2": "EU", "C9": "NA", "IG": "CN", "FNC": "EU"},
		}
		data := `
		century = 21
		[LOL]
		year = 2018
		`

		quarterfinals := map[string]string{
			"G2":  "EU",
			"C9":  "NA",
			"IG":  "CN",
			"FNC": "EU",
		}

		conf, err := ParseString(data)
		if err != nil {
			t.Errorf("parse ini content fail: %v", err)
		}

		conf.SectionUpdate("LOL", quarterfinals)
		if err != nil {
			t.Errorf("could not update exist section:%v", err)
		}

		if !reflect.DeepEqual(conf, expect) {
			t.Errorf("parse ini expect '%#+v' got '%#+v'", expect, conf)
		}

	})

	t.Run("get section map", func(t *testing.T) {
		expect := map[string]string{"year": "2018"}

		data := `
		century = 21
		[LOL]
		year = 2018
		`
		conf, err := ParseString(data)
		if err != nil {
			t.Errorf("parse ini content fail: %v", err)
		}

		got, err := conf.SectionGet("LOL")
		if !reflect.DeepEqual(got, expect) {
			t.Errorf("get section expect '%#+v' got '%#+v'", expect, got)
		}

	})

	t.Run("get not exsit section", func(t *testing.T) {
		data := `
		century = 21
		`
		conf, err := ParseString(data)
		if err != nil {
			t.Errorf("parse ini content fail: %v", err)
		}

		got, err := conf.SectionGet("LOL")
		if got != nil {
			t.Errorf("get not exist section should return nil map: got:%v", got)
		}
		if err != ErrSectionMiss {
			t.Errorf("return wrong type of error: expect %v got %v", ErrSectionMiss, err)
		}

	})
}

func TestWrite(t *testing.T) {
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)

	data := `
	century = 21

	[section1]
	key = val

	[section2]
	key = val

	[section3]
	key = val
	`
	conf, err := ParseString(data)
	if err != nil {
		t.Errorf("parse ini content fail: %v", err)
	}

	if err := conf.Write(w); err != nil {
		t.Errorf("could not write conf into buffer: %v", err)
	}

	// make sure write into buffer
	w.Flush()

	r := bufio.NewReader(&buf)
	readConf, err := Parse(r)
	if err != nil {
		t.Errorf("parse ini content fail: %v", err)
	}
	if !reflect.DeepEqual(conf, readConf) {
		t.Errorf("compare write and read ini expect '%#+v' got '%#+v'", conf, readConf)
	}

}
