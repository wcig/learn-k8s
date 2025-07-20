# lesson 2

# 1. 安装kubebuilder

```shell
# download kubebuilder and install locally.
curl -L -o kubebuilder "https://go.kubebuilder.io/dl/latest/$(go env GOOS)/$(go env GOARCH)"
chmod +x kubebuilder && sudo mv kubebuilder /usr/local/bin/
```

# 2. 创建kubebuilder项目

```shell
mkdir -p kubebuilder/guestbook && cd kubebuilder/guestbook
kubebuilder init --domain my.domain --repo my.domain/guestbook
```

# 参考
* [kubebuilder - Quick Start](https://book.kubebuilder.io/quick-start)