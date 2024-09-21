package ent

import "github.com/google/uuid"

type Output struct {
	Input string
	Sha   string
	UUID  uuid.UUID
}
