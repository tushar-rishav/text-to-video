<script>
  import { createEventDispatcher } from 'svelte';

  const dispatch = createEventDispatcher();

  let prompt = '';
  let isSubmitting = false;
  let error = '';
  let success = '';

  async function submitJob() {
    if (!prompt.trim()) {
      error = 'Please enter a prompt';
      return;
    }

    isSubmitting = true;
    error = '';
    success = '';

    try {
      const response = await fetch('/api/submit', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ prompt: prompt.trim() }),
      });

      const data = await response.json();

      if (response.ok) {
        success = `Job submitted successfully! Job ID: ${data.job_id}`;
        prompt = '';
        dispatch('jobSubmitted', data.job_id);
      } else {
        error = data.error || 'Failed to submit job';
      }
    } catch (err) {
      error = 'Network error. Please try again.';
    } finally {
      isSubmitting = false;
    }
  }
</script>

<div class="job-submission">
  <h2>Submit Video Generation Job</h2>
  
  <form on:submit|preventDefault={submitJob} class="form">
    <div class="form-group">
      <label for="prompt">Text Prompt</label>
      <textarea
        id="prompt"
        bind:value={prompt}
        placeholder="Describe the video you want to generate... (e.g., 'A cat playing in a sunny garden')"
        rows="4"
        disabled={isSubmitting}
        required
      ></textarea>
    </div>

    {#if error}
      <div class="error">{error}</div>
    {/if}

    {#if success}
      <div class="success">{success}</div>
    {/if}

    <button type="submit" disabled={isSubmitting} class="submit-btn">
      {#if isSubmitting}
        <span class="spinner"></span>
        Submitting...
      {:else}
        Generate Video
      {/if}
    </button>
  </form>

  <div class="examples">
    <h3>Example Prompts</h3>
    <ul>
      <li>"A majestic eagle soaring through the clouds"</li>
      <li>"A peaceful lake with rippling water at sunset"</li>
      <li>"A robot dancing in a futuristic city"</li>
      <li>"A flower blooming in time-lapse"</li>
      <li>"A car driving through a mountain road"</li>
    </ul>
  </div>
</div>

<style>
  .job-submission {
    max-width: 600px;
    margin: 0 auto;
  }

  h2 {
    color: #2c3e50;
    margin-bottom: 30px;
    text-align: center;
    font-weight: 400;
  }

  .form {
    margin-bottom: 40px;
  }

  .form-group {
    margin-bottom: 20px;
  }

  label {
    display: block;
    margin-bottom: 8px;
    font-weight: 500;
    color: #34495e;
  }

  textarea {
    width: 100%;
    padding: 12px;
    border: 2px solid #e0e0e0;
    border-radius: 8px;
    font-size: 16px;
    font-family: inherit;
    resize: vertical;
    transition: border-color 0.3s ease;
  }

  textarea:focus {
    outline: none;
    border-color: #3498db;
  }

  textarea:disabled {
    background-color: #f8f9fa;
    cursor: not-allowed;
  }

  .submit-btn {
    width: 100%;
    padding: 15px;
    background: linear-gradient(135deg, #3498db, #2980b9);
    color: white;
    border: none;
    border-radius: 8px;
    font-size: 18px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.3s ease;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 10px;
  }

  .submit-btn:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(52, 152, 219, 0.3);
  }

  .submit-btn:disabled {
    opacity: 0.7;
    cursor: not-allowed;
    transform: none;
  }

  .spinner {
    width: 20px;
    height: 20px;
    border: 2px solid transparent;
    border-top: 2px solid white;
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }

  .error {
    background: #fee;
    color: #c33;
    padding: 12px;
    border-radius: 6px;
    margin-bottom: 20px;
    border-left: 4px solid #c33;
  }

  .success {
    background: #efe;
    color: #363;
    padding: 12px;
    border-radius: 6px;
    margin-bottom: 20px;
    border-left: 4px solid #363;
  }

  .examples {
    background: #f8f9fa;
    padding: 20px;
    border-radius: 8px;
    border-left: 4px solid #3498db;
  }

  .examples h3 {
    margin: 0 0 15px 0;
    color: #2c3e50;
    font-size: 18px;
  }

  .examples ul {
    margin: 0;
    padding-left: 20px;
  }

  .examples li {
    margin-bottom: 8px;
    color: #555;
    font-style: italic;
  }
</style> 