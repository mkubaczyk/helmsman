#!/bin/bash

ytt --ignore-unknown-comments -f - -f $(dirname $0)/overlay.sample.yaml
