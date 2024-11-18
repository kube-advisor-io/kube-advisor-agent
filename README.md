![kube-advisor-logo](https://kube-advisor.io/kube-advisor-logo.png)

# kube-advisor-agent
The cluster agent for kube-advisor.io, written in go.

# Emitted metadata
If you want to check which metadata is sent by the agent, check the resources [here](https://github.com/kube-advisor-io/kube-advisor-agent/tree/main/resources).
The agent is watching these k8s resources and they are unmarshalled into these structs before sending to https://kube-advisor.io.

In addition, there are general data providers (at the moment only one, the K8sVersionProvider). See an complete list [here](https://github.com/kube-advisor-io/kube-advisor-agent/blob/main/data_provider.go#L18) that are providing data sent to the platform.




