package helm

import (
	"bytes"
	"fmt"
	"github.com/howardchn/argus-cli/pkg/conf"
	"os/exec"
)

type Client struct {
	conf *conf.LMConf
}

func NewClient(conf *conf.LMConf) *Client {
	return &Client{conf}
}

func (client *Client) Uninstall() error {
	var err error
	err = client.deleteArgus()
	if err != nil {
		return err
	}

	err = client.deleteCollectorSetController()
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) deleteArgus() error {
	return deleteRelease("argus")
}

func (client *Client) deleteCollectorSetController() error {
	return deleteRelease("collectorset-controller")
}

func deleteRelease(name string) error {
	cmd := exec.Command("helm", "delete", name, "--purge")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}

	fmt.Println(out.String())
	fmt.Printf("%s deleted\n", name)

	return nil
}
