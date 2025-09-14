package user_identity

import (
	"custom-go/generated"
	"custom-go/pkg/plugins"
	"custom-go/pkg/types"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"golang.org/x/exp/slices"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

type (
	importResult struct {
		Succeed []*importResultItem `json:"succeed"`
		Failed  []*importResultItem `json:"failed"`
	}
	importResultItem struct {
		RowIndex int    `json:"rowIndex"`
		Username string `json:"username"`
		Phone    string `json:"phone,omitempty"`
		Error    string `json:"error,omitempty"`
	}
	createUserInput = generated.User__createOneInternalInput
)

func importExcel(hook *types.HttpTransportHookRequest, body *plugins.HttpTransportBody) (resp *types.WunderGraphResponse, err error) {
	if hook.User == nil {
		err = errors.New("请登录后操作")
		return
	}

	reader, err := body.Request.NewRequest().MultipartReader()
	if err != nil {
		return
	}
	var part *multipart.Part
	for {
		part, err = reader.NextPart()
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}
		if part.FileName() != "" {
			break
		}
	}
	if part == nil {
		err = errors.New("no file found")
		return
	}

	excelFile, err := excelize.OpenReader(part)
	if err != nil {
		return
	}
	rows, err := excelFile.GetRows("Sheet1")
	if err != nil {
		return
	}

	result := &importResult{}
	cellFuncs := make(map[int]func(*createUserInput, string))
	cellFuncNames := make(map[int]string)
	for rowIndex, row := range rows {
		input := &createUserInput{}
		var matchCellFuncNames []string
		for cellIndex, cell := range row {
			if rowIndex == 0 {
				if excelFunc, ok := excelHeaderFuncs[cell]; ok {
					cellFuncs[cellIndex] = excelFunc
					cellFuncNames[cellIndex] = cell
				}
				continue
			}

			if cellFunc, ok := cellFuncs[cellIndex]; ok {
				matchCellFuncNames = append(matchCellFuncNames, cellFuncNames[cellIndex])
				cellFunc(input, strings.TrimSpace(cell))
			}
		}
		if rowIndex > 0 {
			resultItem := &importResultItem{RowIndex: rowIndex, Username: input.Name, Phone: input.Phone}
			var missNames []string
			for _, name := range necessaryExcelHeaders {
				if !slices.Contains(matchCellFuncNames, name) {
					missNames = append(missNames, name)
				}
			}
			if len(missNames) > 0 {
				resultItem.Error = fmt.Sprintf("数据缺失: %s", strings.Join(missNames, ","))
				result.Failed = append(result.Failed, resultItem)
				continue
			}

			_, err = generated.User__createOne.Execute(*input, hook.InternalClient)
			if err != nil {
				resultItem.Error = err.Error()
				result.Failed, err = append(result.Failed, resultItem), nil
				continue
			}
			result.Succeed = append(result.Succeed, resultItem)
		}
	}
	resp = &types.WunderGraphResponse{StatusCode: http.StatusOK}
	resp.OriginBody, _ = json.Marshal(result)
	return
}

func init() {
	plugins.RegisterProxyHook(importExcel)
}

var necessaryExcelHeaders = []string{"用户名"}

var excelHeaderFuncs = map[string]func(*createUserInput, string){
	"用户名": func(input *createUserInput, cell string) {
		input.Name = cell
	},
	"手机号": func(input *createUserInput, cell string) {
		input.Phone = cell
	},
}
