#!/bin/bash
#
# This script generates basic "machine" certificate for the FQDN.
#

CAPath="/root/ca"
CACert="${CAPath}/rzCA.crt"
CAKey="${CAPath}/rzCA.key"
OPENSSL="/usr/bin/openssl"
LOOKUPTOOL="/usr/bin/host"

DURATION=365

if [[ -z $1 ]]; then
    echo "Must specify domain FQDN." > /dev/stderr
    exit 42
fi
MachineName="$1"

if [[ ! -f ${LOOKUPTOOL} || ! -x ${LOOKUPTOOL} ]]; then
  echo "The '${LOOKUPTOOL} is not valid." >/dev/stderr
  exit 42
fi
${LOOKUPTOOL} ${MachineName} 2>&1 > /dev/null
if [[ $? -ne 0 ]]; then
    echo "Machine ${MachineName} must be in DNS for a key to be generated for it." > /dev/stderr
    exit 42
fi

MachineKey="${MachineName}.key"
MachineCSR="${MachineName}.csr"
MachineCert="${MachineName}.crt"
MachinePEM="${MachineName}.pem"

${OPENSSL} req \
  -nodes \
  -newkey rsa:2048 \
  -keyout "${MachineKey}" \
  -out "${MachineCSR}" \
  -subj "/C=US/ST=CA/L=Los Gatos/O=Rezidencija/OU=Machine/CN=${MachineName}"
if [[ $? -ne 0 ]]; then
  echo "Certificate '${MachineCSR}' and/or key '${MachineKey}' generation failed." > /dev/stderr
  exit 42
fi

sudo ${OPENSSL} x509 -req \
  -days "${DURATION}" \
  -in "${MachineCSR}" \
  -CA "${CACert}"  \
  -CAkey "${CAKey}" \
  -set_serial "$(date '+%Y%m%d%H%M%S')" \
  -out "${MachineCert}"
if [[ $? -ne 0 ]]; then
  echo "Certificate '${MachineCert}' signing failed." > /dev/stderr
  exit 42
fi

cat "${MachineCert}" "${MachineKey}" > "${MachinePEM}"
if [[ $? -ne 0 ]]; then
  echo "Could not generate '${MachinePEM}'." > /dev/stderr
  exit 42
fi

rm -f "${MachineCSR}" "${MachineKey}" "${MachineCert}"
if [[ $? -ne 0 ]]; then
  echo "Could not remove temporary/intermediate files." > /dev/stderr
  exit 42
fi
