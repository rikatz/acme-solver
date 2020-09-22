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
	"context"
	"fmt"
	"regexp"

	cmclient "github.com/jetstack/cert-manager/pkg/client/clientset/versioned"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	objectRegex = regexp.MustCompile("^(\\S*)-(\\S+)-(\\S+)-(\\S+)$")
)

// GetChallenge picks an object name and a namespace and returns the challengeKey from ACME Challenges
func GetChallenge(certname, namespace string, client cmclient.Interface) (challengeKey []string, err error) {

	challengeList, err := client.AcmeV1().Challenges(namespace).List(context.TODO(), metav1.ListOptions{})

	if apierrors.IsNotFound(err) {
		return challengeKey, fmt.Errorf("No challenge found in Namespace %s", namespace)
	}
	if apierrors.IsForbidden(err) {
		return challengeKey, fmt.Errorf("Permission denied while searching for challenges in namespace %s", namespace)
	}

	if err != nil {
		return challengeKey, fmt.Errorf("Failed to communicate with Kubernetes cluster")
	}

	if len(challengeList.Items) < 1 {
		return challengeKey, fmt.Errorf("Namespace %s does not contain any challenges", namespace)
	}

	for _, challenge := range challengeList.Items {
		name := challenge.GetName()

		challengeObjects := objectRegex.FindStringSubmatch(name)

		// TODO: The idea of splitting the challengeObjects here is that we can use the other identifiers later when
		// evolving the program to validate also the Orders, CertificateRequests and Certificate objects

		if len(challengeObjects) == 5 && challengeObjects[1] == certname {

			challengeTXT := challenge.Spec.Key

			if challengeTXT == "" {
				continue
			}
			challengeKey = append(challengeKey, challengeTXT)

		}
	}

	if len(challengeKey) < 1 {
		return challengeKey, fmt.Errorf("No challenge key found  for namespace %s and certificate %s", namespace, certname)
	}

	return challengeKey, nil
}
