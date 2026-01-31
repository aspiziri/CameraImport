<script>
  const version = '0.1.0';
</script>

<header class="header">
  <div class="logo">
    <svg width="32" height="32" viewBox="0 0 32 32" fill="none">
      <rect width="32" height="32" rx="8" fill="#4F46E5"/>
      <rect x="7" y="12" width="18" height="13" rx="2" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
      <circle cx="16" cy="18.5" r="3.5" stroke="white" stroke-width="2"/>
      <path d="M12 12L13.5 8.5C13.7 8.1 14.1 8 14.5 8H17.5C17.9 8 18.3 8.1 18.5 8.5L20 12" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
      <circle cx="21.5" cy="14.5" r="0.75" fill="white"/>
    </svg>
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
