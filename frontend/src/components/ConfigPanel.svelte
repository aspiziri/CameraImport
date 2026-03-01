<script>
  import { createEventDispatcher } from 'svelte';
  import USBDetectionModal from './USBDetectionModal.svelte';

  export let config;

  const dispatch = createEventDispatcher();
  let scanning = false;
  let showAdvanced = false;
  let showDrivePicker = false;
  let availableDrives = [];
  let loadingDrives = false;

  async function selectFolder(type) {
    try {
      if (!window.go) {
        throw new Error('Wails runtime not available');
      }

      const title = type === 'source'
        ? 'Select Source Folder (SD Card)'
        : 'Select Destination Folder';

      const selected = await window.go.main.App.SelectFolder(title);

      if (selected) {
        if (type === 'source') {
          config.source = selected;
          updateConfig('source', selected);
        } else if (type === 'destination') {
          config.destination = selected;
          updateConfig('destination', selected);
        }
      }
    } catch (error) {
      console.error('Error selecting folder:', error);
      alert('Error selecting folder: ' + error);
    }
  }

  async function openDrivePicker() {
    showDrivePicker = true;
    await refreshDrives();
  }

  async function refreshDrives() {
    loadingDrives = true;
    try {
      availableDrives = await window.go.main.App.GetRemovableDrives() || [];
    } catch (error) {
      console.error('Error getting drives:', error);
      availableDrives = [];
    } finally {
      loadingDrives = false;
    }
  }

  function handleDriveSelect(event) {
    const drive = event.detail;
    config.source = drive.path;
    updateConfig('source', drive.path);
  }

  async function scanSource() {
    if (!config.source) {
      alert('Please select a source folder first');
      return;
    }

    scanning = true;
    dispatch('scanstart');
    console.log('Starting scan of:', config.source);

    try {
      if (!window.go) {
        throw new Error('Wails runtime not available');
      }

      const result = await window.go.importer.ImportService.ScanFiles(
        config.source,
        [...config.imgFormats, ...config.videoFormats, ...config.rawFormats]
      );

      console.log('Scan complete! Found files:', result);

      if (!result || result.length === 0) {
        alert('No media files found in the selected folder.');
      }

      dispatch('scan', result);
    } catch (error) {
      console.error('Scan error:', error);
      alert('Error scanning files: ' + error);
    } finally {
      scanning = false;
    }
  }

  function updateConfig(field, value) {
    dispatch('update', { [field]: value });
  }
</script>

<div class="config-panel">
  <div class="panel-header">
    <span class="step-badge">Step 1</span>
    <h2>Find Your Media</h2>
  </div>

  <div class="panel-content">
    <div class="config-section">
      <label>Source Folder</label>
      <div class="folder-input">
        <input
          type="text"
          bind:value={config.source}
          placeholder="E:\ or /media/sdcard"
          on:change={() => updateConfig('source', config.source)}
        />
        <button on:click={() => selectFolder('source')} class="btn-icon" title="Browse folders">📁</button>
        <button on:click={openDrivePicker} class="btn-icon" title="Select connected device">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="4" y="2" width="16" height="20" rx="2"/>
            <rect x="8" y="6" width="4" height="6" rx="1"/>
            <line x1="8" y1="16" x2="8.01" y2="16"/>
            <line x1="12" y1="16" x2="12.01" y2="16"/>
            <line x1="16" y1="16" x2="16.01" y2="16"/>
          </svg>
        </button>
      </div>
    </div>

    <div class="config-section">
      <label>Destination Folder</label>
      <div class="folder-input">
        <input
          type="text"
          bind:value={config.destination}
          placeholder="D:\Media\Camera Import"
          on:change={() => updateConfig('destination', config.destination)}
        />
        <button on:click={() => selectFolder('destination')} class="btn-icon">📁</button>
      </div>
    </div>

    <button
      class="btn-primary"
      on:click={scanSource}
      disabled={scanning || !config.source}
    >
      {scanning ? 'Scanning...' : '🔍 Scan Files'}
    </button>

    <div class="divider"></div>

    <div class="advanced-container">
      <button
        class="btn-advanced"
        on:click={() => showAdvanced = !showAdvanced}
      >
        {showAdvanced ? '▼' : '▶'} Advanced Settings
      </button>

      {#if showAdvanced}
        <div class="advanced-section">
          <div class="config-section">
            <label>Image Formats</label>
            <input
              type="text"
              bind:value={config.imgFormats}
              placeholder="jpg, png, gif"
              on:change={() => updateConfig('imgFormats', config.imgFormats)}
            />
          </div>

          <div class="config-section">
            <label>Video Formats</label>
            <input
              type="text"
              bind:value={config.videoFormats}
              placeholder="mp4, avi, mov"
              on:change={() => updateConfig('videoFormats', config.videoFormats)}
            />
          </div>

          <div class="config-section">
            <label>RAW Formats</label>
            <input
              type="text"
              bind:value={config.rawFormats}
              placeholder="arw, cr2, nef"
              on:change={() => updateConfig('rawFormats', config.rawFormats)}
            />
          </div>
        </div>
      {/if}
    </div>
  </div>
</div>

<USBDetectionModal
  bind:show={showDrivePicker}
  drives={availableDrives}
  loading={loadingDrives}
  on:select={handleDriveSelect}
  on:refresh={refreshDrives}
/>

<style>
  .config-panel {
    background: rgba(255, 255, 255, 0.05);
    border-radius: 12px;
    border: 1px solid rgba(255, 255, 255, 0.1);
    overflow: hidden;
    display: flex;
    flex-direction: column;
  }

  .panel-header {
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

  .panel-header h2 {
    margin: 0;
    font-size: 18px;
    font-weight: 600;
    color: #ffffff;
  }

  .panel-content {
    padding: 20px;
    display: flex;
    flex-direction: column;
    gap: 16px;
    overflow-y: auto;
  }

  .config-section {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  label {
    font-size: 13px;
    font-weight: 500;
    color: #999;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  input {
    background: rgba(0, 0, 0, 0.3);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 6px;
    padding: 10px 12px;
    color: #ffffff;
    font-size: 14px;
    transition: all 0.2s;
  }

  input:focus {
    outline: none;
    border-color: #4F46E5;
    box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
  }

  .folder-input {
    display: flex;
    gap: 8px;
  }

  .folder-input input {
    flex: 1;
  }

  .btn-icon {
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 6px;
    padding: 10px 14px;
    cursor: pointer;
    transition: all 0.2s;
    font-size: 16px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #ccc;
    flex-shrink: 0;
  }

  .btn-icon svg {
    width: 16px;
    height: 16px;
  }

  .btn-icon:hover {
    background: rgba(255, 255, 255, 0.1);
    color: #fff;
  }

  .btn-primary {
    background: linear-gradient(135deg, #4F46E5 0%, #7C3AED 100%);
    border: none;
    border-radius: 8px;
    padding: 12px 20px;
    color: white;
    font-size: 15px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
    margin-top: 8px;
  }

  .btn-primary:hover:not(:disabled) {
    transform: translateY(-1px);
    box-shadow: 0 8px 16px rgba(79, 70, 229, 0.3);
  }

  .btn-primary:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .divider {
    height: 1px;
    background: rgba(255, 255, 255, 0.1);
    margin: 8px 0;
  }

  .advanced-container {
    background: rgba(0, 0, 0, 0.2);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 8px;
    overflow: hidden;
  }

  .btn-advanced {
    background: transparent;
    border: none;
    border-radius: 0;
    padding: 12px 16px;
    color: #999;
    font-size: 13px;
    cursor: pointer;
    transition: all 0.2s;
    text-align: left;
    width: 100%;
    display: block;
  }

  .btn-advanced:hover {
    background: rgba(255, 255, 255, 0.05);
    color: white;
  }

  .advanced-section {
    display: flex;
    flex-direction: column;
    gap: 16px;
    padding: 16px;
    border-top: 1px solid rgba(255, 255, 255, 0.1);
    background: rgba(0, 0, 0, 0.1);
  }
</style>
