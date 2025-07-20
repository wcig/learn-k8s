# lesson 1

本章节介绍如何在本地使用 kind 创建 k8s 集群，并部署一个简单的 go http 服务。

## 1. 安装kind

MacOS 版本安装如下

```shell
# 1.Maos
# For Intel Macs
[ $(uname -m) = x86_64 ] && curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.29.0/kind-darwin-amd64
# For M1 / ARM Macs
[ $(uname -m) = arm64 ] && curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.29.0/kind-darwin-arm64
chmod +x ./kind
mv ./kind $GOPATH/bin/kind
```

## 2. 创建k8s集群

```shell
kind create cluster --config k8s/kind-config-1c2w.yaml --name 1c2w
```

## 3. 部署go服务

```shell
# 拉取基础镜像
docker pull golang:1.24.5-alpine3.22
docker pull alpine:3.22

# 构建镜像
docker build -t goapp:v1.0 goapp/Dockerfile

# 加载镜像至k8s集群
kind load docker-image goapp:v1.0 --name=1c2w

# 部署go服务
kubectl apply -f k8s/goapp.yaml
```

## 4. 查看 go 服务

```shell
$ kubectl get po,svc                   
NAME                         READY   STATUS    RESTARTS   AGE
pod/goapp-77778d7945-p2qgb   1/1     Running   0          21m

NAME                    TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)        AGE
service/goapp-service   NodePort    10.96.182.51   <none>        80:30080/TCP   21m
service/kubernetes      ClusterIP   10.96.0.1      <none>        443/TCP        25m

$ curl -i 'http://localhost:30080/ready' 
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Sun, 20 Jul 2025 04:32:42 GMT
Content-Length: 16

{"message":"ok"}

$ kubectl logs -f goapp-77778d7945-p2qgb
2025/07/20 04:25:26 >> app run
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /ready                    --> main.readyHandler (3 handlers)
[GIN-debug] [WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.
Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.
[GIN-debug] Listening and serving HTTP on :8080
[GIN] 2025/07/20 - 04:26:01 | 200 |      59.875µs |      172.18.0.4 | GET      "/ready"

```

# 参考
* [Kind - Quick Start](https://kind.sigs.k8s.io/docs/user/quick-start/#installing-from-release-binaries)
* [GitHub - HanFa learn-k8s](https://github.com/HanFa/learn-k8s)
