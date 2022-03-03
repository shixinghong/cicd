package builder

import (
	"fmt"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

func parseImage(img string) (*Image, error) {

	ref, err := name.ParseReference(img)
	if err != nil {
		return nil, err
	}

	desc, err := remote.Get(ref)
	if err != nil {
		return nil, err
	}

	imgBuilder := NewImage(img, desc.Digest)

	if desc.MediaType.IsIndex() { // 如果是index模式 绝大多是都是index模式
		idx, err := desc.ImageIndex()
		if err != nil {
			return nil, err

		}
		mf, err := idx.IndexManifest()
		if err != nil {
			return nil, err

		}
		for _, v := range mf.Manifests {
			img, err := idx.Image(v.Digest)
			if err != nil {
				return nil, err
			}
			cf, err := img.ConfigFile()
			if err != nil {
				return nil, err
			}
			imgBuilder.addCommand(cf.OS, cf.Architecture, cf.Config.Entrypoint, cf.Config.Cmd)
		}
		return imgBuilder, nil
	}

	if desc.MediaType.IsImage() {
		image, err := desc.Image()
		if err != nil {
			return nil, err
		}
		cf, err := image.ConfigFile()
		if err != nil {
			return nil, err
		}
		imgBuilder.addCommand(cf.OS, cf.Architecture, cf.Config.Entrypoint, cf.Config.Cmd)
		return imgBuilder, nil
	}
	return nil, fmt.Errorf("error image url")
}
