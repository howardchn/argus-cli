package uninstaller

import (
	"github.com/howardchn/argus-cli/pkg/conf"
	"github.com/howardchn/argus-cli/pkg/helm"
	"github.com/howardchn/argus-cli/pkg/lm"
	"log"
)

type Client struct {
	LMClient   *lm.Client
	HelmClient *helm.Client
}

func NewClient(conf *conf.LMConf) *Client {
	return &Client{
		lm.NewClient(conf),
		helm.NewClient(conf),
	}
}

func (client *Client) Uninstall() error {
	var err error
	err = client.LMClient.Uninstall()
	if err != nil {
		log.Println("lm uninstall failed")
		return err
	}

	err = client.HelmClient.Uninstall()
	if err != nil {
		log.Println("helm uninstall failed")
		return err
	}

	return nil
}
