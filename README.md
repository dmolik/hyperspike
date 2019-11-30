# HyperSpike

HyperSpike is an attempt at writing a Kubernetes Logging Operator using Rsyslog and Loki. The intention is to automate the drudgery of deploying a logging pipeline (debugging deployments), while also making it easy to scale up from a simple monolith to micro services; read a single instance of Loki to a whole ecosystem.

## Installation

It's simple download the latest manifest and example yaml from the Release page and execute a `kubectl apply`
