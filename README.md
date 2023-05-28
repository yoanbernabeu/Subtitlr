# Subtitlr

A subtitle generator for Youtube (CLI)

## Description

This application, a subtitle generator for YouTube, utilizes OpenAI's Whisper API.

This tool leverages artificial intelligence to efficiently transcribe speech in YouTube videos into text, thereby generating accurate subtitles (in SRT format).

It's designed to improve the accessibility and convenience of video content, ensuring that no matter your language or hearing ability, you can fully engage with and comprehend the material.

## Usage

```bash
Subtitlr generate --id qJpR1NBx4cU --lang fr --output output.srt --apiKey sk-****************************
```

## Requirements

* [OpenAI API key](https://beta.openai.com/)
* [FFmpeg](https://ffmpeg.org/)
* Linux (tested on Ubuntu 22.04)

## Parameters

| Name | Description | Required |
| --- | --- | --- |
| id | Youtube video id | true |
| lang | Language speaking in the video | true |
| output | Output file | true |
| apiKey | OpenAI API key | true |

## Installation

### From binary

* Linux

```bash
# With wget
wget -qO- https://raw.githubusercontent.com/yoanbernabeu/Subtitlr/main/install.sh | bash
```

```bash
# With curl
curl -sL https://raw.githubusercontent.com/yoanbernabeu/Subtitlr/main/install.sh | bash
```

### From source

Subtitlr is written in Go, so you need to install it first.

```bash
git clone git@github.com:yoanbernabeu/Subtitlr.git
cd Subtitlr
go build -o Subtitlr
```

## License

[MIT](LICENSE)

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.