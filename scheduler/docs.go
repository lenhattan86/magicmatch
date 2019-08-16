/*
Package scheduler implements the matching algorithm for Schedulers.
There are 2 matching algorithms 
	1. FirstFitMultiTasks implements FIFO scheduling
	2. LeastFitMultiTasks implements LeastFit, i.e. Pick the host with the least fitness score.

There are two fitness score functions:
	1. CalFenzoCpuFitness: refer Fenzo for details.
	2. CalGenesisCpuFitness: cpu_request/cpu_capacity + cpu_usage
*/
package scheduler