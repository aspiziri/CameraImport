<script>
  const version = '4.0.0';
</script>

<header class="header">
  <div class="logo">
    <svg width="32" height="32" viewBox="0 0 32 32" fill="none">
      <rect width="32" height="32" rx="8" fill="#4F46E5"/>
      <path d="M16 10L8 14V22L16 26L24 22V14L16 10Z" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
      <path d="M8 14L16 18L24 14" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
      <path d="M16 18V26" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
    </svg>
    <h1>CameraImport</h1>
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
