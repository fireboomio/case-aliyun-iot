package hbs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/tidwall/gjson"
)

type (
	TslSpec struct {
		Schema         string                  `json:"schema"`
		Profile        *TslSpecProfile         `json:"profile"`
		Properties     []*TslSpecProperty      `json:"properties"`
		Events         []*TslSpecEvent         `json:"events"`
		Services       []*TslSpecService       `json:"services"`
		FunctionBlocks []*TslSpecFunctionBlock `json:"functionBlocks"`
		*TslSpecFunctionBlock
	}
	TslSpecProfile struct {
		Version    string `json:"version"`
		ProductKey string `json:"productKey"`
	}
	TslSpecProperty struct {
		Identifier string           `json:"identifier"`
		Name       string           `json:"name"`
		AccessMode string           `json:"accessMode"`
		Desc       string           `json:"desc,omitempty"`
		Required   bool             `json:"required"`
		DataType   *TslSpecDataType `json:"dataType"`
	}
	TslSpecEvent struct {
		Identifier string               `json:"identifier"`
		Name       string               `json:"name"`
		Type       string               `json:"type"`
		Required   bool                 `json:"required"`
		Desc       string               `json:"desc,omitempty"`
		Method     string               `json:"method"`
		OutputData []*TslSpecOutputData `json:"outputData"`
	}
	TslSpecService struct {
		Identifier string               `json:"identifier"`
		Name       string               `json:"name"`
		Required   bool                 `json:"required"`
		CallType   string               `json:"callType"`
		Desc       string               `json:"desc,omitempty"`
		Method     string               `json:"method"`
		InputData  []*TslSpecInputData  `json:"inputData"`
		OutputData []*TslSpecOutputData `json:"outputData"`
	}
	TslSpecFunctionBlock struct {
		FunctionBlockId   string `json:"functionBlockId"`
		FunctionBlockName string `json:"functionBlockName"`
		GmtCreated        int64  `json:"gmtCreated"`
		ProductKey        string `json:"productKey"`
	}
	TslSpecDataType struct {
		Type  string                `json:"type"`
		Specs *TslSpecDataTypeSpecs `json:"specs,omitempty"`
	}
	TslSpecInputData struct {
		String *string
		Struct *TslSpecOutputData
	}
	TslSpecOutputData struct {
		Identifier string           `json:"identifier"`
		Name       string           `json:"name"`
		DataType   *TslSpecDataType `json:"dataType"`
	}
	TslSpecDataTypeSpecs struct {
		Number *struct {
			Min      string `json:"min"`
			Max      string `json:"max"`
			Step     string `json:"step"`
			Unit     string `json:"unit,omitempty"`
			UnitName string `json:"unitName,omitempty"`
		}
		Text *struct {
			Length string `json:"length"`
		}
		Map    map[string]string
		Struct []*TslSpecOutputData
		Array  *struct {
			Size string           `json:"size"`
			Item *TslSpecDataType `json:"item"`
		}
	}
)

func (t *TslSpecDataType) UnmarshalJSON(data []byte) (err error) {
	t.Type = gjson.GetBytes(data, "type").String()
	t.Specs = &TslSpecDataTypeSpecs{}
	specsResult := gjson.GetBytes(data, "specs")
	if !specsResult.Exists() {
		return
	}

	specsBytes := []byte(specsResult.String())
	switch t.Type {
	case "bool", "enum":
		err = json.Unmarshal(specsBytes, &t.Specs.Map)
	case "int", "float", "double":
		err = json.Unmarshal(specsBytes, &t.Specs.Number)
	case "text":
		err = json.Unmarshal(specsBytes, &t.Specs.Text)
	case "array":
		err = json.Unmarshal(specsBytes, &t.Specs.Array)
	case "struct":
		err = json.Unmarshal(specsBytes, &t.Specs.Struct)
	default:
		err = fmt.Errorf("unknown data type: %s", t.Type)
	}
	return
}

func (t *TslSpecInputData) UnmarshalJSON(data []byte) (err error) {
	if !bytes.HasPrefix(data, []byte("{")) {
		t.String = tea.String(string(data))
		return
	}

	return json.Unmarshal(data, &t.Struct)
}
