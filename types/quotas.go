package types

type Quotas struct {
	Agents          uint
	CoreInstances   uint
	Pipelines       uint
	PipelineFiles   uint
	PipelineSecrets uint
}
