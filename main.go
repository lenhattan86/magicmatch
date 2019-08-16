/**
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"flag"
	"fmt"
	"strconv"

	"github.paypal.com/PaaS/MagicMatch/web"
)

func main() {
	fmt.Println("usage: ./magicmatch -port [port_number] -enable [true/false] -is_local [true/false]")
	var port string
	var enable bool
	var isLocal bool
	var isDemo bool

	flag.StringVar(&port, "port", "8080", "Web server port")
	flag.BoolVar(&enable, "enable", true, "enable genemtrics")
	flag.BoolVar(&isLocal, "is_local", false, "using docker-container setup.")
	flag.BoolVar(&isDemo, "is_demo", false, "using docker-container setup.")
	flag.Parse()

	if enable {
		fmt.Println("running magicmatch on port " + port + " is_local " + strconv.FormatBool(isLocal))
		if isDemo {
			fmt.Println("running demo mode")
		}
		web.Execute(port, isLocal, isDemo)
	}
}

func parseJsonFile() {

}
