package cloudburst

import (
	"runtime"
)

type Scenario struct {
	TargetManager *TargetManager
}

func NewScenario(schedule *TargetSchedule, factory Factory) *Scenario {
	return &Scenario{NewTargetManager(schedule, factory)}
}

func (scenario *Scenario) Launch() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	joinChannel := make(chan bool)
	go scenario.TargetManager.processSchedule(joinChannel)
	<-joinChannel
}

func (scenario *Scenario) AggregateStatistics() {
	AggregateScoreboards(scenario.TargetManager.Targets, 10)
}
