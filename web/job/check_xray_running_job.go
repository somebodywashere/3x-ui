package job

import "x-ui/web/service"
import "x-ui/logger"

type CheckXrayRunningJob struct {
	xrayService service.XrayService

	checkTime int
}

func NewCheckXrayRunningJob() *CheckXrayRunningJob {
	return new(CheckXrayRunningJob)
}

func (j *CheckXrayRunningJob) Run() {
	if j.xrayService.IsXrayRunning() {
		j.checkTime = 0
		return
	}
	j.checkTime++
	if j.checkTime < 2 {
		return
	}
	err := j.xrayService.RestartXray(false)
			if err != nil {
				logger.Error("Restart xray failed:", err)
			}
}
