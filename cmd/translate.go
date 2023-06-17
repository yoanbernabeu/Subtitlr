/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// translateCmd represents the translate command
var translateCmd = &cobra.Command{
	Use:   "translate",
	Short: "Translate the subtitles in SRT format with DeepL",
	Long:  `The 'translate' command is a crucial feature of our subtitle generator application. Once activated, this command initiates the process of translating subtitles from a provided SRT file`,
	Run: func(cmd *cobra.Command, args []string) {
		/* Variables declaration */
		input, _ := cmd.Flags().GetString("input")
		lang, _ := cmd.Flags().GetString("lang")
		output, _ := cmd.Flags().GetString("output")
		apiKeyDeepL, _ := cmd.Flags().GetString("apiKeyDeepL")

		/* Displaying the values of the flags */
		fmt.Println("---------------------------------------")
		fmt.Println("You have entered the following values:")
		fmt.Println("input:", input)
		fmt.Println("lang:", lang)
		fmt.Println("output:", output)
		fmt.Println("apiKeyDeepL:", apiKeyDeepL)
		fmt.Println("---------------------------------------")

		/* Uploading the file to DeepL */
		deepLResponse := uploadFile(input, lang, apiKeyDeepL)

		/* Checking the status of the translation */
		fmt.Println("=======================================")

		for {
			checkStatusResponse := checkStatus(deepLResponse, apiKeyDeepL)
			if checkStatusResponse.Status == "done" {
				downloadFile(deepLResponse, output, apiKeyDeepL)
				break
			}
		}

	},
}

type DeepLResponse struct {
	DocumentID  string `json:"document_id"`
	DocumentKey string `json:"document_key"`
}

type DeepLCheckStatusResponse struct {
	DocumentID       string `json:"document_id"`
	Status           string `json:"status"`
	SecondsRemaining int    `json:"seconds_remaining"`
}

func init() {
	rootCmd.AddCommand(translateCmd)

	translateCmd.Flags().String("input", "", "Input file")
	translateCmd.Flags().String("lang", "", "Language (in ISO 639-1 format) to translate the subtitles to")
	translateCmd.Flags().String("output", "", "Output file")
	translateCmd.Flags().String("apiKeyDeepL", "", "DeepL API key")

	translateCmd.MarkFlagRequired("input")
	translateCmd.MarkFlagRequired("lang")
	translateCmd.MarkFlagRequired("output")
	translateCmd.MarkFlagRequired("apiKeyDeepL")
}

func uploadFile(input string, lang string, apiKeyDeepL string) DeepLResponse {
	fmt.Println("=======================================")
	fmt.Println("Uploading the file to DeepL...")

	// Create a buffer to store our request body as bytes
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Open the file
	f, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Add the file to the request body
	fw, err := w.CreateFormFile("file", f.Name()+".txt")
	if err != nil {
		log.Fatal(err)
	}
	if _, err = io.Copy(fw, f); err != nil {
		log.Fatal(err)
	}

	// Add our fields to the multipart writer
	if err = w.WriteField("target_lang", lang); err != nil {
		log.Fatal(err)
	}

	// Close the multipart writer
	if err = w.Close(); err != nil {
		log.Fatal(err)
	}

	// Create a new request
	req, _ := http.NewRequest("POST", "https://api-free.deepl.com/v2/document", &b)

	// Set the content type header, as well as the boundary we're going to use
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Set the authorization header (Authorization: DeepL-Auth-Key [yourAuthKey])
	req.Header.Set("Authorization", "DeepL-Auth-Key "+apiKeyDeepL)

	// Set the user agent header (User-Agent: YourApp/1.2.3)
	req.Header.Set("User-Agent", "Subtitlr")

	// Send the request
	client := &http.Client{}
	resp, _ := client.Do(req)

	// Read the response body
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(resp.Body)

	// Close the response body
	_ = resp.Body.Close()

	// Parse the response body into a DeepLResponse struct
	var deepLResponse DeepLResponse
	json.Unmarshal(buf.Bytes(), &deepLResponse)

	// Display the document ID and document key
	fmt.Println("DeepL response:")
	fmt.Println("Document ID:", deepLResponse.DocumentID)
	fmt.Println("Document Key:", deepLResponse.DocumentKey)
	fmt.Println("=======================================")

	return deepLResponse
}

func checkStatus(deepLResponse DeepLResponse, apiKeyDeepL string) DeepLCheckStatusResponse {
	// Create a new request
	req, _ := http.NewRequest("POST", "https://api-free.deepl.com/v2/document/"+deepLResponse.DocumentID, nil)
	req.Header.Set("Authorization", "DeepL-Auth-Key "+apiKeyDeepL)
	req.Header.Set("User-Agent", "Subtitlr")

	// Add the document key to the request body
	q := req.URL.Query()
	q.Add("document_key", deepLResponse.DocumentKey)
	req.URL.RawQuery = q.Encode()

	// Send the request
	client := &http.Client{}
	resp, _ := client.Do(req)

	// Read the response body
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(resp.Body)

	// Close the response body
	_ = resp.Body.Close()

	// Parse the response body into a DeepLCheckStatusResponse struct
	var deepLCheckStatusResponse DeepLCheckStatusResponse
	json.Unmarshal(buf.Bytes(), &deepLCheckStatusResponse)

	fmt.Println("Seconds remaining:", deepLCheckStatusResponse.SecondsRemaining)

	return deepLCheckStatusResponse
}

func downloadFile(deepLResponse DeepLResponse, output string, apiKeyDeepL string) {
	fmt.Println("=======================================")
	fmt.Println("Downloading the translated file...")
	fmt.Println("Document ID:", deepLResponse.DocumentID)
	fmt.Println("Document Key:", deepLResponse.DocumentKey)

	// Create a new request
	req, _ := http.NewRequest("POST", "https://api-free.deepl.com/v2/document/"+deepLResponse.DocumentID+"/result", nil)
	req.Header.Set("Authorization", "DeepL-Auth-Key "+apiKeyDeepL)
	req.Header.Set("User-Agent", "Subtitlr")

	// Add the document key to the request body
	q := req.URL.Query()
	q.Add("document_key", deepLResponse.DocumentKey)
	req.URL.RawQuery = q.Encode()

	// Send the request
	client := &http.Client{}
	resp, _ := client.Do(req)

	// Read the response body
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(resp.Body)

	// Close the response body
	_ = resp.Body.Close()

	// Write the response body to the output file
	err := ioutil.WriteFile(output, buf.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}

	// Display the status of the translation
	fmt.Println("DeepL response:")
	fmt.Println("File downloaded successfully!")
	fmt.Println("=======================================")
}
