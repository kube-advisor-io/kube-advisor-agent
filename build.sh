export VERSION=0.1.21 && docker buildx build ./ --platform linux/386,linux/amd64,linux/arm/v7,linux/arm64 --output "type=image,push=true" -t kubeadvisor/kube-advisor-agent:$VERSION