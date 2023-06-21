package skm

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	mrand "math/rand"
)

func (key *Key) genKey() {

}

const (
	algRsa     = "rsa"
	algEd25519 = "ed25519"
)

var algorithms = []string{algRsa, algEd25519}

// genKey 生成密钥对
// 返回值(私钥,公钥,错误信息)
func genKey(alg string) ([]byte, []byte, error) {
	if alg == algRsa {
		return genRsaKey(3072, "")
	} else if alg == algEd25519 {
		return genEd25519Key()
	}
	return nil, nil, errors.New(fmt.Sprintf("No algorithm: %s", alg))
}

// genRsaKey 生成rsa密钥对
// 返回值(私钥,公钥,错误信息)
func genRsaKey(bits int, passphrase string) ([]byte, []byte, error) {

	if bits%1024 != 0 {
		return nil, nil, errors.New("bits must be a multiple of 1024")
	}

	//生成私钥和公钥
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, errors.New("Failed to generate key pair. ")
	}
	publicKey := &privateKey.PublicKey

	//私钥格式化
	privateKeyPem := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	if passphrase != "" {
		privateKeyPem, err = x509.EncryptPEMBlock(rand.Reader, privateKeyPem.Type, privateKeyPem.Bytes, []byte(passphrase), x509.PEMCipherAES256)
		if err != nil {
			return nil, nil, errors.New("Failed to generate key pair. ")
		}
	}
	privateKeyOutput := pem.EncodeToMemory(privateKeyPem)

	//公钥格式化
	sshPublicKey, err := ssh.NewPublicKey(publicKey)
	publicKeyOutput := ssh.MarshalAuthorizedKey(sshPublicKey)

	return privateKeyOutput, publicKeyOutput, nil

}

// genEd25519Key 生成ed25519密钥对
// 返回值(私钥,公钥,错误信息)
func genEd25519Key() ([]byte, []byte, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, errors.New("Failed to generate key pair. ")
	}

	//私钥格式化
	privateKeyPem := &pem.Block{
		Type:  "OPENSSH PRIVATE KEY",
		Bytes: MarshalED25519PrivateKey(privateKey),
	}
	privateKeyOutput := pem.EncodeToMemory(privateKeyPem)

	//公钥格式化
	sshPublicKey, err := ssh.NewPublicKey(publicKey)
	if err != nil {
		return nil, nil, errors.New("Failed to generate key pair. ")
	}
	publicKeyOutput := ssh.MarshalAuthorizedKey(sshPublicKey)

	return privateKeyOutput, publicKeyOutput, nil

}

// MarshalED25519PrivateKey
// From: https://github.com/mondoohq/cnquery/blob/main/motor/providers/ssh/keypair/edkey.go
func MarshalED25519PrivateKey(key ed25519.PrivateKey) []byte {
	// Add our key header (followed by a null byte)
	magic := append([]byte("openssh-key-v1"), 0)

	var w struct {
		CipherName   string
		KdfName      string
		KdfOpts      string
		NumKeys      uint32
		PubKey       []byte
		PrivKeyBlock []byte
	}

	// Fill out the private key fields
	pk1 := struct {
		Check1  uint32
		Check2  uint32
		Keytype string
		Pub     []byte
		Priv    []byte
		Comment string
		Pad     []byte `ssh:"rest"`
	}{}

	// Set our check ints
	ci := mrand.Uint32()
	pk1.Check1 = ci
	pk1.Check2 = ci

	// Set our key type
	pk1.Keytype = ssh.KeyAlgoED25519

	// Add the pubkey to the optionally-encrypted block
	pk, ok := key.Public().(ed25519.PublicKey)
	if !ok {
		// fmt.Fprintln(os.Stderr, "ed25519.PublicKey type assertion failed on an ed25519 public key. This should never ever happen.")
		return nil
	}
	pubKey := []byte(pk)
	pk1.Pub = pubKey

	// Add our private key
	pk1.Priv = []byte(key)

	// Might be useful to put something in here at some point
	pk1.Comment = ""

	// Add some padding to match the encryption block size within PrivKeyBlock (without Pad field)
	// 8 doesn't match the documentation, but that's what ssh-keygen uses for unencrypted keys. *shrug*
	bs := 8
	blockLen := len(ssh.Marshal(pk1))
	padLen := (bs - (blockLen % bs)) % bs
	pk1.Pad = make([]byte, padLen)

	// Padding is a sequence of bytes like: 1, 2, 3...
	for i := 0; i < padLen; i++ {
		pk1.Pad[i] = byte(i + 1)
	}

	// Generate the pubkey prefix "\0\0\0\nssh-ed25519\0\0\0 "
	prefix := []byte{0x0, 0x0, 0x0, 0x0b}
	prefix = append(prefix, []byte(ssh.KeyAlgoED25519)...)
	prefix = append(prefix, []byte{0x0, 0x0, 0x0, 0x20}...)

	// Only going to support unencrypted keys for now
	w.CipherName = "none"
	w.KdfName = "none"
	w.KdfOpts = ""
	w.NumKeys = 1
	w.PubKey = append(prefix, pubKey...)
	w.PrivKeyBlock = ssh.Marshal(pk1)

	magic = append(magic, ssh.Marshal(w)...)

	return magic
}
