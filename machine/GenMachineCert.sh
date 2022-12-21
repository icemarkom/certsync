#!/bin/bash
#
# This script example generates basic "machine" certificate for the FQDN.
#

###############################################################################
#
# SET THESE VALUES BEFORE USING THE SCRIPT
#
CAPath="/etc/ssl/CA"
CACert="${CAPath}/CA.pem"
CAKey="${CAPath}/CA.key"

CACountry="US"
CAState="California"
CACity="Los Gatos"
CAOrg="Rezidencija"
CAUnit="Machine"
###############################################################################

OPENSSL="/usr/bin/openssl"
LOOKUPTOOL="/usr/bin/host"
MKTEMP="/bin/mktemp"

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

CFGFile="$(${MKTEMP} --tmpdir cfg-XXXXX)"

cat > ${CFGFile} << EOF
[ req ]
default_bits       = 2048
distinguished_name = req_distinguished_name
x509_extensions    = x509_ext
prompt             = no

[ req_distinguished_name ]
C                  = ${CACountry}
ST                 = ${CAState}
L                  = ${CACity}
O                  = ${CAOrg}
OU                 = ${CAUnit}
CN                 = ${MachineName}

[ x509_ext ]
subjectAltName     = @alt_names
keyUsage           = keyEncipherment, dataEncipherment
extendedKeyUsage   = clientAuth

[ alt_names ]
DNS.1              = ${MachineName}
EOF

MachineKey="${MachineName}.key"
MachineCSR="${MachineName}.csr"
MachineCert="${MachineName}.crt"
MachinePEM="${MachineName}.pem"

${OPENSSL} req \
  -nodes \
  -newkey rsa:2048 \
  -keyout "${MachineKey}" \
  -out "${MachineCSR}" \
  -config "${CFGFile}"
if [[ $? -ne 0 ]]; then
  echo "Certificate '${MachineCSR}' and/or key '${MachineKey}' generation failed." > /dev/stderr
  exit 42
fi

sudo ${OPENSSL} x509 \
  -req \
  -days "${DURATION}" \
  -in "${MachineCSR}" \
  -CA "${CACert}"  \
  -CAkey "${CAKey}" \
  -set_serial "$(date '+%Y%m%d%H%M%S')" \
  -out "${MachineCert}" \
  -extensions x509_ext \
  -extfile ${CFGFile}
if [[ $? -ne 0 ]]; then
  echo "Certificate '${MachineCert}' signing failed." > /dev/stderr
  exit 42
fi

cat "${MachineCert}" "${MachineKey}" > "${MachinePEM}"
if [[ $? -ne 0 ]]; then
  echo "Could not generate '${MachinePEM}'." > /dev/stderr
  exit 42
fi

rm -f "${MachineCSR}" "${MachineKey}" "${MachineCert}" "${CFGFile}"
if [[ $? -ne 0 ]]; then
  echo "Could not remove temporary/intermediate files." > /dev/stderr
  exit 42
fi
