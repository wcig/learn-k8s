# lesson 2

本章节介绍使用 kubebuilder 创建并部署一个简单示例 CRD 至 k8s 集群。

# 1. 安装kubebuilder

```shell
# download kubebuilder and install locally.
curl -L -o kubebuilder "https://go.kubebuilder.io/dl/latest/$(go env GOOS)/$(go env GOARCH)"
chmod +x kubebuilder && sudo mv kubebuilder /usr/local/bin/
```

# 2. 创建kubebuilder guestbook项目

创建项目：

```shell
mkdir -p kubebuilder/guestbook && cd kubebuilder/guestbook
kubebuilder init --domain my.domain --repo my.domain/guestbook
```

本地调试：

```shell
$ cd kubebuilder/guestbook

# 安装CRD
$ make install

# 运行控制器
$ make run

# 查看CRD
$ kubectl get crd
NAME                          CREATED AT
guestbooks.webapp.my.domain   2025-07-20T12:58:19Z

$ kubectl api-resources | grep guestbooks
guestbooks                                       webapp.my.domain/v1               true         Guestbook

# 创建guestbook实例
$ kubectl apply -f config/samples/webapp_v1_guestbook.yaml          
guestbook.webapp.my.domain/guestbook-sample created

# 查看guestbook实例
$ kubectl get guestbooks                                  
NAME               AGE
guestbook-sample   9s

$ kubectl get guestbooks guestbook-sample -o yaml
apiVersion: webapp.my.domain/v1
kind: Guestbook
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"webapp.my.domain/v1","kind":"Guestbook","metadata":{"annotations":{},"labels":{"app.kubernetes.io/managed-by":"kustomize","app.kubernetes.io/name":"guestbook"},"name":"guestbook-sample","namespace":"default"},"spec":{"foo":"sample"}}
  creationTimestamp: "2025-07-20T13:11:45Z"
  generation: 1
  labels:
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: guestbook
  name: guestbook-sample
  namespace: default
  resourceVersion: "52475"
  uid: 51a5f612-a205-410b-a672-d1555e7f9558
spec:
  foo: sample

# 删除guestbook实例
$ kubectl delete -f config/samples/webapp_v1_guestbook.yaml 
guestbook.webapp.my.domain "guestbook-sample" deleted

# 卸载CRD
$ make uninstall
```

在集群中运行：

```shell
# 构建本地镜像
$ make docker-build IMG=kubebuilder_guestbook_operator:v1.0

# 加载镜像至k8s集群
$ kind load docker-image kubebuilder_guestbook_operator:v1.0 --name=1c2w

# 部署控制器至集群
$ make deploy IMG=kubebuilder_guestbook_operator:v1.0

# 取消部署控制器
$ make undeploy IMG=kubebuilder_guestbook_operator:v1.0
```

查看CRD安装情况和创建guestbook实例参考本地调试部分。

# 3. 创建kubebuilder memcached-operator项目

创建项目：

```shell
$ mkdir -p kubebuilder/memcached-operator && cd kubebuilder/memcached-operator
$ kubebuilder init --domain=example.com --repo my.domain/memcached-operator
$ kubebuilder create api --group cache --version v1alpha1 --kind Memcached
```

参照 [kubebuilder - Getting Started](https://book.kubebuilder.io/getting-started) 示例修改以下文件：
* memcached-operator/api/v1alpha1/memcached_types.go
* memcached-operator/internal/controller/memcached_controller.go
* memcached-operator/config/samples/cache_v1alpha1_memcached.yaml

测试验证：

```shell
# 生成文件 api/v1alpha1/zz_generated.deepcopy.go
$ make generate

# config/crd/bases 目录下生成 CRD manifests，config/rbac 目录生成权限 CR，config/samples 生成示例
$ make manifests

# 安装CRD
$ make install

# 运行控制器
$ make run
```

```shell
# 镜像拉取
$ docker pull memcached:1.6.26-alpine3.19

# 加载镜像到k8s集群
$ kind load docker-image kubebuilder_guestbook_operator:v1.0 --name=1c2w

# 查看CRD
$ kubectl get crd                                                    
NAME                           CREATED AT
memcacheds.cache.example.com   2025-07-20T14:16:43Z
 
$ kubectl api-resources | grep -i mem  
memcacheds                                       cache.example.com/v1alpha1        true         Memcached

# 创建memcached实例
$ kubectl apply -f config/samples/cache_v1alpha1_memcached.yaml

# 查看memcached实例
$ kubectl get memcacheds                                       
NAME               AGE
memcached-sample   3s

# 查看memcached实例创建的deployment、pod
$ kubectl kb get deployment
NAME               READY   UP-TO-DATE   AVAILABLE   AGE
memcached-sample   1/1     1            1           11s
 
$ kubectl get pod       
NAME                                READY   STATUS    RESTARTS   AGE
memcached-sample-68848467f5-4h5wt   1/1     Running   0          17s

# 删除memcached实例
$ kubectl delete -f config/samples/cache_v1alpha1_memcached.yaml

# 卸载CRD
$ make uninstall
```

# 参考
* [kubebuilder - Quick Start](https://book.kubebuilder.io/quick-start)
* [kubebuilder - Getting Started](https://book.kubebuilder.io/getting-started)