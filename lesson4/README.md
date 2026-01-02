# lesson 4

本章节介绍使用 [controller-runtime](https://github.com/kubernetes-sigs/controller-runtime) 来实现 Foo 自定义资源 Controller。

Controller-runtime 相比传统的 Informer 方式简化了代码，并增加实现了 validation webhook 和 mutation webhook。
