
for secret_name in $(az keyvault secret list --vault-name ${AZURE_KEY_VAULT_NAME} --query '[].[name]' -o tsv | grep ${PREFIX}); do
  new_secret_name=$(echo ${secret_name} | sed -e s!${FIND}!${REPLACE_TO}!)
  echo "Rename ${secret_name} into ${new_secret_name}"
  secret=$(az keyvault secret show --vault-name ${AZURE_KEY_VAULT_NAME} --name ${secret_name} --query value -o tsv)
  az keyvault secret set --vault-name ${AZURE_KEY_VAULT_NAME} --name ${new_secret_name} --value $secret
  # Delete original key if needed
  az keyvault secret delete --vault-name $AZURE_KEY_VAULT_NAME --name ${secret_name}
  # Perminantly
  az keyvault secret purge  --vault-name $AZURE_KEY_VAULT_NAME --name ${secret_name}
done
