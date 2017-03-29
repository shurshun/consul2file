#!/usr/bin/env sh

consul watch -http-addr=${CONSUL_ADDR} -type keyprefix -prefix ${KEY_PREFIX} /bin/consul2file -p ${KEY_PREFIX} -o ${STORE_DIR}