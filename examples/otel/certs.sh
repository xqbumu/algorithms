#!/usr/bin/env bash
# shopt -s nullglob globstar
set -x # have bash print command been ran
set -e # fail if any command fails

SUBJECT="/C=US/ST=CA/O=MyOrg/CN=myOrgCA"
SUBJECT_ALT_NAME="DNS:otel_collector,DNS:localhost,IP:127.0.0.1"

setup_certs() {
  { # create CA.
    openssl req \
      -new \
      -newkey rsa:4096 \
      -nodes \
      -days 1024 \
      -x509 \
      -keyout confs/ca.key \
      -out confs/ca.crt \
      -subj "$SUBJECT"
  }

  { # create server certs.
    openssl req \
      -new \
      -newkey rsa:2048 \
      -nodes \
      -keyout confs/server.key \
      -out confs/server.csr \
      -subj "$SUBJECT"
    openssl x509 \
      -req \
      -days 1000 \
      -in confs/server.csr \
      -CA confs/ca.crt \
      -CAkey confs/ca.key \
      -CAcreateserial \
      -out confs/server.crt \
      -extfile <(echo subjectAltName = $SUBJECT_ALT_NAME)
  }

  { # create client certs.
    openssl req \
      -new \
      -newkey rsa:2048 \
      -days 1000 \
      -nodes \
      -keyout confs/client.key \
      -out confs/client.csr \
      -subj "$SUBJECT"
    openssl x509 \
      -req \
      -in confs/client.csr \
      -days 1000 \
      -CA confs/ca.crt \
      -CAkey confs/ca.key \
      -CAcreateserial \
      -out confs/client.crt
  }

  { # clean
    rm -rf confs/*.csr
    rm -rf confs/*.srl

    chmod 666 confs/server.crt confs/server.key confs/ca.crt
  }
}
setup_certs
