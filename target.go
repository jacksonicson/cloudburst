package cloudburst

import (
	"container/list"
	"fmt"
	"github.com/johanneskross/cloudburst/times"
	"time"
)

const TO_NANO = 1000000000

type Target struct {
	TargetId      int
	Agents        list.List
	AgentChannel  chan bool
	Configuration TargetConfiguration
	Factory       Factory
	Scoreboard    *Scoreboard
	Timing        *Timing
}

func NewTarget(targetConfiguration TargetConfiguration, factory Factory) *Target {
	agents := *list.New()
	channelSize := calcChannelSize(targetConfiguration.TimeSeries.Elements)
	channelSize = 200
	agentChannel := make(chan bool, channelSize)
	scoreboard := NewScoreboard(targetConfiguration.TargetId)
	timing := NewTiming(targetConfiguration.RampUp, targetConfiguration.Duration, targetConfiguration.RampDown)
	return &Target{targetConfiguration.TargetId, agents, agentChannel, targetConfiguration, factory, scoreboard, timing}
}

func calcChannelSize(elements []*times.Element) int {
	channelSize := 0
	for i := 0; i < len(elements); i++ {
		value := int(elements[i].Value)
		if value > channelSize {
			channelSize = value
		}
	}
	return channelSize
}

func (t *Target) RunTimeSeries(c chan bool) {
	fmt.Printf("Running time series on target: %v\n", t.TargetId)

	scoreboardQuitQuannel := make(chan bool)
	go t.Scoreboard.Run(scoreboardQuitQuannel)

	t.Wait(t.Timing.StartSteadyState)
	interval := 0

	for t.Timing.InSteadyState(time.Now().UnixNano()) {
		// wait until next interval is due
		nextInterval := (t.Configuration.TimeSeries.Elements[interval].Timestamp * TO_NANO) + t.Timing.StartSteadyState
		t.Wait(nextInterval)

		runningAgents := len(t.AgentChannel)
		runningNextAgents := int(t.Configuration.TimeSeries.Elements[interval].Value)
		runningNextAgents = 50 // For test reason
		fmt.Printf("Update amount of agents to %v on target%v in interval %v\n", runningNextAgents, t.TargetId, interval)
		interval++

		// update amount of agents for this interval
		switch {
		case runningAgents < runningNextAgents:
			addAgents := runningNextAgents - runningAgents
			startAgents(t, addAgents)
		case runningAgents > runningNextAgents:
			reduceAgents := runningAgents - runningNextAgents
			interruptAgents(t, reduceAgents)
		}
	}
	scoreboardQuitQuannel <- true
	<-scoreboardQuitQuannel
	c <- true
}

func (t *Target) Wait(nextInterval int64) {
	currentTime := time.Now().UnixNano()
	deltaTime := nextInterval - currentTime
	//fmt.Printf("Target %v waits %v seconds for next interval\n", t.TargetId, deltaTime/TO_NANO)
	if deltaTime > 0 {
		time.Sleep(time.Duration(deltaTime))
	}
}

func startAgents(t *Target, amount int) {
	for i := 0; i < amount; i++ {
		agent := NewAgent(t.Agents.Len()+1, t.TargetId, t.Configuration.TargetIp, make(chan bool), t.Factory.CreateGenerator(), t.Scoreboard.OperationResultChannel, t.Timing)
		t.Agents.PushBack(agent)
		go agent.Run(t.AgentChannel)
	}
}

func interruptAgents(t *Target, amount int) {
	for i := 0; i < amount; i++ {
		agentElem := t.Agents.Back()
		agent := agentElem.Value.(*Agent)
		t.Agents.Remove(agentElem)
		go agent.Interrupt(t.AgentChannel)
	}
}
