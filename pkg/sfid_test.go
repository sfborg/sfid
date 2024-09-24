package sfid_test

import (
	"sync"
	"testing"

	"github.com/sfborg/sfid/ent"
	sfid "github.com/sfborg/sfid/pkg"
	"github.com/sfborg/sfid/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		msg  string
		inp  string
		uuid bool
		md5  bool
		res  string
	}{
		{"empty", "some string", false, false, "some string\t\t"},
		{"uuid", "some string", true, false,
			"some string\t5d59011b-e790-5232-8c5c-a13ff1fe88b3\t"},
		{"md5", "some string", false, true,
			"some string\t\t5ac749fbeec93607fc28d666be85e73a"},
		{"uuid+md5", "some string", true, true,
			"some string\t5d59011b-e790-5232-8c5c-a13ff1fe88b3\t5ac749fbeec93607fc28d666be85e73a"},
	}

	for _, v := range tests {
		chOut := make(chan *ent.Output)
		var wg sync.WaitGroup
		wg.Add(1)

		var res []*ent.Output
		cfg := config.New(config.OptWithUUID(v.uuid), config.OptWithMD5(v.md5))
		sf := sfid.New(cfg)
		go func() {
			defer wg.Done()
			for o := range chOut {
				res = append(res, o)
			}
		}()
		err := sf.Process(v.inp, chOut)
		assert.Nil(err)
		close(chOut)
		wg.Wait()
		assert.Equal(1, len(res))
		assert.Equal(v.res, res[0].String())
	}
}

func TestFile(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		msg  string
		inp  string
		uuid bool
		md5  bool
		res  string
	}{
		{"empty", "../testdata/work/meta.md", false, false,
			"../testdata/work/meta.md\t\t"},
		{"uuid", "../testdata/work/meta.md", true, false,
			"../testdata/work/meta.md\tc1893f3b-3056-5515-9660-b69f85eec3f0\t"},
		{"md5", "../testdata/work/meta.md", false, true,
			"../testdata/work/meta.md\t\td8fef66b660717da2352c98215005dca"},
		{"uuid+md5", "../testdata/work/meta.md", true, true,
			"../testdata/work/meta.md\tc1893f3b-3056-5515-9660-b69f85eec3f0\td8fef66b660717da2352c98215005dca"},
	}

	for _, v := range tests {
		chOut := make(chan *ent.Output)
		var wg sync.WaitGroup
		wg.Add(1)
		var res []*ent.Output
		cfg := config.New(config.OptWithUUID(v.uuid), config.OptWithMD5(v.md5))
		sf := sfid.New(cfg)
		go func() {
			defer wg.Done()
			for o := range chOut {
				res = append(res, o)
			}
		}()
		err := sf.Process(v.inp, chOut)
		assert.Nil(err)
		close(chOut)
		wg.Wait()
		assert.Equal(1, len(res))
		o := res[0].String()
		assert.Equal(v.res, o)
	}
}

func TestDir(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		msg  string
		inp  string
		uuid bool
		md5  bool
		res  string
	}{
		{"empty", "../testdata/work", false, false,
			"../testdata/work/friday/breaking-important-stuff-on-friday.jpg\t\t"},
		{"uuid", "../testdata/work", true, false,
			"../testdata/work/friday/breaking-important-stuff-on-friday.jpg\t2f9e3b8e-f9e4-5030-95b1-d5ec3067351d\t"},
		{"md5", "../testdata/work", false, true,
			"../testdata/work/friday/breaking-important-stuff-on-friday.jpg\t\te1d52e661d432badf3ebc2a9cb81a3bc"},
		{"uuid+md5", "../testdata/work", true, true,
			"../testdata/work/friday/breaking-important-stuff-on-friday.jpg\t2f9e3b8e-f9e4-5030-95b1-d5ec3067351d\te1d52e661d432badf3ebc2a9cb81a3bc"},
	}

	for _, v := range tests {
		chOut := make(chan *ent.Output)
		var wg sync.WaitGroup
		wg.Add(1)

		var res []*ent.Output
		cfg := config.New(
			config.OptWithUUID(v.uuid),
			config.OptWithMD5(v.md5),
			config.OptRecursive(true),
			config.OptJobsNum(1),
		)
		sf := sfid.New(cfg)
		go func() {
			defer wg.Done()
			for o := range chOut {
				res = append(res, o)
			}
		}()
		err := sf.Process(v.inp, chOut)
		assert.Nil(err)
		close(chOut)
		wg.Wait()
		assert.Equal(5, len(res))
		o := res[0].String()
		assert.Equal(v.res, o)
	}
}
