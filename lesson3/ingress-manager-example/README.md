# ingress manager example

本示例为参考 [GitHub - ingress manager](https://github.com/baidingtech/operator-lesson-demo/tree/main/11) 实现以下需求:
![prd.png](prd.png)

## 1.kind创建k8s集群并安装ingress
下面以 MacOS 系统为例。

1、安装 kind v0.31.0 版本

```shell
[ $(uname -m) = arm64 ] && curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.31.0/kind-darwin-arm64
chmod +x ./kind
mv ./kind $GOPATH/bin/
```



2、安装 cloud-provider-kind v0.10.0 版本

```shell
$ curl -Lo ./cloud-provider-kind_0.10.0_darwin_arm64.tar.gz https://github.com/kubernetes-sigs/cloud-provider-kind/releases/download/v0.10.0/cloud-provider-kind_0.10.0_darwin_arm64.tar.gz
$ tar -xzf cloud-provider-kind_0.10.0_darwin_arm64.tar.gz
$ chmod +x cloud-provider-kind
$ mv cloud-provider-kind $GOPATH/bin/
```



3、运行 cloud-provider-kind

```shell
$ sudo cloud-provider-kind
```



4、使用 kind 创建 k8s 集群

```shell
$ kind create cluster --config kind-config-1c2w.yaml --name 1c2w
```

kind-config-1c2w.yaml 内容如下：

```yaml
# kind: v0.31.0, node: v1.34.3
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
networking:
  # WARNING: It is _strongly_ recommended that you keep this the default
  # (127.0.0.1) for security reasons. However it is possible to change this.
  apiServerAddress: "127.0.0.1"
  # By default the API server listens on a random open port.
  # You may choose a specific port but probably don't need to in most cases.
  # Using a random port makes it easier to spin up multiple clusters.
  apiServerPort: 6443
# 1 control plane node and 2 workers
nodes: # the control plane node config
  - role: control-plane
    image: kindest/node:v1.34.3@sha256:08497ee19eace7b4b5348db5c6a1591d7752b164530a36f855cb0f2bdcbadd48
    extraPortMappings:
      - containerPort: 80
        hostPort: 80
        listenAddress: "0.0.0.0" # Optional, defaults to "0.0.0.0"
        protocol: tcp # Optional, defaults to tcp
      - containerPort: 443
        hostPort: 443
        listenAddress: "0.0.0.0" # Optional, defaults to "0.0.0.0"
        protocol: tcp # Optional, defaults to tcp
      - containerPort: 30080
        hostPort: 30080
        listenAddress: "0.0.0.0" # Optional, defaults to "0.0.0.0"
        protocol: tcp # Optional, defaults to tcp
      - containerPort: 30777
        hostPort: 30777
        listenAddress: "0.0.0.0" # Optional, defaults to "0.0.0.0"
        protocol: tcp # Optional, defaults to tcp
      - containerPort: 31999
        hostPort: 31999
        listenAddress: "0.0.0.0" # Optional, defaults to "0.0.0.0"
        protocol: tcp # Optional, defaults to tcp
  # the workers
  - role: worker
    image: kindest/node:v1.34.3@sha256:08497ee19eace7b4b5348db5c6a1591d7752b164530a36f855cb0f2bdcbadd48
  - role: worker
    image: kindest/node:v1.34.3@sha256:08497ee19eace7b4b5348db5c6a1591d7752b164530a36f855cb0f2bdcbadd48
```



5、创建示例资源

```shell
$ wget 'https://kind.sigs.k8s.io/examples/ingress/usage.yaml'
$ kubectl apply -f usage.yaml
```

成功后效果如下：

```shell
$ kubectl get pod                            
NAME      READY   STATUS    RESTARTS   AGE
bar-app   1/1     Running   0          11m
foo-app   1/1     Running   0          11m
$ kubectl get svc           
NAME          TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
bar-service   ClusterIP   10.96.53.253    <none>        8080/TCP   11m
foo-service   ClusterIP   10.96.104.213   <none>        8080/TCP   11m
kubernetes    ClusterIP   10.96.0.1       <none>        443/TCP    12m
$ kubectl get ingress                        
NAME              CLASS                 HOSTS   ADDRESS                            PORTS   AGE
example-ingress   cloud-provider-kind   *       172.18.0.5,fc00:f853:ccd:e793::5   80      11m
```

注意 bar-app、foo-app pod 可能起不来，需手动导入 pod image：

```shell
$ docker pull registry.k8s.io/e2e-test-images/agnhost:2.39
$ kind load docker-image registry.k8s.io/e2e-test-images/agnhost:2.39 --name=1c2w
```



6、测试验证

```shell
# get the Ingress IP
$ INGRESS_IP=$(kubectl get ingress example-ingress -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
$ echo $INGRESS_IP
172.18.0.5

# should output "foo-app"
$ curl ${INGRESS_IP}/foo
foo-app

# should output "bar-app"
$ curl ${INGRESS_IP}/bar
bar-app
```


7、资源清理

```shell
$ kubectl delete -f usage.yaml
```


## 2.本地运行

```shell
# 运行go程序
$ go run main.go

# 加载镜像
$ docker pull registry.k8s.io/e2e-test-images/agnhost:2.39 && kind load docker-image registry.k8s.io/e2e-test-images/agnhost:2.39 --name=1c2w

# 创建foo资源
$ kubectl apply -f manifests/foo-example.yaml
deployment.apps/foo-deploy created
service/foo-service created

# 查看foo资源
$ kubectl get deploy,service
NAME                         READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/foo-deploy   1/1     1            1           11s

NAME                  TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
service/foo-service   ClusterIP   10.96.180.214   <none>        8080/TCP   11s
service/kubernetes    ClusterIP   10.96.0.1       <none>        443/TCP    5h2m

# 查看自动创建的ingress资源
$ kubectl get ingress
NAME          CLASS                 HOSTS         ADDRESS                            PORTS   AGE
foo-service   cloud-provider-kind   example.com   172.18.0.5,fc00:f853:ccd:e793::5   80      21s

# 配置host: 添加如下内容
# 172.18.0.5 example.com

# 测试
$ curl http://example.com/foo
foo-deploy-95c975f-4t974

# 测试修改service foo-service annotation ingress/http, 查看ingress变化
# 测试删除ingress foo-service, 查看ingress变化
```

## 3.k8s集群运行

```shell
# 构建镜像
$ docker build --progress=plain --no-cache -t ingress-manager-example:v1.0 .

# 加载镜像
$ kind load docker-image ingress-manager-example:v1.0 --name=1c2w

# 创建ingress-manager-example资源
$ kubectl apply -f manifests/ingress-manager-example.yaml

# 创建foo资源
$ kubectl apply -f manifests/foo-example.yaml

# 测试, 参考1.本地运行
```
