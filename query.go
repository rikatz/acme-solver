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
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/miekg/dns"
	pb "github.com/rikatz/acme-solver/pb"
	k8s "github.com/rikatz/acme-solver/pkg/k8s"
)

// Query implements the Query for the dnsServer Service.
func (d *dnsServer) Query(ctx context.Context, in *pb.DnsPacket) (*pb.DnsPacket, error) {

	m := new(dns.Msg)

	if err := m.Unpack(in.Msg); err != nil {
		return nil, fmt.Errorf("failed to unpack msg: %v", err)
	}

	r := new(dns.Msg)
	r.SetReply(m)
	r.Authoritative = true

	for _, q := range r.Question {
		hdr := dns.RR_Header{Name: q.Name, Rrtype: q.Qtype, Class: q.Qclass}

		if q.Qtype != dns.TypeTXT {
			log.Printf("%s not a TXT query", q.String())
			r.Rcode = dns.RcodeNameError
			break
		}

		log.Printf("Serving TXT Record: %s", q.Name)
		query := strings.SplitN(strings.ToLower(q.Name), ".", 3)
		certname := query[0]
		namespace := query[1]
		if solverDomain != query[2] {
			log.Printf("%s query domain received is different of configured solver domain %s", query[2], solverDomain)
			r.Rcode = dns.RcodeNameError
			break
		}
		challengeKey, err := k8s.GetChallenge(certname, namespace, client)
		if err != nil || len(challengeKey) < 1 {
			log.Printf("ERROR: %v", err)
			r.Rcode = dns.RcodeNameError
			break
		}
		for _, v := range challengeKey {
			r.Answer = append(r.Answer, &dns.TXT{Hdr: hdr, Txt: []string{v}})
		}

	}

	if len(r.Answer) == 0 {
		r.Rcode = dns.RcodeNameError
	}

	out, err := r.Pack()
	if err != nil {
		return nil, fmt.Errorf("failed to pack msg: %v", err)
	}
	return &pb.DnsPacket{Msg: out}, nil
}
