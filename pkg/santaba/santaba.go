package santaba

import (
	"fmt"
	"github.com/howardchn/argus-cli/pkg/conf"
	lmv1 "github.com/logicmonitor/lm-sdk-go"
	"net/url"
)

func newApiClient(conf *conf.LMConf) *lmv1.DefaultApi {
	config := lmv1.NewConfiguration()
	config.APIKey = map[string]map[string]string{
		"Authorization": {
			"AccessID":  conf.AccessId,
			"AccessKey": conf.AccessKey,
		},
	}
	config.BasePath = "https://" + conf.Account + ".logicmonitor.com/santaba/rest"

	api := lmv1.NewDefaultApi()
	api.Configuration = config

	return api
}

type LMClient struct {
	option    *conf.LMConf
	apiClient *lmv1.DefaultApi
}

func NewLMClient(lmConf *conf.LMConf) *LMClient {
	return &LMClient{
		lmConf,
		newApiClient(lmConf),
	}
}

func (lmClient *LMClient) DeleteDeviceGroup() error {
	restDeviceGroup, err := lmClient.findDeviceGroup()
	if err != nil {
		return err
	} else if &restDeviceGroup != nil {
		_, _, deletionErr := lmClient.apiClient.DeleteDeviceGroupById(restDeviceGroup.Id, true)
		return deletionErr
	} else {
		return nil
	}
}

func (lmClient *LMClient) DeleteCollectorGroup() error {
	collectorGroup, err := lmClient.findCollectorGroup()
	if err != nil {
		return err
	}

	collectorIds, err := lmClient.getCollectorIds(collectorGroup)
	if err != nil {
		return err
	}

	for _, id := range collectorIds {
		_, _, err := lmClient.apiClient.DeleteCollectorById(id)
		if err != nil {
			fmt.Printf("delete collector <%d> failed, msg=%v", id, err)
		}
	}

	return nil
}

func (lmClient *LMClient) getCollectorIds(collectorGroup *lmv1.RestCollectorGroup) ([]int32, error) {
	filter := fmt.Sprintf("collectorGroupId:%v", &collectorGroup.Id)
	restRes, _, err := lmClient.apiClient.GetCollectorList("", -1, 0, filter)
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

func (lmClient *LMClient) findDeviceGroup() (*lmv1.RestDeviceGroup, error) {
	api := lmClient.apiClient
	groupName := getGroupName(lmClient.option.Cluster)

	restResp, _, err := api.GetDeviceGroupList("name,id,parentId", -1, 0, fmt.Sprintf("name:%s", groupName))
	if err != nil {
		return nil, fmt.Errorf("get device group <%s> failed. msg: %v", lmClient.option.Cluster, err)
	}

	var deviceGroup *lmv1.RestDeviceGroup
	for _, item := range restResp.Data.Items {
		if item.ParentId == lmClient.option.ParentId {
			deviceGroup = &item
			break
		}
	}

	return deviceGroup, nil
}

func (lmClient *LMClient) findCollectorGroup() (*lmv1.RestCollectorGroup, error) {
	filter := fmt.Sprintf("name:%s", lmClient.option.Cluster)
	restResp, _, err := lmClient.apiClient.GetCollectorGroupList("", -1, 0, filter)
	if err != nil || len(restResp.Data.Items) == 0 {
		return nil, fmt.Errorf("get collector group <%s> failed", lmClient.option.Cluster)
	}

	collectorGroup := &restResp.Data.Items[0]
	return collectorGroup, nil

}
