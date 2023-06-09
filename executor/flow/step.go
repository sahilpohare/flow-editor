package flow

type FlowContext struct {
	FlowId   string `json:"flowId"`
	Data     map[string]interface{}
	ExitChan chan bool
}

func (f *FlowContext) ExitFlow() {
	f.ExitChan <- true
}

type Step interface {
	// Run executes the step.
	Run(ctx *FlowContext, data map[string]interface{}) (map[string]interface{}, error)
}
