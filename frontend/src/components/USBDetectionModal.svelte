<script>
  import { createEventDispatcher } from 'svelte';

  export let show = false;
  export let drives = [];
  export let loading = false;

  const dispatch = createEventDispatcher();

  function handleSelect(drive) {
    dispatch('select', drive);
    show = false;
  }

  function handleRefresh() {
    dispatch('refresh');
  }

  function handleClose() {
    show = false;
  }
</script>

{#if show}
  <div class="modal-overlay" on:click={handleClose}>
    <div class="modal" on:click|stopPropagation>
      <div class="modal-header">
        <h2>Select Source Device</h2>
        <button class="close-btn" on:click={handleClose}>&times;</button>
      </div>

      <div class="modal-body">
        {#if loading}
          <div class="state-message">
            <div class="spinner"></div>
            <p>Scanning for connected devices...</p>
          </div>
        {:else if drives.length === 0}
          <div class="state-message">
            <div class="empty-icon">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <rect x="3" y="6" width="18" height="14" rx="2"/>
                <path d="M7 6V4a1 1 0 0 1 1-1h8a1 1 0 0 1 1 1v2"/>
                <line x1="12" y1="11" x2="12" y2="15"/>
                <line x1="10" y1="13" x2="14" y2="13"/>
              </svg>
            </div>
            <p>No removable drives detected.</p>
            <p class="hint">Connect an SD card or USB drive, then click Refresh.</p>
          </div>
        {:else}
          <div class="drive-list">
            {#each drives as drive}
              <button class="drive-card" on:click={() => handleSelect(drive)}>
                <div class="drive-icon">
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                    <rect x="4" y="2" width="16" height="20" rx="2"/>
                    <rect x="8" y="6" width="4" height="6" rx="1"/>
                    <line x1="8" y1="16" x2="8.01" y2="16"/>
                    <line x1="12" y1="16" x2="12.01" y2="16"/>
                    <line x1="16" y1="16" x2="16.01" y2="16"/>
                  </svg>
                </div>
                <div class="drive-info">
                  <div class="drive-name">{drive.name}</div>
                  <div class="drive-path">{drive.path}</div>
                </div>
                <div class="drive-size">{drive.sizeGB}</div>
              </button>
            {/each}
          </div>
        {/if}
      </div>

      <div class="modal-footer">
        <button class="btn btn-secondary" on:click={handleRefresh} disabled={loading}>
          {loading ? 'Scanning...' : '↻ Refresh'}
        </button>
        <button class="btn btn-secondary" on:click={handleClose}>Cancel</button>
      </div>
    </div>
  </div>
{/if}

<style>
  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.75);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
    backdrop-filter: blur(4px);
    animation: fadeIn 0.2s ease-out;
  }

  @keyframes fadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
  }

  .modal {
    background: linear-gradient(135deg, #2a2a2a 0%, #1f1f1f 100%);
    border-radius: 16px;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
    width: 420px;
    max-width: 90vw;
    overflow: hidden;
    animation: slideUp 0.3s ease-out;
    border: 1px solid rgba(255, 255, 255, 0.1);
  }

  @keyframes slideUp {
    from { transform: translateY(20px); opacity: 0; }
    to { transform: translateY(0); opacity: 1; }
  }

  .modal-header {
    padding: 20px 24px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
    display: flex;
    justify-content: space-between;
    align-items: center;
    background: rgba(255, 255, 255, 0.02);
  }

  .modal-header h2 {
    margin: 0;
    font-size: 18px;
    font-weight: 600;
    color: #fff;
  }

  .close-btn {
    background: none;
    border: none;
    font-size: 28px;
    color: #888;
    cursor: pointer;
    padding: 0;
    width: 28px;
    height: 28px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 6px;
    transition: all 0.2s;
    line-height: 1;
  }

  .close-btn:hover {
    background: rgba(255, 255, 255, 0.1);
    color: #fff;
  }

  .modal-body {
    padding: 16px;
    min-height: 140px;
  }

  .state-message {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 24px 16px;
    gap: 10px;
    text-align: center;
  }

  .state-message p {
    margin: 0;
    color: #aaa;
    font-size: 14px;
  }

  .state-message .hint {
    color: #666;
    font-size: 12px;
  }

  .empty-icon {
    width: 48px;
    height: 48px;
    color: #555;
  }

  .empty-icon svg {
    width: 100%;
    height: 100%;
  }

  .spinner {
    width: 28px;
    height: 28px;
    border: 2px solid rgba(255, 255, 255, 0.1);
    border-top-color: #667eea;
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  .drive-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .drive-card {
    display: flex;
    align-items: center;
    gap: 14px;
    padding: 14px 16px;
    background: rgba(255, 255, 255, 0.04);
    border: 1px solid rgba(255, 255, 255, 0.08);
    border-radius: 10px;
    cursor: pointer;
    transition: all 0.15s;
    text-align: left;
    width: 100%;
  }

  .drive-card:hover {
    background: rgba(102, 126, 234, 0.15);
    border-color: rgba(102, 126, 234, 0.4);
  }

  .drive-icon {
    width: 36px;
    height: 36px;
    flex-shrink: 0;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #fff;
  }

  .drive-icon svg {
    width: 20px;
    height: 20px;
  }

  .drive-info {
    flex: 1;
    min-width: 0;
  }

  .drive-name {
    color: #fff;
    font-size: 14px;
    font-weight: 600;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .drive-path {
    color: #888;
    font-size: 12px;
    font-family: 'Courier New', monospace;
    margin-top: 2px;
  }

  .drive-size {
    color: #aaa;
    font-size: 12px;
    font-weight: 500;
    white-space: nowrap;
    flex-shrink: 0;
  }

  .modal-footer {
    padding: 16px 24px;
    border-top: 1px solid rgba(255, 255, 255, 0.1);
    display: flex;
    gap: 10px;
    justify-content: flex-end;
    background: rgba(0, 0, 0, 0.2);
  }

  .btn {
    padding: 9px 18px;
    border-radius: 8px;
    font-size: 13px;
    font-weight: 600;
    cursor: pointer;
    border: none;
    transition: all 0.2s;
  }

  .btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .btn-secondary {
    background: rgba(255, 255, 255, 0.08);
    color: #ccc;
  }

  .btn-secondary:hover:not(:disabled) {
    background: rgba(255, 255, 255, 0.12);
    color: #fff;
  }
</style>
