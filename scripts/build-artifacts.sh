#!/usr/bin/env bash

set -o errexit
cd cmd

echo "--> Build ENS"

ENS_ARTIFACTS=../builtin/ens/artifacts
go run main.go abigen --source ${ENS_ARTIFACTS}/ENS.abi,${ENS_ARTIFACTS}/Resolver.abi --output ../builtin/ens --package ens

echo "--> Build ERC20"

ERC20_ARTIFACTS=../builtin/erc20/artifacts
go run main.go abigen --source ${ERC20_ARTIFACTS}/ERC20.abi --output ../builtin/erc20 --package erc20

echo "--> Build Testdata"
go run main.go abigen --source ./abigen/testdata/testdata.abi --output ./abigen/testdata --package testdata

CONTRACT_ARTIFACTS=/Users/pagefau1t/work/open-wallet/
echo "--> Build Contract"
# go run main.go abigen --source ${CONTRACT_ARTIFACTS}/abi/erc20.abi --output ${CONTRACT_ARTIFACTS}/internal/treasury/contract --package contract
# go run main.go abigen --source ${CONTRACT_ARTIFACTS}/abi/erc20CloneFactory.abi --output ${CONTRACT_ARTIFACTS}/internal/treasury/contract --package contract
# go run main.go abigen --source ${CONTRACT_ARTIFACTS}/abi/treasury.abi --output ${CONTRACT_ARTIFACTS}/internal/treasury/contract --package contract
# go run main.go abigen --source ${CONTRACT_ARTIFACTS}/abi/treasuryCloneFactory.abi --output ${CONTRACT_ARTIFACTS}/internal/treasury/contract --package contract
#

CONTRACT_ARTIFACTS=/Users/pagefau1t/labs/openwallet-sdk/
# go run main.go abigen --source ${CONTRACT_ARTIFACTS}/abi/erc20.abi --output ${CONTRACT_ARTIFACTS}/internal/client/contract --package contract
go run main.go abigen --source ${CONTRACT_ARTIFACTS}/abi/erc721.abi --output ${CONTRACT_ARTIFACTS}/internal/client/contract --package contract
