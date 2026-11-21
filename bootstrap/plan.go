package bootstrap

import baize "github.com/bpcoder16/baize/bootstrap"

// NewPlan 按启动顺序声明共享资源、业务准备步骤和常驻服务。
func NewPlan() baize.Plan {
	return baize.Plan{
		Resources: resources(),
		Prepare:   prepareSteps(),
		Services:  services(),
	}
}
