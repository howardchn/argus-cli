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
		log.Panicln("-- LM uninstall failed --", err)
		return err
	} else {
		log.Println("-- LM uninstall success --")
	}

	log.Println()
	err = client.HelmClient.Uninstall()
	if err != nil {
		log.Panicln("-- helm uninstall failed --", err)
		return err
	} else {
		log.Panicln("-- helm uninstall success --")
	}

	return nil
}
