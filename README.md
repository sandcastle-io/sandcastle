# 🏰 Kube-Sandcastle

High-Performance, Secure, and Kueue-Native Runtime for AI Agents.

Kube-Sandcastle is an infrastructure layer designed to bridge the 
"Agent Execution Gap." It enables the execution of untrusted, AI-
generated code with the security of a hardened sandbox and the 
speed of a system call, remaining fully native to the Kubernetes 
resource model via Kueue.

## 🚀 The Problem: "The Agent Execution Gap"

AI Agents (AutoGPT, OpenDevin, LangChain) generate code that needs 
to be executed instantly. Current Kubernetes patterns fail to meet 
these demands:

**Security Risk**: Running os.system("rm -rf /") in a standard container 
is a disaster, leading to potential container escapes and cluster 
compromise.

**High Latency**: Spinning up a new Kubernetes Job takes 2–5 seconds. 
Agents require sub-100ms response times to maintain an "intelligent" 
iterative loop.

**Resource Waste**: Maintaining thousands of idle Pods for occasional 
agent tasks results in massive infrastructure overhead and cost.

## 🏗️ The Architecture: "Sandbox-as-a-Service"

Kube-Sandcastle treats compute resources like a financial ledger:

**Wholesale (Kueue)**: We reserve large compute quotas via Kueue Workloads 
to maintain Warm Worker Pools—pre-initialized environments ready for 
action.

**Retail (Sandbox Proxy)**: A high-speed Manager "slices" these reserved 
resources into thousands of ephemeral micro-sandboxes on demand.

**Hard Isolation**: Every task is wrapped in nsjail—a lightweight, 
high-performance sandbox utilizing Linux Namespaces, cgroups, and 
strict seccomp profiles to ensure zero-leak execution.

## ✨ Key Features

⚡ **Sub-100ms** Latency: Bypasses the Kubernetes Control Plane for task execution 
  by leveraging pre-warmed runners.

🔒 **Kernel-Level Isolation**: Powered by nsjail to prevent unauthorized system 
   calls and filesystem access.

♻️ **Instant State Reset**: Every execution starts with a pristine memory and 
   filesystem state using Copy-on-Write (CoW) mechanisms.

📊 **Kueue-Native Accounting**: Real-time reporting of precise CPU/RAM 
   usage (syscall.Getrusage) back to Kueue for granular quota management.

🔌 **Kubernetes API Proxy Integration**: Seamless communication between the CLI 
    and Workers via the native K8s Proxy subresource, eliminating the need 
    for complex Ingress setups.

## 🛠️ Tech Stack

**Language**: Go (Control Plane & Worker logic)

**Isolation**: nsjail (Namespaces, cgroups, seccomp)

**Orchestration**: Kubernetes & Kueue

**Communication**: REST/gRPC via K8s API Proxy

## Why this works for your startup:

**Market Fit**: It targets the "AI Agent" niche specifically, which is currently 
the fastest-growing segment in AI infra.

**Scalability**: It uses Kueue, which is the industry standard for batch/AI 
workloads in K8s.

**Technical Sophistication**: Moving from "just a pod" to "sandboxed process 
inside a managed pool" shows high-level engineering.