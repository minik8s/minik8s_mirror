package function

import (
	"errors"
	"miniK8s/pkg/apiObject"
	"miniK8s/pkg/config"
	"miniK8s/util/executor"
	netrequest "miniK8s/util/netRequest"
	"net/http"
)

type FuncController interface {
	Run()
}

type funcController struct {
	cache map[string]*apiObject.Function
}

func NewFuncController() FuncController {
	return &funcController{
		cache: make(map[string]*apiObject.Function),
	}
}

func (c *funcController) getAllFunc() ([]apiObject.Function, error) {
	url := config.API_Server_URL_Prefix + config.GlobalFunctionsURL

	allFuncs := make([]apiObject.Function, 0)

	code, err := netrequest.GetRequestByTarget(url, &allFuncs, "data")

	if err != nil {
		return nil, err
	}

	if code != http.StatusOK {
		return nil, errors.New("get all functions from apiserver failed, not 200")
	}

	return allFuncs, nil
}

func (c *funcController) routine() {
	res, err := c.getAllFunc()
	if err != nil {
		return
	}

	remoteFuncs := make(map[string]bool)
	for _, f := range res {
		remoteFuncs[f.Metadata.UUID] = true
		// 检查f是否在cache中
		if _, ok := c.cache[f.Metadata.UUID]; !ok {
			// 如果不在cache中，说明是新的function，需要创建
			// 【TODO】
			c.cache[f.Metadata.UUID] = &f
			c.CreateFunction(&f)

		} else {
			// 如果在cache中，说明是已经存在的function，需要检查是否需要更新
			// 【TODO】
			if !c.ComplareTwoFunc(c.cache[f.Metadata.UUID], &f) {
				c.cache[f.Metadata.UUID] = &f
				c.UpdateFunction(&f)
			} else {
				c.cache[f.Metadata.UUID] = &f
			}
		}
	}

	// 检查cache中的function是否在remoteFuncs中
	for uuid, f := range c.cache {
		if _, ok := remoteFuncs[uuid]; !ok {
			// 如果不在remoteFuncs中，说明需要删除
			//
			c.DeleteFunction(f)
			delete(c.cache, uuid)
		}
	}

	// 这样保证这边的缓存和apiserver中的缓存一致
}

// 比较两个function是否相同
func (c *funcController) ComplareTwoFunc(old *apiObject.Function, new *apiObject.Function) bool {
	// 【TODO】
	if old.Spec.UserUploadFilePath != new.Spec.UserUploadFilePath {
		return false
	}
	if len(old.Spec.UserUploadFile) != len(new.Spec.UserUploadFile) {
		return false
	}

	return true
}

func (c *funcController) Run() {
	executor.Period(FuncControllerUpdateDelay, FuncControllerUpdateFrequency, c.routine, FuncControllerUpdateLoop)
}