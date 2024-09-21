package ent

import (
	"encoding/hex"
	"strings"

	"github.com/google/uuid"
)

type Output struct {
	Type  string
	Input string
	Sha   []byte
	UUID  *uuid.UUID
}

func (o Output) String() string {
	var u string
	if o.UUID != nil {
		u = o.UUID.String()
	}
	var sha string
	if len(o.Sha) > 0 {
		sha = hex.EncodeToString(o.Sha)
	}
	inp := o.Input
	if o.Type == "STRING" && len(inp) > 100 {
		inp = o.Input[0:100] + "..."
	}
	res := strings.Join([]string{o.Type, inp, u, sha}, "\t")
	return res
}
