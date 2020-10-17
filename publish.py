#!/usr/bin/env python3
import urllib.request
import subprocess
import json
import yaml
from datetime import datetime

# print the latest tags so we don't have to google our own
# image to check :P
r = urllib.request.urlopen('https://registry.hub.docker.com/v2/repositories/place1/wg-access-server/tags?page_size=10') \
    .read() \
    .decode('utf-8')
tags = json.loads(r).get('results', [])
tags.sort(key=lambda t: datetime.strptime(t.get('last_updated'), '%Y-%m-%dT%H:%M:%S.%f%z'))
tags = [t.get('name') for t in tags]
tags = tags[-4:]
print('current docker tags:')
print('\n'.join(['    ' + t for t in tags]))

# tag the new image
version = input('pick a published tag: ')
docker_tag = f"place1/wg-access-server:{version}"

# update the helm chart and quickstart manifest
with open('deploy/helm/wg-access-server/Chart.yaml', 'r+') as f:
    chart = yaml.load(f)
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
                'https://place1.github.io/wg-access-server'])

# update gh-pages (docs)
subprocess.run(['mkdocs', 'gh-deploy'])

# commit changes
subprocess.run(['git', 'add', '.'])
subprocess.run(['git', 'commit', '-m', f'{version} - helm & docs update'])

# push everything
subprocess.run(['git', 'push'])
subprocess.run(['git', 'push', '--tags'])
