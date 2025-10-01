package manifest

type Manifest struct {
	Version int    `yaml:"version"`
	Service string `yaml:"service"`
	Image   Image  `yaml:"image"`
	Run     Run    `yaml:"run"`
	Deploy  Deploy `yaml:"deploy"`
	Res     *Res   `yaml:"resources,omitempty"` // optional
}

type Image struct {
	Name  string `yaml:"name,omitempty"` // default to service
	Build Build  `yaml:"build"`
}

type Build struct {
	Context    string   `yaml:"context"`
	Dockerfile string   `yaml:"dockerfile"`
	Args       []string `yaml:"args,omitempty"` // KEY=VAL list
}

type Run struct {
	Command     []string `yaml:"command,omitempty"`
	Ports       []Port   `yaml:"ports,omitempty"`
	Volumes     []Volume `yaml:"volumes,omitempty"`
	Env         []string `yaml:"env,omitempty"` // KEY=VAL list (non-secret)
	EnvFromCove []string `yaml:"env_from_cove,omitempty"`
	Res         *Res     `yaml:"resources,omitempty"`
}

type Port struct {
	Name          string `yaml:"name"`                // e.g. http, metrics
	ContainerPort int    `yaml:"container_port"`      // 8080
	Protocol      string `yaml:"protocol,omitempty"`  // tcp|udp (default tcp)
	Public        bool   `yaml:"public"`              // whether to expose on host
	HostPort      int    `yaml:"host_port,omitempty"` // optional explicit host port
}

type Volume struct {
	Name      string `yaml:"name"`       // logical name
	MountPath string `yaml:"mount_path"` // /app/vault
	Type      string `yaml:"type"`       // persistent|bind|tmpfs (we’ll support persistent now)
	Size      string `yaml:"size,omitempty"`
}

type Res struct {
	CPU    string `yaml:"cpu,omitempty"`    // "0.5"
	Memory string `yaml:"memory,omitempty"` // "512Mi"
}

type Deploy struct {
	Strategy string            `yaml:"strategy"` // blue_green|recreate (we’ll do recreate now)
	Restart  string            `yaml:"restart"`
	Labels   map[string]string `yaml:"labels,omitempty"`
}
