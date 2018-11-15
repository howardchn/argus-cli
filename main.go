package main

import (
	"flag"
	"fmt"
	"github.com/howardchn/argus-cli/pkg"
	"github.com/howardchn/argus-cli/pkg/conf"
)

var (
	accessId    string
	accessKey   string
	clusterName string
	account     string
	parentId    int
)

func init() {
	flag.StringVar(&accessId, "accessId", "", "API Access ID")
	flag.StringVar(&accessKey, "accessKey", "", "API Access Key")
	flag.StringVar(&clusterName, "clusterName", "", "Cluster Name")
	flag.StringVar(&account, "account", "", "Account Name")
	flag.IntVar(&parentId, "parentId", 1, "ParentId")
}

func main() {
	flag.Parse()
	conf := &conf.LMConf{AccessId: accessId, AccessKey: accessKey, Account: account, Cluster: clusterName, ParentId: int32(parentId)}
	client := uninstaller.NewClient(conf)
	err := client.Uninstall()
	if err != nil {
		fmt.Printf("uninstall failed. err = %v\n", err)
		return
	}

	fmt.Println("uninstall success")
}
