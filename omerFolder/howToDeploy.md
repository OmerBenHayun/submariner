# how to deploy
* building and pushing

run as root:
```
make build
make images

for IMAGE in submariner-globalnet submariner-route-agent submariner
do
 docker tag quay.io/submariner/${IMAGE}:dev omerbh/${IMAGE}:omer_2n_ver
 docker push omerbh/${IMAGE}:omer_2n_ver
done

```

* deployment
    * on the master:
        * delete the relevant namespaces:
            ```
          kubectl delete ns submariner
          kubectl delete ns submariner-operator
            ```
    * on ther worker(s)
        * `shh` to worker and remove all the images with `docker rmi -f $(docker images -q)` and exit the worker afterwards.
    * on the master again:
        * run the join command 
            ```
            subctl join --repository omerbh --version omer_2n_ver --cable-driver wireguard --disable-nat broker-info.subm --kubeconfig cluster1config --clusterid cluster1
            ```


``
# script for that 

```
# remove all the curently running submariners and remove all the images on the workers
for masterIP in 161.156.166.17 161.156.173.245 161.156.173.199 161.156.166.225
do
	ssh -i /root/.ssh/id_rsa root@${masterIP} 'kubectl delete ns submariner;kubectl delete ns submariner-operator'
done

for WorkerIP in 161.156.170.18 161.156.162.45 161.156.162.156 161.156.163.251 161.156.174.103
do
	ssh -i /root/.ssh/id_rsa root@${WorkerIP} 'docker rmi -f $(docker images -q)'
done

for cluster_num in 1 2
do
    subctl join --repository omerbh --version omer_2n_ver --cable-driver wireguard --disable-nat broker-info.subm --kubeconfig cluster${cluster_num}config --clusterid cluster${cluster_num}
done
```