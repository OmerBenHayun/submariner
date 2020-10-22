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
            subctl join --repository omerbh --version dev --cable-driver wireguard --disable-nat broker-info.subm --kubeconfig cluster2config --clusterid cluster2
            ```
