package defaulted

import (
	"bytes"
	"custom-go/customize"
	"custom-go/iot/base"
	"custom-go/pkg/types"
	"github.com/buger/jsonparser"
	"github.com/labstack/echo/v4"
	"github.com/wundergraph/graphql-go-tools/pkg/lexer/literal"
)

func init() {
	PostEvent.SetHandler(func(logger echo.Logger, client *types.InternalClient, properties *base.ApplicationProperties, output *base.EventTslOutput[PostEventOutput]) error {
		go customize.PushEventData(PostEvent, properties, output)
		return nil
	}, base.TopicTrieOptions{
		ExtractValue:     postEventExtractor,
		CheckRepeatValue: true,
	})
}

func postEventExtractor(data []byte) []byte {
	rewriteBuffer := bytes.NewBuffer(literal.LBRACE)
	_ = jsonparser.ObjectEach(data,
		func(key []byte, value []byte, _ jsonparser.ValueType, _ int) error {
			_value, _valueType, _, _err := jsonparser.Get(value, "value")
			if _err != nil {
				return _err
			}
			if rewriteBuffer.Len() > 1 {
				rewriteBuffer.Write(literal.COMMA)
			}
			rewriteBuffer.Write(literal.QUOTE)
			rewriteBuffer.Write(key)
			rewriteBuffer.Write(literal.QUOTE)
			rewriteBuffer.Write(literal.COLON)
			if _valueType == jsonparser.String {
				rewriteBuffer.Write(literal.QUOTE)
				defer rewriteBuffer.Write(literal.QUOTE)
			}
			rewriteBuffer.Write(_value)
			return nil
		},
		"items")
	rewriteBuffer.Write(literal.RBRACE)
	return rewriteBuffer.Bytes()
}
