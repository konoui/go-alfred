#!/bin/bash

TMP_DIR=$(mktemp -d)
export alfred_workflow_data=${TMP_DIR}
export alfred_workflow_cache=${TMP_DIR}
export alfred_workflow_bundleid=$(date +%s)
