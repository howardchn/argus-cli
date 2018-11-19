package rest

import (
	"fmt"
	"github.com/howardchn/argus-cli/pkg/conf"
	lmv1 "github.com/logicmonitor/lm-sdk-go"
	"log"
	"net/url"
)

type Client struct {
	option    *conf.LMConf
	apiClient *lmv1.DefaultApi
}

func newLMApi(conf *conf.LMConf) *lmv1.DefaultApi {
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

func NewClient(conf *conf.LMConf) *Client {
	return &Client{
		conf,
		newLMApi(conf),
	}
}

func (client *Client) Clean() error {
	var err error
	log.Println("deleting devices and groups")
	err = client.deleteDeviceGroup()
	if err != nil {
		log.Println("delete device group failed")
		return err
	} else {
		log.Println("deleted devices and groups")
	}

	log.Println("deleting collectors and groups")
	err = client.deleteCollectorGroup()
	if err != nil {
		log.Println("delete collector group failed")
		return err
	} else {
		log.Println("deleted collectors and groups")
	}

	return nil
}

func (client *Client) deleteDeviceGroup() error {
	restDeviceGroup, err := client.findDeviceGroup()
	if err != nil {
		return err
	} else if restDeviceGroup != nil {
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
	} else if collectorGroup == nil {
		return nil
	}

	collectorIds, err := client.getCollectorIds(collectorGroup)
	if err != nil {
		return err
	}

	allCollectorDeleted := true
	for _, id := range collectorIds {
		err := client.deleteCollectorById(id)
		if err != nil {
			log.Printf("delete collector <%d> failed, msg=%v\n", id, err)
			allCollectorDeleted = false
		}
	}

	if allCollectorDeleted {
		_, _, err1 := client.apiClient.DeleteCollectorGroupById(collectorGroup.Id)
		return err1
	}

	return nil
}

func (client *Client) deleteCollectorById(id int32) error {
	filter := fmt.Sprintf("currentCollectorId:%d", id)
	restResponse, _, err := client.apiClient.GetDeviceList("id", -1, 0, filter)
	if err != nil {
		log.Printf("find device by collector <%d> failed, err <%v>\n", id, err)
		return err
	}

	deviceIds := getDeviceIds(&restResponse.Data)
	deleteDeviceErr := client.deleteDevicesByIds(deviceIds)
	if deleteDeviceErr != nil {
		log.Println("devices deletion failed, cannot continue to delete its collector", deleteDeviceErr)
		return deleteDeviceErr
	}

	collectorResponse, _, err1 := client.apiClient.DeleteCollectorById(id)
	if err1 != nil {
		log.Printf("delete collector <%d> failed, err <%v>\n", id, err1)
	} else if collectorResponse.Errmsg != "OK" {
		errMsg := fmt.Sprintf("delete collector <%d> failed, err <%v>\n", id, collectorResponse.Errmsg)
		err1 = fmt.Errorf(errMsg)
		log.Printf(errMsg)
	}

	return err1
}

func getDeviceIds(devices *lmv1.RestDevicePagination) []int32 {
	var ids []int32
	for _, d := range devices.Items {
		ids = append(ids, d.Id)
	}

	return ids
}

func (client *Client) deleteDevicesByIds(deviceIds []int32) error {
	if len(deviceIds) == 0 {
		log.Println("no devices to delete")
		return nil
	}

	var errDeviceIds []string
	for _, id := range deviceIds {
		_, _, err := client.apiClient.DeleteDevice(id)
		if err != nil {
			errDeviceIds = append(errDeviceIds, fmt.Sprintf("%d, %v", id, err))
		}
	}

	if len(errDeviceIds) > 0 {
		return fmt.Errorf("delete devices failed, %v", errDeviceIds)
	} else {
		return nil
	}
}

func (client *Client) getCollectorIds(collectorGroup *lmv1.RestCollectorGroup) ([]int32, error) {
	filter := fmt.Sprintf("collectorGroupId:%v", collectorGroup.Id)
	restRes, _, err := client.apiClient.GetCollectorList("", -1, 0, filter)
	if err != nil {
		return nil, fmt.Errorf("get collector ids from group <%v>, group id <%d> failed", collectorGroup.Name, collectorGroup.Id)
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
	filter := fmt.Sprintf("name:%s", groupName)

	restResp, _, err := api.GetDeviceGroupList("name,id,parentId", -1, 0, filter)
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
	collectorGroupName := client.option.Cluster
	filter := fmt.Sprintf("name:%s", collectorGroupName)
	restResp, _, err := client.apiClient.GetCollectorGroupList("", -1, 0, filter)
	if err != nil {
		return nil, fmt.Errorf("get collector group <%s> failed", collectorGroupName)
	}

	var collectorGroup *lmv1.RestCollectorGroup = nil
	if len(restResp.Data.Items) > 0 {
		collectorGroup = &restResp.Data.Items[0]
	} else {
		log.Printf("collector group <%s> not found\n", collectorGroupName)
	}

	return collectorGroup, nil
}
