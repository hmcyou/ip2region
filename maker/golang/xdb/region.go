package xdb

import (
	"fmt"
	"strings"
	"sync"
)

// region manager with:
// 1, content cache.
// 2, util functions

// global cache map
var rcLock sync.Mutex
var regionCache = map[string]*Region{}

type Region struct {
	Str    string   // region string
	fields []string // region fields
}

var EmptyRegion = CacheRegion("")

// Create or get the region from the global cache.
// And it is a thread-safe implementation.
func CacheRegion(str string) *Region {
	// check the cache and return it directly
	// if there is a cache available
	rcLock.Lock()
	defer rcLock.Unlock()

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

// Create a new region without checking cache info
func NewRegion(str string) *Region {
	return &Region{
		Str:    str,
		fields: nil,
	}
}

func (r *Region) Fields() []string {
	if r.fields == nil {
		r.fields = strings.Split(r.Str, "|")
	}

	return r.fields
}

func (r *Region) Join(sep string) string {
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

	new := CacheRegion(strings.Join(sb, "|"))
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
