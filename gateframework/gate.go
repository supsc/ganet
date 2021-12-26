package gateframework

import (
	"github.com/supsc/ganet/network"
	_ "net"
	_ "reflect"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/supsc/ganet/log"
)

type IGateService interface {
	GetAgentActor(Agent) (*actor.PID, error)
	GetRootContext( )* actor.RootContext
}

type Gate struct {
	MaxConnNum      int
	PendingWriteNum int
	MaxMsgLen       uint32
	Processor       network.Processor
	Message func(agent *GFAgent, data []byte)
	// websocket
	WSAddr      string
	HTTPTimeout time.Duration
	CertFile    string
	KeyFile     string

	// tcp
	TCPAddr      string
	LenMsgLen    int
	LittleEndian bool

	//实例
	wsServer  *network.WSServer
	tcpServer *network.TCPServer


}

func (gate *Gate) Run(gs IGateService) {

	var wsServer *network.WSServer
	if gate.WSAddr != "" {
		wsServer = new(network.WSServer)
		wsServer.Addr = gate.WSAddr
		wsServer.MaxConnNum = gate.MaxConnNum
		wsServer.PendingWriteNum = gate.PendingWriteNum
		wsServer.MaxMsgLen = gate.MaxMsgLen
		wsServer.HTTPTimeout = gate.HTTPTimeout
		wsServer.CertFile = gate.CertFile
		wsServer.KeyFile = gate.KeyFile
		wsServer.NewAgent = func(conn *network.WSConn) network.Agent {
			a := &GFAgent{conn: conn, gate: gate,message: gate.Message, netType: WEB_SOCKET}
			if gate.Message == nil {
				ac, err := gs.GetAgentActor(a)
				if err != nil {
					//todo:应该不会发生吧
					log.Error("NewAgent fail:%v", err.Error())
				}
				a.Pid = ac
			}

			a.SetActorSystem(gs.GetRootContext().ActorSystem())
			return a
		}
	}

	var tcpServer *network.TCPServer
	if gate.TCPAddr != "" {
		tcpServer = new(network.TCPServer)
		tcpServer.Addr = gate.TCPAddr
		tcpServer.MaxConnNum = gate.MaxConnNum
		tcpServer.PendingWriteNum = gate.PendingWriteNum
		tcpServer.LenMsgLen = gate.LenMsgLen
		tcpServer.MaxMsgLen = gate.MaxMsgLen
		tcpServer.LittleEndian = gate.LittleEndian
		tcpServer.NewAgent = func(conn *network.TCPConn) network.Agent {
			a := &GFAgent{conn: conn, gate: gate, netType: WEB_SOCKET}
			if gate.Message != nil {
				a.SetMessage(gate.Message)
			} else {
				ac, err := gs.GetAgentActor(a)
				if err != nil {
					//todo:应该不会发生吧
					log.Error("NewAgent fail:%v", err.Error())
				}
				a.Pid = ac
			}
			return a
		}
	}

	if wsServer != nil {
		wsServer.Start()
	}
	if tcpServer != nil {
		tcpServer.Start()
	}

	gate.tcpServer = tcpServer
	gate.wsServer = wsServer
}

func (gate *Gate) OnDestroy() {
	if gate.wsServer != nil {
		gate.wsServer.Close()
	}
	if gate.tcpServer != nil {
		gate.tcpServer.Close()
	}
}
