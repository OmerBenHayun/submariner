# how to make compile and test and image.
on subnet 2 on VPC.
workig on vesion 'tag' in my docker hub
## copy to omervm

## building
on vm `omer-vm` preform:

* git reset (WIP)

and:
To run all linting:
```
make unit
```

To build the Go binaries provided by a repository:
```
make build
```
To package those Go binaries into container images:
```
make images
```

WIP


## install subctl 
```
curl -Ls https://get.submariner.io | bash
export PATH=$PATH:~/.local/bin
echo export PATH=\$PATH:~/.local/bin >> ~/.profile
```

## deploy on broker

### delete old kubectl namespaces
```
kubectl get ns
kubectl delete ns submariner
kubectl delete ns submariner-k8s-broker
kubectl delete ns submariner-operator
```

### Deployment of the Broker
```
subctl deploy-broker --kubeconfig <PATH-TO-KUBECONFIG-BROKER>
```

### join command
on `omer-sub2-master1` preform:

(will use version (tag) `dev` from repo `omerbh`)
make sure to do that on root
```
cd omer/submariner/
subctl join --repository omerbh --version dev --cable-driver wireguard --disable-nat broker-info.subm --kubeconfig cluster2config --clusterid cluster2
```



# another useful things 
#### see metrics
get the endpoint ip
```
subctl show endpoints
```
get the gateway ip
```
kubectl -n submariner-operator get pods -o wide | grep 'gateway\|IP'

```
gather metrics from the ip (`xxx.xxx.xxx.xxx` is the ip address)
```
curl xxx.xxx.xxx.xxx:8080/metrics
```

curl only first 10 lines
```
curl 10.243.64.6:8080/metrics | head -n 10
```
get only getway metrics
```
curl 10.243.64.6:8080/metrics | grep gateway
```

#### delete old broker
```
kubectl delete ns submariner-k8s-broker
```
and delete it's relevant file

#### how to run a connection test
 Switch to the context of one of your clusters, i.e. `kubectl config use-context west`
See the [`subctl verify` docs on Submainer's website](https://submariner.io/deployment/subctl/#verify).
Run an nginx container in this cluster, i.e. `kubectl run -n default nginx --image=nginx`
Retrieve the pod IP of the nginx container, looking under the "Pod IP" column for `kubectl get pod -n default`
Change contexts to your other workload cluster, i.e. `kubectl config use-context east`
Run a busybox pod and ping/curl the nginx pod:
```shell
kubectl run -i -t busybox --image=busybox --restart=Never
```
If you don't see a command prompt, try pressing enter.
```shell
ping <NGINX_POD_IP>
wget -O - <NGINX_POD_IP>
```
after finish the test delte the pod
```
kubectl delete pod busybox
```

#### normal join command
```
cd omer/submariner/
subctl join --disable-nat broker-info.subm --kubeconfig <you-cube-config-file> --clusterid <ID>
```


#### push images to docker hub
```
list = submariner-globalnet submariner-route-agent submariner
for IMAGE in ${list[@]}
do
	docker tag quay.io/submariner/${IMAGE}:dev omerbh/${IMAGE}:dev
	docker push omerbh/${IMAGE}:dev
done
```