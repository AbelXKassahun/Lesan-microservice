from fastapi import (
    FastAPI,
    File,
    UploadFile,
    HTTPException,
)
from fastapi.responses import FileResponse
from speechbrain.inference.ASR import EncoderASR
import edge_tts
from pathlib import Path
import uuid
import logging
import uvicorn
from contextlib import asynccontextmanager
from pydantic import BaseModel, field_validator
import os
from pydub import AudioSegment

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

asr_model = None

# Voice mapping for gender-based selection
VOICE_MAPPING = {
    "male": "am-ET-AmehaNeural",
    "female": "am-ET-MekdesNeural",
}


@asynccontextmanager
async def lifespan(app: FastAPI):
    # Startup
    global asr_model
    try:
        logger.info("Loading ASR model...")
        asr_model = EncoderASR.from_hparams(
            source="speechbrain/asr-wav2vec2-dvoice-amharic",
            savedir="pretrained_models/asr-wav2vec2-dvoice-amharic",
        )
        logger.info("ASR model loaded successfully")
    except Exception as e:
        logger.error(f"Failed to load ASR model: {e}")
        raise

    yield

    logger.info("Shutting down...")


app = FastAPI(title="Amharic STT & TTS", version="1.0.0", lifespan=lifespan)

# Create temp directory for files
TEMP_DIR = Path("temp_files")
TEMP_DIR.mkdir(exist_ok=True)


def convert_audio_to_wav(input_path, output_path=None):
    """Convert audio file to WAV format"""
    if output_path is None:
        base_name = os.path.splitext(input_path)[0]
        output_path = f"{base_name}.wav"

    # Load audio file and convert to WAV
    audio = AudioSegment.from_file(input_path)
    audio.export(output_path, format="wav")
    return output_path


async def speech_to_text(audio_file_path: str) -> str:
    if asr_model is None:
        raise HTTPException(status_code=500, detail="ASR model not loaded")

    # Check if file is already WAV format
    if not audio_file_path.lower().endswith(".wav"):
        logger.info(f"Converting {audio_file_path} to WAV format...")
        # Generate WAV file path
        wav_file_path = convert_audio_to_wav(audio_file_path)
        logger.info(f"Audio converted to: {wav_file_path}")
    else:
        wav_file_path = audio_file_path

    transcription = asr_model.transcribe_file(wav_file_path)
    logger.info(f"STT transcription: {transcription}")

    # Clean up converted file if it was created
    if wav_file_path != audio_file_path and os.path.exists(wav_file_path):
        os.remove(wav_file_path)
        logger.info(f"Cleaned up converted file: {wav_file_path}")

    return transcription


async def text_to_speech(
    text: str, output_path: str, voice: str = "am-ET-MekdesNeural"
):
    # Available Amharic voices:
    # am-ET-MekdesNeural (female)
    # am-ET-AmehaNeural (male)
    communicate = edge_tts.Communicate(text, voice)
    await communicate.save(output_path)
    logger.info(f"TTS audio saved to: {output_path}")


@app.get("/")
async def root():
    return {"message": "Amharic STT & TTS", "status": "running"}


@app.get("/health")
async def health_check():
    return {
        "status": "healthy",
        "models": {
            "asr": "loaded" if asr_model else "not_loaded",
            "tts": "edge-tts (online)",
        },
    }


@app.post("/speech-to-text")
async def stt_endpoint(
    audio_file: UploadFile = File(...),
):
    """Convert uploaded audio to text"""
    try:
        # Validate file type
        if not audio_file.content_type.startswith("audio/"):
            raise HTTPException(status_code=400, detail="File must be an audio file")

        # Check file size (limit to 10MB)
        if audio_file.size > 10 * 1024 * 1024:
            raise HTTPException(
                status_code=400, detail="File size too large. Maximum 10MB allowed."
            )

        file_id = str(uuid.uuid4())
        input_audio_path = TEMP_DIR / f"input_{file_id}.wav"

        with open(input_audio_path, "wb") as buffer:
            content = await audio_file.read()
            buffer.write(content)

        # Speech to Text
        user_text = await speech_to_text(str(input_audio_path))

        return {"transcription": user_text}

    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"STT error: {e}")
        raise HTTPException(status_code=500, detail=f"Speech-to-text failed: {str(e)}")


class TTSRequest(BaseModel):
    text: str
    voice: str = "female"

    @field_validator("voice")
    @classmethod
    def validate_voice(cls, v):
        if v.lower() not in VOICE_MAPPING:
            raise ValueError("Voice must be either 'male' or 'female'")
        return v.lower()


@app.post("/text-to-speech")
async def tts_endpoint(
    request: TTSRequest,
):
    """Convert text to speech audio"""
    try:
        if not request.text.strip():
            raise HTTPException(status_code=400, detail="Text cannot be empty")

        file_id = str(uuid.uuid4())
        output_audio_path = TEMP_DIR / f"output_{file_id}.mp3"

        # Map the voice selection to actual voice ID
        actual_voice = VOICE_MAPPING[request.voice]

        await text_to_speech(request.text, str(output_audio_path), actual_voice)

        return FileResponse(
            path=str(output_audio_path),
            media_type="audio/mpeg",
            filename=f"tts_output_{file_id}.mp3",
        )

    except HTTPException:
        raise
    except ValueError as e:
        raise HTTPException(status_code=400, detail=str(e))
    except Exception as e:
        logger.error(f"TTS error: {e}")
        raise HTTPException(status_code=500, detail=f"Text-to-speech failed: {str(e)}")


if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)
