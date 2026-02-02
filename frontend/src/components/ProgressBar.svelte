<script>
  import { createEventDispatcher } from 'svelte';

  export let progress = 0;
  export let progressData = null;

  const dispatch = createEventDispatcher();

  $: progressPercent = Math.min(Math.max(progress, 0), 100);
  $: current = progressData?.current || progressData?.Current || 0;
  $: total = progressData?.total || progressData?.Total || 0;
  $: currentFile = progressData?.currentFile || progressData?.CurrentFile || '';
  $: status = progressData?.status || progressData?.Status || 'running';
  $: successCount = progressData?.successCount || progressData?.SuccessCount || 0;
  $: failureCount = progressData?.failureCount || progressData?.FailureCount || 0;
  $: destinationPath = progressData?.destinationPath || progressData?.DestinationPath || '';
  $: isCompleted = status === 'completed';

  function handleCancel() {
    if (confirm('Are you sure you want to cancel the import?')) {
      dispatch('cancel');
    }
  }

  function handleClose() {
    dispatch('complete');
  }

  async function openDestination() {
    try {
      await window.go.importer.ImportService.OpenDestinationFolder(destinationPath);
    } catch (error) {
      console.error('Error opening destination folder:', error);
      alert('Failed to open destination folder: ' + error);
    }
  }
</script>

<div class="progress-overlay">
  <div class="progress-modal">
    {#if isCompleted}
      <!-- Completion Summary -->
      <div class="completion-header">
        <div class="completion-icon">✓</div>
        <h3>Import Complete!</h3>
      </div>

      <div class="completion-summary">
        <div class="summary-item success">
          <div class="summary-icon">✓</div>
          <div class="summary-text">
            <div class="summary-count">{successCount}</div>
            <div class="summary-label">{successCount === 1 ? 'file' : 'files'} imported successfully</div>
          </div>
        </div>

        {#if failureCount > 0}
          <div class="summary-item failure">
            <div class="summary-icon">✗</div>
            <div class="summary-text">
              <div class="summary-count">{failureCount}</div>
              <div class="summary-label">{failureCount === 1 ? 'file' : 'files'} failed to import</div>
            </div>
          </div>
        {/if}
      </div>

      <div class="completion-actions">
        <button class="btn-open-folder" on:click={openDestination}>
          📁 Open Destination Folder
        </button>
        <button class="btn-close" on:click={handleClose}>
          Close
        </button>
      </div>
    {:else}
      <!-- Progress Display -->
      <div class="progress-header">
        <h3>Importing Files</h3>
        <div class="progress-percent">{progressPercent.toFixed(0)}%</div>
      </div>

      {#if total > 0}
        <div class="file-count">
          {current} of {total} files
        </div>
      {/if}

      <div class="progress-bar-container">
        <div class="progress-bar" style="width: {progressPercent}%">
          <div class="progress-shine"></div>
        </div>
      </div>

      {#if currentFile}
        <div class="current-file">
          <div class="file-label">Current file:</div>
          <div class="file-name">{currentFile}</div>
        </div>
      {/if}

      <div class="progress-actions">
        <div class="progress-info">
          <div class="spinner">⏳</div>
          <p>Please wait while your files are being imported...</p>
        </div>
        <button class="btn-cancel" on:click={handleCancel}>
          Cancel
        </button>
      </div>
    {/if}
  </div>
</div>

<style>
  .progress-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.8);
    backdrop-filter: blur(8px);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
    animation: fadeIn 0.3s ease;
  }

  @keyframes fadeIn {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }

  .progress-modal {
    background: linear-gradient(135deg, rgba(30, 30, 30, 0.95) 0%, rgba(20, 20, 20, 0.95) 100%);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 16px;
    padding: 32px;
    width: 90%;
    max-width: 500px;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
  }

  .progress-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;
  }

  .progress-header h3 {
    margin: 0;
    font-size: 20px;
    font-weight: 600;
    color: white;
  }

  .progress-percent {
    font-size: 24px;
    font-weight: 700;
    color: #4F46E5;
    font-variant-numeric: tabular-nums;
  }

  .progress-bar-container {
    height: 12px;
    background: rgba(0, 0, 0, 0.3);
    border-radius: 6px;
    overflow: hidden;
    margin-bottom: 20px;
    position: relative;
  }

  .progress-bar {
    height: 100%;
    background: linear-gradient(90deg, #4F46E5 0%, #7C3AED 100%);
    border-radius: 6px;
    transition: width 0.3s ease;
    position: relative;
    overflow: hidden;
  }

  .progress-shine {
    position: absolute;
    top: 0;
    left: -100%;
    width: 100%;
    height: 100%;
    background: linear-gradient(
      90deg,
      transparent,
      rgba(255, 255, 255, 0.3),
      transparent
    );
    animation: shine 2s infinite;
  }

  @keyframes shine {
    to {
      left: 200%;
    }
  }

  .file-count {
    text-align: center;
    font-size: 16px;
    font-weight: 500;
    color: #999;
    margin-bottom: 16px;
    font-variant-numeric: tabular-nums;
  }

  .current-file {
    margin-top: 16px;
    padding: 12px;
    background: rgba(0, 0, 0, 0.3);
    border-radius: 8px;
    margin-bottom: 20px;
  }

  .file-label {
    font-size: 12px;
    color: #666;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    margin-bottom: 6px;
  }

  .file-name {
    font-size: 14px;
    color: #ccc;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    font-family: monospace;
  }

  .progress-actions {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 16px;
  }

  .progress-info {
    display: flex;
    align-items: center;
    gap: 12px;
    color: #999;
    font-size: 14px;
    flex: 1;
  }

  .spinner {
    font-size: 24px;
    animation: spin 2s linear infinite;
  }

  @keyframes spin {
    from {
      transform: rotate(0deg);
    }
    to {
      transform: rotate(360deg);
    }
  }

  .progress-info p {
    margin: 0;
  }

  .btn-cancel {
    background: rgba(220, 38, 38, 0.2);
    border: 1px solid rgba(220, 38, 38, 0.4);
    border-radius: 8px;
    padding: 10px 20px;
    color: #ef4444;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
    white-space: nowrap;
  }

  .btn-cancel:hover {
    background: rgba(220, 38, 38, 0.3);
    border-color: rgba(220, 38, 38, 0.6);
    transform: translateY(-1px);
  }

  /* Completion Summary Styles */
  .completion-header {
    text-align: center;
    margin-bottom: 32px;
  }

  .completion-icon {
    width: 80px;
    height: 80px;
    margin: 0 auto 16px;
    background: linear-gradient(135deg, #10B981 0%, #059669 100%);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 48px;
    color: white;
    animation: scaleIn 0.5s ease;
  }

  @keyframes scaleIn {
    from {
      transform: scale(0);
    }
    to {
      transform: scale(1);
    }
  }

  .completion-header h3 {
    margin: 0;
    font-size: 24px;
    font-weight: 600;
    color: white;
  }

  .completion-summary {
    display: flex;
    flex-direction: column;
    gap: 16px;
    margin-bottom: 32px;
  }

  .summary-item {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 16px;
    border-radius: 12px;
    border: 2px solid;
  }

  .summary-item.success {
    background: rgba(16, 185, 129, 0.1);
    border-color: rgba(16, 185, 129, 0.3);
  }

  .summary-item.failure {
    background: rgba(239, 68, 68, 0.1);
    border-color: rgba(239, 68, 68, 0.3);
  }

  .summary-icon {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 20px;
    font-weight: bold;
  }

  .summary-item.success .summary-icon {
    background: rgba(16, 185, 129, 0.2);
    color: #10B981;
  }

  .summary-item.failure .summary-icon {
    background: rgba(239, 68, 68, 0.2);
    color: #ef4444;
  }

  .summary-text {
    flex: 1;
  }

  .summary-count {
    font-size: 28px;
    font-weight: 700;
    line-height: 1;
    margin-bottom: 4px;
  }

  .summary-item.success .summary-count {
    color: #10B981;
  }

  .summary-item.failure .summary-count {
    color: #ef4444;
  }

  .summary-label {
    font-size: 14px;
    color: #999;
  }

  .completion-actions {
    display: flex;
    gap: 12px;
  }

  .btn-open-folder {
    flex: 1;
    background: linear-gradient(135deg, #4F46E5 0%, #7C3AED 100%);
    border: none;
    border-radius: 8px;
    padding: 14px 20px;
    color: white;
    font-size: 15px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
  }

  .btn-open-folder:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 20px rgba(79, 70, 229, 0.4);
  }

  .btn-close {
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 8px;
    padding: 14px 24px;
    color: white;
    font-size: 15px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
  }

  .btn-close:hover {
    background: rgba(255, 255, 255, 0.1);
    transform: translateY(-1px);
  }
</style>
