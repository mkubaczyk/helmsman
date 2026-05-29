---
version: v3.0.0-beta5
---

# Move charts across namespaces

If you have a workflow for testing a release first in the `staging` namespace then move it to the `production` namespace, Helmsman can help you.

> NOTE: If your chart uses a persistent volume, then you have to read the note on PVs below first.

```toml
#...
[namespaces]
[namespaces.staging]
[namespaces.production]

[apps]
    [apps.jenkins]
    description = "jenkins"
    namespace = "staging" # this is where it is deployed
    enabled = true
    chart = "jenkins/jenkins"
    version = "2.15.1"
    valuesFile = ""
    test = true

#...
```

```yaml
# ...
namespaces:
  staging:
  production:

apps:
  jenkins:
    description: "jenkins"
    namespace: "staging" # this is where it is deployed
    enabled: true
    chart: "jenkins/jenkins"
    version: "2.15.1"
    valuesFile: ""
    test: true

# ...
```

Then if you change the namespace key for jenkins:

```toml
#...
[namespaces]
[namespaces.staging]
[namespaces.production]

[apps]
    [apps.jenkins]
    description = "jenkins"
    namespace = "production" # we want to move it to production
    enabled = true
    chart = "jenkins/jenkins"
    version = "2.15.1"
    valuesFile = ""
    test = true

#...
```

```yaml
# ...
namespaces:
  staging:
  production:

apps:
  jenkins:
    description: "jenkins"
    namespace: "production" # we want to move it to production
    enabled: true
    chart: "jenkins/jenkins"
    version: "2.15.1"
    valuesFile: ""
    test: true

# ...
```

Helmsman will delete the jenkins release from the `staging` namespace and install it in the `production` namespace (default in the above setup).

## Note on Persistent Volumes

Helmsman does not automatically move PVCs across namespaces. You have to follow the steps below to retain your data when moving an app to a different namespace.

Persistent Volumes (PV) are accessed through Persistent Volume Claims (PVC). But **PVCs are namespaced objects** which means moving an application from one namespace to another will result in a new PVC created in the new namespace. The old PV -which possibly contains your application data- will still be mounted to the old PVC (the one in the old namespace) even if you have deleted your application helm release.

Now, the newly created PVC (in the new namespace) will not be able to mount to the old PV and instead it will mount to any other available one or (in the case of dynamic provisioning) will provision a new PV. This means the application in the new namespace does not have the old data. Don't panic, the old PV is still there and contains your old data.

### Mounting the old PV to the new PVC (in the new namespace)

1. You have to make sure the _Reclaim Policy_ of the old PV is set to **Retain**. In dynamic provisioned PVs, the default is Delete.
To change it:

    ```shell
    kubectl patch pv <your-pv-name> -p '{"spec":{"persistentVolumeReclaimPolicy":"Retain"}}'
    ```

2. Once your old helm release is deleted, the old PVC and PV are still there. Go ahead and delete the PVC

    ```shell
    kubectl delete pvc <your-pvc-name> --namespace <the-old-namespace>
    ```

    Since, we changed the Reclaim Policy to Retain, the PV will stay around (with all your data).

3. The PV is now in the **Released** state but not yet available for mounting.

    ```shell
    kubectl get pv
    NAME                                       CAPACITY   ACCESSMODES   RECLAIMPOLICY   STATUS     CLAIM                                                             STORAGECLASS   REASON    AGE
    ...
    pvc-f791ef92-01ab-11e8-8a7e-02412acf5adc   20Gi       RWO           Retain          Released   staging/myapp-persistent-storage-test-old-0       gp2                      5m
    ```

    Now, you need to make it Available, for that we need to remove the `PV.Spec.ClaimRef` from the PV spec:

    ```shell
    kubectl edit pv <pv-name>
    # edit the file and save it
    ```

    Now, the PV should become in the **Available** state:

    ```shell
    kubectl get pv
    NAME                                       CAPACITY   ACCESSMODES   RECLAIMPOLICY   STATUS      CLAIM                                                             STORAGECLASS   REASON    AGE
    ...
    pvc-f791ef92-01ab-11e8-8a7e-02412acf5adc   20Gi       RWO           Retain          Available                                                                     gp2                      7m
    ```

4. Delete the new PVC (and its mounted PV if necessary), then delete your application pod(s) in the new namespace. Assuming you have a deployment/replication controller in place, the pod will be recreated in the new namespace and this time will mount to the old volume and your data will be once again available to your application.

    > NOTE: if there are multiple PVs in the Available state and they match capacity and read access for your application, then your application (in the new namespace) might mount to any of them. In this case, either ensure only the right PV is in the available state or make the PV available to a specific PVC - pre-fill `PV.Spec.ClaimRef` with a pointer to a PVC. Leave the `PV.Spec.ClaimRef,UID` empty, as the PVC does not need to exist at this point and you don't know PVC's UID. This PV can be bound only to the specified PVC

Further details:

* https://github.com/kubernetes/kubernetes/issues/48609
* https://kubernetes.io/docs/tasks/administer-cluster/change-pv-reclaim-policy/

