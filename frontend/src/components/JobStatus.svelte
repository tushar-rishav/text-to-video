<script>
  export let jobs = [];

  $: recentJobs = jobs.filter(job => 
    ['pending', 'processing'].includes(job.job_status)
  ).slice(0, 10);

  $: completedJobs = jobs.filter(job => 
    job.job_status === 'completed'
  ).slice(0, 5);

  function getStatusColor(status) {
    switch (status) {
      case 'pending': return '#f39c12';
      case 'processing': return '#3498db';
      case 'completed': return '#27ae60';
      case 'failed': return '#e74c3c';
      default: return '#95a5a6';
    }
  }

  function getStatusIcon(status) {
    switch (status) {
      case 'pending': return '⏳';
      case 'processing': return '🔄';
      case 'completed': return '✅';
      case 'failed': return '❌';
      default: return '❓';
    }
  }

  function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleString();
  }

  function downloadVideo(videoUrl) {
    const link = document.createElement('a');
    link.href = videoUrl;
    link.download = `video-${Date.now()}.mp4`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  }
</script>

<div class="job-status">
  <h2>Job Status</h2>

  {#if recentJobs.length === 0 && completedJobs.length === 0}
    <div class="empty-state">
      <p>No jobs found. Submit a new job to get started!</p>
    </div>
  {:else}
    {#if recentJobs.length > 0}
      <section class="section">
        <h3>Active Jobs</h3>
        <div class="job-grid">
          {#each recentJobs as job}
            <div class="job-card active">
              <div class="job-header">
                <span class="status-icon" style="color: {getStatusColor(job.job_status)}">
                  {getStatusIcon(job.job_status)}
                </span>
                <span class="status-text" style="color: {getStatusColor(job.job_status)}">
                  {job.job_status}
                </span>
              </div>
              <div class="job-id">ID: {job.job_id}</div>
              <div class="job-prompt">{job.prompt}</div>
              <div class="job-time">Created: {formatDate(job.created_at)}</div>
              {#if job.error}
                <div class="job-error">Error: {job.error}</div>
              {/if}
            </div>
          {/each}
        </div>
      </section>
    {/if}

    {#if completedJobs.length > 0}
      <section class="section">
        <h3>Recently Completed</h3>
        <div class="job-grid">
          {#each completedJobs as job}
            <div class="job-card completed">
              <div class="job-header">
                <span class="status-icon" style="color: {getStatusColor(job.job_status)}">
                  {getStatusIcon(job.job_status)}
                </span>
                <span class="status-text" style="color: {getStatusColor(job.job_status)}">
                  {job.job_status}
                </span>
              </div>
              <div class="job-id">ID: {job.job_id}</div>
              <div class="job-prompt">{job.prompt}</div>
              <div class="job-time">Completed: {formatDate(job.updated_at)}</div>
              {#if job.video_url}
                <button 
                  class="download-btn"
                  on:click={() => downloadVideo(job.video_url)}
                >
                  📥 Download Video
                </button>
              {/if}
            </div>
          {/each}
        </div>
      </section>
    {/if}
  {/if}
</div>

<style>
  .job-status {
    max-width: 1000px;
    margin: 0 auto;
  }

  h2 {
    color: #2c3e50;
    margin-bottom: 30px;
    text-align: center;
    font-weight: 400;
  }

  .section {
    margin-bottom: 40px;
  }

  .section h3 {
    color: #34495e;
    margin-bottom: 20px;
    font-weight: 500;
    border-bottom: 2px solid #ecf0f1;
    padding-bottom: 10px;
  }

  .job-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 20px;
  }

  .job-card {
    background: white;
    border: 1px solid #e0e0e0;
    border-radius: 12px;
    padding: 20px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    transition: all 0.3s ease;
  }

  .job-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  }

  .job-card.active {
    border-left: 4px solid #3498db;
  }

  .job-card.completed {
    border-left: 4px solid #27ae60;
  }

  .job-header {
    display: flex;
    align-items: center;
    gap: 10px;
    margin-bottom: 15px;
  }

  .status-icon {
    font-size: 20px;
  }

  .status-text {
    font-weight: 600;
    text-transform: uppercase;
    font-size: 14px;
  }

  .job-id {
    font-family: monospace;
    font-size: 12px;
    color: #7f8c8d;
    margin-bottom: 10px;
    word-break: break-all;
  }

  .job-prompt {
    font-size: 16px;
    color: #2c3e50;
    margin-bottom: 15px;
    line-height: 1.4;
  }

  .job-time {
    font-size: 12px;
    color: #95a5a6;
    margin-bottom: 10px;
  }

  .job-error {
    background: #fee;
    color: #c33;
    padding: 8px;
    border-radius: 4px;
    font-size: 12px;
    margin-top: 10px;
  }

  .download-btn {
    background: linear-gradient(135deg, #27ae60, #2ecc71);
    color: white;
    border: none;
    padding: 10px 15px;
    border-radius: 6px;
    cursor: pointer;
    font-size: 14px;
    font-weight: 500;
    transition: all 0.3s ease;
    width: 100%;
    margin-top: 10px;
  }

  .download-btn:hover {
    transform: translateY(-1px);
    box-shadow: 0 2px 8px rgba(39, 174, 96, 0.3);
  }

  .empty-state {
    text-align: center;
    padding: 60px 20px;
    color: #7f8c8d;
  }

  .empty-state p {
    font-size: 18px;
    margin: 0;
  }

  @media (max-width: 768px) {
    .job-grid {
      grid-template-columns: 1fr;
    }
    
    .job-card {
      padding: 15px;
    }
  }
</style> 