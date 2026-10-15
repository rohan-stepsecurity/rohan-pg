# Public Ubuntu base image from Docker Hub
FROM docker.io/library/ubuntu:22.04

ARG DEBIAN_FRONTEND=noninteractive
ARG HOOKS_VERSION=v1.4.0

# Base tooling that CodeBuild and most builds expect.
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        ca-certificates \
        curl \
        wget \
        git \
        unzip \
        jq \
        sudo \
        gnupg \
        lsb-release && \
    rm -rf /var/lib/apt/lists/*

# ------------------------------------------------------------------
# 1. StepSecurity GitHub Actions runner job hooks
# ------------------------------------------------------------------
RUN mkdir -p /opt/step-security && \
    curl -fsSL "https://github.com/step-security/gha-runner-job-hooks/releases/download/${HOOKS_VERSION}/pre.js" \
        -o /opt/step-security/pre-job-hook.sh && \
    curl -fsSL "https://github.com/step-security/gha-runner-job-hooks/releases/download/${HOOKS_VERSION}/post.js" \
        -o /opt/step-security/post-job-hook.sh && \
    chmod +x /opt/step-security/pre-job-hook.sh /opt/step-security/post-job-hook.sh

# Picked up by the GitHub Actions runner if this image is used as a runner
# (e.g. CodeBuild-hosted GitHub Actions runners). Harmless otherwise.
ENV ACTIONS_RUNNER_HOOK_JOB_STARTED=/opt/step-security/pre-job-hook.sh \
    ACTIONS_RUNNER_HOOK_JOB_COMPLETED=/opt/step-security/post-job-hook.sh

# ------------------------------------------------------------------
# 2. Harden-Runner agent (download, verify checksum, extract)
# ------------------------------------------------------------------
WORKDIR /home/agent
RUN set -eux; \
    architecture="$(uname -m)"; \
    if [ "$architecture" = "aarch64" ]; then \
        agent_tar="harden-runner-bravo_1.9.3_linux_arm64.tar.gz"; \
        agent_sha="f52555c8ab659a8a9b870fd3b8e820c064e235e3c34e83adf30db4b46e6a2efb"; \
    else \
        agent_tar="harden-runner-bravo_1.9.3_linux_amd64.tar.gz"; \
        agent_sha="995c1157c2764d2b09ba369fc2583104096cfb9a40c09af4b00b77e32651216b"; \
    fi; \
    curl -fsSLO "https://packages.stepsecurity.io/self-hosted/${agent_tar}"; \
    echo "${agent_sha}  ${agent_tar}" | sha256sum -c -; \
    tar -xzf "${agent_tar}"; \
    rm -f "${agent_tar}"; \
    chmod +x /home/agent/agent

# The agent is not run as a service; the pre-job hook starts it.
# agent.json contains the API key, so it is written at runtime from the
# API_KEY env var by configure-agent.sh rather than baked into the image.
COPY configure-agent.sh /opt/step-security/configure-agent.sh
RUN chmod +x /opt/step-security/configure-agent.sh

WORKDIR /codebuild
