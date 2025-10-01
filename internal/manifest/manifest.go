package manifest

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func load(path string) (*Manifest, error) {

	bites, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var manifest Manifest
	if err := yaml.Unmarshal(bites, &manifest); err != nil {
		return nil, err
	}

	if err := manifest.defaults(); err != nil {
		return nil, err
	}

	if err := manifest.validate(); err != nil {
		return nil, err
	}

	// Make build paths absolute relative to manifest file
	base := filepath.Dir(path)
	if !filepath.IsAbs(manifest.Image.Build.Context) {
		manifest.Image.Build.Context = filepath.Clean(filepath.Join(base, manifest.Image.Build.Context))
	}

	return &manifest, nil

}

func (m *Manifest) defaults() error {

	if m.Version == 0 {
		m.Version = 1
	}
	if m.Image.Name == "" {
		m.Image.Name = m.Service
	}

	// Insert default port protocols
	for i := range m.Run.Ports {
		if m.Run.Ports[i].Protocol == "" {
			m.Run.Ports[i].Protocol = "tcp"
		}
	}

	// Reasonable deploy defaults
	if m.Deploy.Strategy == "" {
		m.Deploy.Strategy = "recreate"
	}
	if m.Deploy.Restart == "" {
		m.Deploy.Restart = "unless-stopped"
	}
	if m.Deploy.Labels == nil {
		m.Deploy.Labels = map[string]string{}
	}
	m.Deploy.Labels["managed-by"] = "lighthouse"
	m.Deploy.Labels["service"] = m.Service
	return nil
}

func (m *Manifest) validate() error {
	if m.Service == "" {
		return errors.New("service is required")
	}
	if m.Image.Build.Context == "" {
		return errors.New("image.build.context is required")
	}
	if m.Image.Build.Dockerfile == "" {
		return errors.New("image.build.dockerfile is required")
	}
	seen := map[string]bool{}
	for _, p := range m.Run.Ports {
		if p.Name == "" {
			return errors.New("run.ports[].name is required")
		}
		if seen[p.Name] {
			return fmt.Errorf("duplicate port name %q", p.Name)
		}
		seen[p.Name] = true
		if p.ContainerPort <= 0 {
			return fmt.Errorf("run.ports[%s].container_port must be > 0", p.Name)
		}
		if p.Protocol != "tcp" && p.Protocol != "udp" {
			return fmt.Errorf("run.ports[%s].protocol must be tcp or udp", p.Name)
		}
	}
	for _, v := range m.Run.Volumes {
		switch v.Type {
		case "persistent", "bind", "tmpfs":
		default:
			return fmt.Errorf("run.volumes[%s].type must be persistent|bind|tmpfs", v.Name)
		}
		if v.MountPath == "" {
			return fmt.Errorf("run.volumes[%s].mount_path is required", v.Name)
		}
	}
	// KEY=VAL sanity for env/args
	for _, kv := range append(m.Run.Env, m.Image.Build.Args...) {
		if !strings.Contains(kv, "=") {
			return fmt.Errorf("invalid KEY=VAL: %q", kv)
		}
	}
	return nil
}
