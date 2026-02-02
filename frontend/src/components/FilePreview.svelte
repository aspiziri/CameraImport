<script>
  import { createEventDispatcher } from 'svelte';

  export let files = [];
  export let selectedDates = [];

  const dispatch = createEventDispatcher();

  $: groupedByDate = groupFilesByDate(files);

  function groupFilesByDate(fileList) {
    const groups = {};
    fileList.forEach(file => {
      if (!groups[file.date]) {
        groups[file.date] = [];
      }
      groups[file.date].push(file);
    });
    return groups;
  }

  function toggleDate(date) {
    if (selectedDates.includes(date)) {
      selectedDates = selectedDates.filter(d => d !== date);
    } else {
      selectedDates = [...selectedDates, date];
    }
  }

  function selectAllDates() {
    selectedDates = Object.keys(groupedByDate);
  }

  function deselectAllDates() {
    selectedDates = [];
  }

  function formatFileSize(bytes) {
    if (bytes < 1024) return bytes + ' B';
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB';
    if (bytes < 1024 * 1024 * 1024) return (bytes / (1024 * 1024)).toFixed(1) + ' MB';
    return (bytes / (1024 * 1024 * 1024)).toFixed(2) + ' GB';
  }
</script>

<div class="file-preview">
  <div class="preview-header">
    <span class="step-badge">Step 2</span>
    <h2>Select Dates to Import</h2>
  </div>

  <div class="preview-content">
    {#if files.length === 0}
      <div class="empty-state">
        <div class="empty-icon">📂</div>
        <p>No files scanned yet</p>
        <small>Click "Scan Files" to get started</small>
      </div>
    {:else}
      <div class="date-actions">
        <button class="btn-small" on:click={selectAllDates}>Select All</button>
        <button class="btn-small" on:click={deselectAllDates}>Deselect All</button>
        <span class="file-count">{files.length} files found</span>
      </div>

      <div class="date-list">
        {#each Object.keys(groupedByDate).sort() as date}
          <div
            class="date-item"
            class:selected={selectedDates.includes(date)}
            on:click={() => toggleDate(date)}
          >
            <div class="date-checkbox">
              <input type="checkbox" checked={selectedDates.includes(date)} />
            </div>
            <div class="date-label">{date}</div>
            <div class="date-stats">
              {groupedByDate[date].length} files · {formatFileSize(groupedByDate[date].reduce((sum, f) => sum + f.size, 0))}
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </div>
</div>

<style>
  .file-preview {
    background: rgba(255, 255, 255, 0.05);
    border-radius: 12px;
    border: 1px solid rgba(255, 255, 255, 0.1);
    overflow: hidden;
    display: flex;
    flex-direction: column;
  }

  .preview-header {
    background: rgba(0, 0, 0, 0.2);
    padding: 16px 20px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .step-badge {
    background: linear-gradient(135deg, #4F46E5 0%, #7C3AED 100%);
    color: white;
    font-size: 12px;
    font-weight: 600;
    padding: 6px 12px;
    border-radius: 6px;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    white-space: nowrap;
  }

  .preview-header h2 {
    margin: 0;
    font-size: 18px;
    font-weight: 600;
    color: #ffffff;
  }

  .preview-content {
    flex: 1;
    overflow-y: auto;
    padding: 20px;
  }

  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100%;
    color: #666;
    gap: 12px;
  }

  .empty-icon {
    font-size: 64px;
    opacity: 0.3;
  }

  .date-actions {
    display: flex;
    gap: 12px;
    margin-bottom: 16px;
    align-items: center;
  }

  .btn-small {
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 6px;
    padding: 6px 12px;
    color: white;
    font-size: 13px;
    cursor: pointer;
    transition: all 0.2s;
  }

  .btn-small:hover {
    background: rgba(255, 255, 255, 0.1);
  }

  .file-count {
    margin-left: auto;
    font-size: 13px;
    color: #999;
  }

  .date-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .date-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 10px 14px;
    background: rgba(0, 0, 0, 0.2);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.2s;
  }

  .date-item:hover {
    background: rgba(255, 255, 255, 0.05);
  }

  .date-item.selected {
    border-color: #4F46E5;
    background: rgba(79, 70, 229, 0.1);
  }

  .date-checkbox input {
    cursor: pointer;
  }

  .date-label {
    font-size: 15px;
    font-weight: 600;
    color: white;
    flex: 1;
  }

  .date-stats {
    font-size: 13px;
    color: #999;
    white-space: nowrap;
    margin-right: 8px;
  }
</style>
