package httpapi

import (
	"context"
)

type ecsController struct{}

func newECSController() *ecsController {
	return &ecsController{}
}

// List ECS列表
// @Permission 资源中心:云服务器:服务器列表
// @Router /api/clouds/resources/ecs [get]
func (c *ecsController) List(ctx context.Context) error {
	return nil
}

// Detail ECS详情
// @Permission 资源中心:云服务器:服务器详情
// @Router /api/clouds/resources/ecs/detail [get]
func (c *ecsController) Detail(ctx context.Context) error {
	return nil
}

// Detail ECS详情
// @Summary ECS详情
// @Description ECS详情
func (c *ecsController) Find(ctx context.Context) error {
	return nil
}

