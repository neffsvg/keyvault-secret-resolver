package main

import (
	"errors"
	"flag"
	"github.com/joho/godotenv"
	"log"
	"maps"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

const (
	Prefix = "SECRET:"
)

func main() {
	//configPath := flag.String("configPath", ".setupenvconfig", "config file")
	envPath := flag.String("env", ".env", ".env file where to replace secrets with azure keyvault secrets")
	vault := flag.String("vault", "", "keyvault name to get env variables from")
	flag.Parse()

	if vault == nil || len(*vault) == 0 {
		log.Fatal("Missing vault flag, use -h for help.")
	} else if envPath == nil || len(*vault) == 0 {
		log.Fatal("Missing env flag, use -h for help.")
	}

	configMap := loadEnvFile(*envPath)
	secrets := getSecretsFrom(configMap)
	resolvedSecrets := loadSecrets(*vault, secrets)

	// write resolved secrets back into configmap (unique, because map / key)
	maps.Copy(configMap, resolvedSecrets)
	err := godotenv.Write(configMap, *envPath)
	if err != nil {
		log.Fatal("Failed to write into .env file", err)
	} else {
		log.Println("Done. Resolved " + strconv.Itoa(len(resolvedSecrets)) + " secrets.")
	}
}

func loadEnvFile(envPath string) map[string]string {
	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatal("Unable to load .env file, please check file path: ", envPath, err)
	}
	var config map[string]string
	config, err = godotenv.Read()
	if err != nil {
		log.Fatal("Unable to read variables from .env file: ", envPath)
	}
	return config
}

func getSecretsFrom(configMap map[string]string) map[string]string {
	var secrets = make(map[string]string)
	for key, value := range configMap {
		if strings.HasPrefix(value, Prefix) {
			secrets[key] = strings.Trim(value, Prefix)
		}
	}
	return secrets
}

func loadSecrets(vault string, secrets map[string]string) map[string]string {
	wg := sync.WaitGroup{}
	var resolvedSecrets = make(map[string]string)
	for secretKey, secretName := range secrets {
		wg.Add(1)

		go func() {
			defer wg.Done()
			//log.Println("az", "keyvault", "secret", "show", "--vault-name", vault, "--name", secretName, "--query", "value")
			out, err := exec.Command("az", "keyvault", "secret", "show", "--vault-name", vault, "--name", secretName, "--query", "value").Output()
			//out, err := exec.Command("az", "keyvault", "secret", "show", "--vault", "devops-svg-csi-keyvault", "--name", "m-mysvg-dmaut-db-user-name", "--query", "value").Output()
			if err != nil {
				var ee *exec.ExitError
				errors.As(err, &ee)
				log.Println("Could not load secret " + secretName + " -> " + string(ee.Stderr))
				log.Println("Have you called az login first? You might have insufficient permissions on the given keyvault '" + vault + "'\n")
			} else {
				secretValue := strings.Trim(strings.TrimSpace(string(out)), "\"")
				log.Println("Got secret " + secretName + " " + secretValue)
				resolvedSecrets[secretKey] = secretValue
			}
		}()
	}
	wg.Wait()
	return resolvedSecrets
}
