#!/bin/bash
#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
PATH=build/bin/:${PATH}

# Takes in 4 arguments
# 1. Output doc file
# 2. Preamble Text File
# 3. Postscript File
# 4. Array of commands
generateHelpText(){
    local DOC="$1"
    local preamble="$2"
    local postscript="$3"
    # Shift three times to get to array
    shift
    shift
    shift

    cat "$preamble" > "$DOC"

    local commands=("$@")
    for x in "${commands[@]}" ; do

    cat <<EOF >> "$DOC"

## $x
\`\`\`
$($x --help 2>&1 | sed -E 's/[[:space:]]+$//g')
\`\`\`

EOF
    done
    cat "$postscript" >> "$DOC"
}

commands=("peer version")
generateHelpText \
        docs/source/commands/peerversion.md \
        docs/wrappers/peer_version_preamble.md \
        docs/wrappers/license_postscript.md \
        "${commands[@]}"

commands=("peer chaincode install" "peer chaincode instantiate" "peer chaincode invoke" "peer chaincode list" "peer chaincode package" "peer chaincode query" "peer chaincode signpackage" "peer chaincode upgrade")
generateHelpText \
        docs/source/commands/peerchaincode.md \
        docs/wrappers/peer_chaincode_preamble.md \
        docs/wrappers/peer_chaincode_postscript.md \
        "${commands[@]}"

commands=("peer lifecycle" "peer lifecycle chaincode" "peer lifecycle chaincode package" "peer lifecycle chaincode install" "peer lifecycle chaincode queryinstalled" "peer lifecycle chaincode getinstalledpackage" "peer lifecycle chaincode approveformyorg" "peer lifecycle chaincode queryapproved" "peer lifecycle chaincode checkcommitreadiness" "peer lifecycle chaincode commit" "peer lifecycle chaincode querycommitted")
generateHelpText \
        docs/source/commands/peerlifecycle.md \
        docs/wrappers/peer_lifecycle_chaincode_preamble.md \
        docs/wrappers/peer_lifecycle_chaincode_postscript.md \
        "${commands[@]}"

commands=("peer channel" "peer channel create" "peer channel fetch" "peer channel getinfo" "peer channel join" "peer channel joinbysnapshot" "peer channel joinbysnapshotstatus" "peer channel list" "peer channel signconfigtx" "peer channel update")
generateHelpText \
        docs/source/commands/peerchannel.md \
        docs/wrappers/peer_channel_preamble.md \
        docs/wrappers/peer_channel_postscript.md \
        "${commands[@]}"

commands=("peer node pause" "peer node rebuild-dbs" "peer node reset" "peer node resume" "peer node rollback" "peer node start" "peer node upgrade-dbs")
generateHelpText \
        docs/source/commands/peernode.md \
        docs/wrappers/peer_node_preamble.md \
        docs/wrappers/peer_node_postscript.md \
        "${commands[@]}"

commands=("peer snapshot cancelrequest" "peer snapshot listpending" "peer snapshot submitrequest")
generateHelpText \
        docs/source/commands/peersnapshot.md \
        docs/wrappers/peer_snapshot_preamble.md \
        docs/wrappers/peer_snapshot_postscript.md \
        "${commands[@]}"

commands=("configtxgen")
generateHelpText \
        docs/source/commands/configtxgen.md \
        docs/wrappers/configtxgen_preamble.md \
        docs/wrappers/configtxgen_postscript.md \
        "${commands[@]}"

commands=("cryptogen help" "cryptogen generate" "cryptogen showtemplate" "cryptogen extend" "cryptogen version")
generateHelpText \
        docs/source/commands/cryptogen.md \
        docs/wrappers/cryptogen_preamble.md \
        docs/wrappers/cryptogen_postscript.md \
        "${commands[@]}"

commands=("configtxlator start" "configtxlator proto_encode" "configtxlator proto_decode" "configtxlator compute_update" "configtxlator version")
generateHelpText \
        docs/source/commands/configtxlator.md \
        docs/wrappers/configtxlator_preamble.md \
        docs/wrappers/configtxlator_postscript.md \
        "${commands[@]}"

commands=("osnadmin channel" "osnadmin channel join" "osnadmin channel list" "osnadmin channel remove")
generateHelpText \
        docs/source/commands/osnadminchannel.md \
        docs/wrappers/osnadmin_channel_preamble.md \
        docs/wrappers/osnadmin_channel_postscript.md \
        "${commands[@]}"

exit
