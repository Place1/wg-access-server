#!/usr/bin/env python3

# This script is intended to be run within GitHub Actions, triggered after new tags have been created.
# It updates the version in the Helm Chart, packages it, renders the k8s quickstart.yaml, then commits and pushes everything.
# A separate workflow triggered on pushes should then publish the charts to the GitHub Pages website.

import os
import subprocess
import yaml

version = os.environ.get('GITHUB_REF_NAME')
ref_type = os.environ.get('GITHUB_REF_TYPE')
if not version or ref_type != 'tag':
    print('::error::Aborting, workflow not triggered by tag event')
    exit(1)

# update the helm chart and quickstart manifest
with open('deploy/helm/wg-access-server/Chart.yaml', 'r+') as f:
    chart = yaml.safe_load(f)
    chart['version'] = version
    chart['appVersion'] = version
    f.seek(0)
    yaml.dump(chart, f, default_flow_style=False)
    f.truncate()
with open('deploy/k8s/quickstart.yaml', 'w') as f:
    subprocess.run(['helm', 'template', '--name-template',
                    'quickstart', 'deploy/helm/wg-access-server/'], stdout=f)
subprocess.run(['helm', 'package', 'deploy/helm/wg-access-server/',
                '--destination', 'docs/charts/'])
subprocess.run(['helm', 'repo', 'index', 'docs/', '--url',
                'https://freie-netze.org/wg-access-server'])

# commit changes
subprocess.run(['git', 'add', 'docs/index.yaml', 'docs/charts/', 'deploy/helm/', 'deploy/k8s/'])
subprocess.run(['git', 'commit', '-m', f'{version} - Automated Helm & k8s update'])

# push everything
subprocess.run(['git', 'push'])
