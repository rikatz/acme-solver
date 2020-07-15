# ACME SOLVER

This is a gRPC Backend for CoreDNS that answers for acme-challenges from [cert-manager](https://cert-manager.io/) according to the existing challenges objects in a Kubernetes Cluster.

This is an Alpha / non prod yet program, please use it carefully and report issues :)

The gRPC backend is based in Ahmet Alp Balkan (ahmetb) [coredns-grpc-backend-sample](https://github.com/ahmetb/coredns-grpc-backend-sample)

## The Problem

For Let's Encrypt ACME autnehtication work via DNS, a TXT registry that answers for a challenge is required.

[cert-manager](https://cert-manager.io/) allows nowadays the issuance of those certificates, but it relys on a 'programable' DNS Server so those answers registries might be created.

It supports the [acme-dns](https://github.com/joohoi/acme-dns), a simple DNS server that answers for those challenges but this is not as multi-tenant or automatic as we needed in our case.

Also, we wanted to give to user a change to configure his [DNS CNAME](https://www.eff.org/pt-br/deeplinks/2018/02/technical-deep-dive-securing-automation-acme-dns-challenge-validation) pointing to the solver BEFORE asking for the certificate. 


## Architecture

As this is a gRPC Backend for CoreDNS, you need a CoreDNS server. The configuration is described in the [Quickstart](#quick-start)

So this is the architecture:

CoreDNS -> Acme Solver -> Kubernetes

You may have a better understanding taking a look at [Architecture Draw](assets/architecture.png)

ACME Solver will read challenge objects, created by cert-manager and from a specific namespace and this being equal to the DNS query will be answered with the `keys` stored in the object.


Let's see an example:

* I've requested a certificate for the domains `www.mydomain.example.com` and `www1.mydomain.example.com`.

* The certificate was requested with a `Certificate` object, with the `metadata.name: cert1` and inside the namespace `mysite-prod`, and both the domains above inside the object specification.

* `cert-manager` controller creates a `CertificateRequest`, that originates a `Order`, that originates 2 `Challenges` (one for each domain) in this same namespace. The challenges contains a key `.spec.key` with the expected answer from Let's Encrypt

Considering that acme-solver has been started with the flag `-domain solver.example.com`, this will be the domain that Acme Solver will answer and then I can get the challenges from my domain, as the following:

```
dig txt cert1.mysite-prod.solver.example.com

;cert1.mysite-prod.solver.example.com.        IN      TXT

;; ANSWER SECTION:
cert1.mysite-prod.solver.example.com. 0 IN    TXT     "5OJmcI_gZOb_uJWghi3au9ClSKr2r4wLUgbKaA0FPfg"
cert1.mysite-prod.solver.example.com. 0 IN    TXT     "C-tTHjaMiyJLbllVyNyvybMJyT0CbayhzHxkD9Qfrnk"
```

Moving on, Let's Encrypt DNS Challenges depends on the existence of a "magic" record for each domain to be validated, called `_acme-challenge.domain``.

This way, in the authority domain of `mydomain.example.com` I need to point the magic validation registry to the Acme Solver responsible for the challenge resolution:


```
$ORIGIN mydomain.example.com.
@                      3600 SOA   ns1.mydomain.example.com. (
                              zoneadmin.mydomain.example.com. 
                              2016072701                 ; serial number
                              3600                       ; refresh period
                              600                        ; retry period
                              604800                     ; expire time
                              1800                     ) ; minimum ttl
                      
                           86400 NS      ns1.mydomain.example.com.
_acme-challenge.www        43200 CNAME   cert1.mysite-prod.solver.example.com.
_acme-challenge.www1       43200 CNAME   cert1.mysite-prod.solver.example.com.
```

This way, when Let's Encrypt calls the `_acme-challenge` domains will have the following return:

```
dig txt  _acme-challenge.www.mydomain.example.com.

;; ANSWER SECTION:
_acme-challenge.www.mydomain.example.com. 43200 IN    CNAME   cert1.mysite-prod.solver.example.com.
cert1.mysite-prod.solver.example.com. 0 IN    TXT     "5OJmcI_gZOb_uJWghi3au9ClSKr2r4wLUgbKaA0FPfg"
cert1.mysite-prod.solver.example.com. 0 IN    TXT     "C-tTHjaMiyJLbllVyNyvybMJyT0CbayhzHxkD9Qfrnk"
```

## Quick start

* Define which will be the "solver" domain. Here we will use `solver.example.com`
* Take note of the server that will be configured as the authoritative DNS of the domain `solver.example.com`


### Cluster Kubernetes + Cert Manager

* Deploy a new cluster, [KinD](https://github.com/kubernetes-sigs/kind/) may be used for a demo scenario
* Deploy cert-manager in this cluster
* Deploy the cert-manager [Null Issuer Webhook](webhook/README.md) so attempts of cert-manager to create a new DNS entry will be ignored
* Create a ClusterIssuer
* NOT COVERED YET: Create a User, a certificate and the respective Roles and RoleBindings to allow acme-solver to query only Challenge objects

```
kind create cluster
kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v0.15.2/cert-manager.yaml
kubectl apply -f webhook/deploy/webhook.yaml
kubectl apply -f assets/cluster-issuer.yaml
```

* Optionally you can change cert-manager deployments to use a specific DNS Server
```
spec:
  containers:
  - args:
    - --v=2
    - --cluster-resource-namespace=$(POD_NAMESPACE)
    - --leader-election-namespace=kube-system
    - --dns01-recursive-nameservers="8.8.8.8:53"
```

### CoreDNS + Acme Solver
* Start the `acme-solver` remembering to replace the solver domain for your own, and also pointing to the right kubeconfig file
```
 docker run --net=host --rm -v ~/.kube/config:/etc/kubeconfig rpkatz/acme-solver:v1.0.0 -kubeconfig /etc/kubeconfig -domain solver.example.com
 ```

* Install CoreDNS and start the daemon using the available [Corefile](assets/Corefile). Please remember to change the configuration to reflect your solver domain, and the correct IP in the `A registry`. 
* Still in Corefile configuration, point to the place were the Acme Solver gRPC Backend will be answering

### Testing the installation
* Create an example certificate: ``kubectl apply -f  assets/certificate.yaml``
* Verify if the challenges were created: ``kubectl -n staging get challenges``
* Verify if CoreDNS is answering the TXT registry correctly:
```
dig @127.0.0.1 txt certtest.staging.solver.example.com
[...]
;; ANSWER SECTION:
certtest.staging.solver.example.com. 0 IN TXT   "rTe8yfX86MMZh4MX6q8K8moYLbH6PEua193zgeLXYbM"
certtest.staging.solver.example.com. 0 IN TXT   "6w70LEu4Vzmopj36TrWOIExvlpuVaAk7ixt-r2EbzEU"
```
# TODO, Bugs and issues

* More unit tests
* Some e2e tests
* Better docs (document the RBAC process)
* Keep cert-manager objects updated
* A lot of TODO in the code ;P