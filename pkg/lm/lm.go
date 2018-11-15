package lm

import (
	"fmt"
	"github.com/howardchn/argus-cli/pkg/conf"
	lmv1 "github.com/logicmonitor/lm-sdk-go"
	"net/url"
)

func newLMApi(conf *conf.LMConf) *lmv1.DefaultApi {
	config := lmv1.NewConfiguration()
	config.APIKey = map[string]map[string]string{
		"Authorization": {
			"AccessID":  conf.AccessId,
			"AccessKey": conf.AccessKey,
		},
	}
	config.BasePath = "https://" + conf.Account + ".logicmonitor.com/lm/rest"

	api := lmv1.NewDefaultApi()
	api.Configuration = config

	return api
}

type Client struct {
	option    *conf.LMConf
	apiClient *lmv1.DefaultApi
}

func NewClient(conf *conf.LMConf) *Client {
	return &Client{
		conf,
		newLMApi(conf),
	}
}

func (client *Client) Uninstall() error {
	var err error
	err = client.deleteDeviceGroup()
	if err != nil {
		return err
	}

	err = client.deleteCollectorGroup()
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) deleteDeviceGroup() error {
	restDeviceGroup, err := client.findDeviceGroup()
	if err != nil {
		return err
	} else if &restDeviceGroup != nil {
		_, _, deletionErr := client.apiClient.DeleteDeviceGroupById(restDeviceGroup.Id, true)
		return deletionErr
	} else {
		return nil
	}
}

func (client *Client) deleteCollectorGroup() error {
	collectorGroup, err := client.findCollectorGroup()
	if err != nil {
		return err
	}

	collectorIds, err := client.getCollectorIds(collectorGroup)
	if err != nil {
		return err
	}

	for _, id := range collectorIds {
		_, _, err := client.apiClient.DeleteCollectorById(id)
		if err != nil {
			fmt.Printf("delete collector <%d> failed, msg=%v", id, err)
		}
	}

	return nil
}

func (client *Client) getCollectorIds(collectorGroup *lmv1.RestCollectorGroup) ([]int32, error) {
	filter := fmt.Sprintf("collectorGroupId:%v", &collectorGroup.Id)
	restRes, _, err := client.apiClient.GetCollectorList("", -1, 0, filter)
	if err != nil {
		return nil, fmt.Errorf("get collector ids from group <%v>, group id <%d> failed", &collectorGroup.Name, &collectorGroup.Id)
	}

	var collectorIds []int32
	for _, item := range restRes.Data.Items {
		collectorIds = append(collectorIds, item.Id)
	}

	return collectorIds, nil
}

func getGroupName(cluster string) string {
	groupName := fmt.Sprintf("Kubernetes Cluster: %s", cluster)
	groupName = url.QueryEscape(groupName)
	return groupName
}

func (client *Client) findDeviceGroup() (*lmv1.RestDeviceGroup, error) {
	api := client.apiClient
	groupName := getGroupName(client.option.Cluster)

	restResp, _, err := api.GetDeviceGroupList("name,id,parentId", -1, 0, fmt.Sprintf("name:%s", groupName))
	if err != nil {
		return nil, fmt.Errorf("get device group <%s> failed. msg: %v", client.option.Cluster, err)
	}

	var deviceGroup *lmv1.RestDeviceGroup
	for _, item := range restResp.Data.Items {
		if item.ParentId == client.option.ParentId {
			deviceGroup = &item
			break
		}
	}

	return deviceGroup, nil
}

func (client *Client) findCollectorGroup() (*lmv1.RestCollectorGroup, error) {
	filter := fmt.Sprintf("name:%s", client.option.Cluster)
	restResp, _, err := client.apiClient.GetCollectorGroupList("", -1, 0, filter)
	if err != nil || len(restResp.Data.Items) == 0 {
		return nil, fmt.Errorf("get collector group <%s> failed", client.option.Cluster)
	}

	collectorGroup := &restResp.Data.Items[0]
	return collectorGroup, nil
}
