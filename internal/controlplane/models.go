package controlplane

type RuntimePluginSelector struct {
	PluginUuid string `json:"plugin-uuid"`
}

type RuntimeConfigurationJob struct {
	Uuid              string                   `json:"uuid"`
	ConfigurationUuid string                   `json:"configuration-uuid"`
	RuntimeUuid       string                   `json:"runtime-uuid,omitempty"`
	ActivityId        string                   `json:"activity-id,omitempty"`
	SubjectUuid       string                   `json:"subject-uuid,omitempty"`
	SubjectType       string                   `json:"subject-type,omitempty"`
	Schedule          string                   `json:"schedule,omitempty"`
	Plugins           []*RuntimePluginSelector `json:"plugins,omitempty"`
}
