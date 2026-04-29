package xdb

import (
	"fmt"
	"strings"
)

// region manager with:
// 1, content cache.
// 2, util functions

type Region struct {
	Str    string   // region string
	fields []string // region fields
}

// global cache map
var regionCache = map[string]*Region{}

func RNew(str string) *Region {
	return NewRegion(str)
}

func REmpty() *Region {
	return NewRegion("")
}

func NewRegion(str string) *Region {
	// check the cache and return it directly
	// if there is a cache available
	region, ok := regionCache[str]
	if ok {
		return region
	}

	// cache the new region
	region = &Region{
		Str:    str,
		fields: nil,
	}

	regionCache[str] = region
	return region
}

func (r *Region) Fields() []string {
	if r.fields == nil {
		r.fields = strings.Split(r.Str, "|")
	}

	return r.fields
}

func (r *Region) JoinBy(sep string) string {
	if sep == "|" {
		return r.Str
	}

	return strings.Join(r.Fields(), sep)
}

func (r *Region) Filtering(fields []int) (*Region, error) {
	if len(fields) == 0 {
		return r, nil
	}

	fs := r.Fields()
	var sb []string
	for _, idx := range fields {
		if idx < 0 {
			return r, fmt.Errorf("negative filter index %d", idx)
		}

		if idx >= len(fs) {
			return r, fmt.Errorf("field index %d exceeded the max length of %d", idx, len(fs))
		}

		sb = append(sb, fs[idx])
	}

	new := RNew(strings.Join(sb, "|"))
	if new.fields == nil {
		new.fields = sb
	}

	return new, nil
}

// Equal check ptr (share the same region cache) or the Str is the same.
func (r *Region) Equal(dst *Region) bool {
	return (r == dst || r.Str == dst.Str)
}

func (r *Region) IsEmpty() bool {
	return r.Str == ""
}

func (r *Region) String() string {
	return r.Str
}
