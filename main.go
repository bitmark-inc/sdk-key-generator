package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/spf13/cobra"
)

var (
	randSource = rand.Reader
	filename   string
)

func main() {
	cmdGenerateKey := &cobra.Command{
		Use:   "generate",
		Short: "generate client id and client secret",
		Long:  "generate certificate, allow server to issue sdk token to use with BitmarkSDK",
		Run:   generateClientKeys,
		Args:  cobra.NoArgs,
	}
	cmdGenerateKey.Flags().StringVarP(&filename, "file", "f", "key.pem", "client certificate file")

	cmdIssueSDKToken := &cobra.Command{
		Use:   "issue [account]",
		Short: "issue new jwt token using generated pem file with an account",
		Run:   issueSDKToken,
		Args:  cobra.ExactArgs(1),
	}
	cmdIssueSDKToken.Flags().StringVarP(&filename, "file", "f", "key.pem", "client certificate file")

	var rootCmd = &cobra.Command{}
	rootCmd.AddCommand(cmdGenerateKey, cmdIssueSDKToken)
	rootCmd.Execute()
}

func generateClientKeys(cmd *cobra.Command, args []string) {
	key, err := rsa.GenerateKey(randSource, 2048)
	checkError(err)

	outFile, err := os.Create(filename)
	checkError(err)
	defer outFile.Close()

	var privateKey = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	err = pem.Encode(outFile, privateKey)
	checkError(err)

	pubkey := x509.MarshalPKCS1PublicKey(&key.PublicKey)
	pubkeyBase64 := base64.StdEncoding.EncodeToString(pubkey)
	fmt.Println("Client ID: ", pubkeyBase64)
	fmt.Println("Saved client secret to: " + filename)
}

func issueSDKToken(cmd *cobra.Command, args []string) {
	account := args[0]

	jwtSecretByte, err := ioutil.ReadFile(filename)
	checkError(err)

	jwtPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM(jwtSecretByte)
	checkError(err)

	now := time.Unix(0, time.Now().UnixNano())
	exp := now.Add(1 * time.Hour)

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub": account,
		"exp": exp.Unix(),
		"iat": now.Unix(),
	})

	tokenString, err := token.SignedString(jwtPrivateKey)
	fmt.Println("SDK token: ", tokenString)
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error: ", err.Error())
		os.Exit(1)
	}
}
