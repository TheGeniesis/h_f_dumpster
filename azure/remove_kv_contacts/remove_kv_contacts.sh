validate_vars() {
  if [ "$#" -ne 1 ]; then
    echo "Incorrect arguments"
    echo "Example: ./remove_kv_contacts.sh \"<kv_name>\""
    exit 1
  fi
}

deleteContacts() {
  R=$(az keyvault certificate contact list --vault-name ${KV} --query 'contactList' -o tsv)

  for email in ${R}; do
    echo "Removing ${email}..."
    az keyvault certificate contact delete --vault-name ${KV} --email ${email}
  done
}

validate_vars "$@"
KV="${1}"

deleteContacts
