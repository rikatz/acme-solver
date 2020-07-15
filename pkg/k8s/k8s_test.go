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
package k8s

import (
	"reflect"
	"testing"

	cmclient "github.com/jetstack/cert-manager/pkg/client/clientset/versioned"
	"github.com/jetstack/cert-manager/pkg/client/clientset/versioned/fake"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/jetstack/cert-manager/pkg/apis/acme/v1alpha3"
)

func TestGetChallenge(t *testing.T) {

	var challengeItems []v1alpha3.Challenge

	challenge1 := v1alpha3.Challenge{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "system1",
			Name:      "testecert-example1-2368806360-178418311-1327617168",
		},
		Spec: v1alpha3.ChallengeSpec{
			DNSName: "testcert123.example.com",
			Key:     "fdashfuadsjfuqaj832qFDSAFAEFcsdacadsZZXWSQADASXAS",
		},
	}

	challenge2 := v1alpha3.Challenge{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "system1",
			Name:      "testecert-example1-2368806360-178418311-2519715597",
		},
		Spec: v1alpha3.ChallengeSpec{
			DNSName: "testcert456.example.com",
			Key:     "XXXXXXXXXAAAAAAAaaaaaaaaaaaaaaaaaAAAAAAAAa0000123",
		},
	}

	challenge3 := v1alpha3.Challenge{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "system1",
			Name:      "testecert-example1-2368806360-178418311-2999686616",
		},
		Spec: v1alpha3.ChallengeSpec{
			DNSName: "testcert789.example.com",
			Key:     "123123123VVVVVvvvvvvvvvvvvvvvMMmmmmmMJHGFDSDFSDFSD",
		},
	}

	challenge4 := v1alpha3.Challenge{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "system1",
			Name:      "testecert-example1-2368806360-178418311-324176377",
		},
		Spec: v1alpha3.ChallengeSpec{
			DNSName: "testcert789.example.com",
			Key:     "ThIsIsaS3cretKeyReturnedto4Cm3Challenge",
		},
	}

	challengeItems = append(challengeItems, challenge1, challenge2, challenge3, challenge4)

	challengesList := v1alpha3.ChallengeList{
		Items: challengeItems,
	}

	// TODO: Test unordered array :D
	challenges := []string{"fdashfuadsjfuqaj832qFDSAFAEFcsdacadsZZXWSQADASXAS", "XXXXXXXXXAAAAAAAaaaaaaaaaaaaaaaaaAAAAAAAAa0000123", "123123123VVVVVvvvvvvvvvvvvvvvMMmmmmmMJHGFDSDFSDFSD", "ThIsIsaS3cretKeyReturnedto4Cm3Challenge"}
	data := []struct {
		testName          string
		clientset         cmclient.Interface
		expectedChallenge []string
		inputNamespace    string
		inputCertificate  string
	}{
		{
			clientset:         fake.NewSimpleClientset(&challengesList),
			inputNamespace:    "system1",
			inputCertificate:  "testecert-example1",
			expectedChallenge: challenges,
		},
	}

	for _, single := range data {
		challenge, err := GetChallenge(single.inputCertificate, single.inputNamespace, single.clientset)
		if err != nil {
			t.Errorf("Error in test %v", err)
		}
		if !reflect.DeepEqual(challenge, single.expectedChallenge) {
			t.Errorf("Provided and expected challenges are different - %v XXX %v", challenge, single.expectedChallenge)
		}

	}

}
