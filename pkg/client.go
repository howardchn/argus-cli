package uninstaller

import (
	"github.com/howardchn/argus-cli/pkg/conf"
	"github.com/howardchn/argus-cli/pkg/helm"
	"github.com/howardchn/argus-cli/pkg/rest"
	"log"
)

type Client struct {
	RestClient *rest.Client
	HelmClient *helm.Client
}

func NewClient(conf *conf.LMConf) *Client {
	return &Client{
		rest.NewClient(conf),
		helm.NewClient(conf),
	}
}

func (client *Client) Clean() error {
	var err error
	err = client.RestClient.Clean()
	if err != nil {
		log.Panicln("-- LM uninstall failed --", err)
		return err
	} else {
		log.Println("-- LM uninstall success --")
	}

	log.Println()
	err = client.HelmClient.Clean()
	if err != nil {
		log.Panicln("-- helm uninstall failed --", err)
		return err
	} else {
		log.Panicln("-- helm uninstall success --")
	}

	return nil
}
