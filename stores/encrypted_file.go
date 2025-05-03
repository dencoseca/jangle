package stores

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

// FileStore provides file-based encrypted secret storage functionality.
type FileStore struct {
	filePath string
	key      []byte // 32-byte key for AES-256
	mu       sync.RWMutex
	secrets  map[string]string
}

// NewFileStore initializes a new FileStore with the specified file path and
// 32-byte encryption key.
func NewFileStore(filePath string, encryptionKey []byte) (*FileStore, error) {
	if len(encryptionKey) != 32 {
		return nil, fmt.Errorf("encryption key must be exactly 32 bytes")
	}

	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	fs := &FileStore{
		filePath: filePath,
		key:      encryptionKey,
		secrets:  make(map[string]string),
	}

	// Load existing secrets if the file exists
	if _, err := os.Stat(filePath); err == nil {
		if err := fs.load(); err != nil {
			return nil, fmt.Errorf("failed to load secrets: %w", err)
		}
	}

	return fs, nil
}

// Get retrieves a secret by its name.
func (f *FileStore) Get(name string) (string, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	value, exists := f.secrets[name]
	if !exists {
		return "", fmt.Errorf("secret not found: %s", name)
	}
	return value, nil
}

// Set saves a secret with the given name and value.
func (f *FileStore) Set(name, value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.secrets[name] = value
	return f.save()
}

// Update modifies or creates a secret.
func (f *FileStore) Update(name, value string) error {
	return f.Set(name, value)
}

// Delete removes the specified secret.
func (f *FileStore) Delete(name string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	delete(f.secrets, name)
	return f.save()
}

// List retrieves a list of all secret names.
func (f *FileStore) List() ([]string, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	names := make([]string, 0, len(f.secrets))
	for name := range f.secrets {
		names = append(names, name)
	}
	return names, nil
}

// load decrypts and loads secrets from the file.
func (f *FileStore) load() error {
	data, err := os.ReadFile(f.filePath)
	if err != nil {
		return fmt.Errorf("failed to read secrets file: %w", err)
	}

	if len(data) == 0 {
		f.secrets = make(map[string]string)
		return nil
	}

	decrypted, err := f.decrypt(data)
	if err != nil {
		return fmt.Errorf("failed to decrypt secrets: %w", err)
	}

	if err := json.Unmarshal(decrypted, &f.secrets); err != nil {
		return fmt.Errorf("failed to unmarshal secrets: %w", err)
	}

	return nil
}

// save encrypts and saves secrets to the file.
func (f *FileStore) save() error {
	data, err := json.Marshal(f.secrets)
	if err != nil {
		return fmt.Errorf("failed to marshal secrets: %w", err)
	}

	encrypted, err := f.encrypt(data)
	if err != nil {
		return fmt.Errorf("failed to encrypt secrets: %w", err)
	}

	if err := os.WriteFile(f.filePath, encrypted, 0600); err != nil {
		return fmt.Errorf("failed to write secrets file: %w", err)
	}

	return nil
}

// encrypt and decrypt functions remain the same as in the original implementation
func (f *FileStore) encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(f.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return []byte(base64.StdEncoding.EncodeToString(ciphertext)), nil
}

func (f *FileStore) decrypt(data []byte) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(f.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(decoded) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := decoded[:nonceSize], decoded[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
