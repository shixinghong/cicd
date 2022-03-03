package builder

import (
	"fmt"
	containerregistryv1 "github.com/google/go-containerregistry/pkg/v1"
)

type ImageCommand struct {
	Command []string // image中的command
	Args    []string // image中的args
}

type Image struct {
	Name    string
	Digest  containerregistryv1.Hash
	Command map[string]*ImageCommand // 存储不同的架构的image中的command和args
}

func NewImage(name string, digest containerregistryv1.Hash) *Image {
	return &Image{
		Name:    name,
		Digest:  digest,
		Command: make(map[string]*ImageCommand), // 初始化
	}
}

func (i *Image) addCommand(os, arch string, cmds, args []string) {
	key := fmt.Sprintf("%s/%s", os, arch)
	i.Command[key] = &ImageCommand{
		Command: cmds,
		Args:    args,
	}
}
