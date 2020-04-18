#!/usr/bin/env python3
import urllib.request
import subprocess
import json
import yaml

def is_release_candidate(version):
    return '-rc' in version

# print the latest tags so we don't have to google our own
# image to check :P
r = urllib.request.urlopen('https://registry.hub.docker.com/v2/repositories/place1/wg-access-server/tags?page_size=10') \
    .read() \
    .decode('utf-8')
tags = json.loads(r).get('results', [])
print('current docker tags:', sorted([t.get('name') for t in tags], reverse=True))

# tag the new image
version = input('Version: ')
docker_tag = f"place1/wg-access-server:{version}"
subprocess.run(['docker', 'build', '-t', docker_tag, '.'])

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
subprocess.run(['git', 'commit', '-m', f'{version}'])

# tag the current commit
subprocess.run(['git', 'tag', '-a', f'{version}', '-m', f'{version}'])

# push everything
subprocess.run(['git', 'push'])
subprocess.run(['git', 'push', '--tags'])
subprocess.run(['docker', 'push', docker_tag])
