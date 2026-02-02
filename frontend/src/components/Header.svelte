<script>
  import cameraIcon from '../assets/camera-icon.png';
  const version = '0.1.0';
</script>

<header class="header">
  <div class="logo">
    <img src={cameraIcon} alt="Camera" class="logo-icon" />
    <h1>Camera Import</h1>
    <span class="version">v{version}</span>
  </div>
</header>

<style>
  .header {
    background: rgba(0, 0, 0, 0.3);
    backdrop-filter: blur(10px);
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
    padding: 16px 24px;
  }

  .logo {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .logo-icon {
    width: 32px;
    height: 32px;
  }

  h1 {
    font-size: 24px;
    font-weight: 600;
    color: #ffffff;
    margin: 0;
  }

  .version {
    font-size: 12px;
    color: #999;
    padding: 4px 8px;
    background: rgba(255, 255, 255, 0.05);
    border-radius: 4px;
  }
</style>
