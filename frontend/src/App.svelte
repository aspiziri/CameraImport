<script>
  import { onMount } from 'svelte';
  import Header from './components/Header.svelte';
  import ConfigPanel from './components/ConfigPanel.svelte';
  import FilePreview from './components/FilePreview.svelte';
  import ImportControls from './components/ImportControls.svelte';
  import ProgressBar from './components/ProgressBar.svelte';

  let config = {
    source: '',
    destination: '',
    imgFormats: ['jpg', 'JPG', 'png', 'PNG', 'gif', 'GIF', 'tif', 'TIF', 'jpeg', 'JPEG'],
    videoFormats: ['mp4', 'MP4', 'avi', 'AVI', 'mov', 'MOV'],
    rawFormats: ['arw', 'ARW'],
    imgRelativePath: '/',
    videoRelativePath: '/Videos/',
    rawRelativePath: '/Capture/'
  };

  let files = [];
  let selectedDates = [];
  let importMode = 'full';
  let isImporting = false;
  let isScanning = false;
  let progress = 0;
  let progressData = null;
  let progressInterval = null;

  onMount(async () => {
    console.log('CameraImport UI loaded');

    // Load saved configuration
    try {
      const savedConfig = await window.go.main.App.LoadConfig();
      if (savedConfig) {
        config = savedConfig;
        console.log('Loaded saved config:', config);
      }
    } catch (error) {
      console.error('Error loading config:', error);
    }
  });

  async function handleConfigUpdate(event) {
    config = { ...config, ...event.detail };

    // Save config whenever it changes
    try {
      await window.go.main.App.SaveConfig(config);
      console.log('Config saved');
    } catch (error) {
      console.error('Error saving config:', error);
    }
  }

  function handleScanStart() {
    isScanning = true;
    files = [];
    selectedDates = [];
  }

  async function handleFilesScanned(event) {
    files = event.detail;
    isScanning = false;
    console.log('Files scanned:', files.length, 'files');
  }

  async function handleImportStart(event) {
    isImporting = true;
    importMode = event.detail.mode;
    selectedDates = event.detail.dates || [];
    progress = 0;
    progressData = null;

    try {
      // Start the import process
      const dates = event.detail.dates || [];
      const folderName = event.detail.folderName || '';

      console.log('Starting import with:', { source: config.source, dest: config.destination, dates, folderName });

      await window.go.importer.ImportService.RunImport(
        config.source,
        config.destination,
        dates,
        folderName,
        event.detail.deleteAfter
      );

      // Poll for progress updates
      progressInterval = setInterval(async () => {
        try {
          const data = await window.go.importer.ImportService.GetProgress();
          progressData = data;
          console.log('Progress data:', data);

          // Handle both capital and lowercase field names
          const percentage = data.Percentage || data.percentage || 0;
          const status = data.Status || data.status || 'unknown';

          progress = percentage;

          if (status === 'completed' || percentage >= 100) {
            console.log('Import completed!');
            clearInterval(progressInterval);
            progressInterval = null;
            handleImportComplete();
          } else if (status === 'error' || status === 'cancelled') {
            console.error('Import stopped:', status);
            clearInterval(progressInterval);
            progressInterval = null;
            if (status === 'error') {
              alert('Import failed. Check console for details.');
            }
            isImporting = false;
            progress = 0;
            progressData = null;
          }
        } catch (error) {
          console.error('Error getting progress:', error);
        }
      }, 500);

    } catch (error) {
      console.error('Import error:', error);
      alert('Import failed: ' + error);
      isImporting = false;
      progress = 0;
      progressData = null;
      if (progressInterval) {
        clearInterval(progressInterval);
      }
    }
  }

  async function handleCancel() {
    console.log('Cancelling import...');
    try {
      await window.go.importer.ImportService.CancelImport();
      if (progressInterval) {
        clearInterval(progressInterval);
        progressInterval = null;
      }
      isImporting = false;
      progress = 0;
      progressData = null;
    } catch (error) {
      console.error('Error cancelling import:', error);
    }
  }

  function handleImportComplete() {
    isImporting = false;
    progress = 100;
    setTimeout(() => {
      progress = 0;
      progressData = null;
    }, 1000);
  }

  function handleProgressUpdate(event) {
    progress = event.detail;
  }
</script>

<div class="app">
  <Header />

  <div class="main-content">
    <ConfigPanel
      {config}
      on:update={handleConfigUpdate}
      on:scanstart={handleScanStart}
      on:scan={handleFilesScanned}
    />

    <FilePreview
      {files}
      {config}
      bind:selectedDates
    />

    <ImportControls
      {files}
      {config}
      {selectedDates}
      disabled={isImporting}
      on:import={handleImportStart}
    />
  </div>

  {#if isImporting}
    <ProgressBar
      {progress}
      {progressData}
      on:complete={handleImportComplete}
      on:progress={handleProgressUpdate}
      on:cancel={handleCancel}
    />
  {/if}
</div>

<style>
  .app {
    display: flex;
    flex-direction: column;
    height: 100vh;
    background: linear-gradient(135deg, #1a1a1a 0%, #2d2d2d 100%);
  }

  .main-content {
    flex: 1;
    display: grid;
    grid-template-columns: 1fr 1fr 1fr;
    grid-template-rows: 1fr;
    gap: 20px;
    padding: 20px;
    overflow: hidden;
  }

  @media (max-width: 1024px) {
    .main-content {
      grid-template-columns: 1fr;
      grid-template-rows: auto 1fr auto;
    }
  }
</style>
