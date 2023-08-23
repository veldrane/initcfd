### Init cloudflared for k8s instalations

#### Description

Cloudlfared (https://github.com/cloudflare/cloudflared) is tool for managing tunnels from cloudflare infrastructure to own servers. 
Each tunnel must have own tunnel id (in fact the configuration), so its not too easy to deploy them to kubernetes infrastrcutre, where
all pods in replicas has same cfg. This simple try to solve this issue

#### Build

```cd src
go build cmd/initcfd.go
cd ../build
mv ../src/initcfd
docker build . -f Dockerfile -t czdcm-quay.lx.ifortuna.cz/shared-images/initcfd:<VERSION> --squash
```

#### Prerequisity

You need to have configmap clouddflared with file names config-<number>.yaml and secret files with cloudflare credentials
/etc/cloudflared/creds/credentials-<number>.json mounted in all pods. Initcfd must be part of the cloudflared and must run as 
first process.

example of cloudflared configmap:

```
apiVersion: v1
data:
  config-0.yaml: |
    tunnel: 11111111-1111-1111-1111-111111111111
    credentials-file: /etc/cloudflared/creds/credentials-0.json
    metrics: 0.0.0.0:2000
    no-autoupdate: true
    loglevel: info
    lb-pool: TEST-POOL
    ingress:
      - hostname: api.lab.local
        service: http://localhost:8888
        originRequest:
          noTLSVerify: true
      - service: http_status:403
  config-1.yaml: |
    tunnel: 22222222-2222-2222-2222-222222222222
    credentials-file: /etc/cloudflared/creds/credentials-1.json
    metrics: 0.0.0.0:2000
    no-autoupdate: true
    loglevel: info
    lb-pool: TEST-POOL
    ingress:
      - hostname: api.lab.local
        service: http://localhost:8888
        originRequest:
          noTLSVerify: true
      - service: http_status:403
kind: ConfigMap
```

Tool also needs to have configured aprropiration cluster roles. It must be able to get configmaps and pods from kubernetes api.

and the last POD_NAMESPACE environment variable must be published in the deployment object.

example:

```
       - env:
        - name: POD_NAMESPACE
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: metadata.namespace
```

#### How it works

Cloudflared starts like container with entrypoint initcfd tool. It waits for 5 sec to init all pods in replicaset and then ask
for other members in the certain replicas in current namespace, members are sorted and based on the their order inticfd
runs 

```cloudflared tunnel --config /etc/cloudflared/config/config-<order>.yaml```

so each replica has a own config.


#### Limitation and errors

This tool was a part of the cloudlfared PoC integration and lots of stuff are hardcoded in the code (sorry for that), but feel free
to customize it. 

For example: our cloudflared is part nginx apigw stuff, so there is label filter on the getting pods funtion in file
local/ocp4cli/public.go - function GetList

```
func (pods PodsT) GetList(session *SessionT, namespace *string, replicaSet *string) ([]string, error) {

	var podList []string

	listOptions := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=ng-plus-apigw,pod-template-hash=%s", *replicaSet),
	}
.
.
.
```

So please change your filter based on your needs.