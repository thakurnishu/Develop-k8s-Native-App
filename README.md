# Develop-k8s-Native-App
Basic of Client-go and other library to Create K8S Native Application

## Library Used to Develop K8S Native Application

| Library            | Module                | Used                                            |
| :----------------- | :-------------------- | :---------------------------------------------- |
| **`client-go`**    | `k8s.io/client-go`    | Go clients for talking to a kubernetes cluster. |
| **`api`**          | `k8s.io/api`          | Schema of the external API types that are served by the Kubernetes API server |
| **`apimachinery`** | `k8s.io/apimachinery` | Scheme, typing, encoding, decoding, and conversion packages for Kubernetes and Kubernetes-like API objects | 


### Kubernetes Objects/Resources in GO

- In Go code, when a struct implements the `runtime.Object` `interface` from the **`apimachinery/pkg/runtime`** package, we can classify that Go struct as a **Kubernetes Object**
- Alternatively, we can specify that a Go struct is considered a **Kubernetes Object** if it implements the `DeepCopyObject` **`method`** from the `runtime.Object` **`interface`** and the `SetGroupVersionKind` and `GroupVersionKind` **`method`** from `schema.ObjectKind` **`interface`**

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
![Screenshot from 2024-01-12 02-29-54](https://github.com/thakurnishu/Develop-k8s-Native-App/assets/90508814/c506fc43-50d4-48a7-8694-32c15cfd255b)
