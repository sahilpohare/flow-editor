package flow

type OmitKeysStep struct {
	Keys []string
}

func NewOmitKeysStep(keys []string) *OmitKeysStep {
	return &OmitKeysStep{
		Keys: keys,
	}
}

func (o *OmitKeysStep) Run(ctx *FlowContext, data map[string]interface{}) (map[string]interface{}, error) {
	for _, key := range o.Keys {
		delete(data, key)
	}

	return data, nil
}
