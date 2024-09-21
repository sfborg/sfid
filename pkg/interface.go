package sfid

import (
	"github.com/sfborg/sfid/ent"
)

// SFID interface provides functions to
type SFID interface {
	// Process creates hashes from the given input. Input can be of
	// three types, a simple string, a file, or a directory.
	Process(inp string, chOut chan<- *ent.Output) error
}
