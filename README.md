# cicd
实现cicd的operator

### kubeadm 安装k8s

####  k8s 初始化
```shell
kubeadm init --control-plane-endpoint=172.31.42.220:6443 --apiserver-advertise-address=172.31.42.220 --image-repository registry.aliyuncs.com/google_containers --apiserver-cert-extra-sans="192.168.0.0,69.235.141.221,172.31.42.220" --service-cidr=192.168.0.0/16 --pod-network-cidr=10.244.0.0/16
```

#### 更新k8s认证IP
```shell
rm /etc/kubernetes/pki/apiserver.*
kubeadm init phase certs apiserver --apiserver-advertise-address ${原来的apiserver地址}   --apiserver-cert-extra-sans ${加入认证的地址}
kubeadm alpha certs renew admin.conf
kubectl delete po ${apiserverPod}
```

#### k8s admin安装报错解决方法
```shell
#初始化之后 kubectl get cs 可能会出现scheduler unhealthy
vim /etc/kubernetes/manifests/kube-scheduler.yaml
# 注释 #    - --port=0
```

#### crictl  
```shell
cat /etc/crictl.yaml 
runtime-endpoint: unix:///run/containerd/containerd.sock  ## socket 文件
image-endpoint: unix:///run/containerd/containerd.sock
timeout: 10  ## 超时时间
debug: true  ## 开启debug 
```

####  containerd
```shell
yum install -y containerd
mkdir -p /etc/containerd
containerd config default | sudo tee /etc/containerd/config.toml
sed -i 's#https://registry-1.docker.io#https://docker.mirrors.ustc.edu.cn#g' /etc/containerd/config.toml
vim /etc/containerd/config.toml
[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc]
  ...
  [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc.options]
    SystemdCgroup = true 
    
    ...
    
    [plugins."io.containerd.grpc.v1.cri".registry]
      [plugins."io.containerd.grpc.v1.cri".registry.mirrors]
        [plugins."io.containerd.grpc.v1.cri".registry.mirrors."docker.io"]
          endpoint = ["https://registry-1.docker.io"]
    sandbox_image = "docker.io/juestnow/pause-amd64:3.2"

systemctl enable --now containerd
```

### 使用code-generator 生成crd所需代码
```shell
# 目录层级
.
├──  crd
│   ├── crd.yaml
│   ├── example.yaml
│   └── task.yaml
├── go.mod
├── go.sum
├── LICENSE
├── main.go
├── pkg
│   ├── apis
│   │   └── task
│   │       └── v1alpha1
│   │           ├── doc.go
│   │           ├── register.go
│   │           └── types.go
│   ├── client
│   └── controllers
└── README.md

## 将pkg main.go  go.mod 文件移动至$GOPATH/src/github.com/hongshixing/cicd 无此目录请mkdir
cp  -r pkg main.go  go.mod $GOPATH/src/github.com/hongshixing/cicd
cd $GOPATH/src/github.com/hongshixing/cicd 
go mod tidy
```

#### code-generator 使用
[code-generator github链接](https://github.com/kubernetes/code-generator)

```shell
## 下载源码 此处k8s版本使用的是1.23.4 对应code-generator版本是v0.23.4
cd $GOPATH/src/
wget https://github.com/kubernetes/code-generator/archive/refs/tags/v0.23.4.zip
unzip v0.23.4.zip

## 生成代码
cd $GOPATH/src/github.com/hongshixing/cicd

/home/hsx/go/src/code-generator-0.23.4/generate-groups.sh all github.com/hongshixing/cicd/pkg/client github.com/hongshixing/cicd/pkg/apis task:v1alpha1

```
#### 生成代码后的文件结构是
```shell
.
├── apis
│   └── task
│       └── v1alpha1
│           ├── doc.go
│           ├── register.go
│           ├── types.go
│           └── zz_generated.deepcopy.go
├── client
│   ├── clientset
│   │   └── versioned
│   │       ├── clientset.go
│   │       ├── doc.go
│   │       ├── fake
│   │       │   ├── clientset_generated.go
│   │       │   ├── doc.go
│   │       │   └── register.go
│   │       ├── scheme
│   │       │   ├── doc.go
│   │       │   └── register.go
│   │       └── typed
│   │           └── task
│   │               └── v1alpha1
│   │                   ├── doc.go
│   │                   ├── fake
│   │                   │   ├── doc.go
│   │                   │   ├── fake_task_client.go
│   │                   │   └── fake_task.go
│   │                   ├── generated_expansion.go
│   │                   ├── task_client.go
│   │                   └── task.go
│   ├── informers
│   │   └── externalversions
│   │       ├── factory.go
│   │       ├── generic.go
│   │       ├── internalinterfaces
│   │       │   └── factory_interfaces.go
│   │       └── task
│   │           ├── interface.go
│   │           └── v1alpha1
│   │               ├── interface.go
│   │               └── task.go
│   └── listers
│       └── task
│           └── v1alpha1
│               ├── expansion_generated.go
│               └── task.go
└── controllers
    └── controller.go


## 将pkg目录直接cp回项目目录中 
cp -r $GOPATH/src/github.com/hongshixing/cicd/pkg   $PROJECT  # $PROJECT是项目所在目录
```




#### 几点注意事项
- 在生成的代码文件`$PROJECT/pkg/apis/task/v1alpha1/register.go`中
```go

// 这里默认生成的Version是v1 需要根据实际修改

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: "api.gok8s.fun", Version: "v1alpha1"} 

```
- 在生成的代码文件`$PROJECT/pkg/client/clientset/versioned/scheme/register.go`中
```go
func init() {
	// 这里默认生成的Version是v1 需要根据实际修改
	v1.AddToGroupVersion(Scheme, schema.GroupVersion{Version: "v1alpha1"}) 
	utilruntime.Must(AddToScheme(Scheme))
}

```
- 在生成的代码文件`$PROJECT/pkg/client/clientset/versioned/fake/register.go`中
```go
func init() {
	// 这里默认生成的Version是v1 需要根据实际修改
	v1.AddToGroupVersion(Scheme, schema.GroupVersion{Version: "v1alpha1"}) 
	utilruntime.Must(AddToScheme(Scheme))
}
```

#### 控制pod的顺序执行
- 使用Downward Api [官方链接](https://kubernetes.io/zh/docs/tasks/inject-data-application/downward-api-volume-expose-pod-information/#the-downward-api)

```shell
apiVersion: v1
kind: Pod
metadata:
  name: kubernetes-downwardapi-volume-example
  labels:
    zone: us-est-coast
    cluster: test-cluster1
    rack: rack-22
  annotations:
    build: two
    builder: john-doe
spec:
  containers:
    - name: client-container
      image: k8s.gcr.io/busybox
      command: ["sh", "-c"]
      args:
      - while true; do
          if [[ -e /etc/podinfo/labels ]]; then
            echo -en '\n\n'; cat /etc/podinfo/labels; fi;
          if [[ -e /etc/podinfo/annotations ]]; then
            echo -en '\n\n'; cat /etc/podinfo/annotations; fi;
          sleep 5;
        done;
      volumeMounts:
        - name: podinfo
          mountPath: /etc/podinfo
  volumes:
    - name: podinfo   
      downwardAPI:  ## 使用downward api 将pod本身的变量挂载至容器内
        items:
          - path: "labels"
            fieldRef:
              fieldPath: metadata.labels
          - path: "annotations"
            fieldRef:
              fieldPath: metadata.annotations
```
Downward Api使用中的几个注意事项：
下面这些信息可以通过环境变量和 downwardAPI 卷提供给容器：

- 能通过 fieldRef 获得的：
  - metadata.name - Pod 名称
  - metadata.namespace - Pod 名字空间
  - metadata.uid - Pod 的 UID
  - metadata.labels['<KEY>'] - Pod 标签 <KEY> 的值 (例如, metadata.labels['mylabel']）
  - metadata.annotations['<KEY>'] - Pod 的注解 <KEY> 的值（例如, metadata.annotations['myannotation']）


- 能通过 resourceFieldRef 获得的：
  - 容器的 CPU 约束值
  - 容器的 CPU 请求值
  - 容器的内存约束值
  - 容器的内存请求值
  - 容器的巨页限制值（前提是启用了 DownwardAPIHugePages 特性门控）
  - 容器的巨页请求值（前提是启用了 DownwardAPIHugePages 特性门控）
  - 容器的临时存储约束值
  - 容器的临时存储请求值


- 此外，以下信息可通过 downwardAPI 卷从 fieldRef 获得：
  - metadata.labels - Pod 的所有标签，以 label-key="escaped-label-value" 格式显示，每行显示一个标签
  - metadata.annotations - Pod 的所有注解，以 annotation-key="escaped-annotation-value" 格式显示，每行显示一个标签
   

- 以下信息可通过环境变量获得：
  - status.podIP - 节点 IP
  - spec.serviceAccountName - Pod 服务帐号名称, 版本要求 v1.4.0-alpha.3
  - spec.nodeName - 节点名称, 版本要求 v1.4.0-alpha.3
  - status.hostIP - 节点 IP, 版本要求 v1.7.0-alpha.1

    