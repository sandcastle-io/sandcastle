#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

SOURCE_DIR="$(cd "$(dirname -- "${BASH_SOURCE[0]}")" && pwd -P)"
ROOT_DIR="$SOURCE_DIR/.."
BIN_DIR="$ROOT_DIR/bin"

# File names and constants
RUNSC="${BIN_DIR}/runsc"
SHIM="${BIN_DIR}/containerd-shim-runsc-v1"
KIND_CLUSTER_NAME="kind-control-plane"
KIND_CONFIG="./hack/kind-config.yaml"

echo "🚀 Starting gVisor setup for Kind (ARM64/Apple Silicon)..."

# 1. Download ARM64 binaries if they don't exist
if [ ! -f "$RUNSC" ]; then
    echo "📥 Downloading runsc (arm64)..."
    curl -L https://storage.googleapis.com/gvisor/releases/release/latest/aarch64/runsc -o $RUNSC
    chmod a+rx $RUNSC
fi

if [ ! -f "$SHIM" ]; then
    echo "📥 Downloading containerd-shim-runsc-v1 (arm64)..."
    curl -L https://storage.googleapis.com/gvisor/releases/release/latest/aarch64/containerd-shim-runsc-v1 -o $SHIM
    chmod a+rx $SHIM
fi

# 2. Generate Kind configuration with extraMounts
echo "📝 Generating $KIND_CONFIG..."
cat <<EOF > $KIND_CONFIG
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  extraMounts:
    - hostPath: $RUNSC
      containerPath: /usr/local/bin/runsc
    - hostPath: $SHIM
      containerPath: /usr/local/bin/containerd-shim-runsc-v1
EOF

# 3. Recreate Kind Cluster
echo "🔄 Deleting old cluster and creating a new one..."
kind delete cluster || true
kind create cluster --config $KIND_CONFIG

# 4. Configure Containerd inside the Kind node
echo "⚙️ Configuring containerd inside the node..."

# Ensure binaries are executable inside the container
docker exec -it $KIND_CLUSTER_NAME chmod +x /usr/local/bin/runsc /usr/local/bin/containerd-shim-runsc-v1

# Create runsc.conf in TOML format (required for recent gVisor versions)
docker exec -it $KIND_CLUSTER_NAME bash -c "cat <<EOF > /etc/containerd/runsc.conf
platform = \"systrap\"
network = \"host\"
EOF"

# Rewrite containerd config.toml to include the runsc runtime
docker exec -it $KIND_CLUSTER_NAME bash -c "cat <<EOF > /etc/containerd/config.toml
version = 2
[plugins.\"io.containerd.grpc.v1.cri\".containerd.runtimes.runsc]
  runtime_type = \"io.containerd.runsc.v1\"
[plugins.\"io.containerd.grpc.v1.cri\".containerd.runtimes.runsc.options]
  TypeUrl = \"io.containerd.runsc.v1.options\"
  ConfigPath = \"/etc/containerd/runsc.conf\"
EOF"

# Restart containerd to apply changes
echo "🔄 Restarting containerd..."
docker exec -it $KIND_CLUSTER_NAME systemctl restart containerd

# 5. Create RuntimeClass in Kubernetes
echo "☸️ Creating Kubernetes RuntimeClass 'gvisor'..."
cat <<EOF | kubectl apply -f -
apiVersion: node.k8s.io/v1
kind: RuntimeClass
metadata:
  name: gvisor
handler: runsc
EOF

kind load docker-image kube-sandcastle-worker

echo "✅ Setup complete! You can now use 'runtimeClassName: gvisor' in your Pod specs."
