package imagestore

import (
	"io"
)

// Builder is a means of creating image stores.
type Builder interface {
	// ArchiveDir creates a fresh Builder with the given archive directory.
	ArchiveDir(string) Builder

	// Logs creates a fresh builder with the given log stream.
	Logs(io.Writer) Builder

	// Build creates an image store.
	Build() (Store, error)

	// Create creates an image store with the given options.
	Create(...Option) (Store, error)
}

func Create(b Builder, options ...Option) (Store, error) {
	for _, op := range options {
		b = op(b)
	}
	return b.Build()
}

type Option func(Builder) Builder

func WithArchiveDir(archiveDir string) Option {
	return func(b Builder) Builder {
		return b.ArchiveDir(archiveDir)
	}
}

func WithLogs(logs io.Writer) Option {
	return func(b Builder) Builder {
		return b.Logs(logs)
	}
}
