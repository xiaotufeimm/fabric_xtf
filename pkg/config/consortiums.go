/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package config

import (
	"errors"
	"fmt"

	"github.com/golang/protobuf/proto"
	cb "github.com/hyperledger/fabric-protos-go/common"
	mb "github.com/hyperledger/fabric-protos-go/msp"
)

// Consortium is a group of non-orderer organizations used in channel transactions.
type Consortium struct {
	Name          string
	Organizations []Organization
}

// AddOrgToConsortium adds an org definition to a named consortium in a given
// channel configuration.
func (c *ConfigTx) AddOrgToConsortium(org Organization, consortium string) error {
	var err error

	if consortium == "" {
		return errors.New("consortium is required")
	}

	consortiumsGroup := c.updated.ChannelGroup.Groups[ConsortiumsGroupKey]

	consortiumGroup, ok := consortiumsGroup.Groups[consortium]
	if !ok {
		return fmt.Errorf("consortium '%s' does not exist", consortium)
	}

	if _, ok := consortiumGroup.Groups[org.Name]; ok {
		return fmt.Errorf("org '%s' already defined in consortium '%s'", org.Name, consortium)
	}

	consortiumGroup.Groups[org.Name], err = newOrgConfigGroup(org)
	if err != nil {
		return fmt.Errorf("failed to create consortium org: %v", err)
	}

	return nil
}

// newConsortiumsGroup returns the consortiums component of the channel configuration. This element is only defined for
// the ordering system channel.
// It sets the mod_policy for all elements to "/Channel/Orderer/Admins".
func newConsortiumsGroup(consortiums []Consortium) (*cb.ConfigGroup, error) {
	var err error

	consortiumsGroup := newConfigGroup()
	consortiumsGroup.ModPolicy = ordererAdminsPolicyName

	// acceptAllPolicy always evaluates to true
	acceptAllPolicy := envelope(nOutOf(0, []*cb.SignaturePolicy{}), [][]byte{})

	// This policy is not referenced anywhere, it is only used as part of the implicit meta policy rule at the
	// channel level, so this setting effectively degrades control of the ordering system channel to the ordering admins
	signaturePolicy, err := signaturePolicy(AdminsPolicyKey, acceptAllPolicy)
	if err != nil {
		return nil, err
	}

	consortiumsGroup.Policies[signaturePolicy.key] = &cb.ConfigPolicy{
		Policy:    signaturePolicy.value,
		ModPolicy: ordererAdminsPolicyName,
	}

	for _, consortium := range consortiums {
		consortiumsGroup.Groups[consortium.Name], err = newConsortiumGroup(consortium)
		if err != nil {
			return nil, err
		}
	}

	return consortiumsGroup, nil
}

// newConsortiumGroup returns a consortiums component of the channel configuration.
func newConsortiumGroup(consortium Consortium) (*cb.ConfigGroup, error) {
	var err error

	consortiumGroup := newConfigGroup()
	consortiumGroup.ModPolicy = ordererAdminsPolicyName

	for _, org := range consortium.Organizations {
		consortiumGroup.Groups[org.Name], err = newOrgConfigGroup(org)
		if err != nil {
			return nil, fmt.Errorf("org group '%s': %v", org.Name, err)
		}
	}

	implicitMetaAnyPolicy, err := implicitMetaAnyPolicy(AdminsPolicyKey)
	if err != nil {
		return nil, err
	}

	err = setValue(consortiumGroup, channelCreationPolicyValue(implicitMetaAnyPolicy.value), ordererAdminsPolicyName)
	if err != nil {
		return nil, err
	}

	return consortiumGroup, nil
}

// consortiumValue returns the config definition for the consortium name
// It is a value for the channel group.
func consortiumValue(name string) *standardConfigValue {
	return &standardConfigValue{
		key: ConsortiumKey,
		value: &cb.Consortium{
			Name: name,
		},
	}
}

// channelCreationPolicyValue returns the config definition for a consortium's channel creation policy
// It is a value for the /Channel/Consortiums/*/*.
func channelCreationPolicyValue(policy *cb.Policy) *standardConfigValue {
	return &standardConfigValue{
		key:   ChannelCreationPolicyKey,
		value: policy,
	}
}

// envelope builds an envelope message embedding a SignaturePolicy.
func envelope(policy *cb.SignaturePolicy, identities [][]byte) *cb.SignaturePolicyEnvelope {
	ids := make([]*mb.MSPPrincipal, len(identities))
	for i := range ids {
		ids[i] = &mb.MSPPrincipal{PrincipalClassification: mb.MSPPrincipal_IDENTITY, Principal: identities[i]}
	}

	return &cb.SignaturePolicyEnvelope{
		Version:    0,
		Rule:       policy,
		Identities: ids,
	}
}

// nOutOf creates a policy which requires N out of the slice of policies to evaluate to true.
func nOutOf(n int32, policies []*cb.SignaturePolicy) *cb.SignaturePolicy {
	return &cb.SignaturePolicy{
		Type: &cb.SignaturePolicy_NOutOf_{
			NOutOf: &cb.SignaturePolicy_NOutOf{
				N:     n,
				Rules: policies,
			},
		},
	}
}

// signaturePolicy defines a policy with key policyName and the given signature policy.
func signaturePolicy(policyName string, sigPolicy *cb.SignaturePolicyEnvelope) (*standardConfigPolicy, error) {
	signaturePolicy, err := proto.Marshal(sigPolicy)
	if err != nil {
		return nil, fmt.Errorf("marshaling signature policy: %v", err)
	}

	return &standardConfigPolicy{
		key: policyName,
		value: &cb.Policy{
			Type:  int32(cb.Policy_SIGNATURE),
			Value: signaturePolicy,
		},
	}, nil
}

// implicitMetaPolicy creates a new *cb.Policy of cb.Policy_IMPLICIT_META type.
func implicitMetaPolicy(subPolicyName string, rule cb.ImplicitMetaPolicy_Rule) (*cb.Policy, error) {
	implicitMetaPolicy, err := proto.Marshal(&cb.ImplicitMetaPolicy{
		Rule:      rule,
		SubPolicy: subPolicyName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal implicit meta policy: %v", err)
	}

	return &cb.Policy{
		Type:  int32(cb.Policy_IMPLICIT_META),
		Value: implicitMetaPolicy,
	}, nil
}

// implicitMetaAnyPolicy defines an implicit meta policy whose sub_policy and key is policyname with rule ANY.
func implicitMetaAnyPolicy(policyName string) (*standardConfigPolicy, error) {
	implicitMetaPolicy, err := implicitMetaPolicy(policyName, cb.ImplicitMetaPolicy_ANY)
	if err != nil {
		return nil, fmt.Errorf("failed to make implicit meta ANY policy: %v", err)
	}

	return &standardConfigPolicy{
		key:   policyName,
		value: implicitMetaPolicy,
	}, nil
}
