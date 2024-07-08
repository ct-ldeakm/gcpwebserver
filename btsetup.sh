INSTANCEID=dev-instance
INSTANCENAME=build-testing
CLUSTERID=dev-cluster
ZONE=us-east1-b

gcloud bigtable instances create $INSTANCEID \
    --display-name=$INSTANCENAME \
    --cluster-storage-type=SSD \
    --cluster-config=id=$CLUSTERID,zone=$ZONE,nodes=1 

gcloud bigtable instances describe $INSTANCEID

gcloud bigtable instances delete $INSTANCEID --quiet