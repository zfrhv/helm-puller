source 0#conf.sh

if [ -z "$registry_name" ] || [ -z "$helm_repo" ]
then
  echo "error: you didnt specified the registry name, port or helm_repo in 0#config.sh, please specify the registry+port+helm_repo where to push images and the helm chart"
else
  echo "please login to $registry_name:$registry_port"
  podman login $registry_name:$registry_port
  chart="$(ls *.tgz)"
  curl -u "$(cat $XDG_RUNTIME_DIR/containers/auth.json | jq '.auths["'"$registry_name:$registry_port"'"].auth' --raw-output | base64 -d)" -T $chart "$(echo $helm_repo | sed 's|/$||g')/$chart" 
  cd images
  ls *.tar | while read file; do
    old_image=$(podman load -i $file)
    old_image=${old_image##* }
    new_image="$registry_name:$registry_port/${old_image#*/}"
    podman tag $old_image $new_image
    podman push $new_image
    podman rmi $old_image $new_image
  done
  echo "# done! to use the helm chart type:"
  echo "helm install test $helm_repo/$chart"
  echo "# or type:"
  echo "helm repo add ${helm_repo##*/} $helm_repo"
  echo "helm install test ${helm_repo##*/}/${chart%.tgz}"
fi
