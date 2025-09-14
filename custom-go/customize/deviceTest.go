package customize

import (
	"custom-go/iot/base"
	"custom-go/pkg/plugins"
	"custom-go/pkg/utils"
	"encoding/json"
	"github.com/graphql-go/graphql"
	"strings"
	"sync"
	"time"
)

type (
	DeviceTestArguments struct {
		DeviceName string `json:"deviceName"`
	}
	DeviceTestEventResult struct {
		MessageId       string    `json:"messageId"`
		GeneratedTime   time.Time `json:"generatedTime"`
		EventIdentifier string    `json:"eventIdentifier"`
		EventName       string    `json:"eventName"`
		EventData       string    `json:"eventData,omitempty"`
	}
)

func PushEventData[T any](handler *base.EventTslHandler[T], properties *base.ApplicationProperties, output *base.EventTslOutput[T]) {
	valueBytes, _ := json.Marshal(output)
	eventResult := &DeviceTestEventResult{
		MessageId:       properties.MessageId,
		GeneratedTime:   properties.GenerateTime,
		EventIdentifier: handler.Event.Identifier,
		EventName:       handler.Event.Name,
		EventData:       string(valueBytes),
	}
	eventKeyPrefix := output.DeviceName + "|"
	deviceNameEventChan.Range(func(k, v any) bool {
		if !strings.HasPrefix(k.(string), eventKeyPrefix) {
			return true
		}
		eventChan, ok := v.(chan interface{})
		if !ok {
			deviceNameEventChan.Delete(k)
			return true
		}
		eventChan <- eventResult
		return true
	})
}

var (
	deviceNameEventChan           = &sync.Map{}
	deviceTest_subscriptionFields = graphql.Fields{
		"eventResult": &graphql.Field{
			Args: plugins.BuildGraphqlInput[DeviceTestArguments](),
			Type: plugins.BuildGraphqlOutput[DeviceTestEventResult](),
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return params.Source, nil
			},
			Subscribe: func(params graphql.ResolveParams) (interface{}, error) {
				_, arguments, err := plugins.ResolveArgs[DeviceTestArguments](params)
				if err != nil {
					return nil, err
				}

				uniqueKey := arguments.DeviceName + "|" + utils.RandStr(16)
				resultChan := make(chan interface{})
				deviceNameEventChan.Swap(uniqueKey, resultChan)
				go func() {
					<-params.Context.Done()
					close(resultChan)
					deviceNameEventChan.CompareAndDelete(uniqueKey, resultChan)
				}()
				return resultChan, nil
			},
		},
	}
)

var Device_test_schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: plugins.EmptyRootQuery,
	Subscription: graphql.NewObject(graphql.ObjectConfig{
		Name:   "Subscription",
		Fields: deviceTest_subscriptionFields,
	}),
})

func init() {
	plugins.RegisterGraphql(&Device_test_schema)
}
