# Public Ubuntu base image from Docker Hub
FROM docker.io/library/ubuntu:22.04

ARG DEBIAN_FRONTEND=noninteractive

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
# Harden-Runner installation
# TODO: replace this block with the install steps you'll provide.
# ------------------------------------------------------------------
# RUN curl -fsSL <harden-runner-agent-url> -o /tmp/agent.tar.gz && \
#     <verify checksum> && \
#     <install / configure agent> && \
#     rm -f /tmp/agent.tar.gz
# ------------------------------------------------------------------

WORKDIR /codebuild
