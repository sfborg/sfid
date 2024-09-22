package sfid

import (
	"context"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/sfborg/sfid/ent"
	"github.com/sfborg/sfid/pkg/config"
	"golang.org/x/sync/errgroup"
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
	return s.fromDir(inp)
}

func (s *sfid) fromString(str string) *ent.Output {
	res := ent.Output{Input: str}

	if s.cfg.OutputUUID {
		res.UUID = uuid.NewSHA1(s.ns, []byte(str))
	}

	if s.cfg.OutputSha256 {
		h := sha256.New()
		h.Write([]byte(str))
		res.Sha = string(h.Sum(nil))
	}

	return &res
}

func (s *sfid) fromFile(path string) (*ent.Output, error) {
	res := ent.Output{Input: path}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if s.cfg.OutputUUID {
		h := sha1.New()
		const bufferSize = 65536
		buffer := make([]byte, bufferSize)
		for {
			bsNum, err := f.Read(buffer)
			if err != nil {
				if err != io.EOF {
					return nil, err
				}
				break
			}
			h.Write(buffer[:bsNum])
		}
		res.UUID = uuid.NewSHA1(s.ns, h.Sum(nil))
		res.Sha = string(h.Sum(nil))
	}

	if s.cfg.OutputSha256 {
		h := sha256.New()
		const bufferSize = 65536
		buffer := make([]byte, bufferSize)
		for {
			bsNum, err := f.Read(buffer)
			if err != nil {
				if err != io.EOF {
					return nil, err
				}
				break
			}
			h.Write(buffer[:bsNum])
		}
		res.Sha = string(h.Sum(nil))
	}

	return &res, nil
}

func (s *sfid) fromDir(path string) ([]ent.Output, error) {
	var err error
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	g, ctx := errgroup.WithContext(ctx)
	chIn := make(chan string)
	chOut := make(chan *ent.Output)

	g.Go(func() error {
		for o := range chOut {
			fmt.Printf("%s\t%s\t%s", o.Input, o.UUID.String(), o.Sha)
		}
		return nil
	})

	g.Go(func() error {
		defer close(chOut)
		for range s.cfg.JobsNum {
			err := s.dirWorker(ctx, chIn, chOut)
			if err != nil {
				slog.Error("dirWorker", "error", err)
				return err
			}
		}
		return nil
	})

	g.Go(func() error {
		err = s.loadFiles(ctx, path, chIn)
		if err != nil {
			slog.Error("Get ", "error", err)
		}
		return err
	})
	return nil, nil
}

func (s *sfid) dirWorker(
	ctx context.Context,
	chIn <-chan string,
	chOut chan<- *ent.Output,
) error {
	for path := range chIn {
		out, err := s.fromFile(path)
		if err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			for range chIn {
			}
			return ctx.Err()
		default:
			chOut <- out
		}
	}
	return nil
}

func (s *sfid) loadFiles(
	ctx context.Context,
	path string,
	chIn chan<- string,
) error {
	root := path
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			slog.Error("Problem with directory processing", "error", err)
			return nil
		}

		if !info.IsDir() {
			chIn <- path
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error walking the path:", err)
		return err
	}
	close(chIn)
	return nil
}
