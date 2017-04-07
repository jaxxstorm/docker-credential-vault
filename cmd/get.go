// Copyright © 2017 Lee Briggs <lee@leebriggs.co.uk>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	//"github.com/davecgh/go-spew/spew"

	v "github.com/jaxxstorm/docker-credential-vault/vault"
)

type DockerLoginCredentials struct {
	ServerURL string `mapstructure:"serverurl"`
	Username  string `mapstructure:"username"`
	Secret    string `mapstructure:"password"`
}

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get credentials from vault",
	Run: func(cmd *cobra.Command, args []string) {

		// read in stdin
		bytes, err := ioutil.ReadAll(os.Stdin)

		// if stdin can't be read, bomb
		if err != nil {
			fmt.Printf("Error reading stdin", err)
		}

		// create a base64 encoded url from stdin
		url := b64.StdEncoding.EncodeToString(bytes)

		vaultHost = viper.GetString("vault")
		vaultToken = viper.GetString("token")

		// create a vault client
		client, err := v.New(vaultToken, vaultHost, vaultPort)

		if err != nil {
			fmt.Printf("Error creating vault client", err)
		}

		secret, err := client.Logical().Read("secret/" + url)

		if err != nil {
			fmt.Printf("Error reading vault creds", err)
		}

		if secret == nil {
			fmt.Printf("Error retrieving secret", err)
		}

		var creds DockerLoginCredentials

		if err := mapstructure.Decode(secret.Data, &creds); err != nil {
			fmt.Printf("Error parsing vault response: ", err)
		}

		jsonCreds, _ := json.Marshal(creds)

		fmt.Printf(string(jsonCreds))
	},
}

func init() {
	RootCmd.AddCommand(getCmd)
}
