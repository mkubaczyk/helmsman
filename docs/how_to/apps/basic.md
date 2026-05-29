# Basics

## Install releases

You can run helmsman with the [example.toml](https://github.com/mkubaczyk/helmsman/blob/master/examples/example.toml) or [example.yaml](https://github.com/mkubaczyk/helmsman/blob/master/examples/example.yaml) file.

```shell
$ helmsman --apply -f example.toml
2017/11/19 18:17:57 Parsed [[ example.toml ]] successfully and found [ 2 ] apps.
2017/11/19 18:17:59 WARN: I could not create namespace [staging ]. It already exists. I am skipping this.
2017/11/19 18:17:59 WARN: I could not create namespace [default ]. It already exists. I am skipping this.
2017/11/19 18:18:02 INFO: Executing the following plan ...
---------------
Ok, I have generated a plan for you at: 2017-11-19 18:17:59.347859706 +0100 CET m=+2.255430021
DECISION: release [ jenkins ] is not present in the current k8s context. Will install it in namespace [[ staging ]]
DECISION: release [ artifactory ] is not present in the current k8s context. Will install it in namespace [[ staging ]]
2017/11/19 18:18:02 INFO: attempting: --   installing release [ jenkins ] in namespace [[ staging ]]
2017/11/19 18:18:05 INFO: attempting: --   installing release [ artifactory ] in namespace [[ staging ]]
```

```shell
$ helm list --namespace staging
NAME       	REVISION	UPDATED                 	STATUS  	CHART            	NAMESPACE
artifactory	1       	Sun Nov 19 18:18:06 2017	DEPLOYED	artifactory-6.2.0	staging
jenkins    	1       	Sun Nov 19 18:18:03 2017	DEPLOYED	jenkins-0.9.1    	staging
```

## Delete releases

You can then change your desire, for example to disable the Jenkins release that was created above by setting `enabled = false` :

Then run Helmsman again and it will detect that you want to delete Jenkins:

> Note: Deleting an app entry in the desired state file WILL result in deleting the release. To prevent this, use the `--keep-untracked-releases` flag with your Helmsman command.

```shell
$ helmsman --apply -f example.toml
2017/11/19 18:28:27 Parsed [[ example.toml ]] successfully and found [ 2 ] apps.
2017/11/19 18:28:29 WARN: I could not create namespace [staging ]. It already exists. I am skipping this.
2017/11/19 18:28:29 WARN: I could not create namespace [default ]. It already exists. I am skipping this.
2017/11/19 18:29:01 INFO: Executing the following plan ...
---------------
Ok, I have generated a plan for you at: 2017-11-19 18:28:29.437061909 +0100 CET m=+1.987623555
DECISION: release [ jenkins ] is desired to be deleted . Planning this for you!
DECISION: release [ artifactory ] is desired to be upgraded. Planning this for you!
2017/11/19 18:29:01 INFO: attempting: --   deleting release [ jenkins ]
2017/11/19 18:29:11 INFO: attempting: --   upgrading release [ artifactory ]
```

```shell
$ helm list --namespace staging
NAME       	REVISION	UPDATED                 	STATUS  	CHART            	NAMESPACE
artifactory	2       	Sun Nov 19 18:29:11 2017	DEPLOYED	artifactory-6.2.0	staging
```

```yaml
# ...
apps:
  jenkins:
    description: "jenkins"
    namespace: "staging"
    enabled: false # this tells helmsman to delete it
    chart: "jenkins/jenkins"
    version: "2.15.1"
    valuesFile: ""
    test: false

# ...
```

# Rollback releases

Similarly, if you change `enabled` back to `true`, it will figure out that you would like to roll it back.

```shell
$ helmsman --apply -f example.toml
2017/11/19 18:30:41 Parsed [[ example.toml ]] successfully and found [ 2 ] apps.
2017/11/19 18:30:42 WARN: I could not create namespace [staging ]. It already exists. I am skipping this.
2017/11/19 18:30:43 WARN: I could not create namespace [default ]. It already exists. I am skipping this.
2017/11/19 18:30:49 INFO: Executing the following plan ...
---------------
Ok, I have generated a plan for you at: 2017-11-19 18:30:43.108693039 +0100 CET m=+1.978435517
DECISION: release [ jenkins ] is currently deleted and is desired to be rolledback to namespace [[ staging ]] . No problem!
DECISION: release [ artifactory ] is desired to be upgraded. Planning this for you!
2017/11/19 18:30:49 INFO: attempting: --   rolling back release [ jenkins ]
2017/11/19 18:30:50 INFO: attempting: --   upgrading release [ artifactory ]
```

## Upgrade releases

Every time you run Helmsman, (unless the release is [protected or deployed in a protected namespace](../misc/protect_namespaces_and_releases.md)) it will check if upgrade is necessary (using the helm-diff plugin) and only upgrade if there are changes.

If you change the chart, the existing release will be deleted and a new one with the same name will be created using the new chart.

