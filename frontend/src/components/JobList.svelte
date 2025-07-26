<script>
  import { createEventDispatcher } from 'svelte';

  const dispatch = createEventDispatcher();

  export let jobs = [];

  let selectedStatus = '';
  let searchTerm = '';
  let currentPage = 1;
  let itemsPerPage = 10;

  $: filteredJobs = jobs.filter(job => {
    const matchesStatus = !selectedStatus || job.job_status === selectedStatus;
    const matchesSearch = !searchTerm || 
      job.prompt.toLowerCase().includes(searchTerm.toLowerCase()) ||
      job.job_id.toLowerCase().includes(searchTerm.toLowerCase());
    return matchesStatus && matchesSearch;
  });

  $: totalPages = Math.ceil(filteredJobs.length / itemsPerPage);
  $: paginatedJobs = filteredJobs.slice(
    (currentPage - 1) * itemsPerPage,
    currentPage * itemsPerPage
  );

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

  function refreshJobs() {
    dispatch('refresh');
  }

  function goToPage(page) {
    if (page >= 1 && page <= totalPages) {
      currentPage = page;
    }
  }

  function clearFilters() {
    selectedStatus = '';
    searchTerm = '';
    currentPage = 1;
  }
</script>

<div class="job-list">
  <div class="header">
    <h2>All Jobs</h2>
    <button class="refresh-btn" on:click={refreshJobs}>
      🔄 Refresh
    </button>
  </div>

  <div class="filters">
    <div class="filter-group">
      <label for="status-filter">Status:</label>
      <select id="status-filter" bind:value={selectedStatus}>
        <option value="">All Statuses</option>
        <option value="pending">Pending</option>
        <option value="processing">Processing</option>
        <option value="completed">Completed</option>
        <option value="failed">Failed</option>
      </select>
    </div>

    <div class="filter-group">
      <label for="search">Search:</label>
      <input 
        id="search"
        type="text" 
        bind:value={searchTerm}
        placeholder="Search by prompt or job ID..."
      />
    </div>

    <button class="clear-btn" on:click={clearFilters}>
      Clear Filters
    </button>
  </div>

  <div class="stats">
    <span>Total Jobs: {jobs.length}</span>
    <span>Filtered: {filteredJobs.length}</span>
    <span>Showing: {paginatedJobs.length}</span>
  </div>

  {#if paginatedJobs.length === 0}
    <div class="empty-state">
      <p>No jobs found matching your criteria.</p>
    </div>
  {:else}
    <div class="table-container">
      <table class="job-table">
        <thead>
          <tr>
            <th>Status</th>
            <th>Job ID</th>
            <th>Prompt</th>
            <th>Created</th>
            <th>Updated</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each paginatedJobs as job}
            <tr class="job-row">
              <td>
                <span class="status-badge" style="color: {getStatusColor(job.job_status)}">
                  {getStatusIcon(job.job_status)} {job.job_status}
                </span>
              </td>
              <td class="job-id">{job.job_id}</td>
              <td class="job-prompt">{job.prompt}</td>
              <td class="job-date">{formatDate(job.created_at)}</td>
              <td class="job-date">{formatDate(job.updated_at)}</td>
              <td class="actions">
                {#if job.video_url}
                  <button 
                    class="action-btn download"
                    on:click={() => downloadVideo(job.video_url)}
                    title="Download Video"
                  >
                    📥
                  </button>
                {/if}
                {#if job.error}
                  <span class="error-indicator" title="Error: {job.error}">⚠️</span>
                {/if}
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>

    {#if totalPages > 1}
      <div class="pagination">
        <button 
          class="page-btn"
          disabled={currentPage === 1}
          on:click={() => goToPage(currentPage - 1)}
        >
          ← Previous
        </button>
        
        <span class="page-info">
          Page {currentPage} of {totalPages}
        </span>
        
        <button 
          class="page-btn"
          disabled={currentPage === totalPages}
          on:click={() => goToPage(currentPage + 1)}
        >
          Next →
        </button>
      </div>
    {/if}
  {/if}
</div>

<style>
  .job-list {
    max-width: 1200px;
    margin: 0 auto;
  }

  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 30px;
  }

  h2 {
    color: #2c3e50;
    margin: 0;
    font-weight: 400;
  }

  .refresh-btn {
    background: #3498db;
    color: white;
    border: none;
    padding: 10px 20px;
    border-radius: 6px;
    cursor: pointer;
    font-size: 14px;
    transition: all 0.3s ease;
  }

  .refresh-btn:hover {
    background: #2980b9;
  }

  .filters {
    display: flex;
    gap: 20px;
    margin-bottom: 20px;
    align-items: end;
    flex-wrap: wrap;
  }

  .filter-group {
    display: flex;
    flex-direction: column;
    gap: 5px;
  }

  .filter-group label {
    font-size: 14px;
    font-weight: 500;
    color: #34495e;
  }

  .filter-group select,
  .filter-group input {
    padding: 8px 12px;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 14px;
  }

  .clear-btn {
    background: #95a5a6;
    color: white;
    border: none;
    padding: 8px 16px;
    border-radius: 4px;
    cursor: pointer;
    font-size: 14px;
  }

  .clear-btn:hover {
    background: #7f8c8d;
  }

  .stats {
    display: flex;
    gap: 20px;
    margin-bottom: 20px;
    font-size: 14px;
    color: #7f8c8d;
  }

  .table-container {
    overflow-x: auto;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }

  .job-table {
    width: 100%;
    border-collapse: collapse;
    background: white;
  }

  .job-table th,
  .job-table td {
    padding: 12px;
    text-align: left;
    border-bottom: 1px solid #ecf0f1;
  }

  .job-table th {
    background: #f8f9fa;
    font-weight: 600;
    color: #2c3e50;
    position: sticky;
    top: 0;
  }

  .job-row:hover {
    background: #f8f9fa;
  }

  .status-badge {
    font-weight: 500;
    font-size: 12px;
    text-transform: uppercase;
  }

  .job-id {
    font-family: monospace;
    font-size: 12px;
    color: #7f8c8d;
    word-break: break-all;
  }

  .job-prompt {
    max-width: 300px;
    word-wrap: break-word;
  }

  .job-date {
    font-size: 12px;
    color: #95a5a6;
    white-space: nowrap;
  }

  .actions {
    display: flex;
    gap: 5px;
    align-items: center;
  }

  .action-btn {
    background: none;
    border: none;
    cursor: pointer;
    font-size: 16px;
    padding: 5px;
    border-radius: 4px;
    transition: background 0.2s ease;
  }

  .action-btn:hover {
    background: #ecf0f1;
  }

  .error-indicator {
    font-size: 16px;
    cursor: help;
  }

  .pagination {
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 20px;
    margin-top: 30px;
  }

  .page-btn {
    background: #3498db;
    color: white;
    border: none;
    padding: 10px 20px;
    border-radius: 6px;
    cursor: pointer;
    font-size: 14px;
    transition: all 0.3s ease;
  }

  .page-btn:hover:not(:disabled) {
    background: #2980b9;
  }

  .page-btn:disabled {
    background: #bdc3c7;
    cursor: not-allowed;
  }

  .page-info {
    font-size: 14px;
    color: #7f8c8d;
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
    .filters {
      flex-direction: column;
      align-items: stretch;
    }
    
    .stats {
      flex-direction: column;
      gap: 5px;
    }
    
    .job-table {
      font-size: 12px;
    }
    
    .job-table th,
    .job-table td {
      padding: 8px;
    }
  }
</style> 