# Develop-k8s-Native-App
Basic of Client-go and other library to Create K8S Native Application

## Library Used to Develop K8S Native Application

| Library            | Module                | Used                                            |
| :----------------- | :-------------------- | :---------------------------------------------- |
| **`client-go`**    | `k8s.io/client-go`    | Go clients for talking to a kubernetes cluster. |
| **`api`**          | `k8s.io/api`          | Schema of the external API types that are served by the Kubernetes API server |
| **`apimachinery`** | `k8s.io/apimachinery` | Scheme, typing, encoding, decoding, and conversion packages for Kubernetes and Kubernetes-like API objects | 


### Kubernetes Objects/Resources in GO

- In Go code, when a struct implements the [`runtime.Object`](https://pkg.go.dev/k8s.io/apimachinery/pkg/runtime#Object) `interface` from the **`apimachinery/pkg/runtime`** package, we can classify that Go struct as a **Kubernetes Object**
- Alternatively, we can specify that a Go struct is considered a **Kubernetes Object** if it implements the `DeepCopyObject` **`method`** from the `runtime.Object` **`interface`** and the `SetGroupVersionKind` and `GroupVersionKind` **`method`** from [`schema.ObjectKind`](https://pkg.go.dev/k8s.io/apimachinery@v0.29.0/pkg/runtime/schema#ObjectKind) **`interface`**

```txt
    # Any k8s Objects
        TypeMeta
            Kind
            APIVersion
        
        ObjectMeta
            Name
            Namespace
            Labels
            Annotations
            ResourceVersion
            ...
        Spec
            specs of resource
        Status
```

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: lister-deployment
  labels:
    app: lister
spec:
  replicas: 1
  selector:
    matchLabels:
      app: lister
  template:
    metadata:
      labels:
        app: lister
    spec:
      containers:
      - name: lister
        image: mahakal0510/lister:0.1.2
        resources:
          limits:
            memory: "256Mi"  
            cpu: "100m"      
          requests:
            memory: "128Mi"  
            cpu: "50m"    
```
### Informer, Cache and Queue 

#### **What is custom controller ?**
- It can be any application which listen to specfic resource in kubernetes cluster and perform action based on
what happening with resource in k8s.
- So, for this we need something that watch kubernetes cluster (created, deleted, updated) and tell application to perform further action [`watch` vs `Informer`].

#### **Why don't we use Watch ?**
- *watch* verb is going to query `api-server` repeatedly (created, deleted, updated), which is going to increase load onto `api-server` and use process are going to slow due to use of *watch* verb

#### **Why we use Informer ?**
- It still use *watch* internally but it very efficiently leverages in-memory store, which is going to reduce load onto `api-server` 
- It is very important component while writing kubernetes controller
![Screenshot from 2024-01-12 11-27-52](https://github.com/thakurnishu/Develop-k8s-Native-App/assets/90508814/8dbcb199-57b5-4222-b7ea-e1e90637b0c3)


### Kubernetes API-Machinery
- As in case of **Kubernetes Object**, we discuss what are all property and behaviour particular GO type/struct should have to be called *kubernetes Object*.
#### Kind
- In case of **API-Machinery** the term analogous to kubernetes type is **`Kind`** and are represented in *CamelCase*. Eg:
**`Kind`** for pod -> `Pod`, deployment -> `Deployment`. Kind are not plural.
- **`Kind`** are grouped eg: `apps` group, `authorization` group and they are also versioned.
- `Deployment` are from `apps` group and `v1` version, this is how [**`GroupVersionKind`**](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#GroupVersionKind) comes into picture, here Group is *apps*, Version is *v1* and Kind is *Deployment*
- [**`GroupVersionKind`**](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#GroupVersionKind) is a way to present a particualr *kubernetes Object* using *`API-Machinery`* 
- There is not one-to-one relationship b/w *Kind* and *http-endpoint*.

#### Resource
- We have another contruct similar to `Kind` that is **`resource`**. But *resource* are plural and represented in *lower-case*.
- For example, *resource* for deployment -> `deployments`, pod -> `pods`
- Like *Kind*, *resource* are also grouped and versioned. Eg: `deployments` resource are from `apps` group and `v1` version, this is how [**`GroupVersionResource`**](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#GroupVersionResource) comes into picture.
- Unlike *Kind*, *resource* are mapped to specific *http-endpoint*. For example for `deployments` resource http-endpoint
  - *`https://<api-server-ip>:<api-server-port>/apis/apps/v1/namespaces/default/deployments`*
- We can break it as 
  - *`https://<api-server-ip>:<api-server-port>apis/<group>/<version>/namespaces/<namespace-name>/<resource>`* 
- Note *For `Core` api group http-endpoint is bit different*. Eg: for `services` reource *http-endpoint*
  - `https://<api-server-ip>:<api-server-port>/api/v1/namespaces/default/services`
- We can also break it as 
  - `https://<api-server-ip>:<api-server-port>/api/<version>/namespaces/<namespace-name>/<resource>`
##### From *`resource`* we can very easily figure out *http-endpoint*. 

#### RestMapping
- If we ever need to convert GroupVersionKind to GroupVersionResource, or any orther operation we can use [`RestMapper`](https://pkg.go.dev/k8s.io/apimachinery/pkg/api/meta#RESTMapper) `interface` to run various different method 

#### Scheme
- In [`scheme`](https://pkg.go.dev/k8s.io/apimachinery/pkg/runtime#Scheme) struct, we have a method name [*`ObjectKinds`*](https://pkg.go.dev/k8s.io/apimachinery/pkg/runtime#Scheme.ObjectKinds) which can be use to retrive all possible group,version, kind (`GroupVersionKind` struct) by providing *Kubernetes Object*.
- *`ObjectKinds`* only going to work if passed object is already registered.
- We can register object by using [*`AddKnownTypes`*](https://pkg.go.dev/k8s.io/apimachinery@v0.29.0/pkg/runtime#Scheme.AddKnownTypes) method in [`scheme`](https://pkg.go.dev/k8s.io/apimachinery/pkg/runtime#Scheme) struct.


![Screenshot from 2024-01-12 14-33-32](https://github.com/thakurnishu/Develop-k8s-Native-App/assets/90508814/3f411e72-9770-417e-9fe6-3649c7392c74)
