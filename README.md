![kube-advisor-logo](https://kube-advisor.io/kube-advisor-logo.png)

# kube-advisor-agent
The cluster agent for [kube-advisor.io](https://kube-advisor.io), written in go.

# Emitted metadata
If you want to check which metadata is sent by the agent, check the resources [here](https://github.com/kube-advisor-io/kube-advisor-agent/tree/main/resources).
The agent is watching these k8s resources and they are unmarshalled into these structs before being sent to https://kube-advisor.io.

In addition, there are general data providers that are providing data sent to the platform (at the moment only one, the K8sVersionProvider). See a complete list of them [here](https://github.com/kube-advisor-io/kube-advisor-agent/blob/main/data_provider.go#L18) .




