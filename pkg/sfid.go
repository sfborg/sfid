package sfid

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"github.com/gnames/gnsys"
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

func (s *sfid) Process(inp string, chOut chan<- *ent.Output) error {
	var err error

	dirSt := gnsys.GetDirState(inp)
	if dirSt == gnsys.DirEmpty {
		slog.Warn("Empty directory", "dir", inp)
		return fmt.Errorf("empty directory: '%s'", inp)
	}
	if s.cfg.WithUUID {
		slog.Info("Generating UUID v5", "namespace", s.cfg.NameSpace)
	}
	if s.cfg.WithMD5 {
		slog.Info("Generating MD5 hash")
	}
	if dirSt == gnsys.DirNotEmpty {
		slog.Info("Traversing directory",
			"dir", inp, "recursive", s.cfg.Recursive, "jobs", s.cfg.JobsNum,
		)
		err = s.fromDir(inp, chOut)
		if err != nil {
			return err
		}
		return nil
	}

	isFile, _ := gnsys.FileExists(inp)
	if isFile {
		slog.Info("Processing file")
		out, err := s.fromFile(inp)
		if err != nil {
			return err
		}
		chOut <- out
		return nil
	}

	slog.Info("Processing a string")
	chOut <- s.fromString(inp)
	return nil
}

func (s *sfid) fromString(str string) *ent.Output {
	res := ent.Output{Type: "STRING", Input: str}

	if s.cfg.WithUUID {
		uuidObj := uuid.NewSHA1(s.ns, []byte(str))
		res.UUID = &uuidObj
	}

	if s.cfg.WithMD5 {
		h := md5.New()
		h.Write([]byte(str))
		res.MD5 = h.Sum(nil)
	}

	return &res
}

func (s *sfid) fromFile(path string) (*ent.Output, error) {
	res := ent.Output{Type: "FILE", Input: path}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if s.cfg.WithUUID {
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
		uuidObj := uuid.NewSHA1(s.ns, h.Sum(nil))
		res.UUID = &uuidObj
		if _, err := f.Seek(0, io.SeekStart); err != nil {
			return nil, err
		}
	}

	if s.cfg.WithMD5 {
		h := md5.New()
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
		res.MD5 = h.Sum(nil)
	}

	return &res, nil
}

func (s *sfid) fromDir(path string, chOut chan<- *ent.Output) error {

	var err error
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	g, ctx := errgroup.WithContext(ctx)
	chIn := make(chan string)
	var wg sync.WaitGroup
	wg.Add(s.cfg.JobsNum)

	g.Go(func() error {
		for range s.cfg.JobsNum {
			err := s.dirWorker(ctx, chIn, chOut, &wg)
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
	wg.Wait()

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

func (s *sfid) dirWorker(
	ctx context.Context,
	chIn <-chan string,
	chOut chan<- *ent.Output,
	wg *sync.WaitGroup,
) error {
	defer wg.Done()
	for path := range chIn {
		out, err := s.fromFile(path)
		if err != nil {
			return err
		}
		chOut <- out
		select {
		case <-ctx.Done():
			for range chIn {
			}
			return ctx.Err()
		default:
		}
	}
	return nil
}

func (s *sfid) loadFiles(
	ctx context.Context,
	path string,
	chIn chan<- string,
) error {
	if !s.cfg.Recursive {
		err := s.readDir(ctx, path, chIn)
		if err != nil {
			return err
		}
		return nil
	}

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

func (s *sfid) readDir(
	ctx context.Context,
	path string,
	chIn chan<- string,
) error {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if !file.Type().IsRegular() {
			continue
		}

		chIn <- filepath.Join(path, file.Name())

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}
	close(chIn)
	return nil
}
