package handlers

import (
	"encoding/json"
	"miniK8s/pkg/apiObject"
	"miniK8s/pkg/entity"
	msgutil "miniK8s/pkg/apiserver/msgUtil"
	"miniK8s/pkg/apiserver/serverconfig"
	"miniK8s/pkg/k8log"
	"miniK8s/util/uuid"

	"github.com/gin-gonic/gin"
)

// 添加新的Service
// POST "/api/v1/namespaces/:namespace/services"
func AddService(c *gin.Context) {
	// log
	k8log.InfoLog("APIServer", "AddService: add new service")
	// POST请求，获取请求体
	var service apiObject.Service
	if err := c.ShouldBind(&service); err != nil {
		c.JSON(500, gin.H{
			"error": "parser service failed " + err.Error(),
		})

		k8log.ErrorLog("APIServer", "AddService: parser service failed "+err.Error())
		return
	}

	// 检查name是否重复
	res, err := EtcdStore.PrefixGet(serverconfig.EtcdServicePath + service.Metadata.Name)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "get service failed " + err.Error(),
		})
		k8log.ErrorLog("APIServer", "AddService: get service failed "+err.Error())
		return
	}

	if len(res) != 0 {
		c.JSON(500, gin.H{
			"error": "service name already exist",
		})
		k8log.ErrorLog("APIServer", "AddService: service name already exist")
		return
	}
	// 检查Service的kind是否正确
	if service.Kind != "Service" {
		c.JSON(500, gin.H{
			"error": "service kind is not Service",
		})
		k8log.ErrorLog("APIServer", "AddService: service kind is not Service")
		return
	}

	// 给Service设置UUID, 所以哪怕用户故意设置UUID也会被覆盖
	service.Metadata.UUID = uuid.NewUUID()

	// 将Service转化为ServiceStore
	serviceStore := service.ToServiceStore()

	// 把serviceStore转化为json
	serviceJson, err := json.Marshal(serviceStore)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "service marshal to json failed" + err.Error(),
		})
		return
	}

	// 将Service信息写入etcd
	err = EtcdStore.Put(serverconfig.EtcdServicePath+service.Metadata.Name, serviceJson)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "put service to etcd failed" + err.Error(),
		})
		return
	}
	// 返回201处理成功
	c.JSON(201, gin.H{
		"message": "create service success",
	})

	serviceUpdate := &entity.ServiceUpdate{
		Action: entity.CREATE,
		ServiceTarget: entity.ServiceWithEndpoints{
			Service: service,
			Endpoints: make([]apiObject.Endpoint, 0),
		},
	}

	msgutil.PublishUpdateService(serviceUpdate)

}

// 获取单个Service信息
// 某个特定的Service状态 对应的ServiceSpecURL = "/api/v1/services/:name"
func GetService(c *gin.Context) {
	// 尝试解析请求里面的name
	name := c.Param("name")
	// log
	logStr := "GetSerive: name = " + name
	k8log.InfoLog("APIServer", logStr)

	// 如果解析成功，返回对应的Service信息
	if name != "" {
		res, err := EtcdStore.PrefixGet(serverconfig.EtcdServicePath + name)
		if err != nil {
			c.JSON(400, gin.H{
				"error": "get service failed " + err.Error(),
			})
			return
		}
		// 没找到
		if len(res) == 0 {
			c.JSON(404, gin.H{
				"error": "get service err, not find service",
			})
			return
		}

		// 处理res，如果发现有多个Service，返回错误
		if len(res) != 1 {
			c.JSON(500, gin.H{
				"error": "get service err, find more than one service",
			})
			return
		}
		// 遍历res，返回对应的Service信息
		targetService := res[0].Value
		c.JSON(200, gin.H{
			"data": targetService,
		})
		return
	} else {
		c.JSON(404, gin.H{
			"error": "name is empty",
		})
		return
	}
}

// 获取所有Service信息
func GetServices(c *gin.Context) {
	res, err := EtcdStore.PrefixGet(serverconfig.EtcdServicePath)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "get services failed " + err.Error(),
		})
		return
	}
	// 遍历res，返回对应的Service信息
	var services []string
	for _, service := range res {
		services = append(services, service.Value)
	}
	c.JSON(200, gin.H{
		"data": services,
	})
}

// 删除Service信息
func DeleteService(c *gin.Context) {
	// 尝试解析请求里面的name
	name := c.Params.ByName("name")
	// 如果解析成功，删除对应的Service信息
	if name != "" {
		// log
		logStr := "DeleteService: name = " + name
		k8log.InfoLog("APIServer", logStr)

		err := EtcdStore.Del(serverconfig.EtcdServicePath + name)
		if err != nil {
			c.JSON(400, gin.H{
				"error": "delete service failed " + err.Error(),
			})
			return
		}
		c.JSON(204, gin.H{
			"message": "delete service success",
		})
		return
	} else {
		c.JSON(404, gin.H{
			"error": "name is empty",
		})
		return
	}
}