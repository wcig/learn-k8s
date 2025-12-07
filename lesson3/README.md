# lesson 3

本章节介绍使用 [code-generator](https://github.com/kubernetes/code-generator) 生成代码。

# 1.初始化项目

```shell
mkdir code-generator
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

```