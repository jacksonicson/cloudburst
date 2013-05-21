package benchmark 

import (
	"fmt"
)

type OperationBrowse struct {
	Id                     string
	StartTime, EndTime       int
	Success                  bool
	NumberOfActionsPerformed int
}

func NewOperationBrowse(name string, startTime, endTime int, success bool, numberOfActionsPerformed int) *OperationBrowse{
	return &OperationBrowse{name, startTime, endTime, success, numberOfActionsPerformed}
}

func (o *OperationBrowse) Name() string{
	return o.Id
}

func (o *OperationBrowse) Run() {
	fmt.Printf("running browse operation..\n")
}