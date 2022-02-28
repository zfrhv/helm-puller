# helm-puller
pulls all helm images and compresses into tar.gz

## cloning:
```bash
git clone https://github.com/zfrhv/helm-puller.git
cd helm-puller
```

## usage:
```bash
helm pull <repository>/<chart>
./1#pull-and-compress.sh
```

take the new to-take-<chart>.tar.gz to the restricted network (using usb or any media)

push all images + helm chart to the repository:
```bash
# untar the file
tar -xzvf to-take-<chart>.tar.gz && rm -f to-take-<chart>.tar.gz

cd to-take-<chart>
vim 0#conf.sh
  
./2#extract-and-push.sh
```  
  
## examples:
```
helm repo add grafana https://grafana.github.io/helm-charts
helm pull grafana/grafana
./1#pull-and-compress.sh
```
