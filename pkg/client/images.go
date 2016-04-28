package client

import (
	"fmt"
	"github.com/layer-x/layerx-commons/lxhttpclient"
	"net/http"
	"github.com/layer-x/layerx-commons/lxerrors"
	"encoding/json"
	"strings"
	"github.com/emc-advanced-dev/unik/pkg/types"
)

type images struct {
	unikIP string
}

func (i *images) All() ([]*types.Image, error) {
	resp, body, err := lxhttpclient.Get(i.unikIP, "/images", nil)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, lxerrors.New(fmt.Sprintf("failed with status %v: %s", resp.StatusCode, string(body)), err)
	}
	var images []*types.Image
	if err := json.Unmarshal(body, &images); err != nil {
		return nil, lxerrors.New(fmt.Sprintf("response body %s did not unmarshal to type []*types.Image", string(body)), err)
	}
	return images, nil
}

func (i *images) Get(id string) (*types.Image, error) {
	resp, body, err := lxhttpclient.Get(i.unikIP, "/images/"+id, nil)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, lxerrors.New(fmt.Sprintf("failed with status %v: %s", resp.StatusCode, string(body)), err)
	}
	var image types.Image
	if err := json.Unmarshal(body, &image); err != nil {
		return nil, lxerrors.New(fmt.Sprintf("response body %s did not unmarshal to type *types.Image", string(body)), err)
	}
	return &image, nil
}

func (i *images) Build(name, sourceTar, compiler, provider, args string, mounts []string, force bool) (*types.Image, error) {
	query := fmt.Sprintf("?compiler=%s&provider=%s&args=%s&mounts=%s&force=%v", compiler, provider, args, strings.Join(mounts, ","), force)
	resp, body, err := lxhttpclient.PostFile(i.unikIP, "/images/"+name+query, "tarfile", sourceTar)
	if err != nil || resp.StatusCode != http.StatusCreated {
		return nil, lxerrors.New(fmt.Sprintf("failed with status %v: %s", resp.StatusCode, string(body)), err)
	}
	var image types.Image
	if err := json.Unmarshal(body, &image); err != nil {
		return nil, lxerrors.New(fmt.Sprintf("response body %s did not unmarshal to type *types.Image", string(body)), err)
	}
	return &image, nil
}

func (i *images) Delete(id string, force bool) error {
	query := fmt.Sprintf("?force=%v", force)
	resp, body, err := lxhttpclient.Delete(i.unikIP, "/images/"+id+query, nil)
	if err != nil || resp.StatusCode != http.StatusNoContent {
		return lxerrors.New(fmt.Sprintf("failed with status %v: %s", resp.StatusCode, string(body)), err)
	}
	return nil
}