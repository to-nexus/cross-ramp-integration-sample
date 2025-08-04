package services

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
)

type KeystoreService struct {
	key *ecdsa.PrivateKey
}

func NewKeystoreService() *KeystoreService {
	// TODO: keyStore and passphrase must be loaded from file or env
	keyStore := `{"address":"100cbc7ac2abdb4e75d8e08c6842d1dd8c04df73","crypto":{"cipher":"aes-128-ctr","ciphertext":"ddd3ee2e1eae8a058485146160617d5439f57ab0e900fc68a7632c701315d129","cipherparams":{"iv":"b97e245d56a50673856f3b49a81624a5"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"8c85921c88c4a67c974f4399f046c5ec2dffba9f722e57762508ed161bbe9740"},"mac":"145ca75eb32d366ea108af62ed47f41c04a348ada383304d1995808eb36e9365"},"id":"3b850e08-41a5-49ec-a13e-70a95e1a448e","version":3}`
	passphrase := "strong_password"

	store, err := keystore.DecryptKey([]byte(keyStore), passphrase)
	if err != nil {
		panic(err)
	}

	return &KeystoreService{
		key: store.PrivateKey,
	}
}

func (s *KeystoreService) Sign(digest []byte) ([]byte, error) {
	signature, err := crypto.Sign(digest, s.key)
	if err != nil {
		return nil, err
	}

	signature[64] += 27

	return signature, nil
}
