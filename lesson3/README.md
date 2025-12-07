# lesson 3

本章节介绍使用 [code-generator](https://github.com/kubernetes/code-generator) 生成代码。

# 1.初始化项目

```shell
mkdir code-generator && cd code-generator 
go mod init code-generator
```

# 2.准备类型文件和脚本

这里基于官方 code-generator 项目简单修改，这里我们使用 K8S 1.31 版本。

```shell
# 切换其他目录
git clone https://github.com/kubernetes/code-generator.git
cd code-generator
# 切换至自己需要的版本
git checkout release-1.31
```

切换本地 code-generator 目录，准备类型文件。

```shell
cd code-generator
# 存放脚本文件和模板文件
mkdir hack 
# 存放代码文件
mkdir pkg
# 存放类型文件
mkdir -p pkg/apis/samplecontroller/v1alpha1
# 存放生成代码文件
mkdir pkg/generated
```

准备类型文件 foo_types.go 和 doc.go，拷贝至 pkg/apis/samplecontroller/v1alpha1 目录，从 kubernetes/code-generator 项目下拷贝 hack 目录至本地，简单修改 update-codegen.sh，准备 tools.go 文件并拷贝至 hack 目录。此时项目结构为：

```shell
➜  code-generator git:(master) ✗ tree .                  
.
├── go.mod
├── go.sum
├── hack
│   ├── boilerplate.go.txt
│   ├── kube_codegen.sh
│   ├── tools.go
│   └── update-codegen.sh
├── main.go
└── pkg
    └── apis
        └── samplecontroller
            └── v1alpha1
                ├── doc.go
                └── foo_types.go
```

添加依赖并运行生成文件脚本。

```shell
$ cd code-generator
$ go get k8s.io/apimachinery@v0.31.0
$ go get k8s.io/code-generator@v0.31.0
$ ./hack/update-codegen.sh
Generating deepcopy code for 1 targets
Generating register code for 1 targets
Generating applyconfig code for 1 targets
Generating client code for 1 targets
Generating lister code for 1 targets
Generating informer code for 1 targets
```

脚本运行后项目目录结构为：

```shell
➜  code-generator git:(master) ✗ tree .
.
├── examples
│   ├── crd-status-subresource.yaml
│   ├── crd.yaml
│   └── example-foo.yaml
├── go.mod
├── go.sum
├── hack
│   ├── boilerplate.go.txt
│   ├── kube_codegen.sh
│   ├── tools.go
│   └── update-codegen.sh
├── main.go
└── pkg
    ├── apis
    │   └── samplecontroller
    │       └── v1alpha1
    │           ├── doc.go
    │           ├── foo_types.go
    │           ├── zz_generated.deepcopy.go
    │           └── zz_generated.register.go
    └── generated
        ├── applyconfiguration
        │   ├── internal
        │   │   └── internal.go
        │   ├── samplecontroller
        │   │   └── v1alpha1
        │   │       ├── foo.go
        │   │       ├── foospec.go
        │   │       └── foostatus.go
        │   └── utils.go
        ├── clientset
        │   └── versioned
        │       ├── clientset.go
        │       ├── fake
        │       │   ├── clientset_generated.go
        │       │   ├── doc.go
        │       │   └── register.go
        │       ├── scheme
        │       │   ├── doc.go
        │       │   └── register.go
        │       └── typed
        │           └── samplecontroller
        │               └── v1alpha1
        │                   ├── doc.go
        │                   ├── fake
        │                   │   ├── doc.go
        │                   │   ├── fake_foo.go
        │                   │   └── fake_samplecontroller_client.go
        │                   ├── foo.go
        │                   ├── generated_expansion.go
        │                   └── samplecontroller_client.go
        ├── informers
        │   └── externalversions
        │       ├── factory.go
        │       ├── generic.go
        │       ├── internalinterfaces
        │       │   └── factory_interfaces.go
        │       └── samplecontroller
        │           ├── interface.go
        │           └── v1alpha1
        │               ├── foo.go
        │               └── interface.go
        └── listers
            └── samplecontroller
                └── v1alpha1
                    ├── expansion_generated.go
                    └── foo.go

28 directories, 40 files
```

# 3.测试项目

终端一运行项目：

```shell
➜  code-generator git:(master) ✗ go run main.go -kubeconfig=/Users/yangbo/.kube/config
I1207 18:20:54.995914   77793 main.go:87] Starting Foo controller
Sync/Add/Update for foo default/example-foo
Sync/Add/Update for foo default/example-foo
Sync/Add/Update for foo default/example-foo
I1207 18:21:31.779394   77793 main.go:56] Foo default/example-foo does not exist anymore
```

终端二创建修改和删除Foo对象：

```shell
➜  code-generator git:(master) ✗ k apply -f examples/example-foo.yaml
➜  code-generator git:(master) ✗ k get foo                           
NAME          AGE
example-foo   10s
➜  code-generator git:(master) ✗ k edit foo example-foo        
foo.samplecontroller.k8s.io/example-foo edited
➜  code-generator git:(master) ✗ k delete foo example-foo
foo.samplecontroller.k8s.io "example-foo" deleted
```
