# Thrash is Trash

## Reinstall java 11

```sh
apt remove openjdk-11-jre -y
<path_to_java>

apt update -y && apt upgrade -y; \
apt purge openjdk* -y; \
apt autoremove openjdk* -y; \
rm -Rf /usr/lib/jvm/; \
apt install openjdk-11-jre-headless -y; \
apt install openjdk-11-jre -y; \
java -version \
apt update -y && apt upgrade -y;
```

## Unmount

```shell
systemctl stop <process> && umount <nfs_IP>:<dir>
```

## Install hadoop

```shell
apt install -y zstd lz4 bzip2 openssl libsnappy-dev; \
cd /tmp; \
wget https://archive.apache.org/dist/hadoop/core/hadoop-<version>/hadoop-<version>.tar.gz; \
tar -xvf hadoop-<version>.tar.gz ; \
rm hadoop-<version>.tar.gz ; \
mv hadoop-<version> /usr/local/hadoop; \
ln -s /usr/local/hadoop/bin/hadoop /usr/local/bin/hadoop; \
echo 'export JAVA_HOME=$(readlink -f /usr/bin/java | sed "s:bin/java::")' | tee -a /usr/local/hadoop/etc/hadoop/hadoop-env.sh
```

## Hadoop bashrc

```sh
export HADOOP_HOME=/usr/local/hadoop; \
echo "export HADOOP_HOME=/usr/local/hadoop" >> ${HOME}/.bashrc; \
echo "export HADOOP_INSTALL=$HADOOP_HOME" >> ${HOME}/.bashrc; \
echo "export HADOOP_MAPRED_HOME=$HADOOP_HOME" >> ${HOME}/.bashrc; \
echo "export HADOOP_COMMON_HOME=$HADOOP_HOME" >> ${HOME}/.bashrc; \
echo "export HADOOP_HDFS_HOME=$HADOOP_HOME" >> ${HOME}/.bashrc; \
echo "export YARN_HOME=$HADOOP_HOME" >> ${HOME}/.bashrc; \
echo "export HADOOP_COMMON_LIB_NATIVE_DIR=$HADOOP_HOME/lib/native" >> ${HOME}/.bashrc; \
echo "export PATH=$PATH:$HADOOP_HOME/sbin:$HADOOP_HOME/bin" >> ${HOME}/.bashrc; \
echo "export HADOOP_OPTS=\"-Djava.library.path=$HADOOP_HOME/lib/native\"" >> ${HOME}/.bashrc; \
source ${HOME}/.bashrc;
```

## Update JAVA certs

### PATH

echo $JAVA_HOME
<path_to_java>/lib/security

### locally

Password for example is `changeit`

Download cert

```shell
echo -n | openssl s_client -connect https://<host>:443 | sed -ne '/-BEGIN CERTIFICATE-/,/-END CERTIFICATE-/p' > <filename>.crt

# Copy cert to the server
scp <path_to_downloaded_cert>/<filename>.crt <server>:<path_to_java>/lib/security

# Then ssh to the instance and execute cd <path_to_java>/lib/security
keytool -import -v -trustcacerts -alias <host> -file <filename>.crt -keystore cacerts -keypass changeit -storepass changeit

# Delete alias if needed
keytool -delete -noprompt -trustcacerts -alias "<host>" -keystore "cacerts"

systemctl restart <java_process>
```
