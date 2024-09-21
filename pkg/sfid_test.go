package sfid_test

import (
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
		sha  bool
		res  string
	}{
		{"empty", "some string", false, false, "STRING\tsome string\t\t\n"},
		{"uuid", "some string", true, false,
			"STRING\tsome string\t5d59011b-e790-5232-8c5c-a13ff1fe88b3\t\n"},
		{"sha", "some string", false, true,
			"STRING\tsome string\t\t61d034473102d7dac305902770471fd50f4c5b26f6831a56dd90b5184b3c30fc\n"},
		{"uuid+sha", "some string", true, true,
			"STRING\tsome string\t5d59011b-e790-5232-8c5c-a13ff1fe88b3\t61d034473102d7dac305902770471fd50f4c5b26f6831a56dd90b5184b3c30fc\n"},
	}

	for _, v := range tests {
		chOut := make(chan *ent.Output)
		var res []*ent.Output
		cfg := config.New(config.OptWithUUID(v.uuid), config.OptWithSha(v.sha))
		sf := sfid.New(cfg)
		go func() {
			for o := range chOut {
				res = append(res, o)
			}
		}()
		err := sf.Process(v.inp, chOut)
		assert.Nil(err)
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
		sha  bool
		res  string
	}{
		{"empty", "../testdata/work/meta.md", false, false,
			"FILE\t../testdata/work/meta.md\t\t\n"},
		{"uuid", "../testdata/work/meta.md", true, false,
			"FILE\t../testdata/work/meta.md\tc1893f3b-3056-5515-9660-b69f85eec3f0\t\n"},
		{"sha", "../testdata/work/meta.md", false, true,
			"FILE\t../testdata/work/meta.md\t\tf0bf3b10a2fc9d9361a830022115b3e780dc8850339d2de2900778513e33b003\n"},
		{"uuid+sha", "../testdata/work/meta.md", true, true,
			"FILE\t../testdata/work/meta.md\tc1893f3b-3056-5515-9660-b69f85eec3f0\tf0bf3b10a2fc9d9361a830022115b3e780dc8850339d2de2900778513e33b003\n"},
	}

	for _, v := range tests {
		chOut := make(chan *ent.Output)
		var res []*ent.Output
		cfg := config.New(config.OptWithUUID(v.uuid), config.OptWithSha(v.sha))
		sf := sfid.New(cfg)
		go func() {
			for o := range chOut {
				res = append(res, o)
			}
		}()
		err := sf.Process(v.inp, chOut)
		assert.Nil(err)
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
		sha  bool
		res  string
	}{
		{"empty", "../testdata/work", false, false,
			"FILE\t../testdata/work/friday/breaking-important-stuff-on-friday.jpg\t\t\n"},
		{"uuid", "../testdata/work", true, false,
			"FILE\t../testdata/work/friday/breaking-important-stuff-on-friday.jpg\t2f9e3b8e-f9e4-5030-95b1-d5ec3067351d\t\n"},
		{"sha", "../testdata/work", false, true,
			"FILE\t../testdata/work/friday/breaking-important-stuff-on-friday.jpg\t\t372f1db319519b11907cc649ba7aa35ccd15d1d7c1bf4d4a30b95316aebd0429\n"},
		{"uuid+sha", "../testdata/work", true, true,
			"FILE\t../testdata/work/friday/breaking-important-stuff-on-friday.jpg\t2f9e3b8e-f9e4-5030-95b1-d5ec3067351d\t372f1db319519b11907cc649ba7aa35ccd15d1d7c1bf4d4a30b95316aebd0429\n"},
	}

	for _, v := range tests {
		chOut := make(chan *ent.Output)
		var res []*ent.Output
		cfg := config.New(
			config.OptWithUUID(v.uuid),
			config.OptWithSha(v.sha),
			config.OptJobsNum(1),
		)
		sf := sfid.New(cfg)
		go func() {
			for o := range chOut {
				res = append(res, o)
			}
		}()
		err := sf.Process(v.inp, chOut)
		assert.Nil(err)
		assert.Equal(5, len(res))
		o := res[0].String()
		assert.Equal(v.res, o)
	}
}
