package betypes

import (
	"github.com/bradfitz/gomemcache/memcache"
)

var (
	MemCashed = memcache.New("127.0.0.1:11211")
)