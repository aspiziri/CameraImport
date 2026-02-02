<script>
  import { createEventDispatcher } from 'svelte';

  export let files = [];
  export let config;
  export let selectedDates = [];
  export let disabled = false;

  const dispatch = createEventDispatcher();

  let importMode = 'full';
  let customFolderName = '';
  let deleteAfterImport = false;
  let showDeleteWarning = false;

  function startImport() {
    if (selectedDates.length === 0) {
      alert('Please select at least one date to import in the preview section');
      return;
    }

    if (importMode === 'selective' && !customFolderName.trim()) {
      alert('Please enter a folder name');
      return;
    }

    if (!config.destination) {
      alert('Please set a destination folder');
      return;
    }

    // Show confirmation modal if delete is enabled
    if (deleteAfterImport) {
      showDeleteWarning = true;
      return;
    }

    // Proceed with import
    proceedWithImport();
  }

  function proceedWithImport() {
    showDeleteWarning = false;
    dispatch('import', {
      mode: importMode,
      dates: selectedDates,
      folderName: importMode === 'selective' ? customFolderName : null,
      deleteAfter: deleteAfterImport
    });
  }

  function cancelImport() {
    showDeleteWarning = false;
  }

  $: canImport = files.length > 0 && !disabled;
  $: importCount = files.filter(f => {
    const fileDate = f.Date || f.date;
    return selectedDates.includes(fileDate);
  }).length;
</script>

<div class="import-controls">
  <div class="controls-header">
    <div class="header-left">
      <span class="step-badge">Step 3</span>
      <h2>Pick Your Import Type</h2>
    </div>
  </div>

  <div class="controls-content">
    <div class="mode-selector">
      <label class="radio-card" class:selected={importMode === 'full'}>
        <input
          type="radio"
          bind:group={importMode}
          value="full"
          disabled={disabled}
        />
        <div class="radio-content">
          <div class="radio-icon">📅</div>
          <div class="radio-text">
            <div class="radio-title">Import to Date Folders</div>
            <div class="radio-description">
              Import selected dates to separate date-based folders
            </div>
          </div>
        </div>
      </label>

      <label class="radio-card" class:selected={importMode === 'selective'}>
        <input
          type="radio"
          bind:group={importMode}
          value="selective"
          disabled={disabled}
        />
        <div class="radio-content">
          <div class="radio-icon">📁</div>
          <div class="radio-text">
            <div class="radio-title">Import to Specific Folder</div>
            <div class="radio-description">
              Import selected dates to a single custom folder
            </div>
          </div>
        </div>
      </label>
    </div>

    {#if importMode === 'selective'}
      <div class="selective-options">
        <label>
          <span class="label-text">Folder Name</span>
          <input
            type="text"
            bind:value={customFolderName}
            placeholder="e.g., Wedding Shoot 2024"
            disabled={disabled}
          />
        </label>
        <small class="help-text">
          Selected files will be imported to: {config.destination}\{customFolderName || '...'}
        </small>
      </div>
    {/if}

    <label class="checkbox-option">
      <input
        type="checkbox"
        bind:checked={deleteAfterImport}
        disabled={disabled}
      />
      <span>Delete files from source after successful import</span>
    </label>

    <button
      class="btn-import"
      on:click={startImport}
      disabled={!canImport}
    >
      {disabled ? '⏳ Importing...' : '⚡ Start Import'}
    </button>
  </div>
</div>

{#if showDeleteWarning}
  <div class="warning-overlay">
    <div class="warning-modal">
      <div class="warning-header">
        <div class="warning-icon">⚠️</div>
        <h3>Delete Files After Import</h3>
      </div>

      <div class="warning-content">
        <p class="warning-message">
          You have chosen to <strong>DELETE</strong> source files after they are imported.
        </p>
        <div class="warning-details">
          <div class="detail-item">
            <span class="detail-label">Files to be deleted:</span>
            <span class="detail-value">{importCount} {importCount === 1 ? 'file' : 'files'}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">Source location:</span>
            <span class="detail-value">{config.source}</span>
          </div>
        </div>
        <p class="warning-footer">
          This action <strong>cannot be undone</strong>. Files will be permanently removed from the source location.
        </p>
      </div>

      <div class="warning-actions">
        <button class="btn-cancel-warning" on:click={cancelImport}>
          Cancel
        </button>
        <button class="btn-confirm-delete" on:click={proceedWithImport}>
          Yes, Delete After Import
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
  .import-controls {
    background: rgba(255, 255, 255, 0.05);
    border-radius: 12px;
    border: 1px solid rgba(255, 255, 255, 0.1);
    overflow: hidden;
    display: flex;
    flex-direction: column;
  }

  .controls-header {
    background: rgba(0, 0, 0, 0.2);
    padding: 16px 20px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .header-left {
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

  .controls-header h2 {
    margin: 0;
    font-size: 18px;
    font-weight: 600;
    color: #ffffff;
  }

  .controls-content {
    padding: 20px;
    display: flex;
    flex-direction: column;
    gap: 20px;
  }

  .mode-selector {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .radio-card {
    position: relative;
    display: block;
    cursor: pointer;
  }

  .radio-card input[type="radio"] {
    position: absolute;
    opacity: 0;
    pointer-events: none;
  }

  .radio-content {
    display: flex;
    gap: 12px;
    padding: 16px;
    background: rgba(0, 0, 0, 0.2);
    border: 2px solid rgba(255, 255, 255, 0.1);
    border-radius: 8px;
    transition: all 0.2s;
  }

  .radio-card:hover .radio-content {
    background: rgba(255, 255, 255, 0.05);
  }

  .radio-card.selected .radio-content {
    border-color: #4F46E5;
    background: rgba(79, 70, 229, 0.1);
  }

  .radio-icon {
    font-size: 32px;
  }

  .radio-text {
    flex: 1;
  }

  .radio-title {
    font-size: 15px;
    font-weight: 600;
    color: white;
    margin-bottom: 4px;
  }

  .radio-description {
    font-size: 13px;
    color: #999;
    line-height: 1.4;
  }

  .selective-options {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .label-text {
    display: block;
    font-size: 13px;
    font-weight: 500;
    color: #999;
    margin-bottom: 8px;
  }

  input[type="text"] {
    width: 100%;
    background: rgba(0, 0, 0, 0.3);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 6px;
    padding: 10px 12px;
    color: #ffffff;
    font-size: 14px;
    transition: all 0.2s;
  }

  input[type="text"]:focus {
    outline: none;
    border-color: #4F46E5;
    box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
  }

  .help-text {
    font-size: 12px;
    color: #666;
    font-style: italic;
  }

  .checkbox-option {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 12px;
    background: rgba(255, 100, 100, 0.1);
    border: 1px solid rgba(255, 100, 100, 0.2);
    border-radius: 6px;
    cursor: pointer;
    font-size: 14px;
    color: #ffcccc;
  }

  .checkbox-option input[type="checkbox"] {
    cursor: pointer;
  }

  .btn-import {
    background: linear-gradient(135deg, #10B981 0%, #059669 100%);
    border: none;
    border-radius: 8px;
    padding: 16px 24px;
    color: white;
    font-size: 16px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
    box-shadow: 0 4px 12px rgba(16, 185, 129, 0.3);
  }

  .btn-import:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 8px 20px rgba(16, 185, 129, 0.4);
  }

  .btn-import:disabled {
    opacity: 0.5;
    cursor: not-allowed;
    transform: none;
  }

  @media (max-width: 768px) {
    .mode-selector {
      grid-template-columns: 1fr;
    }
  }

  /* Warning Modal Styles */
  .warning-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.85);
    backdrop-filter: blur(8px);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 2000;
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

  .warning-modal {
    background: linear-gradient(135deg, rgba(30, 30, 30, 0.95) 0%, rgba(20, 20, 20, 0.95) 100%);
    border: 1px solid rgba(255, 165, 0, 0.3);
    border-radius: 16px;
    padding: 32px;
    width: 90%;
    max-width: 500px;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
  }

  .warning-header {
    text-align: center;
    margin-bottom: 24px;
  }

  .warning-icon {
    width: 80px;
    height: 80px;
    margin: 0 auto 16px;
    background: rgba(255, 165, 0, 0.15);
    border: 2px solid rgba(255, 165, 0, 0.3);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 48px;
  }

  .warning-header h3 {
    margin: 0;
    font-size: 24px;
    font-weight: 600;
    color: #ffa500;
  }

  .warning-content {
    margin-bottom: 32px;
  }

  .warning-message {
    font-size: 16px;
    color: #ffffff;
    line-height: 1.6;
    margin-bottom: 20px;
    text-align: center;
  }

  .warning-message strong {
    color: #ffa500;
    font-weight: 700;
  }

  .warning-details {
    background: rgba(0, 0, 0, 0.4);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 8px;
    padding: 16px;
    margin-bottom: 20px;
  }

  .detail-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 8px 0;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  }

  .detail-item:last-child {
    border-bottom: none;
  }

  .detail-label {
    font-size: 14px;
    color: #999;
  }

  .detail-value {
    font-size: 14px;
    color: #ffffff;
    font-weight: 600;
    text-align: right;
    word-break: break-all;
    max-width: 60%;
  }

  .warning-footer {
    font-size: 14px;
    color: #999;
    text-align: center;
    margin: 0;
    font-weight: 500;
  }

  .warning-footer strong {
    color: #ffa500;
    font-weight: 700;
  }

  .warning-actions {
    display: flex;
    gap: 12px;
  }

  .btn-cancel-warning {
    flex: 1;
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid rgba(255, 255, 255, 0.2);
    border-radius: 8px;
    padding: 14px 20px;
    color: white;
    font-size: 15px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
  }

  .btn-cancel-warning:hover {
    background: rgba(255, 255, 255, 0.1);
    transform: translateY(-1px);
  }

  .btn-confirm-delete {
    flex: 1;
    background: linear-gradient(135deg, #ea580c 0%, #c2410c 100%);
    border: none;
    border-radius: 8px;
    padding: 14px 20px;
    color: white;
    font-size: 15px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
    box-shadow: 0 4px 12px rgba(234, 88, 12, 0.3);
  }

  .btn-confirm-delete:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 20px rgba(234, 88, 12, 0.4);
  }
</style>
