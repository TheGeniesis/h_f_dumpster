# kerberos

## Relogin

- klist <- returned incorrect data
- kdestroy  <- remove info related to incorrect data
- klist
- kinit <username> <- check if it works with the proper user after sync
- klist
- hdfs dfs -ls / <- check if commands are working
- ls -la /data/files/scripts/ <- check if keytab exists
- kdestroy <- remove temporary test
- klist
- klist -kt /data/files/scripts/<username>.keytab
- kinit -kt /data/files/scripts/<username>.keytab <username>@<domain> <- do the proper init
- klist <- check if works
- hdfs dfs -ls / <- check if works
