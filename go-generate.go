// +build ignore

package main

//run this as go run go-generate.go $(pwd)

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/haproxytech/config-parser/common"
)

type Data struct {
	ParserMultiple   bool
	ParserName       string
	ParserSecondName string
	StructName       string
	ParserType       string
	NoInit           bool
	NoParse          bool
	TestOK           []string
	TestFail         []string
}

func main() {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}

	log.Println(dir)
	generateTypes(dir)
	generateTypesGeneric(dir)
}

func generateTypesGeneric(dir string) {
	dat, err := ioutil.ReadFile("types/types-generic.go")
	if err != nil {
		log.Println(err)
	}
	lines := common.StringSplitIgnoreEmpty(string(dat), '\n')
	//fmt.Print(lines)

	parsers := map[string]*Data{}
	parserData := &Data{}
	for _, line := range lines {
		//log.Println(parserData)
		//log.Println(line)
		if strings.HasPrefix(line, "//sections:") {
			//log.Println(line)
		}
		if strings.HasPrefix(line, "//name:") {
			data := common.StringSplitIgnoreEmpty(line, ':')
			items := common.StringSplitIgnoreEmpty(data[1], ' ')
			parserData.ParserName = data[1]
			if len(items) > 1 {
				parserData.ParserName = items[0]
				parserData.ParserSecondName = items[1]
			}
		}
		if strings.HasPrefix(line, "//no-init:true") {
			parserData.NoInit = true
		}
		if strings.HasPrefix(line, "//no-parse:true") {
			parserData.NoParse = true
		}
		if strings.HasPrefix(line, "//test:ok") {
			data := strings.SplitN(line, ":", 3)
			parserData.TestOK = append(parserData.TestOK, data[2])
		}
		if strings.HasPrefix(line, "//test:fail") {
			data := strings.SplitN(line, ":", 3)
			parserData.TestFail = append(parserData.TestFail, data[2])
		}
		if strings.HasPrefix(line, "//gen:") {
			data := common.StringSplitIgnoreEmpty(line, ':')
			parserData = &Data{}
			parserData.StructName = data[1]
			parsers[data[1]] = parserData
		}

		if !strings.HasPrefix(line, "type ") {
			continue
		}

		if parserData.ParserName == "" {
			parserData = &Data{}
			continue
		}
		data := common.StringSplitIgnoreEmpty(line, ' ')
		parserType := data[1]

		for _, parserData := range parsers {
			parserData.ParserType = parserType

			filename := parserData.ParserName
			if parserData.ParserSecondName != "" {
				filename = fmt.Sprintf("%s %s", filename, parserData.ParserSecondName)
			}

			filePath := path.Join(dir, "parsers", cleanFileName(filename)+"_autogenerated.go")
			log.Println(filePath)
			//log.Println(parserData)
			//continue
			f, err := os.Create(filePath)
			die(err)
			defer f.Close()

			err = typeTemplate.Execute(f, parserData)
			die(err)

			//parserData.TestFail = append(parserData.TestFail, "") parsers should not get empty line!
			parserData.TestFail = append(parserData.TestFail, "---")
			parserData.TestFail = append(parserData.TestFail, "--- ---")

			filePath = path.Join(dir, "tests", cleanFileName(filename)+"_autogenerated_test.go")
			log.Println(filePath)
			f, err = os.Create(filePath)
			die(err)
			defer f.Close()

			err = testTemplate.Execute(f, parserData)
			die(err)
		}
		parsers = map[string]*Data{}
		parserData = &Data{}
	}
}

func generateTypes(dir string) {
	dat, err := ioutil.ReadFile("types/types.go")
	if err != nil {
		log.Println(err)
	}
	lines := common.StringSplitIgnoreEmpty(string(dat), '\n')
	//fmt.Print(lines)

	parserData := Data{}
	for _, line := range lines {
		if strings.HasPrefix(line, "//sections:") {
			//log.Println(line)
		}
		if strings.HasPrefix(line, "//name:") {
			data := common.StringSplitIgnoreEmpty(line, ':')
			items := common.StringSplitIgnoreEmpty(data[1], ' ')
			parserData.ParserName = data[1]
			if len(items) > 1 {
				parserData.ParserName = items[0]
				parserData.ParserSecondName = items[1]
			}
		}
		if strings.HasPrefix(line, "//is-multiple:true") {
			parserData.ParserMultiple = true
		}
		if strings.HasPrefix(line, "//no-init:true") {
			parserData.NoInit = true
		}
		if strings.HasPrefix(line, "//no-parse:true") {
			parserData.NoParse = true
		}
		if strings.HasPrefix(line, "//test:ok") {
			data := strings.SplitN(line, ":", 3)
			parserData.TestOK = append(parserData.TestOK, data[2])
		}
		if strings.HasPrefix(line, "//test:fail") {
			data := strings.SplitN(line, ":", 3)
			parserData.TestFail = append(parserData.TestFail, data[2])
		}

		if !strings.HasPrefix(line, "type ") {
			continue
		}

		if parserData.ParserName == "" {
			parserData = Data{}
			continue
		}
		data := common.StringSplitIgnoreEmpty(line, ' ')
		parserData.StructName = data[1]
		parserData.ParserType = data[1]

		filename := parserData.ParserName
		if parserData.ParserSecondName != "" {
			filename = fmt.Sprintf("%s %s", filename, parserData.ParserSecondName)
		}

		filePath := path.Join(dir, "parsers", cleanFileName(filename)+"_autogenerated.go")
		log.Println(filePath)
		f, err := os.Create(filePath)
		die(err)
		defer f.Close()

		err = typeTemplate.Execute(f, parserData)
		die(err)

		//parserData.TestFail = append(parserData.TestFail, "") parsers should not get empty line!
		parserData.TestFail = append(parserData.TestFail, "---")
		parserData.TestFail = append(parserData.TestFail, "--- ---")

		filePath = path.Join(dir, "tests", cleanFileName(filename)+"_autogenerated_test.go")
		log.Println(filePath)
		f, err = os.Create(filePath)
		die(err)
		defer f.Close()

		err = testTemplate.Execute(f, parserData)
		die(err)

		parserData = Data{}
	}
}

func cleanFileName(filename string) string {
	return strings.Replace(filename, " ", "-", -1)
}

func die(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var typeTemplate = template.Must(template.New("").Parse(`// Code generated by go generate; DO NOT EDIT.
package parsers

import (
	"fmt"

	"github.com/haproxytech/config-parser/common"
	"github.com/haproxytech/config-parser/errors"
	"github.com/haproxytech/config-parser/types"
)

{{- if not .NoInit }}

func (p *{{ .StructName }}) Init() {
{{- if .ParserMultiple }}
	p.data = []types.{{ .ParserType }}{}
{{- else }}
    p.data = nil
{{- end }}
}
{{- end }}

func (p *{{ .StructName }}) GetParserName() string {
{{- if eq .ParserSecondName "" }}
	return "{{ .ParserName }}"
{{- else }}
	return "{{ .ParserName }} {{ .ParserSecondName }}"
{{- end }}
}

func (p *{{ .StructName }}) Get(createIfNotExist bool) (common.ParserData, error) {
{{- if .ParserMultiple }}
	if len(p.data) == 0 && !createIfNotExist {
		return nil, errors.FetchError
	}
{{- else }}
	if p.data == nil {
		if createIfNotExist {
			p.data = &types.{{ .ParserType }}{}
			return p.data, nil
		}
		return nil, errors.FetchError
	}
{{- end }}
	return p.data, nil
}

func (p *{{ .StructName }}) GetOne(index int) (common.ParserData, error) {
{{- if .ParserMultiple }}
	if len(p.data) == 0 {
		return nil, errors.FetchError
	}
	if index < 0 || index >= len(p.data) {
		return nil, errors.FetchError
	}
	return p.data[index], nil
{{- else }}
	if index != 0 {
		return nil, errors.FetchError
	}
	if p.data == nil {
		return nil, errors.FetchError
	}
	return p.data, nil
{{- end }}
}

func (p *{{ .StructName }}) Set(data common.ParserData, index int) error {
	if data == nil {
		p.Init()
		return nil
	}
{{- if .ParserMultiple }}
	switch newValue := data.(type) {
	case []types.{{ .ParserType }}:
		p.data = newValue
	case *types.{{ .ParserType }}:
		if index > -1 {
			p.data = append(p.data, types.{{ .ParserType }}{})
			copy(p.data[index+1:], p.data[index:])
			p.data[index] = *newValue
		} else {
			p.data = append(p.data, *newValue)
		}
	case types.{{ .ParserType }}:
		if index > -1 {
			p.data = append(p.data, types.{{ .ParserType }}{})
			copy(p.data[index+1:], p.data[index:])
			p.data[index] = newValue
		} else {
			p.data = append(p.data, newValue)
		}
	default:
		return fmt.Errorf("casting error")
	}
{{- else }}
	switch newValue := data.(type) {
	case *types.{{ .ParserType }}:
		p.data = newValue
	case types.{{ .ParserType }}:
		p.data = &newValue
	default:
		return fmt.Errorf("casting error")
	}
{{- end }}
	return nil
}

{{- if and .ParserMultiple (not .NoParse) }}

func (p *{{ .StructName }}) Parse(line string, parts, previousParts []string, comment string) (changeState string, err error) {
{{- if eq .ParserSecondName "" }}
	if parts[0] == "{{ .ParserName }}" {
{{- else }}
	if len(parts) > 1 && parts[0] == "{{ .ParserName }}"  && parts[1] == "{{ .ParserSecondName }}" {
{{- end }}
		data, err := p.parse(line, parts, comment)
		if err != nil {
			return "", &errors.ParseError{Parser: "{{ .StructName }}", Line: line}
		}
		p.data = append(p.data, *data)
		return "", nil
	}
	return "", &errors.ParseError{Parser: "{{ .StructName }}", Line: line}
}
{{- end }}
`))

var testTemplate = template.Must(template.New("").Parse(`// Code generated by go generate; DO NOT EDIT.
package tests

import (
	"testing"

	"github.com/haproxytech/config-parser/parsers"
)

{{ $StructName := .StructName }}
{{- range $index, $val := .TestOK}}
func Test{{ $StructName }}Normal{{$index}}(t *testing.T) {
	err := ProcessLine("{{- $val -}}", &parsers.{{ $StructName }}{})
	if err != nil {
		t.Errorf(err.Error())
	}
}
{{- end }}

{{- range $index, $val := .TestFail}}
func Test{{ $StructName }}Fail{{$index}}(t *testing.T) {
	err := ProcessLine("{{- $val -}}", &parsers.{{ $StructName }}{})
	if err == nil {
		t.Errorf("no data")
	}
}
{{- end }}
`))