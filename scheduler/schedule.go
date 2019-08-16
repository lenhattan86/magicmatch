package scheduler

import (
	"sort"
	"strconv"
	"time"

	realisv2 "github.com/paypal/gorealis/v2/gen-go/apache/aurora"
	"github.paypal.com/PaaS/MagicMatch/mesos"
	hostmetrics "github.paypal.com/PaaS/MagicMatch/metrics"
)

// Match from Task to Host
type Match struct {
	HostOffer string `json:"host_offer"`
	TaskId    string `json:"task_id"`
}

// cached matches
var MatchesMap = make(map[string]Match)

var leastFitnessScore = 0.1 // if fitness score is less than this, the host is picked.
var maxScoreVal = 2.00 // Max value of fitness score.

// structure to wrap all fectched data from tsdb, mesos, aurora.
type HostInfo struct {
	Hostname string
	Capacity *mesos.HostAllocInfo
	Metrics  *hostmetrics.Host
}

// Resource structure
type Resource struct {
	cpu  float64
	mem  int64
	disk int64
}

// Used resources.
var usedResourcesMap = make(map[string]Resource)

// Fitness Score function pointer
type FinessScoreFunc func(task *realisv2.ScheduledTask, hostInfo *HostInfo) float64

// define executor overhead.
var executorOverhead = Resource{cpu: 0.25, mem: 128, disk: 0}

// Scheduler implements scheduling functions.
type Scheduler struct {
}

// FirstFitMultiTasks keeps scheduling tasks to hosts till they are full.
func (this *Scheduler) FirstFitMultiTasks(tasks []*realisv2.ScheduledTask, hostInfoMap map[string]*HostInfo) map[string]Match {
	res := make(map[string]Match)
	usedResourcesMap = make(map[string]Resource)

	for _, task := range tasks {
		taskId := task.GetAssignedTask().GetTaskId()
		var firstOffer *HostInfo
		for hostname, hostInfo := range hostInfoMap {
			// hostoffer must have enough resources.
			if isSufficientResource(task, hostInfo) {
				res[taskId] = Match{TaskId: taskId, HostOffer: hostname}
				firstOffer = hostInfo
				break
			}
		}

		if firstOffer != nil {
			request := GetRequest(task)
			if used, ok := usedResourcesMap[firstOffer.Hostname]; ok {
				usedResourcesMap[firstOffer.Hostname] =
					Resource{cpu: request.cpu + used.cpu,
						mem:  request.mem + used.mem,
						disk: request.disk + used.disk}
			} else {
				usedResourcesMap[firstOffer.Hostname] =
					Resource{cpu: request.cpu,
						mem:  request.mem,
						disk: request.disk}
			}
		} else {
			res[taskId] = Match{TaskId: taskId, HostOffer: "null"}
		}
	}
	return addTimestamp(res)
}

// LeastFitMultiTasks keeps scheduling tasks to hosts till they are full.
func (this *Scheduler) LeastFitMultiTasks(tasks []*realisv2.ScheduledTask, hostInfoMap map[string]*HostInfo, fitnessFunc FinessScoreFunc) map[string]Match {
	res := make(map[string]Match)
	usedResourcesMap = make(map[string]Resource)

	for _, task := range tasks {
		taskId := task.GetAssignedTask().GetTaskId()
		minFitnessScore := 999999.99
		var bestCandidate *HostInfo
		bestCandidate = nil

		var keys []string
		for k := range hostInfoMap {
			keys = append(keys, k)
		}
		// this is used for demo purpose. I we submit 2 jobs continously,
		// they may go to different hosts because of host metric delay.
		sort.Strings(keys)

		// for hostname, hostInfo := range hostInfoMap {
		for _, hostname := range keys {
			hostInfo := hostInfoMap[hostname]
			// hostoffer must have enough resources.
			if isSufficientResource(task, hostInfo) {
				score := fitnessFunc(task, hostInfo)
				// if the score is small enough.
				if score < leastFitnessScore {
					res[taskId] = Match{TaskId: taskId, HostOffer: hostname}
					bestCandidate = hostInfo
					break
				}
				// looking for the smallest fitness score.
				if score < minFitnessScore {
					bestCandidate = hostInfo
					minFitnessScore = score
				}
			}
		}

		if bestCandidate != nil {
			request := GetRequest(task)
			res[taskId] = Match{TaskId: taskId, HostOffer: bestCandidate.Hostname}
			if used, ok := usedResourcesMap[bestCandidate.Hostname]; ok {
				usedResourcesMap[bestCandidate.Hostname] =
					Resource{cpu: request.cpu + used.cpu,
						mem:  request.mem + used.mem,
						disk: request.disk + used.disk}
			} else {
				usedResourcesMap[bestCandidate.Hostname] =
					Resource{cpu: request.cpu,
						mem:  request.mem,
						disk: request.disk}
			}
		} else {
			res[taskId] = Match{TaskId: taskId, HostOffer: "null"}
		}
	}
	return addTimestamp(res)
}

// addTimestamp add timestamp to the matches so that we know when they were created.
func addTimestamp(matchesMap map[string]Match) map[string]Match {
	timestamp := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	matchesMap["timestamp"] = Match{TaskId: timestamp, HostOffer: ""}
	return matchesMap
}

// GetRequest computes the resource request from ScheduledTask
func GetRequest(task *realisv2.ScheduledTask) Resource {
	resources := task.GetAssignedTask().GetTask().GetResources()
	return Resource{cpu: resources[2].GetNumCpus(),
		mem:  int64(resources[1].GetRamMb()),
		disk: int64(resources[0].GetDiskMb())}
}

// CalFenzoCpuFitness computes Fenzo CPU score
func CalFenzoCpuFitness(task *realisv2.ScheduledTask, hostInfo *HostInfo) float64 {
	request := GetRequest(task)
	offer := computeOffer(hostInfo)
	if _, ok := hostmetrics.BadHostMap[hostInfo.Hostname]; ok {
		return maxScoreVal
	}
	if used, ok := usedResourcesMap[hostInfo.Hostname]; ok {
		return (request.cpu + executorOverhead.cpu + hostInfo.Capacity.Cpu - offer.cpu + used.cpu) / hostInfo.Capacity.Cpu
	}
	return (request.cpu + executorOverhead.cpu + hostInfo.Capacity.Cpu - offer.cpu) / hostInfo.Capacity.Cpu
}

// CalGenesisCpuFitness compute a fitness score based on cpu load and requests.
func CalGenesisCpuFitness(task *realisv2.ScheduledTask, hostInfo *HostInfo) float64 {
	request := GetRequest(task)
	if _, ok := hostmetrics.BadHostMap[hostInfo.Hostname]; ok {
		return maxScoreVal
	}
	if used, ok := usedResourcesMap[hostInfo.Hostname]; ok {
		return (request.cpu + executorOverhead.cpu + (hostInfo.Metrics.CpuUsage * hostInfo.Capacity.Cpu / 100) + used.cpu) / hostInfo.Capacity.Cpu
	}
	return (request.cpu + executorOverhead.cpu + (hostInfo.Metrics.CpuUsage * hostInfo.Capacity.Cpu / 100)) / hostInfo.Capacity.Cpu
}

// computeOffer estimates offer from mesos or just use offers from aurora.
func computeOffer(hostInfo *HostInfo) Resource {
	return Resource{cpu: hostInfo.Capacity.AvailCpu,
		mem:  int64(hostInfo.Capacity.AvailMem),
		disk: int64(hostInfo.Capacity.AvailDisk)}
}

// isSufficientResourceBase checks if checks fits the offer.
func isSufficientResourceBase(task *realisv2.ScheduledTask, hostInfo *HostInfo) bool {
	offer := computeOffer(hostInfo)
	request := GetRequest(task)
	if request.cpu+executorOverhead.cpu > offer.cpu {
		return false
	}

	if request.mem+executorOverhead.mem > offer.mem {
		return false
	}

	if request.disk+executorOverhead.disk > offer.disk {
		return false
	}
	return true
}

// isSufficientResource checks if request fits the offer.
func isSufficientResource(task *realisv2.ScheduledTask, hostInfo *HostInfo) bool {
	offer := computeOffer(hostInfo)
	request := GetRequest(task)
	if used, ok := usedResourcesMap[hostInfo.Hostname]; ok {
		if request.cpu+executorOverhead.cpu+used.cpu > offer.cpu {
			return false
		}
		if request.mem+executorOverhead.disk+used.mem > offer.mem {
			return false
		}
		if request.disk+executorOverhead.mem+used.disk > offer.disk {
			return false
		}
		return true
	} else {
		return isSufficientResourceBase(task, hostInfo)
	}
}
