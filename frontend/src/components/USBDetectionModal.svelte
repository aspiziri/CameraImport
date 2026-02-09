<script>
  import { createEventDispatcher } from 'svelte';

  export let device = null;
  export let show = false;

  const dispatch = createEventDispatcher();

  function handleYes() {
    dispatch('accept', device);
    show = false;
  }

  function handleNo() {
    dispatch('decline');
    show = false;
  }
</script>

{#if show && device}
  <div class="modal-overlay" on:click={handleNo}>
    <div class="modal" on:click|stopPropagation>
      <div class="modal-header">
        <h2>USB Device Detected</h2>
        <button class="close-btn" on:click={handleNo}>&times;</button>
      </div>

      <div class="modal-body">
        <div class="device-icon">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 2v6m0 0L9 5m3 3 3-3"/>
            <rect x="3" y="10" width="18" height="12" rx="2"/>
            <path d="M7 15h.01M7 18h.01M11 15h.01M11 18h.01"/>
          </svg>
        </div>

        <div class="device-info">
          <div class="info-row">
            <span class="label">Device Name:</span>
            <span class="value">{device.name}</span>
          </div>
          <div class="info-row">
            <span class="label">Location:</span>
            <span class="value">{device.path}</span>
          </div>
          <div class="info-row">
            <span class="label">Size:</span>
            <span class="value">{device.sizeGB}</span>
          </div>
        </div>

        <p class="question">
          Would you like to set this device as your source location?
        </p>
      </div>

      <div class="modal-footer">
        <button class="btn btn-secondary" on:click={handleNo}>
          No, Thanks
        </button>
        <button class="btn btn-primary" on:click={handleYes}>
          Yes, Use This Device
        </button>
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
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }

  .modal {
    background: linear-gradient(135deg, #2a2a2a 0%, #1f1f1f 100%);
    border-radius: 16px;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
    max-width: 500px;
    width: 90%;
    overflow: hidden;
    animation: slideUp 0.3s ease-out;
    border: 1px solid rgba(255, 255, 255, 0.1);
  }

  @keyframes slideUp {
    from {
      transform: translateY(20px);
      opacity: 0;
    }
    to {
      transform: translateY(0);
      opacity: 1;
    }
  }

  .modal-header {
    padding: 24px 28px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
    display: flex;
    justify-content: space-between;
    align-items: center;
    background: rgba(255, 255, 255, 0.02);
  }

  .modal-header h2 {
    margin: 0;
    font-size: 22px;
    font-weight: 600;
    color: #fff;
  }

  .close-btn {
    background: none;
    border: none;
    font-size: 32px;
    color: #888;
    cursor: pointer;
    padding: 0;
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 6px;
    transition: all 0.2s;
  }

  .close-btn:hover {
    background: rgba(255, 255, 255, 0.1);
    color: #fff;
  }

  .modal-body {
    padding: 28px;
  }

  .device-icon {
    width: 80px;
    height: 80px;
    margin: 0 auto 24px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 8px 24px rgba(102, 126, 234, 0.4);
  }

  .device-icon svg {
    width: 40px;
    height: 40px;
    color: #fff;
  }

  .device-info {
    background: rgba(255, 255, 255, 0.03);
    border-radius: 12px;
    padding: 20px;
    margin-bottom: 20px;
    border: 1px solid rgba(255, 255, 255, 0.05);
  }

  .info-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 10px 0;
  }

  .info-row:not(:last-child) {
    border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  }

  .label {
    color: #888;
    font-size: 14px;
    font-weight: 500;
  }

  .value {
    color: #fff;
    font-size: 14px;
    font-weight: 600;
    font-family: 'Courier New', monospace;
  }

  .question {
    text-align: center;
    color: #ccc;
    font-size: 15px;
    margin: 0;
    line-height: 1.5;
  }

  .modal-footer {
    padding: 20px 28px;
    border-top: 1px solid rgba(255, 255, 255, 0.1);
    display: flex;
    gap: 12px;
    justify-content: flex-end;
    background: rgba(0, 0, 0, 0.2);
  }

  .btn {
    padding: 12px 24px;
    border-radius: 8px;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    border: none;
    transition: all 0.2s;
    flex: 1;
  }

  .btn-secondary {
    background: rgba(255, 255, 255, 0.08);
    color: #ccc;
  }

  .btn-secondary:hover {
    background: rgba(255, 255, 255, 0.12);
    color: #fff;
  }

  .btn-primary {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: #fff;
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
  }

  .btn-primary:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 16px rgba(102, 126, 234, 0.5);
  }

  .btn:active {
    transform: translateY(0);
  }
</style>
