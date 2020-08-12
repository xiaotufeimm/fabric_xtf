// +build pkcs11

/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package pkcs11

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/asn1"
	"testing"

	"github.com/hyperledger/fabric/bccsp"
	"github.com/miekg/pkcs11"
	"github.com/stretchr/testify/require"
)

func TestKeyGenFailures(t *testing.T) {
	var testOpts bccsp.KeyGenOpts
	ki := currentBCCSP
	_, err := ki.KeyGen(testOpts)
	require.Error(t, err)
	require.Contains(t, err.Error(), "Invalid Opts parameter. It must not be nil")
}

func TestLoadLib(t *testing.T) {
	// Setup PKCS11 library and provide initial set of values
	lib, pin, label := FindPKCS11Lib()

	// Test for no specified PKCS11 library
	_, _, _, err := loadLib("", pin, label)
	require.Error(t, err)
	require.Contains(t, err.Error(), "No PKCS11 library default")

	// Test for invalid PKCS11 library
	_, _, _, err = loadLib("badLib", pin, label)
	require.Error(t, err)
	require.Contains(t, err.Error(), "Instantiate failed")

	// Test for invalid label
	_, _, _, err = loadLib(lib, pin, "badLabel")
	require.Error(t, err)
	require.Contains(t, err.Error(), "could not find token with label")

	// Test for no pin
	_, _, _, err = loadLib(lib, "", label)
	require.Error(t, err)
	require.Contains(t, err.Error(), "Login failed: pkcs11")
}

func TestNamedCurveFromOID(t *testing.T) {
	// Test for valid P224 elliptic curve
	namedCurve := namedCurveFromOID(oidNamedCurveP224)
	require.Equal(t, elliptic.P224(), namedCurve, "Did not receive expected named curve for oidNamedCurveP224")

	// Test for valid P256 elliptic curve
	namedCurve = namedCurveFromOID(oidNamedCurveP256)
	require.Equal(t, elliptic.P256(), namedCurve, "Did not receive expected named curve for oidNamedCurveP256")

	// Test for valid P256 elliptic curve
	namedCurve = namedCurveFromOID(oidNamedCurveP384)
	require.Equal(t, elliptic.P384(), namedCurve, "Did not receive expected named curve for oidNamedCurveP384")

	// Test for valid P521 elliptic curve
	namedCurve = namedCurveFromOID(oidNamedCurveP521)
	require.Equal(t, elliptic.P521(), namedCurve, "Did not receive expected named curved for oidNamedCurveP521")

	testAsn1Value := asn1.ObjectIdentifier{4, 9, 15, 1}
	namedCurve = namedCurveFromOID(testAsn1Value)
	if namedCurve != nil {
		t.Fatal("Expected nil to be returned.")
	}
}

func TestPKCS11GetSession(t *testing.T) {
	var sessions []pkcs11.SessionHandle
	for i := 0; i < 3*sessionCacheSize; i++ {
		session, err := currentBCCSP.(*impl).getSession()
		require.NoError(t, err)
		sessions = append(sessions, session)
	}

	// Return all sessions, should leave sessionCacheSize cached
	for _, session := range sessions {
		currentBCCSP.(*impl).returnSession(session)
	}
	sessions = nil

	// Lets break OpenSession, so non-cached session cannot be opened
	oldSlot := currentBCCSP.(*impl).slot
	currentBCCSP.(*impl).slot = ^uint(0)

	// Should be able to get sessionCacheSize cached sessions
	for i := 0; i < sessionCacheSize; i++ {
		session, err := currentBCCSP.(*impl).getSession()
		require.NoError(t, err)
		sessions = append(sessions, session)
	}

	_, err := currentBCCSP.(*impl).getSession()
	require.EqualError(t, err, "OpenSession failed: pkcs11: 0x3: CKR_SLOT_ID_INVALID")

	// Load cache with bad sessions
	for i := 0; i < sessionCacheSize; i++ {
		currentBCCSP.(*impl).returnSession(pkcs11.SessionHandle(^uint(0)))
	}

	// Fix OpenSession so non-cached sessions can be opened
	currentBCCSP.(*impl).slot = oldSlot

	// Request a session, return, and re-acquire. The pool should be emptied
	// before creating a new session so when returned, it should be the only
	// session in the cache.
	sess, err := currentBCCSP.(*impl).getSession()
	require.NoError(t, err)
	currentBCCSP.(*impl).returnSession(sess)
	sess2, err := currentBCCSP.(*impl).getSession()
	require.NoError(t, err)
	require.Equal(t, sess, sess2, "expected to get back the same session")

	// Cleanup
	for _, session := range sessions {
		currentBCCSP.(*impl).returnSession(session)
	}
}

func TestPKCS11ECKeySignVerify(t *testing.T) {
	msg1 := []byte("This is my very authentic message")
	msg2 := []byte("This is my very unauthentic message")
	hash1, _ := currentBCCSP.Hash(msg1, &bccsp.SHAOpts{})
	hash2, _ := currentBCCSP.Hash(msg2, &bccsp.SHAOpts{})

	var oid asn1.ObjectIdentifier
	if currentTestConfig.securityLevel == 256 {
		oid = oidNamedCurveP256
	} else if currentTestConfig.securityLevel == 384 {
		oid = oidNamedCurveP384
	}

	key, pubKey, err := currentBCCSP.(*impl).generateECKey(oid, true)
	if err != nil {
		t.Fatalf("Failed generating Key [%s]", err)
	}

	R, S, err := currentBCCSP.(*impl).signP11ECDSA(key, hash1)

	if err != nil {
		t.Fatalf("Failed signing message [%s]", err)
	}

	_, _, err = currentBCCSP.(*impl).signP11ECDSA(nil, hash1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "Private key not found")

	pass, err := currentBCCSP.(*impl).verifyP11ECDSA(key, hash1, R, S, currentTestConfig.securityLevel/8)
	if err != nil {
		t.Fatalf("Error verifying message 1 [%s]", err)
	}
	if pass == false {
		t.Fatal("Signature should match!")
	}

	pass = ecdsa.Verify(pubKey, hash1, R, S)
	if pass == false {
		t.Fatal("Signature should match with software verification!")
	}

	pass, err = currentBCCSP.(*impl).verifyP11ECDSA(key, hash2, R, S, currentTestConfig.securityLevel/8)
	if err != nil {
		t.Fatalf("Error verifying message 2 [%s]", err)
	}

	if pass != false {
		t.Fatal("Signature should not match!")
	}

	pass = ecdsa.Verify(pubKey, hash2, R, S)
	if pass != false {
		t.Fatal("Signature should not match with software verification!")
	}
}
