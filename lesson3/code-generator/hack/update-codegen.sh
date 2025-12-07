#!/usr/bin/env bash

# Copyright 2017 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

set +x

SCRIPT_DIR="$(dirname "${BASH_SOURCE[0]}")"
CODEGEN_PKG="${SCRIPT_DIR}/../pkg"
source "${SCRIPT_DIR}/kube_codegen.sh"

# 遍历 ${CODEGEN_PKG} 下所有 +k8s:deepcopy-gen 标记的包, 基于 boilerplate.go.txt 模板生成 zz_generated.deepcopy.go
kube::codegen::gen_helpers \
    --boilerplate "${SCRIPT_DIR}/boilerplate.go.txt" \
    "${CODEGEN_PKG}"

# 遍历 ${CODEGEN_PKG} 下所有 +k8s:deepcopy-gen 标记的包, 基于 boilerplate.go.txt 模板生成 zz_generated.register.go
kube::codegen::gen_register \
    --boilerplate "${SCRIPT_DIR}/boilerplate.go.txt" \
    "${CODEGEN_PKG}"

# 扫描指定目录下 +genclient 标记的包, 生成 clientset + informers + listers + applyconfiguration(可选) 代码
## --with-watch: 给生成的 clientset 加上 Watch() 方法, 需要 +genclient:watch 标记或默认开启
## --with-applyconfig: 生成 applyconfiguration 代码
## --output-dir: 生成文件目录
## --output-pkg: 生成的 go 包 import 路径前缀
## --boilerplate: go文件模板, 一般设置为 copyright
## ${PROJECT_ROOT}/pkg/apis: 扫描目录下所有 +genclient 标记的包, 要求目录文件为 <group>/<version>/<types>.go
kube::codegen::gen_client \
    --with-watch \
    --with-applyconfig \
    --output-dir "${CODEGEN_PKG}/generated" \
    --output-pkg "${CODEGEN_PKG}/generated" \
    --boilerplate "${SCRIPT_DIR}/boilerplate.go.txt" \
    "${CODEGEN_PKG}/apis"
