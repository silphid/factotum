package ctx

// RegistryType represents the type of docker registry factotum image should be retrieved from
type RegistryType string

const (
	RegistryGCR       RegistryType = "gcr"
	RegistryECR       RegistryType = "ecr"
	RegistryDockerHub RegistryType = "dockerhub"
)

// Context represents an execution context for factotum (env vars and volumes)
type Context struct {
	Registry RegistryType
	Image    string
	Env      map[string]string
	Volumes  map[string]string
}

// Clone returns a deep-copy of this context
func (c Context) Clone() Context {
	context := Context{
		Registry: c.Registry,
		Image:    c.Image,
		Env:      make(map[string]string),
		Volumes:  make(map[string]string),
	}

	for key, value := range c.Env {
		context.Env[key] = value
	}

	// Volumes
	for key, value := range c.Volumes {
		context.Volumes[key] = value
	}

	return context
}

// Merge creates a deep-copy of this context and copies values from given source context on top of it
func (c Context) Merge(source Context) Context {
	context := c.Clone()

	// Registry
	if source.Registry != "" {
		context.Registry = source.Registry
	}

	// Image
	if source.Image != "" {
		context.Image = source.Image
	}

	// Env
	for key, value := range source.Env {
		context.Env[key] = value
	}

	// Volumes
	for key, value := range source.Volumes {
		context.Volumes[key] = value
	}

	return context
}
