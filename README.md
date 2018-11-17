# Argus Commandline Utility

This is a command utility for [argus project](https://logicmonitor.github.io/k8s-argus/) that help to setup kubernetes monitoring easily.

## Features
The original argus project is only allowing to install resources on your kubernetes environment and sync the resources into Santaba. It is complicated to uninstall it from your env and Santaba. Here are the details:

1. Device groups and resources within that reflect to your k8s resource structure
2. Collector groups and collectors based on the replicas number you set
3. Argus related resources (pods, crd, deployments, statefulset etc.)
4. CollectorSet-Controller related resources (pods, collectorsets etc.)

So this utility helps to clean those resources for you.

## Configurations
An example:
```
$ ./argus-cli --accessId="[Access ID]" --accessKey="[Access Key]" --clusterName="[Cluster Name]" --account="[Company Name]" --parentId=[Parent Group Id]
```

### Require Values:
* **accessId:** the access id that is generate with your own santaba account from UI portal
* **accessKey:** the access key that is generate with your own santaba account from UI portal
* **clusterName:** the cluster name that is used for creating your k8s monitoring. The same configuration as _clusterName_ in [argus configuration](https://logicmonitor.github.io/k8s-argus/docs/configuration/)
* **account:** the company account for your logicmonitor portal. e.g. _mycompany.logicmonitor.com_, _mycompany_ is the value

### Optional Values:
* **parentId:** the parent group id for the k8s cluster node is placed in. Default value is 1 that means the root group

## Build Executable
Feel free to clone this project and build for your specific OS.

### Precondition:
Setup the go environment on your machine. Refer the [installation instruction](https://golang.org/doc/install).

### For macOS (darwin)
```bash
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o argus-cli-darwin main.go
```

### For Linux
```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o argus-cli-linux main.go
```

### For Windows
```bash
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o argus-cli-win main.go
```

### Use Makefile
```bash
# cd to the project
make
```

## Contact
* [Email](mailto:howardch@outlook.com)
* [Linkedin](https://www.linkedin.com/in/howard-chen-328493142/)

