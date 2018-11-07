#!/usr/bin/python
from github import Github
import os, sys
import json

def usage():
    print("usage: {} FILE1 FILE2 ...".format(sys.argv[0]))
    print("\tFILE is output of gosec (https://github.com/securego/gosec) in JSON format")
    exit(1)

debug = True
# debug = False

### export your github access token
### see https://github.com/settings/tokens to get access token
ACCESS_TOKEN = os.environ['GITHUB_ACCESS_TOKEN']

BASE_DIR = os.getcwd()
GOSEC_RESULT_FILES = sys.argv[1:]

if len(GOSEC_RESULT_FILES) == 0:
    usage()

# or using an access token
g = Github(ACCESS_TOKEN)

### FIXME: Hard coded
repo = g.get_repo("K-atc/play-with-gosec")

# master_branch = repo.get_branch("master")
# head = master_branch.commit
# sha = head.sha
# print(head)

### FIXME: Hard coded
sha = "0dc89bcc22163cf3372e79dd5c2b185d3c68a9ff"
commit = repo.get_commit(sha)
if debug: print(commit)
commit_files = []
for x in commit.files:
    commit_files.append(x.filename)
if debug: print(commit_files)

### DANGER: Remove Comments
for x in commit.get_comments():
    if debug: print('deleting {}'.format(x))
    x.delete()

for result_file in GOSEC_RESULT_FILES:
    with open(result_file) as f:
        result = json.load(f)

    issues = result["Issues"]
    for i in issues:
        path = i['file'].replace(BASE_DIR, '').strip('/') # transform to relative path
        if not path in commit_files:
            continue
        if debug: print('path = {}'.format(path))

        if '-' in i['line']:
            position = int(i['line'].split('-')[0])
        else:
            position = int(i['line'])
        if debug: print('position = {}'.format(position))

        body = '### issue reported by gosec\n'
        if True:
            body += '**[{severity}:{confidence}]** {details} ({rule_id})\n'.format(**i)
            body += '{}:{line} `{code}`'.format(path, **i)
        if False:
            body += 'field | description \n'
            body += ':-:|:-:\n'
            body += 'severity | {}\n'.format(i['severity'])
            body += 'confidence | {}\n'.format(i['confidence'])
            body += 'rule_id | {}\n'.format(i['rule_id'])
            body += 'details | {}\n'.format(i['details'])
            # body += 'file | {}\n'.format(i['file']) # NOTE: Should not be appear for security reason
            body += 'code | {}\n'.format(i['code'])
        if debug: print(body)

        commit.create_comment(body, line=0, path=path, position=position) # NOTE: `line` is Deprecated
