package secrets

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/vault/api"
)

const (
	// This is the path were the kubernetes Vault side car container will inject
	// a token into the filesystem of the container this application is running.
	//
	// https://developer.hashicorp.com/vault/docs/platform/k8s/injector/annotations#vault-hashicorp-com-agent-inject-token
	VaultTokenFilePath = "/vault/secrets/token"
)

// A secrets client is used to retrieve, store, and modify key value secrets and certificates.
type SecretsClient interface {
	ReadKeyValue(ctx context.Context, path string) (map[string]any, error)
	WriteKeyValue(ctx context.Context, path string, data map[string]any) error
	DeleteKeyValue(ctx context.Context, path string) error
	RetrieveCert(serialNumber string) (*string, *string, error)
	RetrieveRootCaCert() (*string, error)
	CreateCert(certRequest CertRequest) (serialNumber string, cert string, privateKey string, err error)
	RevokeCert(serialNumber string) error
	GetBoolFromMap(secrets map[string]any, key string) bool
	GetStringFromMap(secrets map[string]any, key string) string
	GetIntFromMap(secrets map[string]any, key string) int
}

type CertRequest struct {
	CommonName string `json:"common_name"`
	AltNames   string `json:"alt_names"`
	IpSans     string `json:"ip_sans"`
	Ttl        string `json:"ttl"`
	KeyType    string `json:"key_type"`
	KeyBits    int    `json:"key_bits"`
}

type Client struct {
	address       string
	username      string
	password      string
	tokenFilePath string
	vaultClient   *api.Client
}

// Creates a new secrets client.
func NewClient(address, username, password, tokenFilePath string) (*Client, error) {
	// Initialize the Vault client
	vaultConnectionConfig := &api.Config{
		Address: address,
	}
	vaultClient, err := api.NewClient(vaultConnectionConfig)
	if err != nil {
		return nil, err
	}

	// Instantiate a secrets client
	secretsClient := &Client{
		address:       address,
		username:      username,
		password:      password,
		tokenFilePath: tokenFilePath,
		vaultClient:   vaultClient,
	}

	// If a token file path is provided, use the token from the file and skip login
	if tokenFilePath != "" {
		token, err := readTokenFromFile(tokenFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read token from file: %v", err)
		}

		// Set the token in the Vault client
		vaultClient.SetToken(token)
		return secretsClient, nil
	}

	// Retry logic for Vault authentication
	const maxRetries = 3
	var loginError error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		vaultClientToken, loginError := secretsClient.login(vaultClient)

		if loginError == nil && vaultClientToken != nil {
			// Successful authentication
			vaultClient.SetToken(*vaultClientToken)
			return secretsClient, nil
		}
	}

	// Return a failed to connect error
	return nil, fmt.Errorf("could not log in to Vault after %d attempts: %v", maxRetries, loginError)
}

// Helper function to read token from file
func readTokenFromFile(tokenFilePath string) (string, error) {
	content, err := os.ReadFile(tokenFilePath)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(content)), nil
}

func (c *Client) refreshTokenFromFile() error {
	if c.tokenFilePath != "" {
		token, err := os.ReadFile(c.tokenFilePath)
		if err != nil {
			return fmt.Errorf("failed to read token from file: %v", err)
		}

		// Set the token in the Vault client
		c.vaultClient.SetToken(strings.TrimSpace(string(token)))
	}

	return nil
}

func (c *Client) login(vaultClient *api.Client) (*string, error) {
	// Authenticate using userpass
	options := map[string]any{
		"password": c.password,
	}

	authPath := fmt.Sprintf("auth/userpass/login/%s", c.username)
	secret, err := vaultClient.Logical().Write(authPath, options)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate with Vault: %v", err)
	}

	if secret == nil || secret.Auth == nil {
		return nil, fmt.Errorf("invalid response from Vault authentication")
	}

	return &secret.Auth.ClientToken, nil
}

func (c *Client) ReadKeyValue(ctx context.Context, path string) (map[string]any, error) {
	// Refresh token before the action
	err := c.refreshTokenFromFile()
	if err != nil {
		return nil, fmt.Errorf("could not refresh token: %s", err)
	}

	// Retrieve the secret
	secret, err := c.vaultClient.Logical().ReadWithContext(ctx, fmt.Sprintf("kv/%s", path))
	if err != nil {
		return nil, fmt.Errorf("could not retrieve secret: %s", err)
	}

	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("secret not found")
	}

	return secret.Data, nil
}

func (c *Client) WriteKeyValue(ctx context.Context, path string, data map[string]any) error {
	// Refresh token before the action
	err := c.refreshTokenFromFile()
	if err != nil {
		return fmt.Errorf("could not refresh token: %s", err)
	}

	// Write the KV secret
	_, err = c.vaultClient.Logical().WriteWithContext(ctx, fmt.Sprintf("kv/%s", path), data)
	if err != nil {
		return fmt.Errorf("could not create secret: %s", err)
	}

	return nil
}

func (c *Client) DeleteKeyValue(ctx context.Context, path string) error {
	// Refresh token before the action
	err := c.refreshTokenFromFile()
	if err != nil {
		return fmt.Errorf("could not refresh token: %s", err)
	}

	// Write the KV secret
	_, err = c.vaultClient.Logical().DeleteWithContext(ctx, fmt.Sprintf("kv/%s", path))
	if err != nil {
		return fmt.Errorf("could not delete secret: %s", err)
	}

	return nil
}

func (c *Client) RetrieveCert(serialNumber string) (*string, *string, error) {
	// Check that the serial number is not empty
	if serialNumber == "" {
		return nil, nil, fmt.Errorf("certificate serial number is empty")
	}

	// Refresh token before the action
	err := c.refreshTokenFromFile()
	if err != nil {
		return nil, nil, fmt.Errorf("could not refresh token: %s", err)
	}

	// Retrieve the certificate
	certPath := fmt.Sprintf("certs/%s", serialNumber)
	certSecretData, err := c.ReadKeyValue(context.Background(), certPath)
	if err != nil {
		return nil, nil, err
	}
	if certSecretData == nil {
		return nil, nil, fmt.Errorf("unknown error when retrieving the certificate")
	}

	certificate, ok := certSecretData["certificate"].(string)
	if !ok {
		return nil, nil, fmt.Errorf("certificate not found in the retrieved secret")
	}

	privateKey, ok := certSecretData["private_key"].(string)
	if !ok {
		return nil, nil, fmt.Errorf("private key not found in the retrieved secret")
	}

	return &certificate, &privateKey, nil
}

func (c *Client) RetrieveRootCaCert() (*string, error) {
	// Refresh token before the action
	err := c.refreshTokenFromFile()
	if err != nil {
		return nil, fmt.Errorf("could not refresh token: %s", err)
	}

	caCertSecret, err := c.vaultClient.Logical().Read("pki/cert/ca")
	if err != nil {
		return nil, fmt.Errorf("could not retrieve CA cert: %s", err)
	}

	certificate := caCertSecret.Data["certificate"].(string)
	return &certificate, nil
}

func (c *Client) CreateCert(certRequest CertRequest) (serialNumber, cert, privateKey string, err error) {
	// Refresh token before the action
	err = c.refreshTokenFromFile()
	if err != nil {
		err = fmt.Errorf("could not refresh token: %s", err)
		return
	}

	// Transform CertRequest to Vault's expected format
	keyBitsString := strconv.FormatInt(int64(certRequest.KeyBits), 10)

	certificateRequest := map[string]any{
		"common_name": certRequest.CommonName,
		"alt_names":   certRequest.AltNames,
		"ip_sans":     certRequest.IpSans,
		"ttl":         certRequest.Ttl,
		"key_bits":    keyBitsString,
		"key_type":    certRequest.KeyType,
	}

	// Create the cert
	certPath := fmt.Sprintf("pki/issue/%s", "mqtt-client")
	certSecret, err := c.vaultClient.Logical().Write(certPath, certificateRequest)
	if err != nil {
		return
	}

	// Extract the certificate, private key, and issuing CA from the response
	serialNumber = certSecret.Data["serial_number"].(string)
	cert = certSecret.Data["certificate"].(string)
	privateKey = certSecret.Data["private_key"].(string)

	// Store the certificate and private key in the KV engine
	err = c.WriteKeyValue(context.Background(), fmt.Sprintf("certs/%s", serialNumber), map[string]any{
		"certificate": cert,
		"private_key": privateKey,
		// "issuing_ca": ca, // TODO: support dynamic CAs in this secret
	})

	return
}

func (c *Client) RevokeCert(path string) error {
	return nil
}

func (c *Client) GetBoolFromMap(secrets map[string]any, key string) bool {
	// Check if the key exists and is not nil
	if val, ok := secrets[key]; ok && val != nil {
		// Attempt to perform a type assertion to string
		if strVal, ok := val.(string); ok {
			// Trim any surrounding quotes from the value
			strVal = strings.Trim(strVal, `"`)

			// Convert the string to a boolean based on common true/false representations
			switch strings.ToLower(strVal) {
			case "true", "1":
				return true
			case "false", "0":
				return false
			default:
				log.Printf("Invalid boolean value for key %s: %s", key, strVal)
			}
		} else {
			log.Printf("Value for key %s is not a string", key)
		}
	}
	return false
}

func (c *Client) GetStringFromMap(secrets map[string]any, key string) string {
	// Check if the key exists and is not nil
	if val, ok := secrets[key]; ok && val != nil {
		strVal := val.(string)

		// Remove surrounding quotes if they exist
		strVal = strings.Trim(strVal, `"`)

		return strVal
	}

	return ""
}

func (c *Client) GetIntFromMap(secrets map[string]any, key string) int {
	// Check if the key exists and is not nil
	if val, ok := secrets[key]; ok && val != nil {
		// Attempt to perform a type assertion to string
		if strVal, ok := val.(string); ok {
			// Trim any surrounding quotes from the value
			strVal = strings.Trim(strVal, `"`)

			// Convert the string to an integer
			if intVal, err := strconv.Atoi(strVal); err == nil {
				return intVal
			} else {
				log.Printf("Error converting %s to integer: %v", key, err)
			}
		} else {
			log.Printf("Value for key %s is not a string", key)
		}
	}

	return 0 // return a default value or handle the missing case
}
