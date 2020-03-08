#!/usr/bin/env python3
import urllib.request
import subprocess
import json
import yaml

# print the latest tags so we don't have to google our own
# image to check :P
r = urllib.request.urlopen('https://registry.hub.docker.com/v2/repositories/place1/wg-access-server/tags?page_size=10') \
    .read() \
    .decode('utf-8')
tags = json.loads(r).get('results', [])
print('current docker tags:', sorted([t.get('name') for t in tags], reverse=True))

# tag and push the new image
version = input('Version: ')
# image_tag=f"place1/wg-access-server:{version}"
# subprocess.run(['docker', 'build', '-t', image_tag, '.'])
# subprocess.run(['docker', 'push', image_tag])

# update the helm chart and quickstart manifest
with open('deploy/helm/wg-access-server/Chart.yaml', 'r+') as f:
    chart = yaml.load(f)
    chart['appVersion'] = version
    yaml.dump(chart, f)
with open('deploy/k8s/quickstart.yaml', 'x') as f:
    subprocess.run(['helm', 'template', '--name-template', 'quickstart', 'deploy/helm/wg-access-server/'], stdout=f)
