# DTree

## Materials 
- Links Kubernetes:
    + Kubernetes JSON file: https://github.com/instrumenta/kubernetes-json-schema/blob/master/v1.18.0/_definitions.json
    + Kubernetes Namespace: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/
- Links Cheops:
    + Cheops Git Repo: https://gitlab.inria.fr/discovery/cheops 
    + Cheops latest Paper: https://hal.inria.fr/hal-03770492/document
- User configs:
    + YAML for a pod: https://kubernetes.io/docs/concepts/configuration/secret/#using-secrets-as-files-from-a-pod
    + YAML for a deployment: https://kubernetes.io/docs/concepts/workloads/controllers/deployment/

## Example

```
go run main.go make -c ./assets/config.yaml -t ./assets/template.json
```

Reference template: [template.json](./assets/template.json)
Reference config: [config.yaml](./assets/config.yaml)

**Example queried config:**

```
apiVersion: v1
kind: Pod
metadata:
  name: mypod
spec:
  containers:
  - name: mypod
    image: redis
    volumeMounts:
    - name: foo
      mountPath: "/etc/foo"
      readOnly: true
  volumes:
  - name: foo
    secret:
      secretName: mysecret
      optional: false
```

### Step 1: Pod

Found template from kind: `io.k8s.api.core.v1.Pod`

Required field in template: `None`

Existing references in template:
  1. metadata: `"#/definitions/io.k8s.apimachinery.pkg.apis.meta.v1.ObjectMeta"`,
  2. spec: `"#/definitions/io.k8s.api.core.v1.PodSpec"`,
  3. status: `"#/definitions/io.k8s.api.core.v1.PodStatus"`,

Found in Config: `1. and 2.`

### Step 2: Pod / ObjectMeta

Referenced template: `io.k8s.apimachinery.pkg.apis.meta.v1.ObjectMeta`

Required field in template: `None`

Existing references in template:
  1. creationTimestamp: `"#/definitions/io.k8s.apimachinery.pkg.apis.meta.v1.Time"`,
  2. deletionTimestamp: `"#/definitions/io.k8s.apimachinery.pkg.apis.meta.v1.Time"`,
  3. managedFields: `"#/definitions/io.k8s.apimachinery.pkg.apis.meta.v1.ManagedFieldsEntry"`,
  4. ownerReferences: `"#/definitions/io.k8s.apimachinery.pkg.apis.meta.v1.OwnerReference"`

Found in Config: `None`

### Step 3: Pod / PodSpec

Referenced template: `io.k8s.api.core.v1.PodSpec`

Required field in template: `"containers"`

Existing references in template:
  1. volumes > items: `"#/definitions/io.k8s.api.core.v1.Volume"`
  2. containers > items: `"#/definitions/io.k8s.api.core.v1.Container"`
  3. ...

Found in Config: `1. and 2.`

### Step 4: Pod / PodSpec / Volume

Referenced template: `io.k8s.api.core.v1.Volume`

Required field in template: `"name"`

Existing references in template:
  1. secret: `io.k8s.api.core.v1.SecretVolumeSource`
  2. ...

Found in Config: `1.`


### Step 5: Pod / PodSpec / Volume / SecretVolumeSource

Referenced template: `io.k8s.api.core.v1.SecretVolumeSource`

Required field in template: `None`

Existing references in template:
  1. ...

Found in Config: `None`

### Step 6: Pod / PodSpec / Container

Referenced template: `io.k8s.api.core.v1.Container`

Required field in template: `"name"`

Existing references in template:
  1. volumeMounts > items: `io.k8s.api.core.v1.VolumeMount`
  2. ...

Found in Config: `1.`

### Step 7: Pod / PodSpec / Container / VolumeMount

Referenced template: `io.k8s.api.core.v1.VolumeMount`

Required field in template: `"name", "mountPath"`

Existing references in template:
 1. ...

Found in Config: `None`

### STOP

results:
 * `io.k8s.api.core.v1.PodSpec.properties.containers`
 * `io.k8s.api.core.v1.Container.properties.name`
 * `io.k8s.api.core.v1.Volume.properties.name`
 * `io.k8s.api.core.v1.VolumeMount.properties.name`
 * `io.k8s.api.core.v1.VolumeMount.properties.mountPath`
