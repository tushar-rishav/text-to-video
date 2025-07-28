import asyncio
import json
import logging
import os
import time
from concurrent.futures import ThreadPoolExecutor
from typing import Dict, Any
import uuid

import redis
import pymysql
from dotenv import load_dotenv
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel

import torch
from diffusers import MochiPipeline
from diffusers.utils import export_to_video


# Load environment variables
load_dotenv()

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# Global variables
redis_client = None
db_connection = None
mochi_pipe = None
executor = ThreadPoolExecutor(max_workers=2)  # Limit concurrent video generation

app = FastAPI(title="Video Generation Service", version="1.0.0")

class JobUpdate(BaseModel):
    job_id: str
    status: str
    video_url: str = None
    error: str = None

def init_redis():
    """Initialize Redis connection"""
    global redis_client
    redis_client = redis.Redis(
        host=os.getenv("REDIS_HOST", "localhost"),
        port=int(os.getenv("REDIS_PORT", 6379)),
        password=os.getenv("REDIS_PASSWORD"),
        decode_responses=True
    )
    logger.info("Redis connection established")

def init_database():
    """Initialize MySQL database connection"""
    global db_connection
    db_connection = pymysql.connect(
        host=os.getenv("DB_HOST", "localhost"),
        user=os.getenv("DB_USER", "root"),
        password=os.getenv("DB_PASSWORD", ""),
        database=os.getenv("DB_NAME", "text_to_video"),
        charset='utf8mb4',
        cursorclass=pymysql.cursors.DictCursor
    )
    logger.info("Database connection established")

def load_model():
    """Load the genmo mochi-1 model"""
    global mochi_pipe
    
    try:
        logger.info("Loading genmo mochi-1 model...")
        
        # Load the model from Hugging Face
        model_name = "genmo/mochi-1-preview"
        
        mochi_pipe = MochiPipeline.from_pretrained("genmo/mochi-1-preview")

        # Enable memory savings
        mochi_pipe.enable_model_cpu_offload()
        mochi_pipe.enable_vae_tiling()

        logger.info("Model loaded successfully")        
    except Exception as e:
        logger.error(f"Failed to load model: {e}")
        raise

def update_job_status(job_id: str, status: str, video_url: str = None, error: str = None):
    """Update job status in database and publish to Redis"""
    try:
        # Update database
        with db_connection.cursor() as cursor:
            if video_url:
                sql = "UPDATE jobs SET status = %s, video_url = %s, updated_at = NOW() WHERE id = %s"
                cursor.execute(sql, (status, video_url, job_id))
            elif error:
                sql = "UPDATE jobs SET status = %s, error = %s, updated_at = NOW() WHERE id = %s"
                cursor.execute(sql, (status, error, job_id))
            else:
                sql = "UPDATE jobs SET status = %s, updated_at = NOW() WHERE id = %s"
                cursor.execute(sql, (status, job_id))
        
        db_connection.commit()
        
        # Publish to Redis for real-time updates
        job_update = JobUpdate(
            job_id=job_id,
            status=status,
            video_url=video_url,
            error=error
        )
        
        redis_client.publish("job_updates", job_update.model_dump_json())
        logger.info(f"Job {job_id} status updated to {status}")
        
    except Exception as e:
        logger.error(f"Failed to update job status: {e}")

def generate_video(prompt: str, job_id: str) -> str:
    """Generate video using genmo mochi-1 model"""
    video_url = f"/videos/{job_id}.mp4"
    try:
        logger.info(f"Starting video generation for job {job_id}")
        update_job_status(job_id, "processing")
        
        # Generate video
        video_path = f"/app/videos/{job_id}.mp4"
        with torch.autocast("cuda", torch.bfloat16, cache_enabled=False):
            frames = mochi_pipe(prompt, num_frames=85).frames[0]

        export_to_video(frames, video_path, fps=30)
        
        logger.info(f"Video generated successfully for job {job_id}")
        update_job_status(job_id, "completed", video_url=video_url)
        
        return video_url            
    except Exception as e:
        error_msg = f"Video generation failed: {str(e)}"
        logger.error(f"Error generating video for job {job_id}: {e}")
        update_job_status(job_id, "failed", error=error_msg)
        raise

def process_job(job_data: Dict[str, Any]):
    """Process a single job"""
    try:
        job_id = job_data.get("job_id")
        prompt = job_data.get("prompt")
        
        if not job_id or not prompt:
            logger.error("Invalid job data received")
            return
        
        logger.info(f"Processing job {job_id} with prompt: {prompt}")
        
        # Submit video generation task to thread pool
        future = executor.submit(generate_video, prompt, job_id)
        future.result()  # Wait for completion
        
    except Exception as e:
        logger.error(f"Error processing job: {e}")

async def listen_for_jobs():
    """Listen for new jobs from Redis"""
    pubsub = redis_client.pubsub()
    pubsub.subscribe("video_jobs")
    
    logger.info("Listening for video generation jobs...")
    
    for message in pubsub.listen():
        if message["type"] == "message":
            try:
                job_data = json.loads(message["data"])
                logger.info(f"Received job: {job_data.get('job_id')}")
                
                # Process job in thread pool to avoid blocking
                loop = asyncio.get_event_loop()
                await loop.run_in_executor(executor, process_job, job_data)
                
            except json.JSONDecodeError as e:
                logger.error(f"Failed to decode job data: {e}")
            except Exception as e:
                logger.error(f"Error processing job: {e}")

@app.on_event("startup")
async def startup_event():
    """Initialize services on startup"""
    logger.info("Starting video generation service...")
    
    # Initialize connections
    init_redis()
    init_database()
    
    # Load model
    load_model()
    
    # Start job listener
    asyncio.create_task(listen_for_jobs())
    
    logger.info("Video generation service started successfully")

@app.on_event("shutdown")
async def shutdown_event():
    """Cleanup on shutdown"""
    logger.info("Shutting down video generation service...")
    
    if db_connection:
        db_connection.close()
    
    if redis_client:
        redis_client.close()
    
    executor.shutdown(wait=True)
    
    logger.info("Video generation service shutdown complete")

@app.get("/health")
async def health_check():
    """Health check endpoint"""
    return {
        "status": "healthy",
        "model_loaded": model is not None,
        "gpu_available": torch.cuda.is_available(),
        "redis_connected": redis_client is not None,
        "db_connected": db_connection is not None
    }

@app.post("/generate")
async def generate_video_endpoint(prompt: str):
    """Direct video generation endpoint (for testing)"""
    job_id = str(uuid.uuid4())
    
    try:
        # Submit to thread pool
        loop = asyncio.get_event_loop()
        video_url = await loop.run_in_executor(executor, generate_video, prompt, job_id)
        
        return {
            "job_id": job_id,
            "video_url": video_url,
            "status": "completed"
        }
        
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

if __name__ == "__main__":
    import uvicorn
    
    port = int(os.getenv("PORT", 8000))
    uvicorn.run(app, host="0.0.0.0", port=port) 