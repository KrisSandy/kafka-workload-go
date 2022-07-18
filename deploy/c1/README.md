## Test cluster

* Start a producer

`export KAFKA_CLUSTER_NAME=test-cluster`
`kubectl run kafka-producer -ti --image=strimzi/kafka:latest-kafka-2.4.0 --rm=true --restart=Never -- bin/kafka-console-producer.sh --broker-list $KAFKA_CLUSTER_NAME-kafka-bootstrap:9092 --topic test-topic`

* Start a consumer

`export KAFKA_CLUSTER_NAME=test-cluster`
`kubectl run kafka-consumer -ti --image=strimzi/kafka:latest-kafka-2.4.0 --rm=true --restart=Never -- bin/kafka-console-consumer.sh --bootstrap-server $KAFKA_CLUSTER_NAME-kafka-bootstrap:9092 --topic test-topic --from-beginning`