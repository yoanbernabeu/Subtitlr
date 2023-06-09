/*
Copyright © 2023 Yoan Bernabeu <contact@yoandev.co>

*/
package cmd

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kkdai/youtube/v2"
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Command to start the translation in SRT format",
	Long: `The 'generate' command is a crucial feature of our subtitle generator application. Once activated, this command initiates the process of generating subtitles from a provided YouTube video ID
	The 'generate' command is a crucial feature of our subtitle generator application. Once activated, this command initiates the process of generating subtitles from a provided YouTube video ID.`,
	Run: func(cmd *cobra.Command, args []string) {
		/* Variables declaration */
		id, _ := cmd.Flags().GetString("id")
		file, _ := cmd.Flags().GetString("file")
		lang, _ := cmd.Flags().GetString("lang")
		output, _ := cmd.Flags().GetString("output")
		apiKey, _ := cmd.Flags().GetString("apiKey")

		/* Displaying the values of the flags */
		fmt.Println("---------------------------------------")
		fmt.Println("You have entered the following values:")
		fmt.Println("id:", id)
		fmt.Println("file:", file)
		fmt.Println("lang:", lang)
		fmt.Println("output:", output)
		fmt.Println("apiKey:", apiKey)
		fmt.Println("---------------------------------------")

		/* Vérification si id et file sont tout les 2 données */
		if (id == "" && file == "") || (id != "" && file != "") {
			fmt.Println("Error: Either 'id' or 'file' flag must be provided, but not both.")
			return
		}
		
		/* Calling the function to generate the subtitles */
		if id != "" {
			generateSubtitles(id, lang, output, apiKey)
			return
		}
		if file != "" {
			verification(file, lang, output, apiKey)
			return
		}

		/* Displaying the success message */
		fmt.Println("---------------------------------------")
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringP("id", "", "", "YouTube video ID (or file)")
	generateCmd.Flags().StringP("file", "", "","Audio file MP3 (or id)")
	generateCmd.Flags().StringP("lang", "l", "fr", "Language (in ISO 639-1 format) speaking in the video")
	generateCmd.Flags().StringP("output", "o", "output.srt", "Output file")

	//generateCmd.MarkFlagRequired("url")
	generateCmd.MarkFlagRequired("lang")
	generateCmd.MarkFlagRequired("output")

	manageApiKeyEnv()
}

func manageApiKeyEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(".env file not found")
		generateCmd.Flags().String("apiKey", "", "OpenAI API key")
		return
	}

	apiKey := os.Getenv("OPENAI_API_KEY")

	if apiKey == "" {
		fmt.Println("Error: OPENAI_API_KEY not found in .env file")
		generateCmd.Flags().String("apiKey", "", "OpenAI API key")
		generateCmd.MarkFlagRequired("apiKey")
		return
	}

	generateCmd.Flags().String("apiKey", apiKey, "OpenAI API key")
}

func generateSubtitles(id string, lang string, output string, apiKey string) {
	/* Downloading the video */
	downloadVideo(id)

	/* Converting the video to audio */
	convertVideoToAudio()

	/* Generating the subtitles */
	generateSubtitlesFromAudio(lang, output, apiKey)

	/* Deleting the temp folder */
	os.RemoveAll("temp")
}

func downloadVideo(id string) {
	fmt.Println("Downloading the video...")
	videoID := id
	client := youtube.Client{}

	video, err := client.GetVideo(videoID)
	if err != nil {
		panic(err)
	}

	formats := video.Formats.WithAudioChannels() // only get videos with audio
	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		panic(err)
	}

	if _, err := os.Stat("temp"); os.IsNotExist(err) {
		os.Mkdir("temp", 0755)
	}

	file, err := os.Create("temp/temp.mp4")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		panic(err)
	}

	fmt.Println("Video downloaded successfully!")
}

func convertVideoToAudio() {
	fmt.Println("Converting the video to audio...")
	/* extract audio from video */
	file, err := os.Open("temp/temp.mp4")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	/* create the output file */
	out, err := os.Create("temp/temp.mp3")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	/* convert to mp3 with ffmpeg with big compression */
	ffmpeg := exec.Command("ffmpeg", "-i", "pipe:0", "-f", "mp3", "-ab", "64k", "-vn", "pipe:1")
	ffmpeg.Stdin = file
	ffmpeg.Stdout = out
	err = ffmpeg.Run()
	if err != nil {
		panic(err)
	}

	fmt.Println("Video converted to audio successfully!")
}

func generateSubtitlesFromAudio(lang string, output string, apiKey string) {
	fmt.Println("Generating the subtitles...")

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Open the file
	file, err := os.Open("temp/temp.mp3")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Add the file to the request
	fw, err := w.CreateFormFile("file", file.Name())
	if err != nil {
		fmt.Println(err)
		return
	}
	if _, err = io.Copy(fw, file); err != nil {
		fmt.Println(err)
		return
	}

	// Add the model to the request
	if err = w.WriteField("model", "whisper-1"); err != nil {
		fmt.Println(err)
		return
	}

	// Add the response format to the request
	if err = w.WriteField("response_format", "srt"); err != nil {
		fmt.Println(err)
		return
	}

	// Add the language to the request
	if err = w.WriteField("language", lang); err != nil {
		fmt.Println(err)
		return
	}

	// Close the request
	if err = w.Close(); err != nil {
		fmt.Println(err)
		return
	}

	// Create the request
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/audio/transcriptions", &b)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Send the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	// Check the response
	if res.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Write the response to the output file
		err = ioutil.WriteFile(output, bodyBytes, 0644)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Subtitles generated successfully!")
	} else {
		fmt.Println("Request failed with status:", res.StatusCode)
	}
}

func generateSubtitlesFromAudioWithMp3(file string, lang string, output string, apiKey string) {
	fmt.Println("Generating the subtitles...")

	audioBytes, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("Error reading MP3 file:", err)
		return
	}

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Add the audio file to the request
	fw, err := w.CreateFormFile("file", "temp.mp3")
	if err != nil {
		fmt.Println("Error creating form file:", err)
		return
	}
	_, err = fw.Write(audioBytes)
	if err != nil {
		fmt.Println("Error writing audio file to form:", err)
		return
	}

	// Add the model to the request
	if err = w.WriteField("model", "whisper-1"); err != nil {
		fmt.Println(err)
		return
	}

	// Add the response format to the request
	if err = w.WriteField("response_format", "srt"); err != nil {
		fmt.Println(err)
		return
	}

	// Add the language to the request
	if err = w.WriteField("language", lang); err != nil {
		fmt.Println(err)
		return
	}

	// Close the request
	if err = w.Close(); err != nil {
		fmt.Println(err)
		return
	}

	// Create the request
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/audio/transcriptions", &b)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Send the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	// Check the response
	if res.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Write the response to the output file
		err = ioutil.WriteFile(output, bodyBytes, 0644)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Subtitles generated successfully!")
	} else {
		fmt.Println("Request failed with status:", res.StatusCode)
	}
}

func verification(file string, lang string, output string, apiKey string) {
	/* vérification de l'existence du fichier */
	if _, err := os.Stat(file); os.IsNotExist(err) {
		fmt.Println("Error: File not found.")
		return
	}
	/* vérification du fichier audio */
	if strings.HasSuffix(file, ".mp3") == true {
		generateSubtitlesFromAudioWithMp3(file, lang, output, apiKey)
		return
	}

}