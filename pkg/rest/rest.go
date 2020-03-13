package rest

import (
	"fmt"
	"github.com/howardchn/argus-cli/pkg/conf"
	"github.com/logicmonitor/lm-sdk-go/client"
	"github.com/logicmonitor/lm-sdk-go/client/lm"
	"github.com/logicmonitor/lm-sdk-go/models"
	"log"
)

type Client struct {
	option    *conf.LMConf
	apiClient *client.LMSdkGo
}

func NewClient(conf *conf.LMConf) *Client {
	config := client.NewConfig()
	config.SetAccessID(&conf.AccessId)
	config.SetAccessKey(&conf.AccessKey)
	domain := conf.Account + ".logicmonitor.com"
	config.SetAccountDomain(&domain)
	return &Client{
		option:    conf,
		apiClient: client.New(config),
	}

}

func cleanTask(name string, action func() error) error {
	log.Println("deleting", name)
	err := action()
	if err != nil {
		log.Println(fmt.Sprintf("delete %s failed", name))
		return err
	} else {
		log.Println("deleted", name)
	}

	return nil
}

func (client *Client) Clean() error {
	err := cleanTask("devices", func() error { return client.deleteDeviceGroup() })
	if err != nil {
		return err
	}

	err = cleanTask("collectors", func() error { return client.deleteCollectorGroup() })
	if err != nil {
		return err
	}

	err = cleanTask("dashboards", func() error { return client.deleteDashboardGroup() })
	if err != nil {
		return err
	}

	err = cleanTask("services", func() error { return client.deleteServiceGroup() })
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) deleteDeviceGroup() error {
	restDeviceGroup, err := client.findDeviceGroup()
	if err != nil {
		return err
	} else if restDeviceGroup != nil {
		//_, _, deletionErr := client.apiClient.LM.DeleteDeviceGroupById(restDeviceGroup.Id, true)
		var params = &lm.DeleteDeviceGroupByIDParams{
			DeleteChildren: nil,
			DeleteHard:     nil,
			ID:             restDeviceGroup.ID,
			Context:        nil,
			HTTPClient:     nil,
		}
		_, deletionErr := client.apiClient.LM.DeleteDeviceGroupByID(params)
		return deletionErr
	} else {
		return nil
	}
}

func (client *Client) deleteServiceGroup() error {
	restServiceGroup, err := client.findServiceGroup()
	if err != nil {
		return err
	} else if restServiceGroup != nil {
		//_, _, deletionErr := client.apiClient.LM.DeleteDeviceGroupById(restServiceGroup.Id, true)
		var params = &lm.DeleteDeviceGroupByIDParams{
			DeleteChildren: nil,
			DeleteHard:     nil,
			ID:             restServiceGroup.ID,
			Context:        nil,
			HTTPClient:     nil,
		}
		_, deletionErr := client.apiClient.LM.DeleteDeviceGroupByID(params)
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

	var params = &lm.DeleteCollectorGroupByIDParams{
		ID:         collectorGroup.ID,
		Context:    nil,
		HTTPClient: nil,
	}
	if allCollectorDeleted {
		_, err1 := client.apiClient.LM.DeleteCollectorGroupByID(params)
		return err1
	}

	return nil
}

func (client *Client) deleteCollectorById(id int32) error {
	filter := fmt.Sprintf("currentCollectorId:%d", id)
	var fields = "id"
	var offset int32 = 0
	var size int32 = -1
	var params = &lm.GetDeviceListParams{
		End:           nil,
		Fields:        &fields,
		Filter:        &filter,
		NetflowFilter: nil,
		Offset:        &offset,
		Size:          &size,
		Start:         nil,
		Context:       nil,
		HTTPClient:    nil,
	}
	restResponse, err := client.apiClient.LM.GetDeviceList(params)
	if err != nil {
		log.Printf("find device by collector <%d> failed, err <%v>\n", id, err)
		return err
	}

	deviceIds := getDeviceIds(restResponse.Payload)
	deleteDeviceErr := client.deleteDevicesByIds(deviceIds)
	if deleteDeviceErr != nil {
		log.Println("devices deletion failed, cannot continue to delete its collector", deleteDeviceErr)
		return deleteDeviceErr
	}

	var dcParams = &lm.DeleteCollectorByIDParams{
		ID:         id,
		Context:    nil,
		HTTPClient: nil,
	}
	collectorResponse, err1 := client.apiClient.LM.DeleteCollectorByID(dcParams)
	if err1 != nil {
		log.Printf("delete collector <%d> failed, err <%v>\n", id, err1)
	} else if collectorResponse.Error() != "OK" {
		errMsg := fmt.Sprintf("delete collector <%d> failed, err <%v>\n", id, collectorResponse.Error())
		err1 = fmt.Errorf(errMsg)
		log.Printf(errMsg)
	}

	return err1
}

func (client *Client) deleteDashboardGroup() error {
	dashboardGroupName := dashboardGroupName(client.option.Cluster)
	filter := fmt.Sprintf("name:%s", dashboardGroupName)
	var fields = "id,name"
	var offset int32 = 0
	var size int32 = -1
	var params = &lm.GetDashboardGroupListParams{
		Fields:     &fields,
		Filter:     &filter,
		Offset:     &offset,
		Size:       &size,
		Context:    nil,
		HTTPClient: nil,
	}
	dashboardGroups, err := client.apiClient.LM.GetDashboardGroupList(params)
	if err != nil {
		log.Printf("dashboard group <%s> found failed\n", dashboardGroupName)
		return err
	}

	for _, d := range dashboardGroups.Payload.Items {
		err := client.deleteDashboardGroupById(d.ID)
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}

func (client *Client) deleteDashboardGroupById(gid int32) error {

	filter := fmt.Sprintf("groupId:%d", gid)
	var fields = "id,name"
	var offset int32 = 0
	var size int32 = -1
	var params = &lm.GetDashboardListParams{
		Fields:     &fields,
		Filter:     &filter,
		Offset:     &offset,
		Size:       &size,
		Context:    nil,
		HTTPClient: nil,
	}
	r, err := client.apiClient.LM.GetDashboardList(params)
	if err != nil {
		log.Printf("get dashboards from group<%d> failed\n", gid)
		return err
	}

	for _, d := range r.Payload.Items {
		r, err := client.apiClient.LM.DeleteDashboardByID(&lm.DeleteDashboardByIDParams{ID: d.ID})
		if err != nil {
			return err
		} else if r.Error() != "OK" {
			return fmt.Errorf("delete dashboard<%d> failed", d.ID)
		}
	}

	deleteGroupResponse, err := client.apiClient.LM.DeleteDashboardGroupByID(&lm.DeleteDashboardGroupByIDParams{ID: gid})
	if err != nil {
		return err
	} else if deleteGroupResponse.Error() != "OK" {
		return fmt.Errorf("delete dashboard group failed, %v", deleteGroupResponse.Error())
	}

	return nil
}

func getDeviceIds(devices *models.DevicePaginationResponse) []int32 {
	var ids []int32
	for _, d := range devices.Items {
		ids = append(ids, d.ID)
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
		_, err := client.apiClient.LM.DeleteDeviceByID(&lm.DeleteDeviceByIDParams{ID:id})
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

func (client *Client) getCollectorIds(collectorGroup *models.CollectorGroup) ([]int32, error) {
	filter := fmt.Sprintf("collectorGroupId:%v", collectorGroup.ID)
	fields := ""
	var offset int32 = 0
	var size int32 = -1
	var params = &lm.GetCollectorListParams{
		Fields:     &fields,
		Filter:     &filter,
		Offset:     &offset,
		Size:       &size,
		Context:    nil,
		HTTPClient: nil,
	}
	restRes, err := client.apiClient.LM.GetCollectorList(params)
	if err != nil {
		return nil, fmt.Errorf("get collector ids from group <%v>, group id <%d> failed", collectorGroup.Name, collectorGroup.ID)
	}

	var collectorIds []int32
	for _, item := range restRes.Payload.Items {
		collectorIds = append(collectorIds, item.ID)
	}

	return collectorIds, nil
}

func (client *Client) findDeviceGroup() (*models.DeviceGroup, error) {
	groupName := deviceGroupName(client.option.Cluster)
	filter := fmt.Sprintf("name:%s", groupName)

	var fields = "name,id,parentId"
	var offset int32 = 0
	var size int32 = -1
	var params = &lm.GetDeviceGroupListParams{
		Fields: &fields,
		Filter: &filter,
		Offset: &offset,
		Size:   &size,
	}
	restResp, err := client.apiClient.LM.GetDeviceGroupList(params)
	if err != nil {
		return nil, fmt.Errorf("get device group <%s> failed. msg: %v", client.option.Cluster, err)
	}

	var deviceGroup *models.DeviceGroup
	for _, item := range restResp.Payload.Items {
		if item.ParentID == client.option.ParentId {
			deviceGroup = item
			break
		}
	}

	return deviceGroup, nil
}

func (client *Client) findServiceGroup() (*models.DeviceGroup, error) {
	groupName := serviceGroupName(client.option.Cluster)
	filter := fmt.Sprintf("name:%s", groupName)

	var fields = "name,id,parentId"
	var offset int32 = 0
	var size int32 = -1
	var params = &lm.GetDeviceGroupListParams{
		Fields: &fields,
		Filter: &filter,
		Offset: &offset,
		Size:   &size,
	}

	restResp, err := client.apiClient.LM.GetDeviceGroupList(params)
	if err != nil {
		return nil, fmt.Errorf("get device group <%s> failed. msg: %v", client.option.Cluster, err)
	}

	var deviceGroup *models.DeviceGroup
	for _, item := range restResp.Payload.Items {
		if item.ParentID == client.option.ParentId {
			deviceGroup = item
			break
		}
	}

	return deviceGroup, nil
}

func (client *Client) findCollectorGroup() (*models.CollectorGroup, error) {
	collectorGroupName := collectorGroupName(client.option.Cluster)
	filter := fmt.Sprintf("name:%s", collectorGroupName)

	var fields = ""
	var offset int32 = 0
	var size int32 = -1
	var params = &lm.GetCollectorGroupListParams{
		Fields: &fields,
		Filter: &filter,
		Offset: &offset,
		Size:   &size,
	}

	restResp, err := client.apiClient.LM.GetCollectorGroupList(params)
	if err != nil {
		return nil, fmt.Errorf("get collector group <%s> failed", collectorGroupName)
	}

	var collectorGroup *models.CollectorGroup = nil
	if len(restResp.Payload.Items) > 0 {
		collectorGroup = restResp.Payload.Items[0]
	} else {
		log.Printf("collector group <%s> not found\n", collectorGroupName)
	}

	return collectorGroup, nil
}
