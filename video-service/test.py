import torch
from diffusers import MochiPipeline, MochiTransformer3DModel
from diffusers.utils import export_to_video

model_id = "genmo/mochi-1-preview"
transformer = MochiTransformer3DModel.from_pretrained(
    model_id,
    subfolder="transformer",
    device_map="auto",
    max_memory={0: "24GB", 1: "24GB"}
)

pipe = MochiPipeline.from_pretrained(model_id,  transformer=transformer)
pipe.enable_model_cpu_offload()
pipe.enable_vae_tiling()

with torch.autocast(device_type="cuda", dtype=torch.bfloat16, cache_enabled=False):
    frames = pipe(
        prompt="Close-up of a chameleon's eye, with its scaly skin changing color. Ultra high resolution 4k.",
        num_frames=85,
    ).frames[0]

export_to_video(frames, "output.mp4", fps=30)
