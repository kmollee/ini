// Package ini :the INI config parser and operate
package ini

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

const (
	// Name for default section.
	defaultSection = ""
	lineBreak      = "\n"
	separator      = "="
	comment        = "#"
)

var (
	commentSign       = []byte(comment)
	keyValueSeparator = []byte(separator)
)
var (
	// ErrFormat :foramt is not correct error
	ErrFormat = errors.New("format is incorrect")
	// ErrSectionMiss :looking section is not exist error
	ErrSectionMiss = errors.New("section is not exsit")
	// ErrKeyMiss :looking section's key is not exist error
	ErrKeyMiss = errors.New("section's key is not exsit")
)

type kvMap map[string]string

// INI :the ini config type
type INI map[string]kvMap

func (ini INI) parse(r io.Reader) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	currentSectionName := defaultSection
	kvmap := make(kvMap)

	ini[currentSectionName] = kvmap

	lines := bytes.Split(data, []byte(lineBreak))

	for _, line := range lines {
		line = bytes.TrimSpace(line)
		// is empty line or comment line
		size := len(line)
		if size == 0 || bytes.HasPrefix(line, commentSign) {
			continue
		}

		// new section
		if line[0] == '[' && line[size-1] == ']' {
			name := string(line[1 : size-1])
			// empty section name
			if len(name) == 0 {
				return ErrFormat
			}
			currentSectionName = name
			kvmap = make(kvMap)
			ini[currentSectionName] = kvmap
			continue
		}

		// key value pair
		var pos int
		if pos = bytes.Index(line, keyValueSeparator); pos == -1 {
			return ErrFormat
		}
		// must trim space before assgin
		// k            =    v => key: k value: v
		key, val := bytes.TrimSpace(line[:pos]), bytes.TrimSpace(line[pos+1:])
		// either key or val is empty
		if len(key) == 0 || len(val) == 0 {
			return ErrFormat
		}
		kvmap[string(key)] = string(val)
	}
	return nil
}

func (ini INI) Write(w io.Writer) error {
	b := bufio.NewWriter(w)

	// wrtie default first, because it has no name
	if kvmap, ok := ini[defaultSection]; ok {
		// write default section
		if err := ini.writeKV(b, kvmap); err != nil {
			return err
		}
	}

	for section, kvmap := range ini {
		if section == defaultSection {
			continue
		}

		// write section
		if _, err := b.WriteString(fmt.Sprintf("[%s]\n", section)); err != nil {
			return err
		}

		// write key val pair
		if err := ini.writeKV(b, kvmap); err != nil {
			return err
		}

	}

	return nil
}

func (ini INI) writeKV(b *bufio.Writer, kv kvMap) error {
	for k, v := range kv {
		if _, err := b.WriteString(k + separator + v + lineBreak); err != nil {
			return err
		}
	}

	return b.Flush()
}

// DefaultSectionSetKey :add key/value pair to default section
func (ini INI) DefaultSectionSetKey(k, v string) {
	ini.SectionSetKey(defaultSection, k, v)
}

// DefaultSectionGetKey :get key's value from default section
func (ini INI) DefaultSectionGetKey(k string) (string, error) {
	return ini.SectionGetKey(defaultSection, k)
}

// DefaultSectionDelKey : delete key from default section
func (ini INI) DefaultSectionDelKey(k string) error {
	return ini.SectionDelKey(defaultSection, k)
}

// DefaultSectionGet :get default section key's value
func (ini INI) DefaultSectionGet() (map[string]string, error) {
	return ini.SectionGet(defaultSection)
}

// SectionSetKey :assign section key value pair
func (ini INI) SectionSetKey(sectionName, k, v string) {
	// if not section not exist, create first
	if _, exist := ini[sectionName]; !exist {
		ini[sectionName] = make(kvMap)
	}

	ini[sectionName][k] = v
}

// SectionGetKey :get section's key val
func (ini INI) SectionGetKey(sectionName, k string) (string, error) {
	// if not section not exist, create first
	s, exist := ini[sectionName]
	if !exist {
		return "", ErrSectionMiss
	}

	v, exist := s[k]
	if !exist {
		return "", ErrKeyMiss
	}
	return v, nil
}

// SectionDelKey :delete setion's key
func (ini INI) SectionDelKey(sectionName, k string) error {
	if _, exist := ini[sectionName]; !exist {
		return ErrSectionMiss
	}
	if _, exist := ini[sectionName][k]; exist {
		delete(ini[sectionName], k)
	}
	return nil
}

// SectionDel :delete section
func (ini INI) SectionDel(sectionName string) error {
	if sectionName == defaultSection {
		return fmt.Errorf("could not remove default section")
	}

	if _, exist := ini[sectionName]; !exist {
		return ErrSectionMiss
	}

	delete(ini, sectionName)
	return nil
}

// SectionGet :get section key/val map
func (ini INI) SectionGet(sectionName string) (map[string]string, error) {
	if _, exist := ini[sectionName]; !exist {
		return nil, ErrSectionMiss
	}

	return ini[sectionName], nil
}

// SectionUpdate :update section key/val map
func (ini INI) SectionUpdate(sectionName string, data map[string]string) {
	if _, exist := ini[sectionName]; !exist {
		ini[sectionName] = make(kvMap)
	}
	for k, v := range data {
		ini[sectionName][k] = v
	}
}

// New :new INI config, with default section
func New() INI {
	return INI{defaultSection: make(kvMap)}
}

// Parse :create INI from io.Reader
func Parse(r io.Reader) (INI, error) {
	conf := New()

	if err := conf.parse(r); err != nil {
		return nil, err
	}
	return conf, nil
}

// ParseString :create INI from string
func ParseString(s string) (INI, error) {
	r := strings.NewReader(s)
	conf := New()
	if err := conf.parse(r); err != nil {
		return nil, err
	}
	return conf, nil
}
