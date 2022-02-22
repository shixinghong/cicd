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