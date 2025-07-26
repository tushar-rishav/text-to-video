<script>
  import { onMount } from 'svelte';
  import JobSubmission from './components/JobSubmission.svelte';
  import JobList from './components/JobList.svelte';
  import JobStatus from './components/JobStatus.svelte';

  let currentView = 'submit';
  let jobs = [];
  let ws = null;

  onMount(() => {
    // Initialize WebSocket connection for real-time updates
    initWebSocket();
    loadJobs();
  });

  function initWebSocket() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/api/ws`;
    
    ws = new WebSocket(wsUrl);
    
    ws.onopen = () => {
      console.log('WebSocket connected');
    };
    
    ws.onmessage = (event) => {
      const jobUpdate = JSON.parse(event.data);
      updateJobStatus(jobUpdate);
    };
    
    ws.onerror = (error) => {
      console.error('WebSocket error:', error);
    };
    
    ws.onclose = () => {
      console.log('WebSocket disconnected');
      // Reconnect after 5 seconds
      setTimeout(initWebSocket, 5000);
    };
  }

  async function loadJobs() {
    try {
      const response = await fetch('/api/list?limit=50');
      const data = await response.json();
      if (data.jobs) {
        jobs = data.jobs;
      }
    } catch (error) {
      console.error('Failed to load jobs:', error);
    }
  }

  function updateJobStatus(jobUpdate) {
    const index = jobs.findIndex(job => job.job_id === jobUpdate.job_id);
    if (index !== -1) {
      jobs[index] = { ...jobs[index], ...jobUpdate };
      jobs = [...jobs]; // Trigger reactivity
    } else {
      // Add new job if not in list
      jobs = [jobUpdate, ...jobs];
    }
  }

  function handleJobSubmitted(jobId) {
    // Switch to status view to monitor the new job
    currentView = 'status';
    // The job will be added to the list via WebSocket
  }

  function navigate(view) {
    currentView = view;
  }
</script>

<main class="container">
  <header class="header">
    <h1>Text-to-Video Generator</h1>
    <nav class="nav">
      <button 
        class="nav-btn" 
        class:active={currentView === 'submit'}
        on:click={() => navigate('submit')}
      >
        Submit Job
      </button>
      <button 
        class="nav-btn" 
        class:active={currentView === 'status'}
        on:click={() => navigate('status')}
      >
        Job Status
      </button>
      <button 
        class="nav-btn" 
        class:active={currentView === 'list'}
        on:click={() => navigate('list')}
      >
        All Jobs
      </button>
    </nav>
  </header>

  <div class="content">
    {#if currentView === 'submit'}
      <JobSubmission on:jobSubmitted={handleJobSubmitted} />
    {:else if currentView === 'status'}
      <JobStatus {jobs} />
    {:else if currentView === 'list'}
      <JobList {jobs} on:refresh={loadJobs} />
    {/if}
  </div>
</main>

<style>
  .container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 20px;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  }

  .header {
    text-align: center;
    margin-bottom: 40px;
  }

  .header h1 {
    color: #2c3e50;
    margin-bottom: 20px;
    font-size: 2.5rem;
    font-weight: 300;
  }

  .nav {
    display: flex;
    justify-content: center;
    gap: 10px;
    margin-bottom: 30px;
  }

  .nav-btn {
    padding: 12px 24px;
    border: 2px solid #3498db;
    background: transparent;
    color: #3498db;
    border-radius: 8px;
    cursor: pointer;
    font-size: 16px;
    font-weight: 500;
    transition: all 0.3s ease;
  }

  .nav-btn:hover {
    background: #3498db;
    color: white;
  }

  .nav-btn.active {
    background: #3498db;
    color: white;
  }

  .content {
    background: white;
    border-radius: 12px;
    padding: 30px;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    min-height: 500px;
  }

  @media (max-width: 768px) {
    .container {
      padding: 10px;
    }
    
    .header h1 {
      font-size: 2rem;
    }
    
    .nav {
      flex-direction: column;
      align-items: center;
    }
    
    .nav-btn {
      width: 200px;
    }
  }
</style> 