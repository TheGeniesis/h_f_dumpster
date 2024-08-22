pat=$(az account get-access-token --query accessToken -o tsv)
url="https://dev.azure.com/${REPO}/_apis/build/builds?statusFilter=notStarted&api-version=7.0"


res=$(curl -u :${pat} --header "ContentType: application/json" $url)

echo "${res}"

names=$(echo "$res" | jq -r .value[].definition.id)

for name in ${names[@]}; do
    echo "deleting job ${name}"
    urlToCancel="https://dev.azure.com/${repo}/_apis/build/builds/${name}?api-version=7.0"
    res=$(curl -u :${pat} -X DELETE --header "ContentType: application/json" ${urlToCancel})

    echo "----------"
    echo "----------"
    echo "----------"
    echo "${res}"

done
