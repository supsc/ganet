package gateframework

import (
	"github.com/supsc/ganet/processor/msgs"
	"github.com/golang/protobuf/proto"
)
type CloseSocket struct {
	Agent *GFAgent
}
type UpdateClient struct {
	Agent *GFAgent
	Msg *msgs.MsgBody
}

type Push struct {
	Msg  msgs.MsgBody
	Mod  msgs.ModType
	Cmd  int32
	Info string
	Data proto.Message
}