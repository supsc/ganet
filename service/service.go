/*
**服务类型，可作为独立进程或线程使用
**
 */
package service

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/supsc/ganet/log"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
)

type Context actor.Context

type IService interface {
	IServiceData
	// OnReceive Receive(context actor.Context)
	OnReceive(context Context)
	OnInit()
	OnStart(as *ActorService)
	// OnRun 正式运行(服务线程)
	OnRun()

	OnDestory()
}

type ServiceRun struct {
}

//interface
//func (s *BaseServer) OnReceive(context Context)            {}
//func (s *BaseServer) OnInit()                              {}
//func (s *BaseServer) OnStart()                             {}

type MessageFunc func(context Context)

// ActorService 服务的代理
type ActorService struct {
	serviceIns IService
	rounter    map[reflect.Type]MessageFunc
}

func (s *ActorService) Receive(context actor.Context) {
	//switch msg := context.Message().(type) {
	//case *hello:
	//	fmt.Printf("Hello %v\n", msg.Who)
	//}
	switch msg := context.Message().(type) {
	case *actor.Started:
		fmt.Println("Started, initialize actor here:",s.serviceIns.GetName())
		s.serviceIns.SetPID(context.Self())
		s.serviceIns.SetActorSystem(context.ActorSystem())
		s.serviceIns.OnStart(s)
		s.serviceIns.OnRun()
	case *actor.Stopping:
		fmt.Println("Stopping, actor is about shut down")
	case *actor.Stopped:
		fmt.Println("Stopped, actor and its children are stopped")
	case *actor.Restarting:
		fmt.Println("Restarting, actor is about restart")
	case *ServiceRun:
		fmt.Println("ServiceRun ", s.serviceIns.GetName())
		//s.serviceIns.OnRun()
	default:
		log.Debug("recv defalult:", msg)
		s.serviceIns.OnReceive(context.(Context))
		fun := s.rounter[reflect.TypeOf(msg)]
		if fun != nil {
			fun(context.(Context))
		}
	}
}

func (s *ActorService) RegisterMsg(t reflect.Type, f MessageFunc) {
	s.rounter[t] = f
}

func StartService(s IService) {
	ac := &ActorService{s, make(map[reflect.Type]MessageFunc)}

	// decider := func(reason interface{}) actor.Directive {
	// 	log.Error("handling failure for child:%v", reason)
	// 	return actor.StopDirective
	// }
	// supervisor := actor.NewOneForOneStrategy(10, 1000, decider)


	system := actor.NewActorSystem()
	rootContext := system.Root

	//sys := actor.NewActorSystem()
	//root := sys.Root
	props := actor.PropsFromProducer(func() actor.Actor { return ac }) //.WithSupervisor(supervisor)
	//props := actor.FromInstance(ac)
	if s.GetAddress() != "" {

		adds := strings.Split(s.GetAddress(),":")
		port,_ := strconv.Atoi(adds[1])
		cfg := remote.Configure(adds[0], port)
		cfg = cfg.WithAdvertisedHost(s.GetAddress())


		//config := remote.Config{
		//	AdvertisedHost:           s.GetAddress(),
		//	DialOptions:              []grpc.DialOption{grpc.WithInsecure()},
		//	EndpointWriterBatchSize:  1000,
		//	EndpointManagerBatchSize: 1000,
		//	EndpointWriterQueueSize:  1000000,
		//	EndpointManagerQueueSize: 1000000,
		//	Kinds:                    make(map[string]*actor.Props),
		//}

		//systemRoot := actor.NewActorSystem()
		remote.NewRemote(system, cfg).Start()
	}

	_, err := rootContext.SpawnNamed(props, s.GetName())
	if err == nil {
		//s.SetPID(pid)
		//s.OnStart(ac)
	} else {
		log.Error("#############actor.SpawnNamed error:%v", err)
	}


}

func DestoryService(s *ActorService) {
	s.serviceIns.OnDestory()
}
