# ansible command

The ansible command allows to download roles, 

sudo ./ansible -playbooks /tmp/v2-ansible/playbooks/docker.yml -inventory /tmp/v2-ansible/local.yml -roles 'geerlingguy.docker geerlingguy.pip' -packages 'ansible sshpass aptitude'

./ansible -base '/tmp/v2-ansible/' -playbooks 'playbooks/upgrade.yml playbooks/data_dir.yml' -inventory 'local.yml' -log trace

```json
{
  "packages": "ansible sshpass aptitude",
  "roles": "geerlingguy.docker geerlingguy.pip",
  "base": "/tmp/v2-ansible/",
  "inventory": "local.yml",
  "playbooks": ["playbooks/upgrade.yml","/tmp/v2-ansible/playbooks/docker.yml"]
}
```


