<script>
  import { createEventDispatcher } from 'svelte';

  export let config;

  const dispatch = createEventDispatcher();
  let scanning = false;

  async function selectFolder(type) {
    try {
      // Check if window.go is available
      if (!window.go) {
        throw new Error('Wails runtime not available');
      }

      const title = type === 'source'
        ? 'Select Source Folder (SD Card)'
        : 'Select Destination Folder';

      // Call the Go backend method for folder selection
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

  async function scanSource() {
    if (!config.source) {
      alert('Please select a source folder first');
      return;
    }

    scanning = true;
    dispatch('scanstart');
    console.log('Starting scan of:', config.source);

    try {
      // Check if window.go is available (will be injected by Wails)
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
    <h2>Configuration</h2>
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
        <button on:click={() => selectFolder('source')} class="btn-icon">📁</button>
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
</div>

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
  }

  .btn-icon:hover {
    background: rgba(255, 255, 255, 0.1);
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
</style>
