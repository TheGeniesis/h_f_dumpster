# databricks

## Update jobs using bash

LOCATION="<azure_region>"
ENV_NAME="preprod"
ADB_PROFILE=${LOCATION}-${ENV_NAME}

### get all the jobs config in json files

```shell
  for jobid in $(databricks --profile ${ADB_PROFILE} jobs list -o json | jq '.[].job_id')
  do
    job_update_json_file="${ADB_PROFILE}_update_job_${jobid}.json"
    echo "json update file: ${job_update_json_file}"
    databricks --profile  ${ADB_PROFILE} jobs get ${jobid} -o json | jq '. |= ( del(.settings.job_clusters[].new_cluster.instance_pool_id) | del(.settings.job_clusters[].new_cluster.driver_instance_pool_id) | del(.settings.job_clusters[].new_cluster.spark_conf) | del(.settings.job_clusters[].new_cluster.single_user_name)) | {job_id: .job_id, new_settings: {job_clusters: .settings.job_clusters}}' >${job_update_json_file}
  done
```

### sample job json file

```json
{
  "job_id": <job_id>,
  "new_settings": {
    "job_clusters": [
      {
        "job_cluster_key": "job-cluster-us-lab",
        "new_cluster": {
          "autoscale": {
            "max_workers": 2,
            "min_workers": 1
          },
          "cluster_name": "",
          "custom_tags": {
            "country": "us"
          },
          "data_security_mode": "SINGLE_USER",
          "policy_id": "<policy>",
          "runtime_engine": "PHOTON",
          "spark_version": "13.3.x-scala2.12"
        }
      }
    ]
  }
}
```

### run update

```shell
for job_update_json in $(ls ${ADB_PROFILE}_update_job_*.json)
do
  if [[ -s "${job_update_json}" ]]
    then
      echo "updating: ${job_update_json}"
      databricks --profile  ${ADB_PROFILE} jobs update --json @${job_update_json}
  fi
done
```
