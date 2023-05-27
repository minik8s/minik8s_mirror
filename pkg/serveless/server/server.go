package server

import (
	"fmt"
	minik8stypes "miniK8s/pkg/minik8sTypes"
	"miniK8s/pkg/serveless/function"
	"miniK8s/util/executor"

	"github.com/gin-gonic/gin"
)

type Server interface {
	Run()
}

type server struct {
	httpServer *gin.Engine
	routeTable map[string][]string
	// routeTable 的key是 namespace/name ，value是一个数组，数组中的每个元素是一个pod的ip地址
	// 当用户的请求到来了之后，首先会根据func的namespace/name找到对应的pod的ip地址，然后再将请求转发到这个pod上
	// 如果发现这个value为空，那就需要创建一个新的pod，然后将请求转发到这个pod上

	// func的controller
	funcController function.FuncController
}

func NewServer() Server {
	return &server{
		httpServer:     gin.Default(),
		routeTable:     make(map[string][]string),
		funcController: function.NewFuncController(),
	}
}

// 周期性的函数都放在这里

// 从API Server获取所有的Pod信息，根据Label筛选出所有的Function Pod
func (s *server) updateRouteTableFromAPIServer() {
	// TODO
	pods, err := GetAllPodFromAPIServer()

	if err != nil {
		return
	}

	// 遍历所有的pod，将其加入到routeTable中
	for _, pod := range pods {
		// 说明是一个Function Pod
		if pod.Metadata.Labels[minik8stypes.Pod_Func_Uuid] != "" {
			funcName := pod.Metadata.Labels[minik8stypes.Pod_Func_Name]
			funcNamespace := pod.Metadata.Labels[minik8stypes.Pod_Func_Namespace]
			key := funcNamespace + "/" + funcName
			// 检查ip是否为空
			if pod.Status.PodIP != "" {
				// ip不为空，说明这个pod已经启动了，可以将其加入到routeTable中
				s.routeTable[key] = append(s.routeTable[key], pod.Status.PodIP)
				fmt.Println("update routeTable: ", s.routeTable)
			}
		}
	}
}

func (s *server) Run() {
	// 周期性的更新routeTable
	go executor.Period(RouterUpdate_Delay, RouterUpdate_WaitTime, s.updateRouteTableFromAPIServer, RouterUpdate_ifLoop)
	// 周期性的检查function的情况，如果有新创建的function，那么就创建一个新的pod
	go s.funcController.Run()

	// 初始化服务器
	s.httpServer.POST("/:funcNamespace/:funcName", s.handleFuncRequest)
	s.httpServer.Run(":28080")
}
