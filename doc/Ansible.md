# Ansible

## encrypt file

```shell
ansible-vault encrypt roles/<role_name>>/vars/main.yml --vault-password-file=~/.ssh/vault_pass
ansible-vault decrypt pass --vault-password-file=~/.ssh/vault_pass

echo -n "<password>" | ansible-vault encrypt_string --vault-password-file=~/.ssh/vault_pass
```

Result to save in yaml looks like

```yaml
password:  !vault |
#  $ANSIBLE_VAULT;1.1;AES256
#  <someID>
```

## Deps update

```shell
ansible-galaxy install -r requirements.yml --force
```
