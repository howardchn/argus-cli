package cmd

import (
	"fmt"
	"github.com/howardchn/argus-cli/pkg"
	"github.com/howardchn/argus-cli/pkg/conf"
	"github.com/spf13/cobra"
)

var (
	cluster  string
	account  string
	parentId int32
	mode     string
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "uninstall argus related resources",
	Long:  "uninstall argus related resources in santaba and k8s",
	Run: func(cmd *cobra.Command, args []string) {
		conf := &conf.LMConf{AccessId: accessId, AccessKey: accessKey, Account: account, Cluster: cluster, ParentId: parentId}
		client := uninstaller.NewClient(conf)
		err := client.Clean(mode)
		if err != nil {
			fmt.Printf("uninstall failed. err = %v\n", err)
			return
		}

		fmt.Println("uninstall success")
	},
}

func init() {
	uninstallCmd.Flags().StringVarP(&cluster, "cluster", "c", "", "cluster name")
	uninstallCmd.Flags().StringVarP(&account, "account", "a", "", "account name")
	uninstallCmd.Flags().Int32VarP(&parentId, "parentId", "g", 1, "parent group id, default: 1")
	uninstallCmd.Flags().StringVarP(&mode, "mode", "m", "all", "uninstall mode: [rest|helm|all], default: all")
	uninstallCmd.MarkFlagRequired("cluster")
	uninstallCmd.MarkFlagRequired("account")
	uninstallCmd.MarkFlagRequired("parentId")
	RootCmd.AddCommand(uninstallCmd)
}
