package web

import (
	"net/http"
	"sync"
	"time"

	realisv2 "github.com/paypal/gorealis/v2/gen-go/apache/aurora"
	"github.com/sirupsen/logrus"
	"github.paypal.com/PaaS/MagicMatch/aurora"
	"github.paypal.com/PaaS/MagicMatch/mesos"
	"github.paypal.com/PaaS/MagicMatch/metrics"
	"github.paypal.com/PaaS/MagicMatch/scheduler"
)

var log = logrus.New()
var capacitiesMap = make(map[string]*mesos.HostAllocInfo) // mesos data
var hostMetricsMap = make(map[string]*metrics.Host) // TSDB data
var tasks = make([]*realisv2.ScheduledTask, 0) // aurora tasks

var fetcher = metrics.Fetcher{}
var port = "8080" // server port
var dnsSufix string // this is used to map mesos slaves and tsdb hosts.
var isLocalhost = false // test on local without using TSDB.

// initalize does the initialization for MagicMatch.
func initalize(webPort string, isLocal bool, isDemo bool) {
	port = webPort
	isLocalhost = isLocal

	// for local box.
	if isLocal {
		mesos.Host = "192.168.33.3"
		dnsSufix = ""
	} else {
		// for gcp
		mesos.Host = "localhost"
		// clean sensitive data.
	}

	if isDemo {
		// clean sensitive data ...
	}

	fetcher.Init()
}

// fetchHostMetricsData fetches the metrics data from TSDB
func fetchHostMetricsData() {
	for true {
		if isLocalhost {
			// hostMetricsMap = fetcher.FakeHostDataFromFile()
			hostMetricsMap = fetcher.FakeHostData()
		} else {
			hostMetricsMap = fetcher.FetchHostData()
		}
		time.Sleep(60 * time.Second)
	}
}

// fetchAuroraThriftData fetches tasks and offers from Aurora via Thrift call
func fetchAuroraThriftData() {
	for true {
		aurora.Connect()
		tasks = aurora.GetPendingTasks()
		log.Infof("There are %v sorted pending tasks", len(tasks))
		for _, task := range tasks {
			log.Infof("==== task: %v", task.GetAssignedTask().GetTaskId())
		}
	}
}

// fetchMesosData pulls host data from mesos server
func fetchMesosData() {
	log.Infof("FetchMesosData")
	for true {
		capacitiesMap = mesos.GetHostAllocInfo()
		time.Sleep(1 * time.Second)
	}
}

// runJsonServer starts MagicMatch as a web server.
func runJsonServer() {
	router := newRouter()
	log.Fatal(http.ListenAndServe(":"+port, router))
}

var hostMapCache = make(map[string]*scheduler.HostInfo)

// schedule does scheduling based on cached data.
func schedule() {
	for true {
		hostMap := make(map[string]*scheduler.HostInfo)

		for _, c := range capacitiesMap {
			hostMap[c.Hostname] = &scheduler.HostInfo{Hostname: c.Hostname, Capacity: c}
		}

		metrics.BadHostMap = make(map[string]bool)
		for k, _ := range hostMap {
			shortHostName := k[:len(k)-len(dnsSufix)]
			if v, ok := hostMetricsMap[shortHostName]; ok {
				hostMap[k].Metrics = v
				if metrics.IsBadHost(*v) {
					metrics.BadHostMap[k] = true
				}
			} else {
				log.Errorf("cannot find %v host on metric data", shortHostName)
				hostMap[k].Metrics = &metrics.Host{Name: k, CpuUsage: 0}
			}
		}

		sched := &scheduler.Scheduler{}
		log.Infof("Scheduling ....")
		start := time.Now()
		scheduler.MatchesMap = sched.LeastFitMultiTasks(tasks, hostMap, scheduler.CalGenesisCpuFitness)
		elapsed := time.Since(start)
		log.Printf("Scheduling took %s", elapsed)
		hostMapCache = hostMap
		time.Sleep(1 * time.Second)
	}
}

//Execute is called by main.go to run MagicMatch.
func Execute(port string, isLocal bool, isDemo bool) {
	initalize(port, isLocal, isDemo)
	parallelize(runJsonServer, fetchMesosData, fetchHostMetricsData, fetchAuroraThriftData, schedule)
}

// Parallelize runs multiple functions in parallel
func parallelize(functions ...func()) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(functions))

	defer waitGroup.Wait()

	for _, function := range functions {
		go func(copy func()) {
			defer waitGroup.Done()
			copy()
		}(function)
	}
}
