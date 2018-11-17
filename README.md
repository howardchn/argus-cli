# Argus Commandline Utility

This is a command utility for argus project that help to setup kubernetes monitoring easily.

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

Require Values:
* accessId:

