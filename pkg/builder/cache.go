package builder

import (
	lru "github.com/hashicorp/golang-lru"
	"k8s.io/klog/v2"
)

var ImageCache *lru.Cache

func InitCache(size int) {
	c, err := lru.New(size)
	if err != nil {
		klog.Exit(err)
	}
	ImageCache = c
}
