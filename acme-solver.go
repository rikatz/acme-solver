/*
Copyright [2020] [The Acme Solver Authors]

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/
package main

import (
	"flag"
	"log"
	"net"
	"path/filepath"
	"regexp"

	pb "github.com/rikatz/acme-solver/pb"
	"google.golang.org/grpc"

	cmclient "github.com/jetstack/cert-manager/pkg/client/clientset/versioned"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

const (
	defaultBind = ":18853"
)

var (
	solverDomain, solverBind, solverIP, kubeconfig string
	client                                         cmclient.Interface
	err                                            error
)

type dnsServer struct {
	pb.UnimplementedDnsServiceServer
}

func main() {

	if solverDomain == "" {
		log.Fatalf("The -domain argument should not be empty")
	}

	// If the domain does not finish with ".", then insert
	if matched, _ := regexp.MatchString("^\\S*\\.$", solverDomain); !matched {
		solverDomain = solverDomain + "."
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}
	client, err = cmclient.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	log.Printf("Starting server on %s for domain %s", solverDomain, solverBind)
	lis, err := net.Listen("tcp", solverBind)
	if err != nil {
		log.Fatalf("failed to open port: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterDnsServiceServer(grpcServer, &dnsServer{})
	panic(grpcServer.Serve(lis))
}

func init() {

	if home := homedir.HomeDir(); home != "" {
		flag.StringVar(&kubeconfig, "kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		flag.StringVar(&kubeconfig, "kubeconfig", "", "absolute path to the kubeconfig file")
	}

	flag.StringVar(&solverDomain, "domain", "", "What's the domain of this solver. Can not be empty")
	flag.StringVar(&solverBind, "bind", defaultBind, "What's the bind address of the daemon")

	flag.Parse()
}
