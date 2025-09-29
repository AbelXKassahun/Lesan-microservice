# official Python image (pick a stable version like 3.11)
FROM python:3.12-slim

# install system dependencies (if some packages need gcc, ffmpeg, etc.)
RUN apt-get update && apt-get install -y curl \
build-essential \
ffmpeg \
git \
&& curl -LsSf https://astral.sh/uv/install.sh | sh \
&& rm -rf /var/lib/apt/lists/*

# Add uv to PATH
ENV PATH="/root/.cargo/bin:${PATH}"

# Set working directory
WORKDIR /app

# Copy requirements 
COPY requirements.txt .

# install Python dependencies
RUN uv pip install -r requirements.txt
# RUN pip install --no-cache-dir --timeout=1200 -r requirements.txt

# copy the rest of the app
COPY . .

# Default command
CMD ["uv", "run", "python", "main.py"]