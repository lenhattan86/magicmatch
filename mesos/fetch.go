package mesos

import (
	"bytes"
	"net/http"

	"github.com/Jeffail/gabs"
	"github.com/sirupsen/logrus"
)

var host = "localhost" // mesos host, it is localhost because we deploy MagicMath in the same node with aurora.
var port = "5050" // mesos port
var log = logrus.New()

// mesos hostinfo
type HostAllocInfo struct {
	Hostname  string
	Cpu       float64
	Mem       float64
	Disk      float64
	AvailCpu  float64
	AvailMem  float64
	AvailDisk float64
}

// GetHostAllocInfo pulls host info from mesos master.
func GetHostAllocInfo() map[string]*HostAllocInfo {
	capacitiesMap := make(map[string]*HostAllocInfo)
	HostCapacityUrl := "http://" + host + ":" + port + "/master/slaves"
	resp, err := http.Get(HostCapacityUrl)
	if err != nil {
		log.Error("cannot fetch data from mesos master")
		log.Error(err.Error())
		return capacitiesMap
	}

	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	respByte := buf.Bytes()
	hostAllocJson, err := gabs.ParseJSON(respByte)

	// S is shorthand for Search
	children, err := hostAllocJson.S("slaves").Children()
	if err != nil {
		log.Error(err.Error())
		return capacitiesMap
	}

	for _, child := range children {
		hostname := child.Path("hostname").Data().(string)

		resources, err := child.Path("resources").ChildrenMap()
		if err != nil {
			panic(err)
		}
		cpu := resources["cpus"].Data().(float64)
		mem := resources["mem"].Data().(float64)
		disk := resources["disk"].Data().(float64)

		resources, err = child.Path("used_resources").ChildrenMap()
		if err != nil {
			panic(err)
		}
		availCpu := cpu - resources["cpus"].Data().(float64)
		availMem := mem - resources["mem"].Data().(float64)
		availDisk := disk - resources["disk"].Data().(float64)
		capacity := &HostAllocInfo{Hostname: hostname, Cpu: cpu, Mem: mem, Disk: disk, AvailCpu: availCpu, AvailMem: availMem, AvailDisk: availDisk}
		capacitiesMap[hostname] = capacity
	}
	log.Infof("There are %v hosts", len(capacitiesMap))
	return capacitiesMap
}