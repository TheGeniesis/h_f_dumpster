# Vault

Generate new secrets for role

```shell
export VAULT_ADDR=<addr> && vault login -namespace=<ns> -method=ldap username=$(whoami)

vault write -namespace <ns>  -force auth/approle/role/<role_name>/secret-id ttl=30d
```
