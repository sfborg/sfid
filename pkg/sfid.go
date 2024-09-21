package sfid

import (
	"github.com/google/uuid"
	"github.com/sfborg/sfid/ent"
	"github.com/sfborg/sfid/pkg/config"
)

type sfid struct {
	// config
	cfg config.Config
	// name space uuid
	ns uuid.UUID
}

func New(cfg config.Config) SFID {
	res := sfid{
		cfg: cfg,
		ns:  uuid.NewSHA1(uuid.NameSpaceDNS, []byte(cfg.NameSpace)),
	}
	return &res
}

func (s *sfid) Process(inp string) ([]ent.Output, error) {
}

func (s *sfid) FromString(str string) ent.Output {
	res := ent.Output{}
	return res
}

func (s *sfid) FromFile(path string) (ent.Output, error) {
	res := ent.Output{}
	return res, nil
}

func (s *sfid) FromDir(path string) ([]ent.Output, error) {
	return nil, nil
}
