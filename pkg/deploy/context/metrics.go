package context

type Metrics struct {
	CheClusterUpTime float64      `json:"che_cluster_up_time"`
	ChePods          []Containers `json:"che_pods"`
	Workspaces       `json:"workspaces"`
}

type ChePods struct {
	DevfileRegistry Containers `json:"devfile_registry"`
	PluginRegistry  Containers `json:"plugin_registry"`
	CheServer       Containers `json:"che_server"`
}

type Workspaces struct {
	JavaMaven `json:"java_maven"`
	NodeJS    `json:"nodejs"`
	Sample    `json:"sample"`
}

type JavaMaven struct {
	Containers []Containers `json:"containers"`
}

type NodeJS struct {
	Containers []Containers `json:"containers"`
}

type Sample struct {
	Containers []Containers `json:"containers"`
}

type Containers struct {
	ContainerName  string  `json:"container_name"`
	StartupLatency float64 `json:"startup_latency"`
}
