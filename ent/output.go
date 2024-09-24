package ent

import (
	"encoding/hex"
	"strings"

	"github.com/google/uuid"
)

type Output struct {
	Type  string
	Input string
	MD5   []byte
	UUID  *uuid.UUID
}

func (o Output) String() string {
	var u string
	if o.UUID != nil {
		u = o.UUID.String()
	}
	var md5 string
	if len(o.MD5) > 0 {
		md5 = hex.EncodeToString(o.MD5)
	}
	inp := o.Input
	if o.Type == "STRING" && len(inp) > 100 {
		inp = o.Input[0:100] + "..."
	}
	res := strings.Join([]string{inp, md5, u}, "\t")
	return res
}
