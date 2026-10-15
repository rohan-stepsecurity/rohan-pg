# Public Ubuntu base image from Docker Hub
FROM docker.io/library/ubuntu:22.04

ARG DEBIAN_FRONTEND=noninteractive
ARG HOOKS_VERSION=v1.4.0

# Base tooling that CodeBuild and most builds expect, plus the GitHub
# Actions runner's .NET dependencies (libicu70, libkrb5-3, zlib1g,
# libssl3, liblttng-ust1).
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
        lsb-release \
        tar \
        libicu70 \
        libkrb5-3 \
        zlib1g \
        libssl3 \
        liblttng-ust1 && \
    rm -rf /var/lib/apt/lists/*

# Node.js is required to run the StepSecurity job hooks (pre.js / post.js).
# Ubuntu 22.04's apt nodejs is too old, so use NodeSource.
RUN curl -fsSL https://deb.nodesource.com/setup_20.x | bash - && \
    apt-get install -y --no-install-recommends nodejs && \
    rm -rf /var/lib/apt/lists/*

# ------------------------------------------------------------------
# 1. StepSecurity GitHub Actions runner job hooks
# ------------------------------------------------------------------
RUN mkdir -p /opt/step-security/hooks && \
    curl -fsSL "https://github.com/step-security/gha-runner-job-hooks/releases/download/${HOOKS_VERSION}/pre.js" \
        -o /opt/step-security/hooks/pre.js && \
    curl -fsSL "https://github.com/step-security/gha-runner-job-hooks/releases/download/${HOOKS_VERSION}/post.js" \
        -o /opt/step-security/hooks/post.js && \
    for hook in pre post; do \
        printf '%s\n' \
            '#!/bin/bash' \
            'export STEP_HOOK_MODE="vm"' \
            'export STEP_AGENT_ROOT="/home/agent"' \
            'preserved_env="$(compgen -e | paste -sd, -)"' \
            "sudo --preserve-env=\"\$preserved_env\" env node /opt/step-security/hooks/${hook}.js" \
            > "/opt/step-security/${hook}-job-hook.sh"; \
        chmod 0755 "/opt/step-security/${hook}-job-hook.sh"; \
    done

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
    echo "${agent_sha}  ${agent_tar}" | sha256sum -c -
# The agent is not run as a service; the pre-job hook starts it.
# agent.json (with API_KEY from the environment) is handled at runtime,
# not baked into the image.

WORKDIR /codebuild
