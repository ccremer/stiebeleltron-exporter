# This file is managed by greposync.
# Do not modify manually.
# Adjust variables in `.sync.yml`.

# These are some common variables for Make

BIN_FILENAME ?= stiebeleltron-exporter

# Image URL to use all building/pushing image targets
IMG_TAG ?= latest
LOCAL_IMG ?= local.dev/ccremer/stiebeleltron-exporter:$(IMG_TAG)
