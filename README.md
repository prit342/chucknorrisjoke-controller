# Chuck Norris Joke Kubernetes Controller

## Overview

This Kubernetes controller is a fun project that manages Chuck Norris jokes within your cluster by leveraging custom resources of type `ChuckNorris`. It serves as a wrapper around the Chuck Norris API available at [https://api.chucknorris.io/jokes/](https://api.chucknorris.io/jokes/), fetching jokes based on specified categories. The controller watches for `ChuckNorris` custom resources, fetches jokes from the aforementioned API, and updates the resource status with a new joke and its observed generation. Additionally, it handles conditions to inform users about the success or failure of fetching jokes.

## Features

- **Custom Joke Categories**: Supports various joke categories including **animal**, **career**, **celebrity**, **dev**, **explicit**, **fashion**, **food**, **history**, **money**, **movie**, **music**, **political**, **religion**, **science**, **sport**, and **travel**.
- **Automatic Joke Updates**: Watches for `ChuckNorris` custom resource creations and updates, fetching new jokes as specified.
- **Status and Condition Reporting**: Updates custom resources with jokes and fetch status, offering transparency and insights into operation outcomes.

## Prerequisites

- go version v1.20.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.

## Local Testing and Installation

- Create a new Kubernetes cluster, I am using `k3d` for this purpose. You can install `k3d` by following the instructions [here](https://k3d.io/).

```bash
$ k3d cluster create demo --agents 3

INFO[0000] Prep: Network                                
INFO[0000] Created network 'k3d-demo'                   
INFO[0000] Created image volume k3d-demo-images         
INFO[0000] Starting new tools node...                   
INFO[0000] Pulling image 'ghcr.io/k3d-io/k3d-tools:5.6.0' 
INFO[0001] Creating node 'k3d-demo-server-0'            
INFO[0002] Pulling image 'docker.io/rancher/k3s:v1.27.4-k3s1' 
INFO[0003] Starting Node 'k3d-demo-tools'               
```

- Ensure that the cluster is ready.

```bash
$ kubectl cluster-info
Kubernetes control plane is running at https://0.0.0.0:38957
CoreDNS is running at https://0.0.0.0:38957/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy
Metrics-server is running at https://0.0.0.0:38957/api/v1/namespaces/kube-system/services/https:metrics-server:https/proxy

```

- Install the CRDs into the cluster:

```sh

$ make deploy 

```

You should be able to read the new CRDs:

```sh
$ kubectl get customresourcedefinitions.apiextensions.k8s.io chucknorris.jokes.example.com -oyaml

```

- Run the controller locally, which used your Kubeconfig to connect to the `k3d` cluster
You can apply the samples (examples) from the config/sample:

```sh
$ make run
```

- Apply the sample CR using the following command from a differnet terminal windows:
 
```sh

$ kubectl apply -f config/samples/jokes_v1alpha1_chucknorris.yaml
chucknorris.jokes.example.com/chucknorris-dev-joke-example created
chucknorris.jokes.example.com/chucknorris-music-joke-example created

```

- Read the CRs, the joke will be updated in the status field:

```sh
$ kubectl get chucknorris.jokes.example.com/chucknorris-dev-joke-example                                                                                   
NAME                           CATEGORY   JOKE                                                                   AGE
chucknorris-dev-joke-example   dev        Chuck Norris can unit test entire applications with a single assert.   54s

$ kubectl get chucknorris.jokes.example.com/chucknorris-music-joke-example
NAME                             CATEGORY   JOKE                                AGE
chucknorris-music-joke-example   music      Chuck Norris can touch MC Hammer.   2m24s
```
