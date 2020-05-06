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
	access_key   string      `json:"access_key"`
	secret_key   string      `json:"secret_key"`
}

func main() {
	var vaultEnv = os.Getenv("VAULT_TOKEN")
	responseData, err := makeApiCall(vaultEnv)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	var jsonData map[string]interface{}
	if err := json.Unmarshal(responseData, &jsonData); err != nil{
		fmt.Println(err)
		os.Exit(1)
	}

	var accessKey string  = fmt.Sprintf("%v", jsonData["data"].(map[string]interface{})["access_key"])
	var secretKey string = fmt.Sprintf("%v", jsonData["data"].(map[string]interface{})["secret_key"])
	fmt.Println("Calling makeAwsCall")
	if makeAwsCall(accessKey, secretKey) == true {

	}

}

func formatJson(accessKey string, secretKey string) {
	keys := &awsKeys{
		access_key:   accessKey,
		secret_key: secretKey}
	jsonKeys, _ := json.Marshal(keys)
	fmt.Println(jsonKeys)
}

func makeAwsCall( accessKey string, secretKey string ) bool {

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey,""),
	})

	if err != nil {
		return false
	}
	if  sess == nil {
		return false
	}

	for i := 0; i < 100; i++ {
		// Create a IAM service client.
		svc := iam.New(sess)
		result, err := svc.ListRoles(&iam.ListRolesInput{})
		if  err != nil {
			fmt.Println(err.Error())
			time.Sleep(time.Second)
			continue
		}
		if result != nil{
			fmt.Println("Success")
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
	req.Header.Add("X-Vault-Token", vaultEnv)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
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
