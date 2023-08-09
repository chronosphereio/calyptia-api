package types

import "time"

// PipelinePortProtocol valid protocol types defined for a pipeline.
type PipelinePortProtocol string

const (
	PipelineProtocolTCP PipelinePortProtocol = "tcp"
	PipelineProtocolUDP PipelinePortProtocol = "udp"
)

// AllPipelinePortProtocols all valid protocol types for a pipeline.
var AllPipelinePortProtocols = [...]PipelinePortProtocol{
	PipelineProtocolTCP,
	PipelineProtocolUDP,
}

// PipelinePort model.
type PipelinePort struct {
	ID           string    `json:"id" yaml:"id"`
	Protocol     string    `json:"protocol" yaml:"protocol"`
	FrontendPort uint      `json:"frontendPort" yaml:"frontendPort"`
	BackendPort  uint      `json:"backendPort" yaml:"backendPort"`
	Endpoint     string    `json:"endpoint" yaml:"endpoint"`
	PluginID     *string   `json:"pluginID" yaml:"pluginID"`
	PluginName   *string   `json:"pluginName" yaml:"pluginName"`
	PluginAlias  *string   `json:"pluginAlias" yaml:"pluginAlias"`
	CreatedAt    time.Time `json:"createdAt" yaml:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt" yaml:"updatedAt"`
}

// PipelinePorts paginated list.
type PipelinePorts struct {
	Items     []PipelinePort `json:"items" yaml:"items"`
	EndCursor *string        `json:"endCursor" yaml:"endCursor"`
}

// CreatePipelinePort request payload for creating a pipeline port.
type CreatePipelinePort struct {
	Protocol     string `json:"protocol"`
	FrontendPort uint   `json:"frontendPort"`
	BackendPort  uint   `json:"backendPort"`
	Endpoint     string `json:"endpoint"`

	pluginID    *string
	pluginName  *string
	pluginAlias *string
}

func CreatePipelinePortWithOpts(base CreatePipelinePort, pluginID, pluginName, pluginAlias *string) CreatePipelinePort {
	base.pluginID = pluginID
	base.pluginName = pluginName
	base.pluginAlias = pluginAlias
	return base
}

func (in *CreatePipelinePort) SetPluginID(pluginID string) {
	in.pluginID = &pluginID
}

func (in *CreatePipelinePort) SetPluginName(pluginName string) {
	in.pluginName = &pluginName
}

func (in *CreatePipelinePort) SetPluginAlias(pluginAlias string) {
	in.pluginAlias = &pluginAlias
}

func (in CreatePipelinePort) PluginID() *string {
	return in.pluginID
}

func (in CreatePipelinePort) PluginName() *string {
	return in.pluginName
}

func (in CreatePipelinePort) PluginAlias() *string {
	return in.pluginAlias
}

// CreatedPipelinePort response payload after creating a pipeline port successfully.
type CreatedPipelinePort struct {
	ID          string    `json:"id"`
	PluginID    *string   `json:"pluginID" yaml:"pluginID"`
	PluginName  *string   `json:"pluginName" yaml:"pluginName"`
	PluginAlias *string   `json:"pluginAlias" yaml:"pluginAlias"`
	CreatedAt   time.Time `json:"createdAt"`
}

// PipelinePortsParams request payload for querying the pipeline ports.
type PipelinePortsParams struct {
	Last   *uint
	Before *string
}

// UpdatePipelinePort request payload for updating a pipeline port.
type UpdatePipelinePort struct {
	Protocol     *string `json:"protocol"`
	FrontendPort *uint   `json:"frontendPort"`
	BackendPort  *uint   `json:"backendPort"`
	Endpoint     *string `json:"endpoint"`
}
