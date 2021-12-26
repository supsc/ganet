package processor

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/supsc/ganet/processor/msgs"
)

//json pb
var protoType string = "json"

//json pb 默认json
func SetProtoType(t string) {
	protoType = t
}

//消息组帧
func Marshal(msgType msgs.MsgType,bodyBuf []byte) (buf []byte,  err error){
	headData := &msgs.MsgHead{}
	headData.Type = msgType
	headData.CheckNum =12
	headData.Encrypt=1
	buf, err = proto.Marshal(headData)
	if err !=nil{
		return
	}
	buf=append(buf,bodyBuf...)
	return
}

func ProtoParse(buf []byte) (body *msgs.MsgBody, err error) {
	if len(buf) < 9 {
		err = errors.New("消息格式错误")
		return
	}

	headBuf := buf[:9]
	headData := &msgs.MsgHead{}
	err = proto.Unmarshal(headBuf, headData)
	if err != nil {
		return
	}
	switch headData.Type {
	case msgs.MsgType_REQUEST:
		bodyBuf := buf[9:]
		body = &msgs.MsgBody{}
		err = proto.Unmarshal(bodyBuf, body)
		if err != nil {
			return
		}
		return
	default:
		err = errors.New("不支持的消息类型")
		return
	}
}