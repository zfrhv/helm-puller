


# This script pulls image only by tag, it ignores the images by digest
ls *.tgz | while read -r chart; do
  appVersion=$(tar -xOf $chart */Chart.yaml | yq e '.appVersion')
  all_images=$(tar -xOf $chart */values.yaml | yq e '[.. | select(has("repository"))]')
  images=$(echo "$all_images" | yq e '.[] | select(.tag != null and .tag != "") | .repository + ":" + .tag')
  images="$images $(echo "$all_images" | yq e '.[] | select(.tag == "" or .tag == null) | .repository + ":" + "'"$appVersion"'"')" # when the tag wasnt defined at all (null), the tag was the same as Chart appVersion. when the tag is epty string(""), then idk, possible that its the same, or latest

  folder="to-take-$(echo $chart | grep -Po '.*(?=.tgz)')"
  rm -rf $folder
  mkdir -p $folder/images
  mv $chart $folder/
  cp .0#conf.sh $folder/0#conf.sh
  cp .2#extract-and-push.sh $folder/2#extract-and-push.sh
  cd $folder/images

  echo -n "$images" > images.txt
  counter=1
  for image_name in $images; do
    echo "pulling image $image_name"
    podman pull $image_name
    file=${image_name##*/}
    file=${file%%:*}
    file="$counter-${file}.tar"
    podman save $image_name -o $file
    podman rmi $image_name
    ((counter++))
  done
  cd ../..
  rm -f $folder.tar.gz
  tar -czvf $folder.tar.gz $folder
  rm -rf $folder
  echo "take the ${folder}.tar.gz to restricted environment."
done
