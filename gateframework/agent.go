package gateframework

import (
	"github.com/supsc/ganet/log"
	"github.com/supsc/ganet/network"

	"net"
	"reflect"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/supsc/ganet/processor"
)

type NetType byte

const (
	TCP        NetType = 0
	WEB_SOCKET NetType = 1
)

type Agent interface {
	WriteMsg(msg []byte)
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close()
	Destroy()
	UserData() interface{}
	SetUserData(data interface{})
	SetDead()
	GetNetType() NetType
	SetActorSystem(system *actor.ActorSystem)
	GetRootContext() *actor.RootContext
	SetMessage(func(*GFAgent, []byte))
	GetUid() float64
	GetConn() network.Conn
}

type GFAgent struct {
	Uid              float64
	conn             network.Conn
	gate             *Gate
	agentActorSystem *actor.ActorSystem
	RootContext      *actor.RootContext
	Pid              *actor.PID
	userData         interface{}
	dead             bool
	netType          NetType
	message          func(agent *GFAgent, data []byte)
}

func (a *GFAgent) GetNetType() NetType {
	return a.netType
}
func (a *GFAgent) SetDead() {
	a.dead = true
}

func (a *GFAgent) Run() {

	for {
		data, err := a.conn.ReadMsg()
		if err != nil {
			log.Debug("read message: %v", err)
			break
		}

		if a.message != nil {
			a.message(a, data)
		} else {

			//
			result, parseerr := processor.ProtoParse(data)
			if parseerr != nil {
				log.Error("消息解析失败,err:%s msg: %s", err, data)
				break
			}

			if result.Mod == 0 || result.Cmd == 0 {
				log.Error("Mod和Cmd是0")
				break
			}
			//
			//msg := &msgs.MsgBody{}
			//msg := .(msgs.MsgBody)
			a.RootContext.Send(a.Pid, result)
			//a.rootContext.Send(a.agentActor, &msgs.MsgBody{Data: data})

			//break
			//
			//if result.Mod == server.ModType_BASE && int32(result.Cmd) == int32(base.Cmd_TOKEN) {
			//
			//}

			//if a.gate.Processor != nil {
			//	msg, err := a.gate.Processor.Unmarshal(data)
			//	if err != nil {
			//		log.Debug("unmarshal message error: %v", err)
			//		break
			//	}
			//	err = a.gate.Processor.Route(msg, a)
			//	if err != nil {
			//		log.Debug("route message error: %v", err)
			//		break
			//	}
			//
			//} else {
			//todo:not safe
			//a.agentActor.Tell(&msgs.ReceviceClientMsg{data})
			//a.agentActorRoot.Send(a.agentActor,&msgs.MsgBody)
			//if err != nil {
			//	log.Error("ReceviceClientMsg message error: %v", err)
			//	break
			//}
			//}
		}

	}
}

func (a *GFAgent) OnClose() {
	//todo:not safe
	if a.Pid != nil && !a.dead {
		//a.agentActor.Tell(&msgs.ClientDisconnect{})
		a.RootContext.Send(a.Pid, &CloseSocket{Agent: a})
	}

}

func (a *GFAgent) WriteMsg(data []byte) {
	err := a.conn.WriteMsg(data)
	if err != nil {
		log.Error("write message %v error: %v", reflect.TypeOf(data), err)
	}

}

func (a *GFAgent) LocalAddr() net.Addr {
	return a.conn.LocalAddr()
}

func (a *GFAgent) RemoteAddr() net.Addr {
	return a.conn.RemoteAddr()
}

func (a *GFAgent) Close() {
	a.conn.Close()
}

func (a *GFAgent) Destroy() {
	a.conn.Destroy()

}

func (a *GFAgent) UserData() interface{} {
	return a.userData
}

func (a *GFAgent) SetUserData(data interface{}) {
	a.userData = data
}

func (a *GFAgent) SetActorSystem(system *actor.ActorSystem) {
	a.agentActorSystem = system
	a.RootContext = system.Root
}

func (a GFAgent) GetRootContext() *actor.RootContext {
	return a.RootContext
}


func (a GFAgent) SetMessage(message func(*GFAgent, []byte)) {
	a.message = message
}

func (a GFAgent) GetUid() float64 {
	return a.Uid
}

func (a GFAgent) GetConn() network.Conn {
	return a.conn
}
