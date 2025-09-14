package base

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"custom-go/generated"
	"custom-go/operation/dict/item/upsertMany"
	"custom-go/pkg/types"
	"encoding/base64"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
	"os"
	"pack.ag/amqp"
	"strings"
	"time"
)

const (
	// 表示客户端ID，需您自定义，长度不可超过64个字符。建议使用您的AMQP客户端所在服务器UUID、MAC地址、IP等唯一标识
	cloudClientIdEnv = "ALIBABA_CLOUD_CLIENT_ID"
	// 当前物联网平台对应实例中的消费组ID
	cloudConsumerGroupIdEnv = "ALIBABA_CLOUD_CONSUMER_GROUP_ID"
	// 当前物联网平台实例的ID
	cloudIotInstanceIdEnv = "ALIBABA_CLOUD_IOT_INSTANCE_ID"
	// AMQP接入域名
	cloudHostEnv = "ALIBABA_CLOUD_HOST"
	// 取值为阿里云主账号的AccessKey ID
	cloudAccessKeyEnv = "ALIBABA_CLOUD_ACCESS_KEY"
	// 取值为阿里云主账号的AccessKey Secret
	cloudAccessSecretEnv = "ALIBABA_CLOUD_ACCESS_SECRET"
	// 忽略接收AMQP消息
	cloudIgnoreReceiveAMQPEnv = "IGNORE_RECEIVE_AMQP"

	cloudAmqpManagerSettingDictCode              = "amqp_manager_setting"
	cloudAmqpManagerSettingDictKeyIncludeDevices = "include_devices"
	cloudAmqpManagerSettingDictKeyExcludeDevices = "exclude_devices"
)

func init() {
	types.AddRegisteredHook(func(logger echo.Logger) {
		if os.Getenv(cloudIgnoreReceiveAMQPEnv) == "1" {
			logger.Debugf("amqp ignore receive because of existing env [%s=1]", cloudIgnoreReceiveAMQPEnv)
			return
		}

		accessKey, accessSecret := os.Getenv(cloudAccessKeyEnv), os.Getenv(cloudAccessSecretEnv)
		address := fmt.Sprintf("amqps://%s:5671", os.Getenv(cloudHostEnv))
		timestamp := time.Now().Nanosecond() / int(time.Millisecond)
		userName := fmt.Sprintf("%s|authMode=aksign,signMethod=Hmacsha1,consumerGroupId=%s,authId=%s,iotInstanceId=%s,timestamp=%d|",
			os.Getenv(cloudClientIdEnv), os.Getenv(cloudConsumerGroupIdEnv), accessKey, os.Getenv(cloudIotInstanceIdEnv), timestamp)
		stringToSign := fmt.Sprintf("authId=%s&timestamp=%d", accessKey, timestamp)
		hmacKey := hmac.New(sha1.New, []byte(accessSecret))
		hmacKey.Write([]byte(stringToSign))
		password := base64.StdEncoding.EncodeToString(hmacKey.Sum(nil))
		amqpManager := &AmqpManager{
			address:  address,
			userName: userName,
			password: password,
			logger:   logger,
		}

		amqpManager.refreshSetting(types.NewEmptyInternalClient())
		upsertMany.Subscribe(cloudAmqpManagerSettingDictCode, amqpManager.refreshSetting)
		ctx := context.Background()
		amqpManager.startReceiveMessage(ctx)
	})
}

type AmqpManager struct {
	logger             echo.Logger
	address            string
	userName           string
	password           string
	includeDeviceNames []string
	excludeDeviceNames []string

	client   *amqp.Client
	session  *amqp.Session
	receiver *amqp.Receiver
}

func (am *AmqpManager) refreshSetting(client *types.InternalClient) {
	amqpManagerDictInput := generated.Dict__item__findManyInternalInput{DictCode: cloudAmqpManagerSettingDictCode, Enabled: true}
	amqpManagerDictResp, _ := generated.Dict__item__findMany.Execute(amqpManagerDictInput, client)
	for _, item := range amqpManagerDictResp.Data {
		itemValue := strings.TrimSpace(item.Value)
		if len(itemValue) == 0 {
			continue
		}
		itemValueSlice := strings.Split(itemValue, ",")
		switch item.Key {
		case cloudAmqpManagerSettingDictKeyIncludeDevices:
			am.includeDeviceNames = itemValueSlice
		case cloudAmqpManagerSettingDictKeyExcludeDevices:
			am.excludeDeviceNames = itemValueSlice
		}
	}
}

type ApplicationProperties struct {
	Qos          int
	Topic        string
	MessageId    string
	GenerateTime time.Time
	AcceptedAt   time.Time
	Repeated     bool
}

// 业务函数。用户自定义实现，该函数被异步执行，请考虑系统资源消耗情况。
func (am *AmqpManager) processMessage(message *amqp.Message) {
	appProperties := &ApplicationProperties{
		Qos:          cast.ToInt(message.ApplicationProperties["qos"]),
		Topic:        cast.ToString(message.ApplicationProperties["topic"]),
		MessageId:    cast.ToString(message.ApplicationProperties["messageId"]),
		GenerateTime: castToTime(cast.ToInt64(message.ApplicationProperties["generateTime"])),
		AcceptedAt:   time.Now(),
	}
	if err := am.executeHandler(appProperties, message.GetData()); err != nil {
		am.logger.Errorf("process message for topic [%s] with error: %v", appProperties.Topic, err)
	}
}

func (am *AmqpManager) startReceiveMessage(ctx context.Context) {
	if err := am.generateReceiverWithRetry(ctx); err != nil {
		return
	}
	defer func() {
		_ = am.receiver.Close(ctx)
		_ = am.session.Close(ctx)
		_ = am.client.Close()
	}()

	for {
		// 阻塞接受消息，如果ctx是background则不会被打断。
		message, err := am.receiver.Receive(ctx)
		if err != nil {
			am.logger.Errorf("amqp receive data error: %v", err)
			// 如果是主动取消，则退出程序。
			if ctx.Err() != nil {
				return
			}

			// 非主动取消，则重新建立连接。
			if err = am.generateReceiverWithRetry(ctx); nil != err {
				return
			}

			continue
		}

		go am.processMessage(message)
		_ = message.Accept()
	}
}

func (am *AmqpManager) generateReceiverWithRetry(ctx context.Context) error {
	//退避重连，从10ms依次x2，直到20s。
	duration := 10 * time.Millisecond
	maxDuration := 20000 * time.Millisecond
	times := 1

	//异常情况，退避重连。
	for {
		if ctx.Err() != nil {
			return amqp.ErrConnClosed
		}

		if err := am.generateReceiver(); err == nil {
			am.logger.Debug("amqp connect init success")
			return nil
		}

		time.Sleep(duration)
		if duration < maxDuration {
			duration *= 2
		}
		am.logger.Debugf("amqp connect retry, times: %d, duration: %d", times, duration)
		times++
	}
}

// 由于包不可见，无法判断Connection和Session状态，重启连接获取。
func (am *AmqpManager) generateReceiver() error {
	if am.session != nil {
		receiver, err := am.session.NewReceiver(
			amqp.LinkSourceAddress("/queue-name"),
			amqp.LinkCredit(20),
		)
		//如果断网等行为发生，Connection会关闭导致Session建立失败，未关闭连接则建立成功。
		if err == nil {
			am.receiver = receiver
			return nil
		}
	}

	//清理上一个连接。
	if am.client != nil {
		_ = am.client.Close()
	}

	client, err := amqp.Dial(am.address, amqp.ConnSASLPlain(am.userName, am.password), amqp.ConnConnectTimeout(time.Second*10))
	if err != nil {
		return err
	}
	am.client = client

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	am.session = session

	receiver, err := am.session.NewReceiver(
		amqp.LinkSourceAddress("/queue-name"),
		amqp.LinkCredit(20),
	)
	if err != nil {
		return err
	}
	am.receiver = receiver

	return nil
}
