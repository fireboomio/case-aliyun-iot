package hbs

import (
	"custom-go/pkg/utils"
	"encoding/json"
	"fmt"
	"github.com/flowchartsman/handlebars/v3"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/cast"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type (
	templateContext struct {
		enumFields   map[string]*enumField
		objectFields map[string]*objectField

		Tsl                        *TslSpec
		Package                    string
		OnceMap                    map[string]any
		EnumFieldArray             []*enumField   // 枚举类型定义
		ObjectFieldArray           []*objectField // 对象类型定义
		SimpleFieldPointerRequired bool           // 普通字段是否指针
	}
	enumField struct {
		Name                string            // 枚举名称
		ValueType           string            // 枚举值类型
		Values              []interface{}     // 枚举值列表
		ValueDescriptionMap map[string]string // 枚举变量名描述
	}
	objectFieldType struct {
		TypeName      string       // 类型名(为字段时使用)
		TypeRefObject *objectField // 类型引用(为字段时使用)
		TypeRefEnum   *enumField   // 枚举引用(为字段时使用)
		Required      bool         // 是否必须(为字段时使用)
		IsArray       bool         // 是否数组(为字段时使用)
	}
	objectField struct {
		objectFieldType
		Name         string         // 对象/字段名
		Description  string         // 对象/字段描述
		IsDefinition bool           // 是否全局定义
		DocumentPath []string       // 文档路径(建议拼接后用来做对象名/字段类型名)
		Fields       []*objectField // 字段列表(为对象时使用)
	}
)

func newParentObjectField(name string) *objectField {
	return &objectField{Name: name, DocumentPath: []string{name}, IsDefinition: true}
}

func GenerateTslTemplate(suffix string) {
	initHelpers()
	ctx := &templateContext{
		enumFields:                 make(map[string]*enumField),
		objectFields:               make(map[string]*objectField),
		Package:                    suffix,
		OnceMap:                    make(map[string]any),
		SimpleFieldPointerRequired: true,
	}
	if err := utils.ReadStructAndCacheFile(fmt.Sprintf("iot/%s/tsl.json", ctx.Package), &ctx.Tsl); err != nil {
		fmt.Printf("ReadStructAndCacheFile failed: %v\n", err)
		return
	}

	ctx.buildProperties()
	ctx.buildServices()
	ctx.buildEvents()
	ctx.optimizeFieldInfo()
	if err := ctx.generateModels(); err != nil {
		fmt.Printf("generateModels failed: %v\n", err)
	}
	if err := ctx.generateOthers(); err != nil {
		fmt.Printf("generateOthers failed: %v\n", err)
	}
}

func (t *templateContext) buildServices() {
	for _, item := range t.Tsl.Services {
		var inputData []*TslSpecOutputData
		for _, itemInput := range item.InputData {
			if itemInput.Struct != nil {
				inputData = append(inputData, itemInput.Struct)
			}
		}
		item.Desc = strings.ReplaceAll(item.Desc, "\n", "")
		t.buildOutputData(inputData, item.Identifier+"ServiceInput", item.Identifier)
		t.buildOutputData(item.OutputData, item.Identifier+"ServiceOutput", item.Identifier)
	}
}

func (t *templateContext) buildEvents() {
	for _, item := range t.Tsl.Events {
		item.Desc = strings.ReplaceAll(item.Desc, "\n", "")
		t.buildOutputData(item.OutputData, item.Identifier+"EventOutput", item.Identifier)
	}
}

func (t *templateContext) buildProperties() {
	property := newParentObjectField("properties")
	for _, item := range t.Tsl.Properties {
		itemField := &objectField{
			Name:        item.Identifier,
			Description: fmt.Sprintf("%s[accessMode:%s]", item.Name, item.AccessMode),
		}
		if item.Desc != "" {
			itemField.Description += fmt.Sprintf("(%s)", item.Desc)
		}
		fieldType, fieldDesc := t.buildDataType(item.DataType, item.Identifier, property.Name)
		fieldType.Required = item.Required
		itemField.objectFieldType = fieldType
		if fieldDesc != "" {
			itemField.Description += fieldDesc
		}
		property.Fields = append(property.Fields, itemField)
	}
	t.objectFields[property.Name] = property
}

// *   **int**: integer
// *   **float**: single-precision floating-point number
// *   **double**: double-precision floating-point number
// *   **enum**: enumeration
// *   **bool**: Boolean
// *   **text**: character
// *   **date**: time (string-type UTC timestamp in milliseconds)
// *   **array**: array
// *   **struct**: structure
func (t *templateContext) buildDataType(dataType *TslSpecDataType, parent, root string) (fieldType objectFieldType, desc string) {
	switch dataType.Type {
	case "int":
		fieldType.TypeName = openapi3.TypeInteger
		number, _ := json.Marshal(dataType.Specs.Number)
		desc = fmt.Sprintf("【%s】", string(number))
	case "float", "double":
		fieldType.TypeName = openapi3.TypeNumber
		number, _ := json.Marshal(dataType.Specs.Number)
		desc = fmt.Sprintf("【%s】", string(number))
	case "enum", "bool":
		enum := &enumField{
			Name:                root + "_" + parent,
			ValueDescriptionMap: dataType.Specs.Map,
		}
		for k := range dataType.Specs.Map {
			enum.Values = append(enum.Values, k)
		}
		fieldType.TypeRefEnum = enum
		if dataType.Type == "bool" {
			enum.ValueType = openapi3.TypeInteger
		}
		t.enumFields[enum.Name] = enum
	case "text":
		fieldType.TypeName = openapi3.TypeString
		text, _ := json.Marshal(dataType.Specs.Text)
		desc = fmt.Sprintf("【%s】", string(text))
	case "array":
		fieldType, _ = t.buildDataType(dataType.Specs.Array.Item, parent, root)
		fieldType.IsArray = true
	case "struct":
		fieldType.TypeName = openapi3.TypeObject
		fieldType.TypeRefObject = t.buildOutputData(dataType.Specs.Struct, parent, root)
	}
	return
}

func (t *templateContext) buildOutputData(data []*TslSpecOutputData, parent, root string) (output *objectField) {
	output, ok := t.objectFields[parent]
	if ok {
		return
	}

	output = newParentObjectField(parent)
	for _, item := range data {
		itemField := &objectField{Name: item.Identifier, Description: item.Name}
		fieldType, fieldDesc := t.buildDataType(item.DataType, item.Identifier, root)
		itemField.objectFieldType = fieldType
		if fieldDesc != "" {
			itemField.Description += fieldDesc
		}
		output.Fields = append(output.Fields, itemField)
	}
	t.objectFields[output.Name] = output
	return
}

// 格式化对象/枚举定义，对字段按名称进行排序
func (t *templateContext) optimizeFieldInfo() {
	for _, enum := range t.enumFields {
		slices.SortFunc(enum.Values, func(a, b any) bool { return cast.ToString(a) < cast.ToString(b) })
	}
	t.EnumFieldArray = maps.Values(t.enumFields)
	slices.SortFunc(t.EnumFieldArray, func(a *enumField, b *enumField) bool { return a.Name < b.Name })

	t.ObjectFieldArray = maps.Values(t.objectFields)
	slices.SortFunc(t.ObjectFieldArray, func(a *objectField, b *objectField) bool { return a.Name < b.Name })
}

func (t *templateContext) generateModels() (err error) {
	// 读取并加载片段函数
	partialsDirname := "../template/golang-server/partials"
	var partials []string
	partialFiles, _ := os.ReadDir(partialsDirname)
	for _, item := range partialFiles {
		partials = append(partials, filepath.Join(partialsDirname, item.Name()))
	}

	modelsHbsFilepath := "../template/golang-server/files/generated/models.go.hbs"
	fileBytes, err := os.ReadFile(modelsHbsFilepath)
	if err != nil {
		return
	}
	template, err := handlebars.Parse(string(fileBytes))
	if err != nil {
		return
	}
	if err = template.RegisterPartialFiles(partials...); err != nil {
		return
	}
	content, err := template.Exec(t)
	if err != nil {
		return
	}
	content = strings.ReplaceAll(content, "package generated", "package "+t.Package)
	content = strings.ReplaceAll(content, "&quot;", "\"")
	writePath := fmt.Sprintf("iot/%s/models.go", t.Package)
	err = os.WriteFile(writePath, []byte(content), os.ModePerm)
	return
}

func (t *templateContext) generateOthers() (err error) {
	hbsDir := "iot/hbs"
	dirEntry, _ := os.ReadDir(hbsDir)
	var (
		hbsFileBytes []byte
		hbsTemplate  *handlebars.Template
		writeContent string
	)
	for _, item := range dirEntry {
		before, ok := strings.CutSuffix(item.Name(), ".hbs")
		if !ok {
			continue
		}
		if hbsFileBytes, err = os.ReadFile(filepath.Join(hbsDir, item.Name())); err != nil {
			return
		}
		if hbsTemplate, err = handlebars.Parse(string(hbsFileBytes)); err != nil {
			return
		}
		if writeContent, err = hbsTemplate.Exec(t); err != nil {
			return
		}
		writePath := strings.ReplaceAll(filepath.Join(hbsDir, before), "hbs", t.Package)
		if err = os.WriteFile(writePath, []byte(writeContent), os.ModePerm); err != nil {
			return
		}
	}
	return
}

var initOnce sync.Once

func initHelpers() {
	initOnce.Do(func() {
		// 去除前置空格
		handlebars.RegisterHelper("trimPrefix", func(str string, cut string) string {
			cutLength := len(cut)
			for {
				if !strings.HasPrefix(str, cut) {
					break
				}

				str = str[cutLength:]
			}
			return str
		})
		// 首字母大写
		handlebars.RegisterHelper("upperFirst", func(str string) string {
			strLen := len(str)
			if strLen == 0 {
				return ""
			}

			result := strings.ToUpper(str[:1])
			if strLen > 1 {
				result += str[1:]
			}
			return result
		})
		// 以指定字符串连接数组
		handlebars.RegisterHelper("joinString", func(sep string, strArr []string) string {
			return strings.Join(strArr, sep)
		})
		// 判断任意等于，target为','拼接的多个变量
		handlebars.RegisterHelper("equalAny", func(source string, target string) bool {
			return slices.Contains(strings.Split(target, ","), source)
		})
		// 判断任意等于，target为','拼接的多个变量
		handlebars.RegisterHelper("getMapValue", func(data map[string]string, key string) string {
			return data[key]
		})
		// 判断是否唯一，可用于判断import等唯一性
		handlebars.RegisterHelper("isAbsent", func(onceMap map[string]any, name string, val any) bool {
			if utils.IsZeroValue(val) {
				return false
			}

			if _, ok := onceMap[name]; ok {
				return false
			}

			onceMap[name] = val
			return true
		})
	})
}
