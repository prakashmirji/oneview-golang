package main

import (
	"fmt"
	"github.com/HewlettPackard/oneview-golang/ov"
	"github.com/HewlettPackard/oneview-golang/utils"
	"os"
	"strconv"
	"strings"
)

func main() {
	var (
		ClientOV            *ov.OVClient
		logical_enclosure   = "TestLE"
		logical_enclosure_1 = "TestLE"
		logical_enclosure_2 = "log_enclosure88"
		scope_name          = "testing"
		li_name             = "<logical_interconnect_name>"
	)
	apiversion, _ := strconv.Atoi(os.Getenv("ONEVIEW_APIVERSION"))

	ovc := ClientOV.NewOVClient(
		os.Getenv("ONEVIEW_OV_USER"),
		os.Getenv("ONEVIEW_OV_PASSWORD"),
		os.Getenv("ONEVIEW_OV_DOMAIN"),
		os.Getenv("ONEVIEW_OV_ENDPOINT"),
		false,
		apiversion,
		"*")

	fmt.Println("#................... Create Logical Enclosure ...............#")
	enclosureUris := new([]utils.Nstring)
	*enclosureUris = append(*enclosureUris, utils.NewNstring("/rest/enclosures/0000000000A66101"))
	*enclosureUris = append(*enclosureUris, utils.NewNstring("/rest/enclosures/0000000000A66102"))
	*enclosureUris = append(*enclosureUris, utils.NewNstring("/rest/enclosures/0000000000A66103"))

	enc_grp, err := ovc.GetEnclosureGroupByName("EG-Synergy-Local")

	logicalEnclosure := ov.LogicalEnclosure{Name: logical_enclosure_1,
		EnclosureUris:     *enclosureUris,
		EnclosureGroupUri: enc_grp.URI}

	er := ovc.CreateLogicalEnclosure(logicalEnclosure)
	if er != nil {
		fmt.Println("............... Logical Enclosure Creation Failed:", er)
	} else {
		fmt.Println(".... Logical Enclosure Created Success")
	}

	logicalInterconnect, _ := ovc.GetLogicalInterconnects("", "", "")
	li := ov.LogicalInterconnect{}
	for i := 0; i < len(logicalInterconnect.Members); i++ {
		if logicalInterconnect.Members[i].Name == li_name {
			li = logicalInterconnect.Members[i]
		}
	}

	fmt.Println("#................... Create Logical Enclosure Support Dumps ...............#")

	supportdmp := ov.SupportDumps{ErrorCode: "MyDump16",
		ExcludeApplianceDump:    false,
		LogicalInterconnectUris: []utils.Nstring{li.URI}}

	li_id := strings.Replace(string(li.URI), "/rest/logical-interconnects/", "", 1)

	data, er := ovc.CreateSupportDump(supportdmp, li_id)

	if er != nil {
		fmt.Println("............... Logical Enclosure Support Dump Creation Failed:", er)
	} else {
		fmt.Println(".... Logical Enclosure Support Dump Created Successfully", data)
		fmt.Println(data["URI"])
		id := strings.Trim(data["URI"], "/rest/tasks/")
		task, err := ovc.GetTasksById("", "", "", "", id)
		if err != nil {
			fmt.Println("Error getting the task details ", err)
		}
		fmt.Println(task)
	}

	log_enc, _ := ovc.GetLogicalEnclosureByName(logical_enclosure_1)

	fmt.Println("#................... Logical Enclosure by Name ...............#")
	log_en, err := ovc.GetLogicalEnclosureByName(logical_enclosure)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(log_en)
	}

	// Update From Group

	err = ovc.UpdateFromGroupLogicalEnclosure(log_en)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("#............. Update From Group Logical Enclosure Successfully .....#")
	}

	scope1, err := ovc.GetScopeByName(scope_name)
	scope_uri := scope1.URI
	scope_Uris := new([]string)
	*scope_Uris = append(*scope_Uris, scope_uri.String())

	sort := "name:desc"
	log_en_list, err := ovc.GetLogicalEnclosures("", "", "", *scope_Uris, sort)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("# ................... Logical Enclosures List .................#")
		for i := 0; i < len(log_en_list.Members); i++ {
			fmt.Println(log_en_list.Members[i].Name)
		}
	}

	log_enc.Name = logical_enclosure_2
	err = ovc.UpdateLogicalEnclosure(log_enc)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("#.................... Logical Enclosure after Updating ...........#")
		log_en_after_update, err := ovc.GetLogicalEnclosures("", "", "", *scope_Uris, sort)
		if err != nil {
			fmt.Println(err)
		} else {
			for i := 0; i < len(log_en_after_update.Members); i++ {
				fmt.Println(log_en_after_update.Members[i].Name)
			}
		}
	}

	err = ovc.DeleteLogicalEnclosure(logical_enclosure_2)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("#...................... Deleted Logical Enclosure Successfully .....#")
	}

}
