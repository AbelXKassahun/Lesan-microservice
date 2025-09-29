# Amharic Speech Service

A FastAPI-based web service for Amharic Text-to-Speech (TTS) and Speech-to-Text (STT) conversion.

## Features

- **Speech-to-Text**: Convert Amharic audio files to text using SpeechBrain's Wav2Vec2 model
- **Text-to-Speech**: Convert Amharic text to speech using Microsoft Edge TTS

## Installation

0. (Recommended) Create and activate a virtual environment:

```bash
python -m venv venv
venv\Scripts\activate
```

1. Install dependencies:

```bash
pip install -r requirements.txt
```

2. Run the application:

```bash
python main.py
```

The API will be available at `http://localhost:8000`

## API Endpoints

### Health Check

- `GET /health` - Check service status and model availability

### Speech-to-Text

- `POST /speech-to-text` - Upload audio file (max 10MB) for transcription

### Text-to-Speech

- `POST /text-to-speech` - Convert text to speech audio
- Supported voices: `am-ET-MekdesNeural` (female), `am-ET-AmehaNeural` (male)

## Usage Example

```bash
# Convert text to speech
curl -X POST "http://localhost:8000/text-to-speech" \
     -H "Content-Type: application/json" \
     -d '{"text": "ሰላም ወንድሜ", "voice": "am-ET-MekdesNeural"}'

# Convert speech to text
curl -X POST "http://localhost:8000/speech-to-text" \
     -F "audio_file=@your_audio.wav"
```

## Requirements

- Python 3.7+
- Internet connection (for TTS model)
- Audio files in supported formats (WAV, MP3, etc.)
