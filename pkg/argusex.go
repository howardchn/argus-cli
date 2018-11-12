package argusex

import (
	"github.com/howardchn/argus-cli/pkg/conf"
	"github.com/howardchn/argus-cli/pkg/santaba"
)

type ArgusClient struct {
	LMClient *santaba.LMClient
}

func NewArgusClient(lmConf *conf.LMConf) *ArgusClient {
	return &ArgusClient{
		santaba.NewLMClient(lmConf),
	}
}
