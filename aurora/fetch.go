package aurora

import (
	"strings"
	"time"

	realis "github.com/paypal/gorealis/v2"
	realisv2 "github.com/paypal/gorealis/v2/gen-go/apache/aurora"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var configFile string
var client *realis.Client

// authentication purpose.
var username, password, zkAddr, schedAddr string
var clientKey, clientCert string
var skipCertVerification bool
var caCertsPath string
var cmdInterval = time.Second * 5
var cmdTimeout = time.Minute * 10
var log = logrus.New()

// Init initalizes a aurora.fetcher
func init() {
	username = "aurora"
	password = "secret"
	schedAddr = "localhost:8081"
}

// Connect to aurora
func Connect() {
	var err error

	zkAddrSlice := strings.Split(zkAddr, ",")

	viper.SetConfigFile(configFile)
	err = viper.ReadInConfig()
	if err == nil {
		// Best effort load configuration. Will only set config values when flags have not set them already.
		if viper.IsSet("zk") && len(zkAddrSlice) == 1 && zkAddrSlice[0] == "" {
			zkAddrSlice = viper.GetStringSlice("zk")
		}

		if viper.IsSet("username") && username == "" {
			username = viper.GetString("username")
		}

		if viper.IsSet("password") && password == "" {
			password = viper.GetString("password")
		}

		if viper.IsSet("clientKey") && clientKey == "" {
			clientKey = viper.GetString("clientKey")
		}

		if viper.IsSet("clientCert") && clientCert == "" {
			clientCert = viper.GetString("clientCert")
		}

		if viper.IsSet("caCertsPath") && caCertsPath == "" {
			caCertsPath = viper.GetString("caCertsPath")
		}

		if viper.IsSet("skipCertVerification") && !skipCertVerification {
			skipCertVerification = viper.GetBool("skipCertVerification")
		}
	}

	realisOptions := []realis.ClientOption{realis.BasicAuth(username, password),
		realis.ThriftJSON(),
		realis.Timeout(20 * time.Second),
		realis.BackOff(realis.Backoff{
			Steps:    2,
			Duration: 10 * time.Second,
			Factor:   2.0,
			Jitter:   0.1,
		}),
		realis.SetLogger(log)}

	// Prefer zookeeper if both ways of connecting are provided
	if len(zkAddrSlice) > 0 && zkAddrSlice[0] != "" {
		// Configure Zookeeper to connect
		zkOptions := []realis.ZKOpt{realis.ZKEndpoints(zkAddrSlice...), realis.ZKPath("/aurora/scheduler")}
		realisOptions = append(realisOptions, realis.ZookeeperOptions(zkOptions...))
	} else if schedAddr != "" {
		realisOptions = append(realisOptions, realis.SchedulerUrl(schedAddr))
	} else {
		logrus.Fatalln("Zookeeper address or Scheduler URL must be provided.")
	}

	// Client certificate configuration if available
	if clientKey != "" || clientCert != "" || caCertsPath != "" {
		realisOptions = append(realisOptions,
			realis.CertsPath(caCertsPath),
			realis.ClientCerts(clientKey, clientCert),
			realis.InsecureSkipVerify(skipCertVerification))
	}

	// Connect to Aurora Scheduler and create a client object
	client, err = realis.NewClient(realisOptions...)

	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("connect to %v sucessfully")
	}

}

// Close aurora connection
func Close() {
	client.Close()
	log.Infof("closed Thrift connection to Aurora...")
}

// get pending tasks from aurora
func GetPendingTasks() []*realisv2.ScheduledTask {
	PENDING_STATES := []realisv2.ScheduleStatus{realisv2.ScheduleStatus_PENDING}
	taskQuery := &realisv2.TaskQuery{Statuses: PENDING_STATES}
	tasks, err := client.GetTasksWithoutConfigs(taskQuery)

	if err != nil {
		log.Errorf("error: %+v\n", err)
		tasks = make([]*realisv2.ScheduledTask, 0)
	}

	return tasks
}