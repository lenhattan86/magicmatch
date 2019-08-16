package metrics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

// metrics.Fetcher is used to fetch metrics from TSDB
type Fetcher struct {
	cpuUsageUrl string
	memUsageUrl string
	diskFreeUrl string
}

// for development purpose, you can create a tsdb json file and place it in jsonFilePath
var jsonFilePath = "/usr/share/magicmatch/data" 

var log = logrus.New()

// Init initialize Fetcher.
// clean up sensitive data.
func (fetcher *Fetcher) Init() {
	fetcher.cpuUsageUrl = TSDB_URL + "start=" + TIME_AGO + "&m=" + AGREGATOR + ":" + CPU_METRIC + "{host=" + HOST + ",colo=" + COLO + "}&format=json"
	fetcher.memUsageUrl = TSDB_URL + "start=" + TIME_AGO + "&m=" + AGREGATOR + ":" + MEM_USAGE_METRIC + "{host=" + HOST + ",colo=" + COLO + "}&format=json"
	fetcher.diskFreeUrl = TSDB_URL + "start=" + TIME_AGO + "&m=" + AGREGATOR + ":" + DISK_FREE_METRIC + "{host=" + HOST + ",colo=" + COLO + "}&format=json"
}

// FetchHostData fetch metrics data in json format.
func (fetcher *Fetcher) FetchHostData() map[string]*Host {
	hosts := make(map[string]*Host)

	// CPU usage
	cpuUsageMetricList, err := fetcher.readJSONFromUrlForCpu()
	if err != nil {
		log.Errorf("error: %+v\n", err)
	}
	for _, cpuUsageMetric := range cpuUsageMetricList {
		host := &Host{ CpuUsage:0.0, MemUsage: 0, DiskUsage: 0}
		host.Name = cpuUsageMetric.Tags.Host
		host.CpuUsage = cpuUsageMetric.DataPoints[len(cpuUsageMetric.DataPoints)-1][1]
		hosts[host.Name] = host
	}

	// memory
	//TODO(nhatle) add pulling memory usage
	// disk
	//TODO(nhatle) add pulling disk usage
	log.Infof("number of hosts from tsdb: %v ", len(hosts))
	return hosts
}

// FakeHostDataFromFile for development purpose
func (fetcher *Fetcher) FakeHostDataFromFile() map[string]*Host {
	hosts := make(map[string]*Host)

	// CPU usage
	cpuUsageMetricList, err := fetcher.readJSONFromFileForCpu()
	if err != nil {
		log.Errorf("error: %+v\n", err) //todo(tanle) remove Fatalf
	}
	log.Infof("cpuUsageMetricList: %v", len(cpuUsageMetricList))

	for _, cpuUsageMetric := range cpuUsageMetricList {
		host := &Host{CpuUsage:0.0, MemUsage: 0, DiskUsage: 0}
		host.Name = cpuUsageMetric.Tags.Host
		host.CpuUsage = cpuUsageMetric.DataPoints[len(cpuUsageMetric.DataPoints)-1][1]
		hosts[host.Name] = host
	}

	return hosts
}

// FakeHostData for development purpose
func (fetcher *Fetcher) FakeHostData() map[string]*Host {
	hosts := make(map[string]*Host)
	hosts["agent1"] = &Host{Name: "agent1", CpuUsage: 0.0}
	hosts["agent2"] = &Host{Name: "agent2", CpuUsage: 0.0}
	return hosts
}


// json format of tsdb
type Tags struct {
	Colo        string `json:"colo"`
	Environment string `json:"environment"`
	Servergroup string `json:"servergroup"`
	Host string `json:"host"`
}

type MetricData struct {
	MetricName string      `json:"MetricName"`
	Tags       Tags        `json:"Tags"`
	DataPoints [][]float64 `json:"DataPoints"`
}

// readJSONFromFileForCpu for development purpose
func (fetcher *Fetcher) readJSONFromFileForCpu() ([]MetricData, error) {
	filePath := jsonFilePath + "/hostCpu.json"
	respByte, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Errorf("error reading file " + filePath)
		log.Errorf(err.Error())
	}
	metricList := make([]MetricData, 0)
	if err := json.Unmarshal(respByte, &metricList); err != nil {
		fmt.Println(respByte)
		return nil, err
	}

	return metricList, nil
}

// readJSONFromUrlForCpu fetches cpu usage
func (fetcher *Fetcher) readJSONFromUrlForCpu() ([]MetricData, error) {
	log.Infof(fetcher.cpuUsageUrl)
	resp, err := http.Get(fetcher.cpuUsageUrl)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	metricList := make([]MetricData, 0)

	buf := new(bytes.Buffer)

	buf.ReadFrom(resp.Body)

	respByte := buf.Bytes()
	if err := json.Unmarshal(respByte, &metricList); err != nil {
		return nil, err
	}

	return metricList, nil
}
