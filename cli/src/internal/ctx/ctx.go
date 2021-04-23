package ctx

// Context represents an execution context for factotum (env vars and volumes)
type Context struct {
	Env     map[string]string
	Volumes map[string]string
}

// Clone returns a deep-copy of this context
func (c Context) Clone() Context {
	var context Context

	// Env
	context.Env = make(map[string]string)
	for key, value := range c.Env {
		context.Env[key] = value
	}

	// Volumes
	context.Volumes = make(map[string]string)
	for key, value := range c.Volumes {
		context.Volumes[key] = value
	}

	return context
}

// Merge creates a deep-copy of this context and copies values from given source context on top of it
func (c Context) Merge(source Context) Context {
	context := c.Clone()

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
