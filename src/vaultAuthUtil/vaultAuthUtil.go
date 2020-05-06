package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type awsKeys struct {
	Access_key string `json:"access_key"`
	Secret_key string `json:"secret_key"`
}

func main() {
	var vaultEnv = os.Getenv("VAULT_TOKEN")
	responseData, err := makeApiCall(vaultEnv)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	var jsonData map[string]interface{}
	if err := json.Unmarshal(responseData, &jsonData); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var accessKey = fmt.Sprintf("%v", jsonData["data"].(map[string]interface{})["access_key"])
	var secretKey = fmt.Sprintf("%v", jsonData["data"].(map[string]interface{})["secret_key"])
	if makeAwsCall(accessKey, secretKey) == true {
		formatJson(accessKey, secretKey)
	}

}

func formatJson(accessKey string, secretKey string) {

	keys := awsKeys{
		Access_key: fmt.Sprintf("%s", accessKey),
		Secret_key: fmt.Sprintf("%s", secretKey)}

	prettyJSON, err := json.Marshal(keys)
	if err != nil {
		os.Exit(1)
	}
	fmt.Printf("%s\n", string(prettyJSON))

}

func makeAwsCall(accessKey string, secretKey string) bool {

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})

	if err != nil {
		return false
	}
	if sess == nil {
		return false
	}

	for i := 0; i < 100; i++ {
		// Create a IAM service client.
		svc := iam.New(sess)
		result, err := svc.ListRoles(&iam.ListRolesInput{})
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		if result != nil {
			break
		}
	}

	return true
}

func makeApiCall(vaultEnv string) ([]byte, error) {
	tr := &http.Transport{
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", "https://vault.rbcloudsandbox.com:8200/v1/aws-sandboxops/creds/poweruser", nil)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	} else {
		req.Header.Add("X-Vault-Token", vaultEnv)
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return responseData, err

}
