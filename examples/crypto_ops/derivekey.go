package main

import (
	"context"
	"fmt"
	"log"

	"github.com/fortanix/sdkms-client-go/sdkms"
)

func sobjectToString(sobject *sdkms.Sobject) string {
	created, err := sobject.CreatedAt.Std()
	if err != nil {
		log.Fatalf("Failed to convert sobject.CreatedAt: %v", err)
	}
	return fmt.Sprintf("{ %v %#v group(%v) enabled: %v created: %v }",
		*sobject.Kid, *sobject.Name, *sobject.GroupID, sobject.Enabled,
		created.Local())
}

func sample_derive_key(client *sdkms.Client, objId string) {
	var derivedName string = "Derived-Key"
	deriveKeyMec := sdkms.DeriveKeyMechanismHkdf{
		HashAlg: sdkms.DigestAlgorithmSha224,
		Info:    someBytes([]byte("signing")),
		Salt:    someBytes(generateRandom(16)),
	}
	ctx := context.Background()
	var key sdkms.SobjectDescriptor

	if objId != "" {
		key = *sdkms.SobjectByID(objId)
	} else {
		key = *sdkms.SobjectByName(keyName)
	}
	deriveKeyReq := sdkms.DeriveKeyRequest{
		Name:      &derivedName,
		KeyType:   sdkms.ObjectTypeAria,
		KeySize:   128,
		Key:       &key,
		Mechanism: sdkms.DeriveKeyMechanism{Hkdf: &deriveKeyMec},
	}
	deriveKeyResp, err := client.Derive(ctx, deriveKeyReq)
	if err != nil {
		log.Fatalf("Exception in Key Derivation: %v", err)
		return
	}
	log.Printf("Sobject after Derivation: %v", sobjectToString(deriveKeyResp))
}
