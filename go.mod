module github.paypal.com/PaaS/MagicMatch

require (
	github.com/Jeffail/gabs v1.4.0
	github.com/gorilla/mux v1.7.2
	github.com/paypal/gorealis v1.21.1
	github.com/paypal/gorealis/v2 v2.0.1
	github.com/pkg/errors v0.8.0
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.0-20180115160933-0c34d16c3123
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.4.0
	gopkg.in/yaml.v2 v2.2.2
)

//replace github.com/paypal/gorealis/v2 => ../../projects/gorealis/

replace github.paypal.com/PaaS/MagicMatch => ./
