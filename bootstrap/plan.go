package bootstrap

import baize "github.com/bpcoder16/baize/bootstrap"

// NewPlan 将共享资源和业务准备步骤组装成 demo 的启动计划。
func NewPlan() baize.Plan {
	return baize.Plan{
		Resources: resources(),
		Prepare:   prepareSteps(),
	}
}
