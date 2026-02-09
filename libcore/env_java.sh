#!/bin/bash

# Set JAVA_HOME if not already set
if [ -z "$JAVA_HOME" ]; then
  # Try to find Java installation
  if command -v java >/dev/null 2>&1; then
    # Use java_home command on macOS if available
    if command -v /usr/libexec/java_home >/dev/null 2>&1; then
      export JAVA_HOME=$(/usr/libexec/java_home)
    # Try common Linux paths
    elif [ -d "/usr/lib/jvm/default-java" ]; then
      export JAVA_HOME="/usr/lib/jvm/default-java"
    elif [ -d "/usr/lib/jvm/java-11-openjdk-amd64" ]; then
      export JAVA_HOME="/usr/lib/jvm/java-11-openjdk-amd64"
    elif [ -d "/usr/lib/jvm/java-17-openjdk-amd64" ]; then
      export JAVA_HOME="/usr/lib/jvm/java-17-openjdk-amd64"
    fi
  fi
fi
