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

    dispatch('import', {
      mode: importMode,
      dates: selectedDates,
      folderName: importMode === 'selective' ? customFolderName : null,
      deleteAfter: deleteAfterImport
    });
  }

  $: canImport = files.length > 0 && !disabled;
  $: importCount = files.filter(f => {
    const fileDate = f.Date || f.date;
    return selectedDates.includes(fileDate);
  }).length;
</script>

<div class="import-controls">
  <div class="controls-header">
    <h2>Step 3) Pick Your Import Type</h2>
    <span class="import-count">
      {importCount} {importCount === 1 ? 'file' : 'files'} to import
    </span>
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

  .controls-header h2 {
    margin: 0;
    font-size: 18px;
    font-weight: 600;
    color: #ffffff;
  }

  .import-count {
    font-size: 14px;
    color: #999;
    padding: 4px 12px;
    background: rgba(79, 70, 229, 0.2);
    border-radius: 12px;
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
</style>
