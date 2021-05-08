sam_host="http://localhost:3000"
echo "Waiting ${sam_host} to be up"

while true
do
  curl ${sam_host} &> /dev/null
  if [[ $? -eq 0 ]]
  then
    echo "${sam_host} is up"
    break
  else
    echo "try again in 1s ..."
    sleep 1
  fi
done
