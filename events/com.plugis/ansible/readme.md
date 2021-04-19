# com.plugis.ansible cloud events

Cloud events :
- Type : com.plugis.ansible.playbook
- topic : com.plugis.ansible.8E:82:45:0E:A4:6F
- parameters :
* packages : packages to install with apt install
* roles : roles to install with ansible-galaxy install
* base : prefix for inventory and playbooks, can be a folder or an http url
* inventory : one inventory file
* playbooks : array of playbook files

- parameter sample :

```json
{
  "packages": "ansible sshpass aptitude",
  "roles": "geerlingguy.docker geerlingguy.pip",
  "base": "/tmp/v2-ansible/",
  "inventory": "local.yml",
  "playbooks": ["playbooks/upgrade.yml","/tmp/v2-ansible/playbooks/docker.yml"]
}
```

./event-sender -type 'com.plugis.ansible.playbook' -topic "com.plugis.ansible.8E:82:45:0E:A4:6F" -request -timeout 600 \
    -data '{"packages": "ansible sshpass aptitude","roles": "geerlingguy.docker geerlingguy.pip","base": "https://update.plugis.com/ansible/","inventory": "local.yml","playbooks": ["upgrade.yml","site.yml"]}' \

