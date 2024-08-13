package main

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"github.com/joho/godotenv"
	"log"
	"sync"
)
import flag "github.com/spf13/pflag"

func main() {
	secretsSpecFile := flag.StringP("secrets-env-file", "s", "", "secrets.env file where secrets are specified with their keyvault names")
	resultPath := flag.StringP("result-env-file-path", "r", ".env", ".env file to put in config and resolved secrets")
	vault := flag.StringP("keyvault", "k", "", "keyvault name to get env variables from")
	flag.Parse()

	// check if all flags are set
	if vault == nil || len(*vault) == 0 {
		log.Fatal("Missing vault flag, use -h for help.")
	} else if secretsSpecFile == nil || len(*secretsSpecFile) == 0 {
		log.Fatal("Missing secrets-env flag, use -h for help.")
	} else if resultPath == nil || len(*resultPath) == 0 {
		log.Fatal("Missing result-path flag, use -h for help.")
	}

	var secretsMap map[string]string
	secretsMap, readError := godotenv.Read(*secretsSpecFile)
	if readError != nil {
		log.Fatal("Unable to read variables from .env file:", *secretsSpecFile, "\n", readError)
	}

	client, clientError := createAzureClient(*vault)
	if clientError != nil {
		log.Fatal(clientError)
	}
	resolvedSecrets := loadSecrets(client, secretsMap)

	writeError := godotenv.Write(resolvedSecrets, *resultPath)
	if writeError != nil {
		log.Fatal("Failed to write into .env file. " + writeError.Error())
	} else {
		log.Printf("Done. Resolved %v of %v secrets.\n", len(resolvedSecrets), len(secretsMap))
	}
}

func loadSecrets(client *azsecrets.Client, secrets map[string]string) map[string]string {
	wg := sync.WaitGroup{}
	var resolvedSecrets = make(map[string]string)
	for secretKey, secretName := range secrets {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// empty version, will omit version parameter in api call and defaults to latest version
			secretResponse, err := client.GetSecret(context.TODO(), secretName, "", nil)
			if err != nil {
				log.Println("Could not resolve secret:", secretName, "\n", err)
			} else {
				secretValue := *secretResponse.Value
				log.Println("Resolved secret:", secretName)
				resolvedSecrets[secretKey] = secretValue
			}
		}()
	}
	wg.Wait()
	return resolvedSecrets
}

func createAzureClient(vaultName string) (*azsecrets.Client, error) {
	credentials, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, errors.New("Unable to authenticate with azure, run `az login` in terminal. " + err.Error())
	}
	vaultURL := "https://" + vaultName + ".vault.azure.net/"
	client, err := azsecrets.NewClient(vaultURL, credentials, nil)
	if err != nil {
		return nil, errors.New("Unable to create azure client / connect to key vault. " + err.Error())
	} else {
		return client, nil
	}
}
